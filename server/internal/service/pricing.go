package service

import (
	_ "embed"
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/mizanmahi/aiusage/types"
)

//go:embed litellm_prices.json
var litellmPricesJSON []byte

type litellmPrice struct {
	InputCostPerToken         float64 `json:"input_cost_per_token"`
	OutputCostPerToken        float64 `json:"output_cost_per_token"`
	CacheCreationCostPerToken float64 `json:"cache_creation_input_token_cost"`
	CacheReadCostPerToken     float64 `json:"cache_read_input_token_cost"`
}

type pricingResolver struct {
	prices  map[string]litellmPrice
	aliases map[string]string
}

var defaultPricing = newPricingResolver()

func calculateCost(event types.UsageEvent) float64 {
	cost, ok := defaultPricing.calculate(event)
	if !ok {
		slog.Warn("missing model pricing", "model", event.Model)
	}
	return cost
}

func newPricingResolver() *pricingResolver {
	prices := map[string]litellmPrice{}
	if err := json.Unmarshal(litellmPricesJSON, &prices); err != nil {
		slog.Error("failed to load embedded pricing", "error", err)
	}
	prices = validPrices(prices)

	return &pricingResolver{
		prices: prices,
		aliases: map[string]string{
			"gpt-5-codex":   "gpt-5",
			"gpt-5.3-codex": "gpt-5.2-codex",
		},
	}
}

func validPrices(prices map[string]litellmPrice) map[string]litellmPrice {
	valid := map[string]litellmPrice{}
	for model, price := range prices {
		if price.InputCostPerToken == 0 && price.OutputCostPerToken == 0 {
			continue
		}
		valid[model] = price
	}
	return valid
}

func (r *pricingResolver) calculate(event types.UsageEvent) (float64, bool) {
	price, ok := r.resolve(event.Model)
	if !ok {
		return 0, false
	}

	cacheCreateCost := price.CacheCreationCostPerToken
	if cacheCreateCost == 0 {
		cacheCreateCost = price.InputCostPerToken
	}

	cacheReadCost := price.CacheReadCostPerToken
	if cacheReadCost == 0 {
		cacheReadCost = price.InputCostPerToken
	}

	return float64(event.InputTokens)*price.InputCostPerToken +
		float64(event.OutputTokens)*price.OutputCostPerToken +
		float64(event.CacheCreateTokens)*cacheCreateCost +
		float64(event.CacheReadTokens)*cacheReadCost, true
}

func (r *pricingResolver) resolve(model string) (litellmPrice, bool) {
	model = strings.TrimSpace(model)
	if model == "" {
		return litellmPrice{}, false
	}

	if price, ok := r.prices[model]; ok {
		return price, true
	}

	if alias, ok := r.aliases[model]; ok {
		if price, found := r.prices[alias]; found {
			return price, true
		}
	}

	if family := familyFallback(model); family != model {
		if price, ok := r.prices[family]; ok {
			return price, true
		}
	}

	return litellmPrice{}, false
}

func familyFallback(model string) string {
	switch {
	case strings.HasPrefix(model, "gpt-5.5"):
		return "gpt-5"
	case strings.HasPrefix(model, "gpt-5.2-codex"):
		return "gpt-5.2-codex"
	case strings.HasPrefix(model, "claude-sonnet-4-5"):
		return "claude-sonnet-4-5"
	default:
		return model
	}
}

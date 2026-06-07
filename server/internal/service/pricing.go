package service

import "github.com/mizanmahi/aiusage/types"

type modelPrice struct {
	inputPerMillion  float64
	outputPerMillion float64
	cachePerMillion  float64
}

var modelPrices = map[string]modelPrice{
	"claude-sonnet-4-5": {inputPerMillion: 3.00, outputPerMillion: 15.00, cachePerMillion: 0.30},
	"gpt-5.5":           {inputPerMillion: 2.00, outputPerMillion: 8.00, cachePerMillion: 0.50},
}

func calculateCost(event types.UsageEvent) float64 {
	price, ok := modelPrices[event.Model]
	if !ok {
		return 0
	}

	return (float64(event.InputTokens)*price.inputPerMillion +
		float64(event.OutputTokens)*price.outputPerMillion +
		float64(event.CacheTokens)*price.cachePerMillion) / 1_000_000
}

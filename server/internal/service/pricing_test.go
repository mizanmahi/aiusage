package service

import (
	"testing"

	"github.com/mizanmahi/aiusage/types"
)

func TestPricingResolverExactMatch(t *testing.T) {
	resolver := pricingResolver{
		prices: map[string]litellmPrice{
			"exact-model": {
				InputCostPerToken:         0.000001,
				OutputCostPerToken:        0.000002,
				CacheCreationCostPerToken: 0.000003,
				CacheReadCostPerToken:     0.0000005,
			},
		},
	}

	cost, ok := resolver.calculate(types.UsageEvent{
		Model:             "exact-model",
		InputTokens:       1_000_000,
		OutputTokens:      2_000_000,
		CacheCreateTokens: 3_000_000,
		CacheReadTokens:   4_000_000,
	})
	if !ok {
		t.Fatal("calculate() ok = false, want true")
	}

	want := 16.0
	if cost != want {
		t.Fatalf("cost = %.2f, want %.2f", cost, want)
	}
}

func TestPricingResolverAliasFallback(t *testing.T) {
	resolver := pricingResolver{
		prices: map[string]litellmPrice{
			"target-model": {InputCostPerToken: 0.000001},
		},
		aliases: map[string]string{"alias-model": "target-model"},
	}

	cost, ok := resolver.calculate(types.UsageEvent{
		Model:       "alias-model",
		InputTokens: 1_000_000,
	})
	if !ok {
		t.Fatal("calculate() ok = false, want true")
	}
	if cost != 1 {
		t.Fatalf("cost = %.2f, want 1.00", cost)
	}
}

func TestPricingResolverFamilyFallback(t *testing.T) {
	resolver := pricingResolver{
		prices: map[string]litellmPrice{
			"gpt-5": {InputCostPerToken: 0.000002},
		},
	}

	cost, ok := resolver.calculate(types.UsageEvent{
		Model:       "gpt-5.5-preview",
		InputTokens: 1_000_000,
	})
	if !ok {
		t.Fatal("calculate() ok = false, want true")
	}
	if cost != 2 {
		t.Fatalf("cost = %.2f, want 2.00", cost)
	}
}

func TestPricingResolverUnknownModel(t *testing.T) {
	resolver := pricingResolver{prices: map[string]litellmPrice{}}

	cost, ok := resolver.calculate(types.UsageEvent{
		Model:       "unknown-model",
		InputTokens: 1_000_000,
	})
	if ok {
		t.Fatal("calculate() ok = true, want false")
	}
	if cost != 0 {
		t.Fatalf("cost = %.2f, want 0", cost)
	}
}

package domain

import "testing"

func TestLoanLimitInterestAndRanking(t *testing.T) {
	g, _ := NewGame("x", []string{"a", "b"})
	p := g.Players[0]
	for i := 0; i < MaxLoans; i++ {
		if err := g.TakeLoan(p.ID); err != nil {
			t.Fatal(err)
		}
	}
	if err := g.TakeLoan(p.ID); err != ErrLoanLimit {
		t.Fatalf("expected loan limit, got %v", err)
	}
	p.Cash = 50
	if err := g.SettleInterest(); err != nil {
		t.Fatalf("expected automatic interest loan, got %v", err)
	}
	if p.Loans != MaxLoans+1 {
		t.Fatalf("interest loan count=%d, want %d", p.Loans, MaxLoans+1)
	}
	g.Players[1].BrandAwareness = 3
	p.BrandAwareness = 3
	g.Players[1].Order = 0
	p.Order = 1
	if got := g.Rank()[0].ID; got != "b" {
		t.Fatalf("rank=%s", got)
	}
}

func TestFinalScoreConvertsCashToScoreUnits(t *testing.T) {
	p := Player{Cash: 250, Tableau: []Card{{ID: "score-product", Kind: "product"}, {ID: "score-value", Kind: "value"}}, TotalRevenue: 750, SelectedKPIs: []string{"products", "quality"}}
	got := p.Cash/CashScoreDivisor + p.metricScore()
	if got != 31 {
		t.Fatalf("final score=%d, want 31", got)
	}

	p.SelectedKPIs = []string{"surplus", "products"}
	if got := p.Cash/CashScoreDivisor + p.metricScore(); got != 36 {
		t.Fatalf("surplus final score=%d, want 36", got)
	}
}

func TestMetricScoreUsesRequestedCardAndRevenueValues(t *testing.T) {
	p := Player{
		Partner:       Card{ID: "partner", Function: "channel", ColorKey: "channel"},
		StarterShop:   Card{ID: "shop", Icons: []string{"coffee"}},
		Tableau:       []Card{{ID: "marketing-1", Kind: "marketing", Function: "marketing", BrandAwareness: 4}, {ID: "product-1", Kind: "product"}, {ID: "value-1", Kind: "value"}, {ID: "resource-1", Kind: "resource"}},
		RetainedCards: []Card{{ID: "marketing-2", Kind: "marketing", Function: "marketing", BrandAwareness: 3}},
		TotalRevenue:  90,
	}

	p.SelectedKPIs = []string{"channel", "awareness"}
	if got := p.metricScore(); got != 8 {
		t.Fatalf("channel and awareness score=%d, want 8", got)
	}
	p.SelectedKPIs = []string{"products", "quality"}
	if got := p.metricScore(); got != 6 {
		t.Fatalf("products and values score=%d, want 6", got)
	}
	p.SelectedKPIs = []string{"cost", "surplus"}
	if got := p.metricScore(); got != 3 {
		t.Fatalf("resources and surplus score=%d, want 3", got)
	}

	p.GourmetSatisfaction = 2
	p.RegularSatisfaction = 3
	p.SelectedKPIs = []string{"gourmet_satisfaction", "regular_satisfaction"}
	if got := p.metricScore(); got != 22 {
		t.Fatalf("satisfaction score=%d, want 22", got)
	}
}

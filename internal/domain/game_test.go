package domain

import (
	"testing"

	fixturedata "cafestartups/data"
)

func TestLoadCatalogRejectsDuplicateCards(t *testing.T) {
	if _, err := LoadCatalog([]byte(`{"cards":[{"id":"x","period":1,"cost":{"cash":0,"icons":[]}}, {"id":"x","period":1,"cost":{"cash":0,"icons":[]}}]}`)); err == nil {
		t.Fatal("expected duplicate card error")
	}
}

func TestCatalogDealsPeriodCards(t *testing.T) {
	catalog, err := LoadCatalog(fixturedata.MVPFixture)
	if err != nil {
		t.Fatal(err)
	}
	g, err := NewGame("catalog-seed", []string{"a", "b", "c", "d"})
	if err != nil {
		t.Fatal(err)
	}
	g.SetCatalog(catalog)
	for _, p := range g.Players {
		if err := g.SetKPIs(p.ID, "brand_awareness", "products"); err != nil {
			t.Fatal(err)
		}
	}
	if err := g.BeginExperiment(); err != nil {
		t.Fatal(err)
	}
	for _, p := range g.Players {
		for _, card := range p.Hand {
			if card.Period != PeriodOne || card.Name == "Management Card" {
				t.Fatalf("unexpected dealt card: %+v", card)
			}
		}
	}
}

func gameForTest(t *testing.T) *Game {
	t.Helper()
	g, err := NewGame("fixed-seed", []string{"a", "b", "c", "d"})
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range g.Players {
		if err := g.SetKPIs(p.ID, "brand_awareness", "products"); err != nil {
			t.Fatal(err)
		}
	}
	if err := g.BeginExperiment(); err != nil {
		t.Fatal(err)
	}
	return g
}

func TestDraftEndsWithOneCardAndPeriodDirection(t *testing.T) {
	g := gameForTest(t)
	for round := 0; round < 6; round++ {
		for _, p := range g.Players {
			if err := g.SelectCard(p.ID, p.Hand[0].ID); err != nil {
				t.Fatal(err)
			}
			if err := g.DiscardSelectedCard(p.ID); err != nil {
				t.Fatal(err)
			}
		}
		if err := g.PassHands(); err != nil {
			t.Fatal(err)
		}
	}
	if g.Phase != PhaseLearning {
		t.Fatalf("phase=%s", g.Phase)
	}
	if g.Center.ID == "" {
		t.Fatal("expected final covered card in center")
	}
	for _, p := range g.Players {
		if len(p.Hand) != 1 {
			t.Fatalf("player %s has %d cards", p.ID, len(p.Hand))
		}
	}
}

func TestPlayingCardAppliesMetricAndMarketEffects(t *testing.T) {
	g, _ := NewGame("effects", []string{"a", "b"})
	p := g.Players[0]
	g.Phase = PhaseExperiment
	card := Card{ID: "marketing-card", Kind: "marketing", Cost: Cost{Cash: 10}, MarketChange: map[string]int{"gourmet": 2}}
	p.Hand = []Card{card}
	if err := g.SelectCard(p.ID, card.ID); err != nil {
		t.Fatal(err)
	}
	if err := g.PlaySelectedCard(p.ID); err != nil {
		t.Fatal(err)
	}
	if p.BrandAwareness != 1 || g.DemandBoard["gourmet"] != 2 {
		t.Fatalf("effects not applied: brand=%d demand=%d", p.BrandAwareness, g.DemandBoard["gourmet"])
	}
}

func TestCostAndMissingIcons(t *testing.T) {
	g, _ := NewGame("x", []string{"a", "b"})
	p := g.Players[0]
	p.Cash = 60
	c := Card{ID: "c", Cost: Cost{Cash: 30, Icons: []string{"coffee", "operations"}}}
	g.selected[p.ID] = c
	p.Hand = []Card{c}
	g.Phase = PhaseExperiment
	if err := g.PlaySelectedCard(p.ID); err != ErrInsufficientCash {
		t.Fatalf("expected cash error, got %v", err)
	}
	p.Cash = 100
	g.selected[p.ID] = c
	if err := g.PlaySelectedCard(p.ID); err != nil {
		t.Fatal(err)
	}
	if p.Cash != 30 {
		t.Fatalf("cash=%d, want 30", p.Cash)
	}
}

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
	if err := g.SettleInterest(); err != ErrInsufficientCash {
		t.Fatalf("expected interest error, got %v", err)
	}
	g.Players[1].BrandAwareness = 3
	p.BrandAwareness = 3
	g.Players[1].Order = 0
	p.Order = 1
	if got := g.Rank()[0].ID; got != "b" {
		t.Fatalf("rank=%s", got)
	}
}

func TestInterestCanTakeLoanWhenCashIsShort(t *testing.T) {
	g, _ := NewGame("x", []string{"a", "b"})
	p := g.Players[0]
	p.Loans, p.Cash = 1, 0
	if err := g.SettleInterest(); err != nil {
		t.Fatal(err)
	}
	if p.Loans != 2 || p.Cash != 40 {
		t.Fatalf("loans=%d cash=%d", p.Loans, p.Cash)
	}
}

func TestSeedAndRevenueAreDeterministic(t *testing.T) {
	a := gameForTest(t)
	b := gameForTest(t)
	for i := range a.Players {
		if a.Players[i].Hand[0].ID != b.Players[i].Hand[0].ID {
			t.Fatal("same seed produced different deal")
		}
	}
	a.Players[0].Tableau = []Card{{Demand: map[string]int{"gourmet": 1}}}
	a.DistributeCustomers([]Customer{{Kind: "gourmet", Demand: "gourmet", Count: 2}, {Kind: "regular", Demand: "gourmet", Count: 1}})
	a.SettleRevenue()
	if a.Players[0].Revenue != 20 || a.Players[0].Cash != 170 {
		t.Fatalf("revenue=%d cash=%d", a.Players[0].Revenue, a.Players[0].Cash)
	}
}

func TestPeriodsAdvanceThreeTimes(t *testing.T) {
	g := gameForTest(t)
	for round := 0; round < 6; round++ {
		for _, p := range g.Players {
			if err := g.SelectCard(p.ID, p.Hand[0].ID); err != nil {
				t.Fatal(err)
			}
			if err := g.DiscardSelectedCard(p.ID); err != nil {
				t.Fatal(err)
			}
		}
		if err := g.PassHands(); err != nil {
			t.Fatal(err)
		}
	}
	if err := g.AdvancePeriod(); err != nil {
		t.Fatal(err)
	}
	if g.Period != PeriodTwo || g.Phase != PhaseHypothesis {
		t.Fatalf("period=%d phase=%s", g.Period, g.Phase)
	}
	for _, p := range g.Players {
		if err := g.SetKPIs(p.ID, "brand_awareness", "products"); err != nil {
			t.Fatal(err)
		}
	}
	if err := g.BeginExperiment(); err != nil {
		t.Fatal(err)
	}
}

func TestResolveLearningCompletesThreePeriods(t *testing.T) {
	g := gameForTest(t)
	for period := 1; period <= 3; period++ {
		for round := 0; round < 6; round++ {
			for _, p := range g.Players {
				if err := g.SelectCard(p.ID, p.Hand[0].ID); err != nil {
					t.Fatal(err)
				}
				if err := g.DiscardSelectedCard(p.ID); err != nil {
					t.Fatal(err)
				}
			}
			if err := g.PassHands(); err != nil {
				t.Fatal(err)
			}
		}
		if err := g.ResolveLearning(); err != nil {
			t.Fatal(err)
		}
		if period < 3 {
			if g.Phase != PhaseHypothesis || int(g.Period) != period+1 {
				t.Fatalf("after period %d: period=%d phase=%s", period, g.Period, g.Phase)
			}
			for _, p := range g.Players {
				if err := g.SetKPIs(p.ID, "brand_awareness", "products"); err != nil {
					t.Fatal(err)
				}
			}
			if err := g.BeginExperiment(); err != nil {
				t.Fatal(err)
			}
		}
	}
	if g.Phase != PhaseFinished || g.Period != PeriodThree {
		t.Fatalf("final state: period=%d phase=%s", g.Period, g.Phase)
	}
	for _, p := range g.Players {
		if p.Score <= 0 {
			t.Fatalf("player %s has no final score", p.ID)
		}
	}
}

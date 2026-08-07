package domain

import "testing"

func TestDemandCardsRevealByRoundAndScoreCustomerTypesSeparately(t *testing.T) {
	g := gameForTest(t)
	if !g.DemandCards["gourmet"][0].Revealed || g.DemandCards["gourmet"][1].Revealed {
		t.Fatal("expected only the initial demand position to be revealed")
	}
	for _, kind := range []string{"gourmet", "regular"} {
		for index, card := range g.DemandCards[kind] {
			if card.Position != index {
				t.Fatalf("%s demand card %s position=%d, want %d", kind, card.ID, card.Position, index)
			}
		}
	}
	g.Round = 1
	g.revealDemandCards()
	if !g.DemandCards["regular"][1].Revealed {
		t.Fatal("expected the demand position to reveal after the round ends")
	}
	g.DemandCards = map[string][]DemandCard{
		"gourmet": {{Position: 0, Icons: []string{"coffee"}, Revealed: true}},
		"regular": {{Position: 0, Icons: []string{"taste"}, Revealed: true}},
	}
	g.Players[0].Tableau = []Card{{Icons: []string{"coffee"}}}
	g.updateSatisfactionScores()
	if g.Players[0].GourmetSatisfaction != 1 || g.Players[0].RegularSatisfaction != 0 {
		t.Fatalf("satisfaction gourmet=%d regular=%d, want 1 and 0", g.Players[0].GourmetSatisfaction, g.Players[0].RegularSatisfaction)
	}
}

func TestDemandCardsUseProductAndValueIconsOnly(t *testing.T) {
	allowed := map[string]bool{"coffee": true, "dessert": true, "beans": true, "taste": true, "service": true, "value": true}
	for _, icon := range demandIconTypes {
		if !allowed[icon] {
			t.Fatalf("unexpected demand icon %q", icon)
		}
	}
	if len(demandIconTypes) != 6 {
		t.Fatalf("demand icon types=%d, want 6", len(demandIconTypes))
	}
}

func TestRevenueUsesSatisfiedDemandValuePerCustomerType(t *testing.T) {
	g := gameForTest(t)
	p := g.Players[0]
	p.Tableau = []Card{{Icons: []string{"coffee", "coffee", "taste"}, Demand: map[string]int{"gourmet": 1}}}
	g.DemandCards = map[string][]DemandCard{
		"gourmet": {{Position: 0, Icons: []string{"coffee"}, Revealed: true}, {Position: 1, Icons: []string{"coffee"}, Revealed: true}},
		"regular": {{Position: 0, Icons: []string{"taste"}, Revealed: true}},
	}
	p.Customers = []Customer{{Kind: "gourmet", Count: 2}, {Kind: "regular", Count: 1}}
	g.SettleRevenue()
	if p.Revenue != 80 || p.TotalRevenue != 80 || p.cashFlowGourmetRevenue != 60 || p.cashFlowRegularRevenue != 20 {
		t.Fatalf("revenue=%d total=%d gourmet=%d regular=%d, want 80, 80, 60, 20", p.Revenue, p.TotalRevenue, p.cashFlowGourmetRevenue, p.cashFlowRegularRevenue)
	}
}

func TestChannelCustomerCountsAreIncludedInRevenueCustomers(t *testing.T) {
	g := gameForTest(t)
	p := g.Players[0]
	p.Partner = Card{}
	p.StarterShop = Card{}
	p.Tableau = []Card{
		{Kind: "channel", CustomerCount: map[string]int{"regular": 2}},
		{Kind: "channel", CustomerCount: map[string]int{"regular": 1}},
	}
	g.DistributeCustomers(nil)
	if len(p.Customers) != 2 || p.Customers[0].Count != 2 || p.Customers[1].Count != 1 {
		t.Fatalf("channel customers=%+v, want regular counts 2 and 1", p.Customers)
	}
}

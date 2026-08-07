package domain

import "testing"

func TestDraftEndsWithOneCardAndPeriodDirection(t *testing.T) {
	g := gameForTest(t)
	for round := InitialRound; round < ExperimentRounds; round++ {
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
		if g.Round != round+1 {
			t.Fatalf("after draft round %d: round=%d", round, g.Round)
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

func TestRoundEndRevealsThatRoundsDemandCard(t *testing.T) {
	g := gameForTest(t)
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
	if !g.DemandCards["gourmet"][1].Revealed || !g.DemandCards["regular"][1].Revealed {
		t.Fatal("expected the first completed round's demand cards to be revealed")
	}
}

func TestSelectCardCanReplaceBeforeAction(t *testing.T) {
	g := gameForTest(t)
	p := g.Players[0]
	first, second := p.Hand[0], p.Hand[1]

	if err := g.SelectCard(p.ID, first.ID); err != nil {
		t.Fatal(err)
	}
	if err := g.SelectCard(p.ID, second.ID); err != nil {
		t.Fatalf("reselecting before action should succeed: %v", err)
	}
	if g.selected[p.ID].ID != second.ID {
		t.Fatalf("selected card=%s, want %s", g.selected[p.ID].ID, second.ID)
	}
}

func TestPlayingCardAppliesMetricAndMarketEffects(t *testing.T) {
	g, _ := NewGame("effects", []string{"a", "b"})
	p := g.Players[0]
	g.Phase = PhaseExperiment
	card := Card{ID: "marketing-card", Kind: "marketing", Icons: []string{"marketing"}, Cost: Cost{Cash: 10}, MarketChange: map[string]int{"gourmet": 2}}
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

func TestPlayingCardAddsOneMetricPerIcon(t *testing.T) {
	g, _ := NewGame("multi-icon-effects", []string{"a", "b"})
	p := g.Players[0]
	g.Phase = PhaseExperiment
	card := Card{ID: "multi-icon-card", Kind: "resource", Icons: []string{"operations", "coffee", "taste", "marketing"}, Cost: Cost{Cash: 10}}
	p.Hand = []Card{card}
	if err := g.SelectCard(p.ID, card.ID); err != nil {
		t.Fatal(err)
	}
	if err := g.PlaySelectedCard(p.ID); err != nil {
		t.Fatal(err)
	}
	if p.Resources != 1 || p.Products != 1 || p.Values != 1 || p.BrandAwareness != 1 {
		t.Fatalf("metrics resources=%d products=%d values=%d awareness=%d, want 1 each", p.Resources, p.Products, p.Values, p.BrandAwareness)
	}
	if p.IconValues["operations"] != 1 || p.IconValues["coffee"] != 1 || p.IconValues["taste"] != 1 || p.IconValues["marketing"] != 1 || len(p.IconValues) != 4 {
		t.Fatalf("icon values=%v, want one count for each of the four icons", p.IconValues)
	}
}

func TestMarketingCardUsesPrintedStarValue(t *testing.T) {
	g, _ := NewGame("marketing-stars", []string{"a", "b"})
	p := g.Players[0]
	g.Phase = PhaseExperiment
	card := Card{ID: "five-star-card", Kind: "marketing", Icons: []string{"marketing"}, BrandAwareness: 5, Cost: Cost{Cash: 10}}
	p.Hand = []Card{card}
	if err := g.SelectCard(p.ID, card.ID); err != nil {
		t.Fatal(err)
	}
	if err := g.PlaySelectedCard(p.ID); err != nil {
		t.Fatal(err)
	}
	if p.BrandAwareness != 5 || p.IconValues["marketing"] != 5 {
		t.Fatalf("marketing values awareness=%d icons=%v, want 5", p.BrandAwareness, p.IconValues)
	}
}

func TestPlayingCardCostUsesInitialCardsAndCountsDuplicateIcons(t *testing.T) {
	g, err := NewGame("cost-icons", []string{"a", "b"})
	if err != nil {
		t.Fatal(err)
	}
	g.Phase = PhaseExperiment
	p := g.Players[0]
	p.Cash = 100
	p.Partner = Card{Icons: []string{"coffee"}}
	p.StarterShop = Card{Icons: []string{"coffee"}}
	p.Tableau = []Card{{Icons: []string{"operations"}}}

	card := Card{ID: "cost-card", Kind: "resource", Cost: Cost{Cash: 10, Icons: []string{"coffee", "coffee", "beans"}}}
	p.Hand = []Card{card}
	if err := g.SelectCard(p.ID, card.ID); err != nil {
		t.Fatal(err)
	}
	if err := g.PlaySelectedCard(p.ID); err != nil {
		t.Fatal(err)
	}
	if p.Cash != 70 {
		t.Fatalf("cash=%d, want 70 after $10 base cost and one $20 missing icon", p.Cash)
	}
}

func TestCostAndMissingIcons(t *testing.T) {
	g, _ := NewGame("x", []string{"a", "b"})
	p := g.Players[0]
	p.Partner = Card{}
	p.StarterShop = Card{}
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

func TestDiscardedCardIsNotPassedToNextPlayer(t *testing.T) {
	g, err := NewGame("discard-pass", []string{"a", "b"})
	if err != nil {
		t.Fatal(err)
	}
	g.Phase = PhaseExperiment
	g.Players[0].Hand = []Card{{ID: "a-selected"}, {ID: "a-kept"}}
	g.Players[1].Hand = []Card{{ID: "b-selected"}, {ID: "b-kept"}}

	for _, p := range g.Players {
		if err := g.SelectCard(p.ID, p.Hand[0].ID); err != nil {
			t.Fatal(err)
		}
		if err := g.DiscardSelectedCard(p.ID); err != nil {
			t.Fatal(err)
		}
	}
	passed, err := g.PassHandsIfReady()
	if err != nil {
		t.Fatal(err)
	}
	if !passed || g.Round != 1 {
		t.Fatalf("round was not passed automatically: passed=%v round=%d", passed, g.Round)
	}
	wantHands := map[string]string{"a": "b-kept", "b": "a-kept"}
	for _, p := range g.Players {
		if len(p.Hand) != 1 || p.Hand[0].ID != wantHands[p.ID] {
			t.Fatalf("player %s received discarded card or wrong hand: %+v", p.ID, p.Hand)
		}
		if len(p.Discard) != 1 || p.Discard[0].ID != p.ID+"-selected" {
			t.Fatalf("player %s discard pile=%+v", p.ID, p.Discard)
		}
	}
}

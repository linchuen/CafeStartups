package domain

import "testing"

func TestExperimentStartsAtRoundZero(t *testing.T) {
	g := gameForTest(t)
	if err := g.SetKPIs(g.Players[0].ID, "brand_awareness", "products"); err != ErrInvalidAction {
		t.Fatalf("expected KPI selection to be locked during period one, got %v", err)
	}
	if g.Round != InitialRound {
		t.Fatalf("round=%d, want initial round %d", g.Round, InitialRound)
	}
	if len(g.Players[0].Hand) != 7 {
		t.Fatalf("initial hand=%d, want 7", len(g.Players[0].Hand))
	}
}

func TestInitialCardsCanBeSelectedInLobby(t *testing.T) {
	g, err := NewGame("setup-seed", []string{"a", "b"})
	if err != nil {
		t.Fatal(err)
	}
	if g.Period != PeriodZero {
		t.Fatalf("new game period=%d, want period zero", g.Period)
	}
	if err := g.SetInitialCards("a", "partner-service", "starter-station"); err != nil {
		t.Fatal(err)
	}
	if g.Players[0].Partner.ID != "partner-service" || g.Players[0].StarterShop.ID != "starter-station" {
		t.Fatalf("selected cards not saved: partner=%s shop=%s", g.Players[0].Partner.ID, g.Players[0].StarterShop.ID)
	}
	if err := g.SetInitialCards("b", "partner-service", "starter-songshan"); err != nil {
		t.Fatalf("expected duplicate partner selection to be allowed, got %v", err)
	}
}

func TestInitialCardsCannotBeChangedAfterPeriodOne(t *testing.T) {
	g, err := NewGame("setup-period-seed", []string{"a", "b"})
	if err != nil {
		t.Fatal(err)
	}
	g.Period = PeriodTwo
	if err := g.SetInitialCards("a", "partner-service", "starter-station"); err != ErrInvalidAction {
		t.Fatalf("expected initial card selection to be locked after period one, got %v", err)
	}
}

func TestBeginExperimentMovesInitialSetupToPeriodOne(t *testing.T) {
	g, err := NewGame("begin-period-seed", []string{"a", "b"})
	if err != nil {
		t.Fatal(err)
	}
	if err := g.BeginExperiment(); err != nil {
		t.Fatal(err)
	}
	if g.Period != PeriodOne || g.Phase != PhaseExperiment {
		t.Fatalf("after initial setup period=%d phase=%s, want period one experiment", g.Period, g.Phase)
	}
	if g.Players[0].Cash != InitialCash-g.Players[0].StarterShop.Cost.Cash {
		t.Fatalf("starter shop cost was not deducted: cash=%d shop cost=%d", g.Players[0].Cash, g.Players[0].StarterShop.Cost.Cash)
	}
}

func TestChannelPartnerAddsInitialCash(t *testing.T) {
	g, err := NewGame("partner-cash-seed", []string{"a", "b"})
	if err != nil {
		t.Fatal(err)
	}
	if err := g.SetInitialCards("a", "partner-community", "starter-station"); err != nil {
		t.Fatal(err)
	}
	if err := g.SetInitialCards("b", "partner-service", "starter-songshan"); err != nil {
		t.Fatal(err)
	}
	if err := g.BeginExperiment(); err != nil {
		t.Fatal(err)
	}
	want := InitialCash + 30 - g.Players[0].StarterShop.Cost.Cash
	if g.Players[0].Cash != want {
		t.Fatalf("cash=%d, want %d after initial cash bonus and shop cost", g.Players[0].Cash, want)
	}
}

func TestChannelPartnerAddsConfiguredCustomers(t *testing.T) {
	p := &Player{
		Partner: Card{Kind: "partner", Function: "channel", CustomerCount: map[string]int{"gourmet": 0, "regular": 1}},
		Tableau: []Card{{Demand: map[string]int{"regular": 1}}},
	}
	g := &Game{Players: []*Player{p}}
	g.DistributeCustomers(nil)
	if len(p.Customers) != 1 || p.Customers[0].Kind != "regular" || p.Customers[0].Count != 1 {
		t.Fatalf("partner customers=%+v, want one regular customer", p.Customers)
	}
}

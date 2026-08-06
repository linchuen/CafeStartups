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

func TestCatalogPreservesExplicitCostIcons(t *testing.T) {
	catalog, err := LoadCatalog(fixturedata.MVPFixture)
	if err != nil {
		t.Fatal(err)
	}

	for _, card := range catalog {
		if len(card.Cost.Icons) == 0 {
			continue
		}
		for index, icon := range card.Cost.Icons {
			if icon == "" {
				t.Fatalf("card %q cost icon[%d] is empty", card.ID, index)
			}
		}
	}
}

func TestCatalogKeepsBlankCostIconsBlank(t *testing.T) {
	catalog, err := LoadCatalog([]byte(`{"cards":[{"id":"multi-icon","name":"Multi icon","kind":"product","period":3,"cost":{"cash":35,"icons":[]},"icons":["coffee","operations"],"marketChange":{},"source":"mvp-fixture"}]}`))
	if err != nil {
		t.Fatal(err)
	}

	got := catalog[0].Cost.Icons
	if len(got) != 0 {
		t.Fatalf("cost icons=%v, want blank cost", got)
	}
}

func TestCatalogRejectsMarketingAndChannelCostIcons(t *testing.T) {
	for _, icon := range []string{"marketing", "channel"} {
		fixture := []byte(`{"cards":[{"id":"invalid-cost","name":"Invalid cost","kind":"resource","period":1,"cost":{"cash":10,"icons":["` + icon + `"]},"icons":["operations"],"marketChange":{},"source":"mvp-fixture"}]}`)
		if _, err := LoadCatalog(fixture); err == nil {
			t.Fatalf("expected %q cost icon to be rejected", icon)
		}
	}
}

func gameForTest(t *testing.T) *Game {
	t.Helper()
	g, err := NewGame("fixed-seed", []string{"a", "b", "c", "d"})
	if err != nil {
		t.Fatal(err)
	}
	if err := g.BeginExperiment(); err != nil {
		t.Fatal(err)
	}
	return g
}

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
	// Two coffee icons are covered by the partner and starter shop; beans is missing.
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
	if a.Players[0].Revenue != 40 || a.Players[0].Cash != InitialCash-a.Players[0].StarterShop.Cost.Cash+40 {
		t.Fatalf("revenue=%d cash=%d", a.Players[0].Revenue, a.Players[0].Cash)
	}
}

func TestPeriodsAdvanceThreeTimes(t *testing.T) {
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
	}
	if !g.MarketBagReady || g.MarketBag["difficult"] != 1 {
		t.Fatalf("market bag was not prepared automatically: ready=%t bag=%v", g.MarketBagReady, g.MarketBag)
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
		}
		if err := g.DrawMarket(); err != nil {
			t.Fatal(err)
		}
		if err := g.ResolveLearning(); err != nil {
			t.Fatal(err)
		}
		if period < 3 {
			if g.Phase != PhaseHypothesis || int(g.Period) != period+1 {
				t.Fatalf("after period %d: period=%d phase=%s", period, g.Period, g.Phase)
			}
			kpis := []string{"brand_awareness", "products"}
			if period == 2 {
				kpis = []string{"values", "resources"}
			}
			for _, p := range g.Players {
				if err := g.SetKPIs(p.ID, kpis...); err != nil {
					t.Fatal(err)
				}
			}
			if period == 2 {
				if err := g.SetKPIs(g.Players[0].ID, "brand_awareness", "products"); err != ErrInvalidAction {
					t.Fatalf("expected only one KPI reselection after period two, got %v", err)
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
		if len(p.CashFlow) != 3 {
			t.Fatalf("player %s cash-flow statements=%d, want 3", p.ID, len(p.CashFlow))
		}
		if len(p.CashFlowRounds) != 3*ExperimentRounds {
			t.Fatalf("player %s round cash-flow statements=%d, want %d", p.ID, len(p.CashFlowRounds), 3*ExperimentRounds)
		}
		for index, statement := range p.CashFlow {
			if int(statement.Period) != index+1 {
				t.Fatalf("player %s statement %d has period %d", p.ID, index, statement.Period)
			}
			calculated := statement.BeginningCash + statement.OperatingRevenue + statement.OtherIncome - statement.OperatingExpenses - statement.InterestPaid - statement.PrincipalRepayment + statement.NewLoans
			if calculated != statement.EndingCash {
				t.Fatalf("player %s period %d ending cash=%d, calculated=%d", p.ID, statement.Period, statement.EndingCash, calculated)
			}
		}
	}
}

func TestDrawMarketUsesReferencePeriodRules(t *testing.T) {
	want := map[Period][]int{PeriodOne: {3, 2, 1, 1}, PeriodTwo: {4, 3, 2, 1}, PeriodThree: {5, 3, 2, 1}}
	for period, expected := range want {
		g := gameForTest(t)
		g.Period = period
		g.Phase = PhaseLearning
		if err := g.DrawMarket(); err != nil {
			t.Fatal(err)
		}
		for index, count := range expected {
			if g.MarketRanking[index] != count {
				t.Fatalf("period %d rank %d=%d, want %d", period, index+1, g.MarketRanking[index], count)
			}
		}
	}
}

func TestDrawMarketUsesFinalHandsToBuildBag(t *testing.T) {
	g := gameForTest(t)
	g.Phase = PhaseLearning
	g.Players[0].Tableau = []Card{{MarketChange: map[string]int{"gourmet": 99}}}
	g.Players[0].Hand = []Card{{MarketChange: map[string]int{"gourmet": 2}}}
	g.Players[1].Hand = []Card{{MarketChange: map[string]int{"regular": 1}}}
	g.Players[2].Hand = []Card{{MarketChange: map[string]int{"difficult": -5}}}
	g.Players[3].Hand = []Card{{MarketChange: map[string]int{}}}

	if err := g.PrepareMarketBag(); err != nil {
		t.Fatal(err)
	}
	if total := g.MarketBag["gourmet"] + g.MarketBag["regular"] + g.MarketBag["difficult"]; total != 4 {
		t.Fatalf("prepared bag=%d, want 4 including one difficult customer", total)
	}
	if err := g.DrawMarket(); err != nil {
		t.Fatal(err)
	}
	if total := g.MarketBag["gourmet"] + g.MarketBag["regular"] + g.MarketBag["difficult"]; total != 0 {
		t.Fatalf("remaining bag=%d, want 0 after drawing all four customers", total)
	}
	totalDrawn := 0
	for _, draw := range g.MarketDraws {
		totalDrawn += draw.Total
	}
	if totalDrawn != 4 {
		t.Fatalf("drawn=%d, want 4 customers including the difficult customer", totalDrawn)
	}
}

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
	if p.Revenue != 80 || p.cashFlowGourmetRevenue != 60 || p.cashFlowRegularRevenue != 20 {
		t.Fatalf("revenue=%d gourmet=%d regular=%d, want 80, 60, 20", p.Revenue, p.cashFlowGourmetRevenue, p.cashFlowRegularRevenue)
	}
}

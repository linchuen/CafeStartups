package domain

import "testing"

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
		if err := g.SetKPIs(p.ID, "products", "brand_awareness"); err != nil {
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
			kpis := []string{"products", "brand_awareness"}
			if period == 2 {
				kpis = []string{"products", "resources"}
			}
			for _, p := range g.Players {
				if err := g.SetKPIs(p.ID, kpis...); err != nil {
					t.Fatal(err)
				}
			}
			if period == 2 {
				if err := g.SetKPIs(g.Players[0].ID, "products", "brand_awareness"); err != ErrInvalidAction {
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
		if len(p.RetainedCards) != 3 {
			t.Fatalf("player %s retained cards=%d, want 3", p.ID, len(p.RetainedCards))
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
	if total := g.MarketBag["gourmet"] + g.MarketBag["regular"] + g.MarketBag["difficult"]; total != 13 {
		t.Fatalf("prepared bag=%d, want 13 including initial customers and one difficult customer", total)
	}
	if err := g.DrawMarket(); err != nil {
		t.Fatal(err)
	}
	if total := g.MarketBag["gourmet"] + g.MarketBag["regular"] + g.MarketBag["difficult"]; total != 6 {
		t.Fatalf("remaining bag=%d, want 6 after drawing seven customers", total)
	}
	totalDrawn := 0
	for _, draw := range g.MarketDraws {
		totalDrawn += draw.Total
	}
	if totalDrawn != 7 {
		t.Fatalf("drawn=%d, want 7 customers", totalDrawn)
	}
}

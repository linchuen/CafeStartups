package domain

import "fmt"

var marketCustomersByPeriod = map[Period][]int{
	PeriodOne:   {3, 2, 1, 1},
	PeriodTwo:   {4, 3, 2, 1},
	PeriodThree: {5, 3, 2, 1},
}

// PrepareMarketBag fills the market bag from each player's final card without
// drawing. The full bag is exposed so players can verify how many customers
// will be available before the ranking draw.
func (g *Game) PrepareMarketBag() error {
	if g.Phase != PhaseLearning || g.marketDrawn || g.MarketBagReady {
		return ErrInvalidAction
	}
	counts := marketCustomersByPeriod[g.Period]
	if len(counts) < len(g.Players) {
		return fmt.Errorf("%w: market customer rule", ErrInvalidAction)
	}
	g.MarketRanking = append([]int(nil), counts[:len(g.Players)]...)
	g.MarketBag = map[string]int{"gourmet": 0, "regular": 0, "difficult": 0}
	// Every market draw starts with one difficult customer in the bag.
	g.MarketBag["difficult"] = 1
	for _, p := range g.Players {
		if len(p.Hand) == 0 {
			continue
		}
		for kind, count := range p.Hand[0].MarketChange {
			if _, ok := g.MarketBag[kind]; ok && count > 0 {
				g.MarketBag[kind] += count
			}
		}
	}
	g.MarketBagReady = true
	return nil
}

// DrawMarket assigns the reference-board draw amounts in ranking order.
func (g *Game) DrawMarket() error {
	if g.Phase != PhaseLearning || g.marketDrawn {
		return ErrInvalidAction
	}
	// Keep direct domain callers compatible while the UI uses the explicit
	// PREPARE_MARKET_BAG command first.
	if !g.MarketBagReady {
		if err := g.PrepareMarketBag(); err != nil {
			return err
		}
	}
	market := make([]string, 0)
	for _, kind := range []string{"gourmet", "regular", "difficult"} {
		for i := 0; i < g.MarketBag[kind]; i++ {
			market = append(market, kind)
		}
	}
	g.rng.Shuffle(len(market), func(i, j int) { market[i], market[j] = market[j], market[i] })
	g.MarketDraws = make([]MarketDraw, 0, len(g.Players))
	for i, p := range g.Players {
		countsByType := map[string]int{"gourmet": 0, "regular": 0, "difficult": 0}
		for draw := 0; draw < g.MarketRanking[i] && len(market) > 0; draw++ {
			kind := market[0]
			market = market[1:]
			countsByType[kind]++
			g.MarketBag[kind]--
		}
		total := countsByType["gourmet"] + countsByType["regular"] + countsByType["difficult"]
		g.MarketDraws = append(g.MarketDraws, MarketDraw{Rank: i + 1, PlayerID: p.ID, CustomerCounts: countsByType, Total: total})
	}
	g.marketDrawn = true
	return nil
}

// ResolveLearning settles the period after the market ranking has been drawn.
func (g *Game) ResolveLearning() error {
	if g.Phase != PhaseLearning {
		return ErrInvalidPhase
	}
	if !g.marketDrawn || len(g.MarketRanking) != len(g.Players) {
		return ErrInvalidAction
	}
	if err := g.SettleInterest(); err != nil {
		return err
	}
	customers := make([]Customer, 0, len(g.Players))
	for _, draw := range g.MarketDraws {
		for kind, count := range draw.CustomerCounts {
			if count <= 0 {
				continue
			}
			demand := ""
			if kind != "difficult" {
				demand = kind
			}
			customers = append(customers, Customer{Kind: kind, Demand: demand, Count: count})
		}
	}
	g.DistributeCustomers(customers)
	g.SettleRevenue()
	g.updateSatisfactionScores()
	g.recordCashFlow()
	g.recordCashFlowRound()
	for _, p := range g.Players {
		p.Score = p.Cash + p.metricScore()
	}
	if g.Period == PeriodThree {
		g.Phase = PhaseFinished
		return nil
	}
	g.Period++
	g.Round = 0
	g.Phase = PhaseHypothesis
	g.MarketRanking = nil
	g.MarketDraws = nil
	g.MarketBag = nil
	g.MarketBagReady = false
	g.marketDrawn = false
	g.selected = map[string]Card{}
	g.acted = map[string]bool{}
	return nil
}

// AdvancePeriod moves from learning to the next period after the current
// period's market, revenue and interest have been resolved.
func (g *Game) AdvancePeriod() error {
	if g.Phase != PhaseLearning {
		return ErrInvalidPhase
	}
	if g.Period == PeriodThree {
		g.Phase = PhaseFinished
		return nil
	}
	g.Period++
	g.Round = 0
	g.Phase = PhaseHypothesis
	g.selected = map[string]Card{}
	g.acted = map[string]bool{}
	return nil
}

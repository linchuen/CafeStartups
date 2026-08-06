package domain

import "math/rand"

func (g *Game) SetCatalog(cards []Card) { g.Catalog = append([]Card(nil), cards...) }

// NewGame creates a deterministic, offline game. Two to four players are required to start.
func NewGame(seed string, playerIDs []string) (*Game, error) {
	g := &Game{Seed: seed, Period: PeriodZero, Phase: PhaseHypothesis, PartnerOptions: MVPPartnerCards(), StarterShopOptions: MVPStarterShopCards(), DemandBoard: map[string]int{"gourmet": 0, "regular": 0, "difficult": 0}, selected: map[string]Card{}, acted: map[string]bool{}, rng: rand.New(rand.NewSource(seedValue(seed)))}
	for _, id := range playerIDs {
		if err := g.AddPlayer(id, id); err != nil {
			return nil, err
		}
	}
	return g, nil
}

func (g *Game) AddPlayer(id, displayName string) error {
	if g.Phase != PhaseHypothesis || len(g.Players) >= 4 || id == "" {
		return ErrInvalidAction
	}
	for _, p := range g.Players {
		if p.ID == id {
			return ErrInvalidAction
		}
	}
	p := &Player{ID: id, DisplayName: displayName, Cash: InitialCash, Order: len(g.Players)}
	if len(g.PartnerOptions) > 0 {
		p.Partner = g.PartnerOptions[len(g.Players)%len(g.PartnerOptions)]
	}
	if len(g.StarterShopOptions) > 0 {
		p.StarterShop = g.StarterShopOptions[len(g.Players)%len(g.StarterShopOptions)]
	}
	g.Players = append(g.Players, p)
	return nil
}

func (g *Game) SetInitialCards(playerID, partnerID, starterShopID string) error {
	if g.Phase != PhaseHypothesis || g.Period != PeriodZero || partnerID == "" || starterShopID == "" {
		return ErrInvalidAction
	}
	p, err := g.player(playerID)
	if err != nil {
		return err
	}
	partner, ok := findCard(g.PartnerOptions, partnerID)
	if !ok {
		return ErrCardNotFound
	}
	starterShop, ok := findCard(g.StarterShopOptions, starterShopID)
	if !ok {
		return ErrCardNotFound
	}
	p.Partner = partner
	p.StarterShop = starterShop
	p.InitialCardsSelected = true
	return nil
}

func findCard(cards []Card, id string) (Card, bool) {
	for _, card := range cards {
		if card.ID == id {
			return card, true
		}
	}
	return Card{}, false
}

func (g *Game) Start() error {
	if len(g.Players) < 2 {
		return ErrNotEnoughPlayers
	}
	if g.Phase != PhaseHypothesis {
		return ErrInvalidPhase
	}
	return nil
}

var validKPI = map[string]bool{
	"gourmet_satisfaction": true,
	"regular_satisfaction": true,
	"total_satisfaction":   true,
	"channel":              true,
	"awareness":            true,
	"products":             true,
	"quality":              true,
	"cost":                 true,
	"surplus":              true,
	"brand_awareness":      true,
	"values":               true,
	"resources":            true,
}

func (g *Game) SetKPIs(playerID string, kpis ...string) error {
	if g.Phase != PhaseHypothesis || g.Period <= PeriodOne || len(kpis) != 2 || kpis[0] == kpis[1] {
		return ErrInvalidAction
	}
	for _, kpi := range kpis {
		if !validKPI[kpi] {
			return ErrInvalidAction
		}
	}
	p, err := g.player(playerID)
	if err != nil {
		return err
	}
	if p.KPISelectionPeriod == g.Period {
		return ErrInvalidAction
	}
	p.SelectedKPIs = append([]string(nil), kpis...)
	p.KPISelectionPeriod = g.Period
	return nil
}

func (g *Game) BeginExperiment() error {
	if g.Phase != PhaseHypothesis || len(g.Players) < 2 {
		return ErrInvalidPhase
	}
	initialSetup := g.Period == PeriodZero
	if initialSetup {
		g.Period = PeriodOne
	}
	if g.Period != PeriodOne {
		for _, p := range g.Players {
			if len(p.SelectedKPIs) != 2 || p.KPISelectionPeriod != g.Period {
				return ErrInvalidAction
			}
		}
	}
	for _, p := range g.Players {
		if g.Period == PeriodOne && len(p.SelectedKPIs) != 0 {
			return ErrInvalidAction
		}
	}
	if initialSetup {
		for _, p := range g.Players {
			if p.Cash+partnerInitialCashBonus(p) < p.StarterShop.Cost.Cash {
				return ErrInsufficientCash
			}
		}
		for _, p := range g.Players {
			p.Cash += partnerInitialCashBonus(p)
			p.Cash -= p.StarterShop.Cost.Cash
		}
	}
	for _, p := range g.Players {
		p.cashFlowBeginning = p.Cash
		p.cashFlowRevenue = 0
		p.cashFlowGourmetCount = 0
		p.cashFlowRegularCount = 0
		p.cashFlowOtherIncome = 0
		p.cashFlowExpenses = 0
		p.cashFlowInterest = 0
		p.cashFlowPrincipal = 0
		p.cashFlowNewLoans = 0
	}
	g.Round, g.Phase = InitialRound, PhaseExperiment
	g.deal()
	return nil
}

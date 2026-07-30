package domain

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
)

type Period int

const (
	PeriodOne   Period = 1
	PeriodTwo   Period = 2
	PeriodThree Period = 3
)

type Phase string

const (
	PhaseHypothesis Phase = "hypothesis"
	PhaseExperiment Phase = "experiment"
	PhaseLearning   Phase = "learning"
	PhaseFinished   Phase = "finished"
)

const (
	InitialCash   = 150
	LoanAmount    = 50
	LoanInterest  = 10
	MaxLoans      = 6
	DiscardRefund = 20
)

var (
	ErrInvalidAction    = errors.New("INVALID_ACTION")
	ErrNotEnoughPlayers = errors.New("NOT_ENOUGH_PLAYERS")
	ErrNotYourTurn      = errors.New("NOT_YOUR_TURN")
	ErrInsufficientCash = errors.New("INSUFFICIENT_CASH")
	ErrLoanLimit        = errors.New("LOAN_LIMIT")
	ErrInvalidPhase     = errors.New("INVALID_PHASE")
	ErrCardNotFound     = errors.New("CARD_NOT_FOUND")
	ErrAlreadySelected  = errors.New("CARD_ALREADY_SELECTED")
)

type Cost struct {
	Cash  int
	Icons []string
}
type Card struct {
	ID, Name, Kind                              string
	Period                                      Period
	Cost                                        Cost
	Icons                                       []string
	BrandAwareness, Products, Values, Resources int
	Demand                                      map[string]int
}
type Player struct {
	ID, DisplayName        string
	Cash, Loans            int
	BrandAwareness         int
	Products               int
	Values                 int
	Resources              int
	Partner, StarterShop   Card
	Hand, Tableau, Discard []Card
	SelectedKPIs           []string
	Customers              []Customer
	Revenue, Score         int
	Order                  int
}
type Customer struct {
	Kind, Demand     string
	UnitPrice, Count int
}
type Game struct {
	Seed        string
	Period      Period
	Phase       Phase
	Round       int
	Players     []*Player
	Center      Card
	Market      []Card
	DemandBoard map[string]int
	selected    map[string]Card
	acted       map[string]bool
	rng         *rand.Rand
}

// NewGame creates a deterministic, offline game. Two to four players are required to start.
func NewGame(seed string, playerIDs []string) (*Game, error) {
	g := &Game{Seed: seed, Period: PeriodOne, Phase: PhaseHypothesis, DemandBoard: map[string]int{"gourmet": 0, "regular": 0, "difficult": 0}, selected: map[string]Card{}, acted: map[string]bool{}, rng: rand.New(rand.NewSource(seedValue(seed)))}
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
	p.Partner = Card{ID: "partner-" + id, Name: "Partner", Kind: "partner"}
	p.StarterShop = Card{ID: "starter-shop-" + id, Name: "Founding Shop", Kind: "starter_shop"}
	g.Players = append(g.Players, p)
	return nil
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

func (g *Game) SetKPIs(playerID string, kpis ...string) error {
	if g.Phase != PhaseHypothesis || len(kpis) != 2 || kpis[0] == kpis[1] {
		return ErrInvalidAction
	}
	p, err := g.player(playerID)
	if err != nil {
		return err
	}
	p.SelectedKPIs = append([]string(nil), kpis...)
	return nil
}

func (g *Game) BeginExperiment() error {
	if g.Phase != PhaseHypothesis || len(g.Players) < 2 {
		return ErrInvalidPhase
	}
	for _, p := range g.Players {
		if len(p.SelectedKPIs) != 2 {
			return ErrInvalidAction
		}
	}
	g.Round, g.Phase = 0, PhaseExperiment
	g.deal()
	return nil
}

func (g *Game) SelectCard(playerID, cardID string) error {
	if g.Phase != PhaseExperiment {
		return ErrInvalidPhase
	}
	p, err := g.player(playerID)
	if err != nil {
		return err
	}
	if _, ok := g.selected[playerID]; ok {
		return ErrAlreadySelected
	}
	for _, c := range p.Hand {
		if c.ID == cardID {
			g.selected[playerID] = c
			g.acted[playerID] = false
			return nil
		}
	}
	return ErrCardNotFound
}

// PassHands resolves one simultaneous drafting round. The selected card is played or
// discarded first, then the remaining six cards move to the next player.
func (g *Game) PassHands() error {
	if g.Phase != PhaseExperiment {
		return ErrInvalidPhase
	}
	if len(g.selected) != len(g.Players) {
		return ErrInvalidAction
	}
	for _, p := range g.Players {
		if !g.acted[p.ID] {
			return ErrInvalidAction
		}
	}
	remaining := make([][]Card, len(g.Players))
	for i, p := range g.Players {
		remaining[i] = append([]Card(nil), p.Hand...)
	}
	direction := 1
	if g.Period == PeriodTwo {
		direction = -1
	}
	for i, p := range g.Players {
		target := (i + direction + len(g.Players)) % len(g.Players)
		g.Players[target].Hand = append(g.Players[target].Hand[:0], remaining[i]...)
		_ = p
	}
	g.Round++
	g.selected, g.acted = map[string]Card{}, map[string]bool{}
	if g.Round == 6 {
		for _, p := range g.Players {
			if len(p.Hand) != 1 {
				return fmt.Errorf("%w: hand size", ErrInvalidAction)
			}
		}
		g.Phase = PhaseLearning
	}
	return nil
}

func (g *Game) PlaySelectedCard(playerID string) error {
	p, err := g.player(playerID)
	if err != nil {
		return err
	}
	if g.Phase != PhaseExperiment {
		return ErrInvalidPhase
	}
	c, ok := g.selected[playerID]
	if !ok {
		return ErrInvalidAction
	}
	if err := g.pay(p, c); err != nil {
		return err
	}
	p.Tableau = append(p.Tableau, c)
	removeCard(&p.Hand, c.ID)
	g.acted[playerID] = true
	return nil
}

func (g *Game) DiscardSelectedCard(playerID string) error {
	p, err := g.player(playerID)
	if err != nil {
		return err
	}
	if g.Phase != PhaseExperiment {
		return ErrInvalidPhase
	}
	c, ok := g.selected[playerID]
	if !ok {
		return ErrInvalidAction
	}
	p.Cash += DiscardRefund
	p.Discard = append(p.Discard, c)
	removeCard(&p.Hand, c.ID)
	g.acted[playerID] = true
	return nil
}

func (g *Game) TakeLoan(playerID string) error {
	p, e := g.player(playerID)
	if e != nil {
		return e
	}
	if p.Loans >= MaxLoans {
		return ErrLoanLimit
	}
	p.Loans++
	p.Cash += LoanAmount
	return nil
}
func (g *Game) RepayLoan(playerID string, count int) error {
	p, e := g.player(playerID)
	if e != nil {
		return e
	}
	if count < 1 || count > p.Loans || p.Cash < count*LoanAmount {
		return ErrInvalidAction
	}
	p.Loans -= count
	p.Cash -= count * LoanAmount
	return nil
}
func (g *Game) SettleInterest() error {
	newLoans := make([]int, len(g.Players))
	newCash := make([]int, len(g.Players))
	for _, p := range g.Players {
		interest := LoanInterest * p.Loans
		loans, cash := p.Loans, p.Cash
		for cash < interest && loans < MaxLoans {
			loans++
			cash += LoanAmount
		}
		if cash < interest {
			return ErrInsufficientCash
		}
		idx, _ := g.playerIndex(p.ID)
		newLoans[idx], newCash[idx] = loans, cash-interest
	}
	for i, p := range g.Players {
		p.Loans, p.Cash = newLoans[i], newCash[i]
	}
	return nil
}

// DistributeCustomers assigns customers in stable player order. A customer only
// contributes revenue when the player's tableau contains a card satisfying its demand.
func (g *Game) DistributeCustomers(customers []Customer) {
	for _, p := range g.Players {
		p.Customers = nil
	}
	for i, customer := range customers {
		p := g.Players[i%len(g.Players)]
		if customer.UnitPrice == 0 {
			customer.UnitPrice = basePrice(customer.Kind)
		}
		if !g.satisfies(p, customer.Demand) {
			customer.UnitPrice = 0
		}
		p.Customers = append(p.Customers, customer)
	}
}

func (g *Game) satisfies(p *Player, demand string) bool {
	for _, c := range p.Tableau {
		if c.Demand[demand] > 0 {
			return true
		}
	}
	return demand == "" // a blank demand is useful for fixture customers
}

func basePrice(kind string) int {
	if kind == "difficult" {
		return 0
	}
	return 10
}

func (g *Game) SettleRevenue() {
	for _, p := range g.Players {
		revenue := 0
		for _, c := range p.Customers {
			revenue += c.Count * c.UnitPrice
		}
		p.Revenue = revenue
		p.Cash += revenue
	}
}

// AdvancePeriod moves from learning to the next period after the current period's
// market, revenue and interest have been resolved.
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
func (g *Game) Rank() []*Player {
	out := append([]*Player(nil), g.Players...)
	sort.SliceStable(out, func(i, j int) bool { return less(out[i], out[j]) })
	return out
}
func (g *Game) Finish() error {
	if g.Phase != PhaseLearning {
		return ErrInvalidPhase
	}
	if err := g.SettleInterest(); err != nil {
		return err
	}
	g.SettleRevenue()
	for _, p := range g.Players {
		p.Score = p.Cash + p.metricScore()
	}
	g.Phase = PhaseFinished
	return nil
}

func (g *Game) Deal() { g.deal() }
func (g *Game) deal() {
	cards := make([]Card, 0, len(g.Players)*7)
	for i := 0; i < len(g.Players)*7; i++ {
		cards = append(cards, Card{ID: fmt.Sprintf("p%d-r%d-c%d", i/7, g.Round, i%7), Name: "Management Card", Kind: "resource", Period: g.Period, Cost: Cost{Cash: 10}})
	}
	g.rng.Shuffle(len(cards), func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })
	for i, p := range g.Players {
		p.Hand = append([]Card(nil), cards[i*7:(i+1)*7]...)
	}
}
func (g *Game) pay(p *Player, c Card) error {
	missing := 0
	owned := map[string]bool{}
	for _, x := range p.Tableau {
		for _, i := range x.Icons {
			owned[i] = true
		}
	}
	for _, i := range c.Cost.Icons {
		if !owned[i] {
			missing++
		}
	}
	cost := c.Cost.Cash + missing*20
	if p.Cash < cost {
		return ErrInsufficientCash
	}
	p.Cash -= cost
	return nil
}
func (g *Game) player(id string) (*Player, error) {
	for _, p := range g.Players {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, ErrInvalidAction
}

func (g *Game) playerIndex(id string) (int, error) {
	for i, p := range g.Players {
		if p.ID == id {
			return i, nil
		}
	}
	return -1, ErrInvalidAction
}
func (p *Player) metricScore() int {
	s := 0
	for _, k := range p.SelectedKPIs {
		switch k {
		case "brand_awareness":
			s += p.BrandAwareness
		case "products":
			s += p.Products
		case "values":
			s += p.Values
		case "resources":
			s += p.Resources
		}
	}
	return s
}
func less(a, b *Player) bool {
	for _, v := range [5]func(*Player) int{func(p *Player) int { return p.BrandAwareness }, func(p *Player) int { return p.Products }, func(p *Player) int { return p.Values }, func(p *Player) int { return p.Resources }, func(p *Player) int { return p.Cash }} {
		if v(a) != v(b) {
			return v(a) > v(b)
		}
	}
	return a.Order < b.Order
}
func removeCard(cards *[]Card, id string) {
	for i, c := range *cards {
		if c.ID == id {
			*cards = append((*cards)[:i], (*cards)[i+1:]...)
			return
		}
	}
}
func seedValue(seed string) int64 {
	var n int64
	for _, r := range seed {
		n = n*31 + int64(r)
	}
	return n
}

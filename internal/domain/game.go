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
	InitialCash      = 150
	LoanAmount       = 50
	LoanInterest     = 10
	MaxLoans         = 6
	DiscardRefund    = 20
	InitialRound     = 0
	ExperimentRounds = 6
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
	Cash  int      `json:"cash"`
	Icons []string `json:"icons"`
}
type Card struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Kind           string         `json:"kind"`
	Description    string         `json:"description,omitempty"`
	Effect         string         `json:"effect,omitempty"`
	Period         Period         `json:"period"`
	Cost           Cost           `json:"cost"`
	Icons          []string       `json:"icons"`
	MarketChange   map[string]int `json:"marketChange"`
	Source         string         `json:"source"`
	BrandAwareness int            `json:"brandAwareness,omitempty"`
	Products       int            `json:"products,omitempty"`
	Values         int            `json:"values,omitempty"`
	Resources      int            `json:"resources,omitempty"`
	Demand         map[string]int `json:"demand,omitempty"`
}
type Player struct {
	ID, DisplayName        string
	IsBot                  bool
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
	KPISelectionPeriod     Period
	CashFlow               []CashFlowStatement
	CashFlowRounds         []CashFlowStatement
	cashFlowBeginning      int
	cashFlowRevenue        int
	cashFlowOtherIncome    int
	cashFlowExpenses       int
	cashFlowInterest       int
	cashFlowPrincipal      int
	cashFlowNewLoans       int
}
type Customer struct {
	Kind, Demand     string
	UnitPrice, Count int
}
type CashFlowStatement struct {
	Period             Period `json:"period"`
	Round              int    `json:"round,omitempty"`
	BeginningCash      int    `json:"beginningCash"`
	OperatingRevenue   int    `json:"operatingRevenue"`
	OtherIncome        int    `json:"otherIncome"`
	OperatingExpenses  int    `json:"operatingExpenses"`
	InterestPaid       int    `json:"interestPaid"`
	PrincipalRepayment int    `json:"principalRepayment"`
	NewLoans           int    `json:"newLoans"`
	EndingCash         int    `json:"endingCash"`
}
type Game struct {
	Seed                               string
	Period                             Period
	Phase                              Phase
	Round                              int
	Players                            []*Player
	Center                             Card
	Market                             []Card
	Catalog                            []Card
	PartnerOptions, StarterShopOptions []Card
	DemandBoard                        map[string]int
	selected                           map[string]Card
	acted                              map[string]bool
	rng                                *rand.Rand
}

func (g *Game) SetCatalog(cards []Card) { g.Catalog = append([]Card(nil), cards...) }

// NewGame creates a deterministic, offline game. Two to four players are required to start.
func NewGame(seed string, playerIDs []string) (*Game, error) {
	g := &Game{Seed: seed, Period: PeriodOne, Phase: PhaseHypothesis, PartnerOptions: MVPPartnerCards(), StarterShopOptions: MVPStarterShopCards(), DemandBoard: map[string]int{"gourmet": 0, "regular": 0, "difficult": 0}, selected: map[string]Card{}, acted: map[string]bool{}, rng: rand.New(rand.NewSource(seedValue(seed)))}
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

// SetInitialCards changes a player's lobby selections. The default cards make
// old clients and deterministic tests compatible while new clients can choose
// explicit MVP cards before the game starts.
func (g *Game) SetInitialCards(playerID, partnerID, starterShopID string) error {
	if g.Phase != PhaseHypothesis || partnerID == "" || starterShopID == "" {
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
	for _, other := range g.Players {
		if other.ID != playerID && other.Partner.ID == partner.ID {
			return ErrInvalidAction
		}
	}
	p.Partner = partner
	p.StarterShop = starterShop
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

func (g *Game) SetKPIs(playerID string, kpis ...string) error {
	if g.Phase != PhaseHypothesis || g.Period == PeriodOne || len(kpis) != 2 || kpis[0] == kpis[1] {
		return ErrInvalidAction
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
	// Round 0 is the initial experiment state: cards are dealt, but no
	// select/pass/action cycle has been completed yet.
	for _, p := range g.Players {
		p.cashFlowBeginning = p.Cash
		p.cashFlowRevenue = 0
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

// HasSelected and HasActed expose read-only action state to the server's bot
// adapter without allowing callers to mutate the domain internals.
func (g *Game) HasSelected(playerID string) bool {
	_, ok := g.selected[playerID]
	return ok
}

func (g *Game) HasActed(playerID string) bool { return g.acted[playerID] }

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
	g.recordCashFlowRound()
	if g.Round == ExperimentRounds {
		for _, p := range g.Players {
			if len(p.Hand) != 1 {
				return fmt.Errorf("%w: hand size", ErrInvalidAction)
			}
		}
		// The final unselected card becomes the central covered card in the
		// digital MVP. It remains visible only through its public count/state.
		g.Center = g.Players[0].Hand[0]
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
	g.applyCardEffects(p, c)
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
	p.cashFlowOtherIncome += DiscardRefund
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
	p.cashFlowNewLoans += LoanAmount
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
	p.cashFlowPrincipal += count * LoanAmount
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
		newLoansTaken := newLoans[i] - p.Loans
		p.Loans, p.Cash = newLoans[i], newCash[i]
		p.cashFlowInterest += LoanInterest * (p.Loans - newLoansTaken)
		p.cashFlowNewLoans += newLoansTaken * LoanAmount
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
		if c.Demand[demand] > 0 || c.MarketChange[demand] > 0 {
			return true
		}
	}
	return demand == "" // a blank demand is useful for fixture customers
}

func (g *Game) applyCardEffects(p *Player, c Card) {
	switch c.Kind {
	case "marketing":
		p.BrandAwareness++
	case "product":
		p.Products++
	case "value":
		p.Values++
	case "resource", "channel":
		p.Resources++
	}
	for demand, change := range c.MarketChange {
		if demand == "difficult" {
			continue
		}
		g.DemandBoard[demand] += change
		if g.DemandBoard[demand] < 0 {
			g.DemandBoard[demand] = 0
		}
	}
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
		p.cashFlowRevenue += revenue
	}
}

func (g *Game) recordCashFlow() {
	for _, p := range g.Players {
		p.CashFlow = append(p.CashFlow, CashFlowStatement{
			Period:             g.Period,
			Round:              g.Round,
			BeginningCash:      p.cashFlowBeginning,
			OperatingRevenue:   p.cashFlowRevenue,
			OtherIncome:        p.cashFlowOtherIncome,
			OperatingExpenses:  p.cashFlowExpenses,
			InterestPaid:       p.cashFlowInterest,
			PrincipalRepayment: p.cashFlowPrincipal,
			NewLoans:           p.cashFlowNewLoans,
			EndingCash:         p.Cash,
		})
	}
}

func (g *Game) recordCashFlowRound() {
	for _, p := range g.Players {
		statement := CashFlowStatement{
			Period:             g.Period,
			Round:              g.Round,
			BeginningCash:      p.cashFlowBeginning,
			OperatingRevenue:   p.cashFlowRevenue,
			OtherIncome:        p.cashFlowOtherIncome,
			OperatingExpenses:  p.cashFlowExpenses,
			InterestPaid:       p.cashFlowInterest,
			PrincipalRepayment: p.cashFlowPrincipal,
			NewLoans:           p.cashFlowNewLoans,
			EndingCash:         p.Cash,
		}
		if len(p.CashFlowRounds) > 0 {
			last := &p.CashFlowRounds[len(p.CashFlowRounds)-1]
			if last.Period == statement.Period && last.Round == statement.Round {
				*last = statement
				continue
			}
		}
		p.CashFlowRounds = append(p.CashFlowRounds, statement)
	}
}

// ResolveLearning applies the MVP learning phase in one server-side step.
// Customer generation is intentionally simple and deterministic; the full
// market/customer workshop can be added without changing the phase contract.
func (g *Game) ResolveLearning() error {
	if g.Phase != PhaseLearning {
		return ErrInvalidPhase
	}
	if err := g.SettleInterest(); err != nil {
		return err
	}
	market := make([]string, 0)
	for _, kind := range []string{"gourmet", "regular", "difficult"} {
		count := g.DemandBoard[kind]
		if count < 1 {
			count = 1
		}
		for i := 0; i < count; i++ {
			market = append(market, kind)
		}
	}
	g.rng.Shuffle(len(market), func(i, j int) { market[i], market[j] = market[j], market[i] })
	customers := make([]Customer, 0, len(g.Players))
	for i := range g.Players {
		kind := market[(int(g.Period)+i)%len(market)]
		demand := ""
		if kind != "difficult" {
			demand = kind
		}
		customers = append(customers, Customer{Kind: kind, Demand: demand, Count: 1})
	}
	g.DistributeCustomers(customers)
	g.SettleRevenue()
	g.recordCashFlow()
	g.recordCashFlowRound()
	if g.Period == PeriodThree {
		for _, p := range g.Players {
			p.Score = p.Cash + p.metricScore()
		}
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
	if len(g.Catalog) > 0 {
		cards := make([]Card, 0, len(g.Catalog))
		for _, card := range g.Catalog {
			if card.Period == g.Period {
				cards = append(cards, card)
			}
		}
		g.rng.Shuffle(len(cards), func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })
		if len(cards) < len(g.Players)*7 {
			return
		}
		for i, p := range g.Players {
			p.Hand = append([]Card(nil), cards[i*7:(i+1)*7]...)
		}
		return
	}
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
	p.cashFlowExpenses += cost
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

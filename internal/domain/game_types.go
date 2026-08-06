package domain

import (
	"errors"
	"math/rand"
)

type Period int

const (
	PeriodZero  Period = 0
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
	ID                string         `json:"id"`
	Name              string         `json:"name"`
	Kind              string         `json:"kind"`
	Description       string         `json:"description,omitempty"`
	Function          string         `json:"function,omitempty"`
	ColorKey          string         `json:"colorKey,omitempty"`
	StartingCashBonus int            `json:"startingCashBonus,omitempty"`
	StarterShopID     string         `json:"starterShopId,omitempty"`
	Period            Period         `json:"period"`
	Cost              Cost           `json:"cost"`
	Icons             []string       `json:"icons"`
	MarketChange      map[string]int `json:"marketChange"`
	CustomerCount     map[string]int `json:"customerCount,omitempty"`
	Source            string         `json:"source"`
	BrandAwareness    int            `json:"brandAwareness,omitempty"`
	Products          int            `json:"products,omitempty"`
	Values            int            `json:"values,omitempty"`
	Resources         int            `json:"resources,omitempty"`
	Demand            map[string]int `json:"demand,omitempty"`
}

type DemandCard struct {
	ID       string   `json:"id"`
	Position int      `json:"position"`
	Icons    []string `json:"icons"`
	Revealed bool     `json:"revealed"`
}

type Player struct {
	ID, DisplayName        string
	IsBot                  bool
	Cash, Loans            int
	BrandAwareness         int
	Products               int
	Values                 int
	Resources              int
	GourmetSatisfaction    int
	RegularSatisfaction    int
	IconValues             map[string]int
	Partner, StarterShop   Card
	InitialCardsSelected   bool
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
	cashFlowGourmetCount   int
	cashFlowRegularCount   int
	cashFlowGourmetRevenue int
	cashFlowRegularRevenue int
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

type MarketDraw struct {
	Rank           int            `json:"rank"`
	PlayerID       string         `json:"playerId"`
	CustomerCounts map[string]int `json:"customerCounts"`
	Total          int            `json:"total"`
}

type CashFlowStatement struct {
	Period               Period `json:"period"`
	Round                int    `json:"round,omitempty"`
	BeginningCash        int    `json:"beginningCash"`
	OperatingRevenue     int    `json:"operatingRevenue"`
	GourmetRevenue       int    `json:"gourmetRevenue"`
	RegularRevenue       int    `json:"regularRevenue"`
	GourmetCustomerCount int    `json:"gourmetCustomerCount"`
	RegularCustomerCount int    `json:"regularCustomerCount"`
	OtherIncome          int    `json:"otherIncome"`
	OperatingExpenses    int    `json:"operatingExpenses"`
	InterestPaid         int    `json:"interestPaid"`
	PrincipalRepayment   int    `json:"principalRepayment"`
	NewLoans             int    `json:"newLoans"`
	EndingCash           int    `json:"endingCash"`
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
	DemandCards                        map[string][]DemandCard `json:"demandCards,omitempty"`
	MarketRanking                      []int                   `json:"marketRanking,omitempty"`
	MarketDraws                        []MarketDraw            `json:"marketDraws,omitempty"`
	MarketBag                          map[string]int          `json:"marketBag,omitempty"`
	marketDrawn                        bool
	selected                           map[string]Card
	acted                              map[string]bool
	rng                                *rand.Rand
}

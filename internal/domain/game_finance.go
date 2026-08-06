package domain

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

func (g *Game) recordCashFlow() {
	for _, p := range g.Players {
		p.CashFlow = append(p.CashFlow, CashFlowStatement{
			Period:               g.Period,
			Round:                g.Round,
			BeginningCash:        p.cashFlowBeginning,
			OperatingRevenue:     p.cashFlowRevenue,
			GourmetRevenue:       p.cashFlowGourmetCount * basePrice("gourmet"),
			RegularRevenue:       p.cashFlowRegularCount * basePrice("regular"),
			GourmetCustomerCount: p.cashFlowGourmetCount,
			RegularCustomerCount: p.cashFlowRegularCount,
			OtherIncome:          p.cashFlowOtherIncome,
			OperatingExpenses:    p.cashFlowExpenses,
			InterestPaid:         p.cashFlowInterest,
			PrincipalRepayment:   p.cashFlowPrincipal,
			NewLoans:             p.cashFlowNewLoans,
			EndingCash:           p.Cash,
		})
	}
}

func (g *Game) recordCashFlowRound() {
	for _, p := range g.Players {
		statement := CashFlowStatement{
			Period:               g.Period,
			Round:                g.Round,
			BeginningCash:        p.cashFlowBeginning,
			OperatingRevenue:     p.cashFlowRevenue,
			GourmetRevenue:       p.cashFlowGourmetCount * basePrice("gourmet"),
			RegularRevenue:       p.cashFlowRegularCount * basePrice("regular"),
			GourmetCustomerCount: p.cashFlowGourmetCount,
			RegularCustomerCount: p.cashFlowRegularCount,
			OtherIncome:          p.cashFlowOtherIncome,
			OperatingExpenses:    p.cashFlowExpenses,
			InterestPaid:         p.cashFlowInterest,
			PrincipalRepayment:   p.cashFlowPrincipal,
			NewLoans:             p.cashFlowNewLoans,
			EndingCash:           p.Cash,
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

func (g *Game) pay(p *Player, c Card) error {
	owned := map[string]int{}
	addOwnedIcons := func(card Card) {
		for _, icon := range card.Icons {
			owned[icon]++
		}
	}
	addOwnedIcons(p.Partner)
	addOwnedIcons(p.StarterShop)
	for _, x := range p.Tableau {
		addOwnedIcons(x)
	}
	missing := 0
	for _, i := range c.Cost.Icons {
		if owned[i] > 0 {
			owned[i]--
		} else {
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

func partnerInitialCashBonus(p *Player) int {
	if p.Partner.Function == "channel" && len(p.Partner.CustomerCount) == 0 {
		return p.Partner.Cost.Cash
	}
	return 0
}

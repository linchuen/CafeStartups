package domain

import "sort"

func (g *Game) Rank() []*Player {
	out := append([]*Player(nil), g.Players...)
	sort.SliceStable(out, func(i, j int) bool { return less(out[i], out[j]) })
	return out
}

func (g *Game) Finish() error {
	if g.Phase != PhaseLearning {
		return ErrInvalidPhase
	}
	g.SettleRevenue()
	if err := g.SettleInterest(); err != nil {
		return err
	}
	g.updateSatisfactionScores()
	for _, p := range g.Players {
		p.Score = p.Cash/CashScoreDivisor - p.Loans*6 + p.metricScore()
	}
	g.rememberFinalHands()
	g.Phase = PhaseFinished
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
	cards := p.scoreCardCounts()
	for _, k := range p.SelectedKPIs {
		switch k {
		case "gourmet_satisfaction":
			s += p.GourmetSatisfaction * 5
		case "regular_satisfaction":
			s += p.RegularSatisfaction * 4
		case "total_satisfaction":
			s += (p.GourmetSatisfaction + p.RegularSatisfaction) * 2
		case "channel":
			s += cards.channelCards * 4
		case "awareness":
			s += cards.marketingStars
		case "brand_awareness":
			s += cards.marketingStars
		case "products":
			s += cards.productCards * 3
		case "quality":
			s += cards.valueCards * 3
		case "values":
			s += cards.valueCards * 3
		case "cost":
			s += cards.resourceCards * 3
		case "surplus":
			cashAfterLoans := p.Cash - p.Loans*LoanAmount
			if cashAfterLoans > 0 {
				s += cashAfterLoans / 30
			}
		case "resources":
			s += cards.resourceCards * 3
		}
	}
	return s
}

type scoreCardCounts struct {
	channelCards, marketingStars, productCards, valueCards, resourceCards int
}

func (p *Player) scoreCardCounts() scoreCardCounts {
	counts := scoreCardCounts{}
	seen := map[string]bool{}
	cards := append([]Card{p.Partner, p.StarterShop}, p.Tableau...)
	for _, card := range cards {
		if card.ID != "" && seen[card.ID] {
			continue
		}
		if card.ID != "" {
			seen[card.ID] = true
		}
		category := card.Function
		if category == "" {
			category = card.ColorKey
		}
		if category == "" {
			category = card.Kind
		}
		if category == "marketing" {
			counts.marketingStars += card.BrandAwareness
		}
		switch category {
		case "channel":
			counts.channelCards++
		case "product":
			counts.productCards++
		case "value":
			counts.valueCards++
		case "resource":
			counts.resourceCards++
		}
	}
	return counts
}

func (p *Player) satisfiedCustomerCount(kind string) int {
	count := 0
	for _, customer := range p.Customers {
		if customer.Kind == kind && customer.UnitPrice > 0 {
			count += customer.Count
		}
	}
	return count
}

func less(a, b *Player) bool {
	aCards, bCards := a.scoreCardCounts(), b.scoreCardCounts()
	for _, values := range [][2]int{{aCards.marketingStars, bCards.marketingStars}, {aCards.productCards, bCards.productCards}, {aCards.valueCards, bCards.valueCards}, {aCards.resourceCards, bCards.resourceCards}, {a.Cash, b.Cash}} {
		if values[0] != values[1] {
			return values[0] > values[1]
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

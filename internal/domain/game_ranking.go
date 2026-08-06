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
	if err := g.SettleInterest(); err != nil {
		return err
	}
	g.SettleRevenue()
	g.updateSatisfactionScores()
	for _, p := range g.Players {
		p.Score = p.Cash + p.metricScore()
	}
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
	for _, k := range p.SelectedKPIs {
		switch k {
		case "gourmet_satisfaction":
			s += p.GourmetSatisfaction * 5
		case "regular_satisfaction":
			s += p.RegularSatisfaction * 4
		case "total_satisfaction":
			s += (p.GourmetSatisfaction + p.RegularSatisfaction) * 2
		case "channel":
			s += p.Resources
		case "awareness":
			s += p.BrandAwareness
		case "brand_awareness":
			s += p.BrandAwareness
		case "products":
			s += p.Products
		case "quality":
			s += p.Values
		case "values":
			s += p.Values
		case "cost":
			s += p.cashFlowExpenses
		case "surplus":
			s += p.Cash
		case "resources":
			s += p.Resources
		}
	}
	return s
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

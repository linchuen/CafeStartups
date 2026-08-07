package domain

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
	for _, p := range g.Players {
		g.appendCardCustomers(p)
	}
}

func (g *Game) appendCardCustomers(p *Player) {
	cards := []Card{p.Partner, p.StarterShop}
	cards = append(cards, p.Tableau...)
	for _, card := range cards {
		for kind, count := range card.CustomerCount {
			if count <= 0 {
				continue
			}
			demand := kind
			customer := Customer{Kind: kind, Demand: demand, Count: count, UnitPrice: basePrice(kind)}
			if !g.satisfies(p, demand) {
				customer.UnitPrice = 0
			}
			p.Customers = append(p.Customers, customer)
		}
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
	if p.IconValues == nil {
		p.IconValues = map[string]int{}
	}
	for _, icon := range c.Icons {
		switch icon {
		case "marketing":
			stars := c.BrandAwareness
			if stars < 1 {
				stars = 1
			}
			p.BrandAwareness += stars
			p.IconValues[icon] += stars
		case "coffee", "dessert", "beans":
			p.IconValues[icon]++
			p.Products++
		case "value", "taste", "service":
			p.IconValues[icon]++
			p.Values++
		case "data", "procurement", "operations", "marketing_resource", "channel":
			p.IconValues[icon]++
			p.Resources++
		}
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
		gourmetRate := g.demandRevenuePerCustomer("gourmet", p)
		regularRate := g.demandRevenuePerCustomer("regular", p)
		gourmetCount, regularCount := 0, 0
		for index := range p.Customers {
			c := &p.Customers[index]
			switch c.Kind {
			case "gourmet":
				c.UnitPrice = gourmetRate
				gourmetCount += c.Count
				p.cashFlowGourmetCount += c.Count
			case "regular":
				c.UnitPrice = regularRate
				regularCount += c.Count
				p.cashFlowRegularCount += c.Count
			}
		}
		p.cashFlowGourmetRevenue = gourmetRate * gourmetCount
		p.cashFlowRegularRevenue = regularRate * regularCount
		revenue := p.cashFlowGourmetRevenue + p.cashFlowRegularRevenue
		p.Revenue = revenue
		p.TotalRevenue += revenue
		p.Cash += revenue
		p.cashFlowRevenue += revenue
	}
}

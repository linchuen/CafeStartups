package domain

var demandIconTypes = []string{"coffee", "dessert", "beans", "taste", "service", "price"}

func (g *Game) prepareDemandCards() {
	first := make([]DemandCard, len(demandIconTypes))
	second := make([]DemandCard, len(demandIconTypes))
	for index, icon := range demandIconTypes {
		first[index] = DemandCard{ID: "demand-" + twoDigit(index+1), Position: index, Icons: []string{icon}}
		next := demandIconTypes[(index+1)%len(demandIconTypes)]
		second[index] = DemandCard{ID: "demand-" + twoDigit(index+7), Position: index, Icons: []string{icon, next}}
	}
	shuffleDemandCards(first, demandSeed(g.Seed))
	shuffleDemandCards(second, demandSeed(g.Seed)+97)
	for index := range first {
		first[index].Revealed = index == 0
		second[index].Revealed = index == 0
	}
	g.DemandCards = map[string][]DemandCard{"gourmet": first, "regular": second}
}

func (g *Game) revealDemandCards() {
	for _, cards := range g.DemandCards {
		for index := range cards {
			cards[index].Revealed = index <= g.Round
		}
	}
}

func demandSeed(seed string) uint32 {
	value := uint32(1)
	for index := 0; index < len(seed); index++ {
		value = value*31 + uint32(seed[index])
	}
	return value
}

func twoDigit(value int) string {
	if value < 10 {
		return "0" + string(rune('0'+value))
	}
	return string(rune('0'+value/10)) + string(rune('0'+value%10))
}

func shuffleDemandCards(cards []DemandCard, seed uint32) {
	value := seed
	for index := len(cards) - 1; index > 0; index-- {
		value = value*1664525 + 1013904223
		target := int(value % uint32(index+1))
		cards[index], cards[target] = cards[target], cards[index]
	}
}

func demandQuantity(kind string, position int) int {
	quantities := map[string][]int{"gourmet": {1, 2, 2, 2}, "regular": {1, 1, 1, 2}}
	values := quantities[kind]
	if position >= len(values) {
		return values[len(values)-1]
	}
	return values[position]
}

func (g *Game) satisfactionFor(kind string, p *Player) int {
	count := 0
	for _, card := range g.DemandCards[kind] {
		if !card.Revealed {
			continue
		}
		if g.satisfiesDemandCard(kind, card, p) {
			count++
		}
	}
	return count
}

func (g *Game) satisfiesDemandCard(kind string, card DemandCard, p *Player) bool {
	quantity := demandQuantity(kind, card.Position)
	available := map[string]int{}
	for _, owned := range []Card{p.Partner, p.StarterShop} {
		for _, icon := range owned.Icons {
			available[icon]++
		}
	}
	for _, owned := range p.Tableau {
		for _, icon := range owned.Icons {
			available[icon]++
		}
	}
	for index := 0; index < quantity; index++ {
		icon := card.Icons[index%len(card.Icons)]
		if available[icon] == 0 {
			return false
		}
		available[icon]--
	}
	return true
}

func demandValue(kind string, position int) int {
	values := map[string][]int{"gourmet": {10, 10, 20, 30}, "regular": {10, 10, 10, 10}}
	amounts := values[kind]
	if len(amounts) == 0 {
		return 0
	}
	if position >= len(amounts) {
		return amounts[len(amounts)-1]
	}
	return amounts[position]
}

func (g *Game) demandRevenuePerCustomer(kind string, p *Player) int {
	if len(g.DemandCards[kind]) == 0 {
		return 10
	}
	revenue := 0
	for _, card := range g.DemandCards[kind] {
		if card.Revealed && g.satisfiesDemandCard(kind, card, p) {
			revenue += demandValue(kind, card.Position)
		}
	}
	return revenue
}

func (g *Game) updateSatisfactionScores() {
	for _, p := range g.Players {
		p.GourmetSatisfaction = g.satisfactionFor("gourmet", p)
		p.RegularSatisfaction = g.satisfactionFor("regular", p)
	}
}

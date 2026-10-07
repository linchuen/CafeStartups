package domain

var demandIconTypes = []string{"coffee", "dessert", "beans", "taste", "service", "value"}

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
	g.DemandCards = map[string][]DemandCard{"gourmet": make([]DemandCard, 4), "regular": make([]DemandCard, 4)}
	ordinary, advanced := 0, 0
	for _, kind := range []string{"gourmet", "regular"} {
		for position := 0; position < 4; position++ {
			var card DemandCard
			if demandQuantity(kind, position) == 1 {
				card = first[ordinary]
				ordinary++
			} else {
				card = second[advanced]
				advanced++
			}
			card.Position = position
			card.Revealed = position == 0
			g.DemandCards[kind][position] = card
		}
	}
}

func (g *Game) revealDemandCards() {
	for _, cards := range g.DemandCards {
		for index := range cards {
			cards[index].Revealed = index <= int(g.Period)
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
	return len(g.matchedDemandCards(kind, p))
}

// A played card can support one demand in each customer row. Matching is
// performed left to right so repeated icons cannot reuse the same card.
func (g *Game) matchedDemandCards(kind string, p *Player) map[int]bool {
	owned := make([]Card, 0, len(p.Tableau))
	for _, card := range p.Tableau {
		category := card.Function
		if category == "" {
			category = card.ColorKey
		}
		if category == "" {
			category = card.Kind
		}
		if category == "product" || category == "value" {
			owned = append(owned, card)
		}
	}
	used := make([]bool, len(owned))
	matched := map[int]bool{}
	for _, demand := range g.DemandCards[kind] {
		if !demand.Revealed || len(demand.Icons) == 0 {
			continue
		}
		selected := map[int]map[string]int{}
		var assign func(int) bool
		assign = func(index int) bool {
			if index == len(demand.Icons) {
				return true
			}
			for cardIndex, card := range owned {
				if used[cardIndex] {
					continue
				}
				available := 0
				for _, icon := range card.Icons {
					if icon == demand.Icons[index] {
						available++
					}
				}
				if selected[cardIndex] == nil {
					selected[cardIndex] = map[string]int{}
				}
				if selected[cardIndex][demand.Icons[index]] >= available {
					continue
				}
				selected[cardIndex][demand.Icons[index]]++
				if assign(index + 1) {
					return true
				}
				selected[cardIndex][demand.Icons[index]]--
			}
			return false
		}
		if assign(0) {
			matched[demand.Position] = true
			for cardIndex, icons := range selected {
				for _, count := range icons {
					if count > 0 {
						used[cardIndex] = true
					}
				}
			}
		}
	}
	return matched
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
	revenue := basePrice(kind)
	matched := g.matchedDemandCards(kind, p)
	for _, card := range g.DemandCards[kind] {
		if matched[card.Position] {
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

package domain

import "fmt"

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
		marketChange := map[string]int{"gourmet": 1}
		if i%2 == 1 {
			marketChange = map[string]int{"regular": 1}
		}
		cards = append(cards, Card{ID: fmt.Sprintf("p%d-r%d-c%d", i/7, g.Round, i%7), Name: "Management Card", Kind: "resource", Period: g.Period, Cost: Cost{Cash: 10}, MarketChange: marketChange})
	}
	g.rng.Shuffle(len(cards), func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })
	for i, p := range g.Players {
		p.Hand = append([]Card(nil), cards[i*7:(i+1)*7]...)
	}
}

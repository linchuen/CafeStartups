package domain

import (
	"encoding/json"
	"fmt"
)

type fixtureFile struct {
	Source  string `json:"source"`
	Version int    `json:"version"`
	Seed    string `json:"seed"`
	Cards   []Card `json:"cards"`
}

// LoadCatalog validates and loads a data-driven MVP card fixture.
func LoadCatalog(data []byte) ([]Card, error) {
	var file fixtureFile
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("decode card fixture: %w", err)
	}
	if len(file.Cards) == 0 {
		return nil, fmt.Errorf("card fixture is empty")
	}
	seen := make(map[string]bool, len(file.Cards))
	for i, card := range file.Cards {
		if card.ID == "" || seen[card.ID] {
			return nil, fmt.Errorf("invalid or duplicate card id: %q", card.ID)
		}
		if card.Period < PeriodOne || card.Period > PeriodThree {
			return nil, fmt.Errorf("invalid period for card %q", card.ID)
		}
		if card.Cost.Cash < 0 {
			return nil, fmt.Errorf("negative cost for card %q", card.ID)
		}
		file.Cards[i] = withSchemaCostIcons(card)
		seen[card.ID] = true
	}
	return file.Cards, nil
}

// withSchemaCostIcons uses the card's schema icons as the resources required
// to obtain the card. Explicit cost icons remain a fallback for older data.
func withSchemaCostIcons(card Card) Card {
	if len(card.Cost.Icons) > 0 || len(card.Icons) == 0 {
		return card
	}

	card.Cost.Icons = append([]string(nil), card.Icons...)
	return card
}

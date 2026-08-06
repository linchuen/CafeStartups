package domain

import (
	"encoding/json"
	"fmt"

	fixturedata "cafestartups/data"
)

// MVPPartnerCards returns the selectable founder cards used during lobby setup.
// The card text is an MVP fixture and must be checked against the physical cards
// before the complete official card set is enabled.
func MVPPartnerCards() []Card {
	cards, err := LoadPartnerCards(fixturedata.MVPPartnerFixture)
	if err != nil {
		panic(fmt.Sprintf("invalid partner card fixture: %v", err))
	}
	return cards
}

// LoadPartnerCards loads partner cards using the same card-shaped fixture
// contract as the regular card catalog.
func LoadPartnerCards(data []byte) ([]Card, error) {
	var file fixtureFile
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("decode partner card fixture: %w", err)
	}
	if len(file.Cards) == 0 {
		return nil, fmt.Errorf("partner card fixture is empty")
	}

	seen := make(map[string]bool, len(file.Cards))
	for index, card := range file.Cards {
		if card.ID == "" || seen[card.ID] {
			return nil, fmt.Errorf("invalid or duplicate partner card id: %q", card.ID)
		}
		if card.Kind != "partner" || card.Function == "" || card.ColorKey != card.Function {
			return nil, fmt.Errorf("invalid partner card classification for %q", card.ID)
		}
		if len(card.Icons) == 0 || card.Cost.Cash < 0 {
			return nil, fmt.Errorf("invalid partner card schema fields for %q", card.ID)
		}
		for kind, count := range card.CustomerCount {
			if (kind != "gourmet" && kind != "regular") || count < 0 {
				return nil, fmt.Errorf("invalid partner customer count for %q", card.ID)
			}
		}
		if card.Function == "channel" {
			hasInitialCashBonus := card.Cost.Cash > 0
			hasCustomerCount := card.CustomerCount["gourmet"] > 0 || card.CustomerCount["regular"] > 0
			if hasInitialCashBonus == hasCustomerCount {
				return nil, fmt.Errorf("channel partner %q must have exactly one effect", card.ID)
			}
		}
		file.Cards[index] = card
		seen[card.ID] = true
	}
	return file.Cards, nil
}

// MVPStarterShopCards returns the selectable founding shop cards used during
// lobby setup.
func MVPStarterShopCards() []Card {
	cards, err := LoadStarterShopCards(fixturedata.MVPStarterShopFixture)
	if err != nil {
		panic(fmt.Sprintf("invalid starter shop fixture: %v", err))
	}
	return cards
}

// LoadStarterShopCards loads founding shop cards from the shared card-shaped
// fixture contract.
func LoadStarterShopCards(data []byte) ([]Card, error) {
	var file fixtureFile
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("decode starter shop fixture: %w", err)
	}
	if len(file.Cards) == 0 {
		return nil, fmt.Errorf("starter shop fixture is empty")
	}

	seen := make(map[string]bool, len(file.Cards))
	for index, card := range file.Cards {
		if card.ID == "" || seen[card.ID] || card.Kind != "starter_shop" || card.Cost.Cash < 0 {
			return nil, fmt.Errorf("invalid starter shop card %q", card.ID)
		}
		if card.Cost.Cash%10 != 0 {
			return nil, fmt.Errorf("starter shop %q cost must use 10-unit increments", card.ID)
		}
		for kind, count := range card.CustomerCount {
			if (kind != "gourmet" && kind != "regular") || count < 0 {
				return nil, fmt.Errorf("invalid starter shop customer count for %q", card.ID)
			}
		}
		file.Cards[index] = card
		seen[card.ID] = true
	}
	return file.Cards, nil
}

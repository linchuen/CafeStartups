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
// lobby setup. Costs and detailed card effects remain fixture data for now.
func MVPStarterShopCards() []Card {
	return []Card{
		{ID: "starter-songshan", Name: "松山店", Kind: "starter_shop", Description: "吸引高消費力的饕客。", Effect: "每期 +2 位饕客", Cost: Cost{Cash: 20}, Demand: map[string]int{"gourmet": 2}, Source: "mvp-fixture"},
		{ID: "starter-minsheng", Name: "民生店", Kind: "starter_shop", Description: "吸引高消費力的饕客與一般客。", Effect: "每期 +1 位饕客、+1 位一般客", Cost: Cost{Cash: 25}, Demand: map[string]int{"gourmet": 1, "regular": 1}, Source: "mvp-fixture"},
		{ID: "starter-xinyi", Name: "信義店", Kind: "starter_shop", Description: "開設分店，同時吸引饕客與一般客群。", Effect: "每期 +1 位饕客、+2 位一般客", Cost: Cost{Cash: 30}, Demand: map[string]int{"gourmet": 1, "regular": 2}, Source: "mvp-fixture"},
		{ID: "starter-station", Name: "站前品牌旗艦店", Kind: "starter_shop", Description: "開設旗艦店，吸引大批一般客。", Effect: "每期 +3 位一般客", Cost: Cost{Cash: 35}, Demand: map[string]int{"regular": 3}, Source: "mvp-fixture"},
		{ID: "starter-daan", Name: "大安巷口店", Kind: "starter_shop", Description: "深耕社區客群，帶來穩定的一般客。", Effect: "每期 +2 位一般客", Cost: Cost{Cash: 22}, Demand: map[string]int{"regular": 2}, Source: "mvp-fixture"},
		{ID: "starter-beitou", Name: "北投溫泉店", Kind: "starter_shop", Description: "吸引休閒旅客與高消費力饕客。", Effect: "每期 +2 位饕客", Cost: Cost{Cash: 28}, Demand: map[string]int{"gourmet": 2}, Source: "mvp-fixture"},
		{ID: "starter-neihu", Name: "內湖辦公店", Kind: "starter_shop", Description: "服務辦公商圈，兼顧饕客與一般客。", Effect: "每期 +1 位饕客、+2 位一般客", Cost: Cost{Cash: 32}, Demand: map[string]int{"gourmet": 1, "regular": 2}, Source: "mvp-fixture"},
		{ID: "starter-banqiao", Name: "板橋轉運店", Kind: "starter_shop", Description: "掌握交通人流，快速累積一般客。", Effect: "每期 +3 位一般客", Cost: Cost{Cash: 34}, Demand: map[string]int{"regular": 3}, Source: "mvp-fixture"},
		{ID: "starter-ximen", Name: "西門潮流店", Kind: "starter_shop", Description: "連結年輕客群，帶來饕客與一般客。", Effect: "每期 +2 位饕客、+1 位一般客", Cost: Cost{Cash: 38}, Demand: map[string]int{"gourmet": 2, "regular": 1}, Source: "mvp-fixture"},
		{ID: "starter-gongguan", Name: "公館學府店", Kind: "starter_shop", Description: "經營學生與社群客群，擴大來客數。", Effect: "每期 +1 位饕客、+2 位一般客", Cost: Cost{Cash: 26}, Demand: map[string]int{"gourmet": 1, "regular": 2}, Source: "mvp-fixture"},
	}
}

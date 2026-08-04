package domain

// MVPPartnerCards returns the selectable founder cards used during lobby setup.
// The card text is an MVP fixture and must be checked against the physical cards
// before the complete official card set is enabled.
func MVPPartnerCards() []Card {
	return []Card{
		{ID: "partner-barista", Name: "國際認證咖啡師", Kind: "partner", Description: "熟悉咖啡館營運，協助降低營運成本。", Effect: "營運成本 -10", Source: "mvp-fixture"},
		{ID: "partner-roaster", Name: "精品烘豆顧問", Kind: "partner", Description: "協助挑選咖啡豆，建立特色產品方向。", Effect: "產品能力 +1", Source: "mvp-fixture"},
		{ID: "partner-marketer", Name: "品牌行銷顧問", Kind: "partner", Description: "協助規劃行銷活動，建立品牌知名度。", Effect: "品牌知名度 +1", Source: "mvp-fixture"},
		{ID: "partner-service", Name: "服務設計顧問", Kind: "partner", Description: "協助改善服務流程，提升顧客體驗。", Effect: "價值主張 +1", Source: "mvp-fixture"},
	}
}

// MVPStarterShopCards returns the selectable founding shop cards used during
// lobby setup. Costs and detailed card effects remain fixture data for now.
func MVPStarterShopCards() []Card {
	return []Card{
		{ID: "starter-songshan", Name: "松山店", Kind: "starter_shop", Description: "吸引高消費力的饕客。", Effect: "每期 +2 位精品顧客", Cost: Cost{Cash: 20}, Demand: map[string]int{"gourmet": 2}, Source: "mvp-fixture"},
		{ID: "starter-minsheng", Name: "民生店", Kind: "starter_shop", Description: "吸引高消費力的饕客與一般客。", Effect: "每期 +1 位精品顧客、+1 位一般顧客", Cost: Cost{Cash: 25}, Demand: map[string]int{"gourmet": 1, "regular": 1}, Source: "mvp-fixture"},
		{ID: "starter-xinyi", Name: "信義店", Kind: "starter_shop", Description: "開設分店，同時吸引饕客與一般客群。", Effect: "每期 +1 位精品顧客、+2 位一般顧客", Cost: Cost{Cash: 30}, Demand: map[string]int{"gourmet": 1, "regular": 2}, Source: "mvp-fixture"},
		{ID: "starter-station", Name: "站前品牌旗艦店", Kind: "starter_shop", Description: "開設旗艦店，吸引大批一般消費客群。", Effect: "每期 +3 位一般顧客", Cost: Cost{Cash: 35}, Demand: map[string]int{"regular": 3}, Source: "mvp-fixture"},
	}
}

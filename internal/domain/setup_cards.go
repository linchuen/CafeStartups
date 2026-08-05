package domain

// MVPPartnerCards returns the selectable founder cards used during lobby setup.
// The card text is an MVP fixture and must be checked against the physical cards
// before the complete official card set is enabled.
func MVPPartnerCards() []Card {
	return []Card{
		{ID: "partner-barista", Name: "國際認證咖啡師", Kind: "partner", Function: "resource", ColorKey: "resource", Description: "熟悉咖啡館營運，協助降低營運成本。", Source: "mvp-fixture"},
		{ID: "partner-roaster", Name: "精品烘豆顧問", Kind: "partner", Function: "product", ColorKey: "product", Description: "協助挑選咖啡豆，建立特色產品方向。", Source: "mvp-fixture"},
		{ID: "partner-marketer", Name: "品牌行銷顧問", Kind: "partner", Function: "marketing", ColorKey: "marketing", Description: "協助規劃行銷活動，建立品牌知名度。", Source: "mvp-fixture"},
		{ID: "partner-service", Name: "服務設計顧問", Kind: "partner", Function: "value", ColorKey: "value", Description: "協助改善服務流程，提升顧客體驗。", Source: "mvp-fixture"},
		{ID: "partner-finance", Name: "財務管理顧問", Kind: "partner", Function: "resource", ColorKey: "resource", Description: "協助規劃資金運用，控制貸款與利息成本。", Source: "mvp-fixture"},
		{ID: "partner-pastry", Name: "甜點開發師", Kind: "partner", Function: "product", ColorKey: "product", Description: "開發搭配甜點，豐富咖啡館的產品品項。", Source: "mvp-fixture"},
		{ID: "partner-supply", Name: "供應鏈管理師", Kind: "partner", Function: "resource", ColorKey: "resource", Description: "整合原料供應與採購，降低營運浪費。", Source: "mvp-fixture"},
		{ID: "partner-community", Name: "社區合作顧問", Kind: "partner", Function: "channel", ColorKey: "channel", Description: "經營社區關係，為開局提供額外預算。", StartingCashBonus: 10, Source: "mvp-fixture"},
		{ID: "partner-hr", Name: "人才培訓顧問", Kind: "partner", Function: "value", ColorKey: "value", Description: "培育團隊服務能力，穩定顧客體驗。", Source: "mvp-fixture"},
		{ID: "partner-analytics", Name: "數據營運顧問", Kind: "partner", Function: "channel", ColorKey: "channel", Description: "分析商圈資料，開局提供一間大安巷口店。", StarterShopID: "starter-daan", Source: "mvp-fixture"},
	}
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

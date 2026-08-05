package main

func (room *gameRoom) view(token string) map[string]any {
	players := make([]map[string]any, 0, len(room.Domain.Players))
	for _, p := range room.Domain.Players {
		players = append(players, map[string]any{"id": p.ID, "displayName": p.DisplayName, "bot": p.IsBot, "cash": p.Cash, "loans": p.Loans, "revenue": p.Revenue, "score": p.Score, "selectedKPIs": p.SelectedKPIs, "brandAwareness": p.BrandAwareness, "products": p.Products, "values": p.Values, "resources": p.Resources, "handCount": len(p.Hand)})
	}
	result := map[string]any{"id": room.ID, "status": room.Status, "seed": room.Seed, "gameVersion": room.Version, "period": room.Domain.Period, "phase": room.Domain.Phase, "round": room.Domain.Round, "demandBoard": room.Domain.DemandBoard, "marketRanking": room.Domain.MarketRanking, "center": room.Domain.Center, "partnerOptions": room.Domain.PartnerOptions, "starterShopOptions": room.Domain.StarterShopOptions, "players": players}
	if token == room.Token {
		for _, p := range room.Domain.Players {
			if p.ID == room.PlayerID {
				result["me"] = map[string]any{"id": p.ID, "hand": p.Hand, "tableau": p.Tableau, "discardCount": len(p.Discard), "partner": p.Partner, "starterShop": p.StarterShop, "initialCardsSelected": p.InitialCardsSelected, "cash": p.Cash, "loans": p.Loans, "customers": p.Customers, "revenue": p.Revenue, "score": p.Score, "selectedKPIs": p.SelectedKPIs, "cashFlow": p.CashFlow, "cashFlowRounds": p.CashFlowRounds, "brandAwareness": p.BrandAwareness, "products": p.Products, "values": p.Values, "resources": p.Resources}
			}
		}
	}
	return result
}

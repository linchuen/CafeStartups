package main

import (
	"cafestartups/internal/application"
	"cafestartups/internal/domain"
)

func applyCommand(room *gameRoom, playerID string, input commandRequest) error {
	return application.ExecuteCommand(room.Domain, playerID, application.Command{
		Type:  input.Type,
		Card:  input.CardID,
		KPIs:  input.KPIs,
		Count: input.Count,
	})
}

func beginNextPeriod(room *gameRoom) error {
	return application.BeginNextPeriod(room.Domain, room.Seed, &room.Version)
}

func allKPIsSelected(room *gameRoom) bool {
	return application.AllKPIsSelected(room.Domain)
}

func setRandomBotKPIs(room *gameRoom, player *domain.Player) error {
	// This wrapper keeps existing server tests and call sites stable while the
	// application package owns the actual bot policy.
	return application.BeginNextPeriodForBot(room.Domain, room.Seed, room.Version, player)
}

func addBots(room *gameRoom) error {
	for len(room.Domain.Players) < 4 {
		index := len(room.Domain.Players) + 1
		if err := room.Domain.AddPlayer("bot-"+itoa(int64(index)), "?餉?拙振 "+itoa(int64(index))); err != nil {
			return err
		}
		room.Domain.Players[len(room.Domain.Players)-1].IsBot = true
		room.Version++
	}
	return nil
}

func runBots(room *gameRoom) {
	application.RunBots(room.Domain, room.Seed, &room.Version)
}

package application

import "cafestartups/internal/domain"

// AllKPIsSelected reports whether every player completed the KPI choice for
// the current hypothesis phase. Period one is the initial experiment.
func AllKPIsSelected(game *domain.Game) bool {
	if game.Phase != domain.PhaseHypothesis || game.Period == domain.PeriodOne {
		return false
	}
	for _, player := range game.Players {
		if len(player.SelectedKPIs) != 2 || player.KPISelectionPeriod != game.Period {
			return false
		}
	}
	return len(game.Players) > 0
}

// BeginNextPeriod completes bot KPI setup, starts the next experiment, and
// advances the optimistic concurrency version owned by the application layer.
func BeginNextPeriod(game *domain.Game, seed string, version *int64) error {
	for _, player := range game.Players {
		if player.IsBot && player.KPISelectionPeriod != game.Period {
			if err := setRandomBotKPIs(game, seed, *version, player); err != nil {
				return err
			}
		}
	}
	if err := game.BeginExperiment(); err != nil {
		return err
	}
	*version++
	return nil
}

// BeginNextPeriodForBot keeps the server compatibility adapter small while
// the bot policy remains owned by this application package.
func BeginNextPeriodForBot(game *domain.Game, seed string, version int64, player *domain.Player) error {
	return setRandomBotKPIs(game, seed, version, player)
}

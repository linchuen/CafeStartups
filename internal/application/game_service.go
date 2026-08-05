package application

import (
	"crypto/sha1"
	"fmt"

	"cafestartups/internal/domain"
)

// Command is the application-layer input for a player action. HTTP details
// such as tokens and request versions stay outside this package.
type Command struct {
	Type  string
	Card  string
	KPIs  []string
	Count int
}

// ExecuteCommand translates an application command into a domain operation.
func ExecuteCommand(game *domain.Game, playerID string, command Command) error {
	switch command.Type {
	case "BEGIN_EXPERIMENT":
		return game.BeginExperiment()
	case "SET_KPI":
		return game.SetKPIs(playerID, command.KPIs...)
	case "SELECT_CARD":
		return game.SelectCard(playerID, command.Card)
	case "PLAY_SELECTED_CARD":
		return game.PlaySelectedCard(playerID)
	case "DISCARD_SELECTED_CARD":
		return game.DiscardSelectedCard(playerID)
	case "PASS_HAND":
		return game.PassHands()
	case "TAKE_LOAN":
		return game.TakeLoan(playerID)
	case "REPAY_LOAN":
		return game.RepayLoan(playerID, command.Count)
	case "CONFIRM_INTEREST":
		return game.SettleInterest()
	case "CONFIRM_REVENUE":
		game.SettleRevenue()
		return nil
	case "RESOLVE_LEARNING":
		return game.ResolveLearning()
	case "DRAW_MARKET":
		return game.DrawMarket()
	default:
		return domain.ErrInvalidAction
	}
}

// AllKPIsSelected reports whether every player completed the KPI choice for
// the current hypothesis phase. Period one is the initial experiment and does
// not require this transition.
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

// BeginNextPeriodForBot applies the same deterministic KPI policy to one bot.
// It is exported for the compatibility adapter in cmd/server.
func BeginNextPeriodForBot(game *domain.Game, seed string, version int64, player *domain.Player) error {
	return setRandomBotKPIs(game, seed, version, player)
}

// RunBots performs deterministic, intentionally simple actions for solo MVP
// games. Domain methods remain responsible for legality and state changes.
func RunBots(game *domain.Game, seed string, version *int64) {
	if game.Phase == domain.PhaseHypothesis {
		for _, player := range game.Players {
			if player.IsBot && player.KPISelectionPeriod != game.Period {
				kpis := validKPIs()
				i := botChoice(seed, player.ID, *version, len(kpis))
				j := botChoice(seed, player.ID+"-second", *version, len(kpis)-1)
				if j >= i {
					j++
				}
				if err := game.SetKPIs(player.ID, kpis[i], kpis[j]); err == nil {
					*version++
				}
			}
		}
	}
	if game.Phase != domain.PhaseExperiment {
		return
	}
	for _, player := range game.Players {
		if !player.IsBot {
			continue
		}
		if !game.HasSelected(player.ID) {
			if len(player.Hand) == 0 {
				continue
			}
			card := player.Hand[botChoice(seed, player.ID, *version, len(player.Hand))]
			if err := game.SelectCard(player.ID, card.ID); err == nil {
				*version++
			}
		}
		if game.HasActed(player.ID) {
			continue
		}
		if botChoice(seed, player.ID+"-action", *version, 2) == 0 {
			if err := game.PlaySelectedCard(player.ID); err == nil {
				*version++
				continue
			}
		}
		if err := game.DiscardSelectedCard(player.ID); err == nil {
			*version++
		}
	}
}

func setRandomBotKPIs(game *domain.Game, seed string, version int64, player *domain.Player) error {
	kpis := validKPIs()
	i := botChoice(seed, player.ID+"-kpi", version, len(kpis))
	j := botChoice(seed, player.ID+"-kpi-second", version, len(kpis)-1)
	if j >= i {
		j++
	}
	return game.SetKPIs(player.ID, kpis[i], kpis[j])
}

func validKPIs() []string {
	return []string{"gourmet_satisfaction", "regular_satisfaction", "total_satisfaction", "channel", "awareness", "products", "quality", "cost", "surplus"}
}

func botChoice(seed, key string, version int64, size int) int {
	if size <= 1 {
		return 0
	}
	sum := sha1.Sum([]byte(fmt.Sprintf("%s|%s|%d", seed, key, version)))
	return int(sum[0]) % size
}

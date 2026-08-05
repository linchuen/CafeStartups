package application

import (
	"crypto/sha1"
	"fmt"

	"cafestartups/internal/domain"
)

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

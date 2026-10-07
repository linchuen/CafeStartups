package application

import "cafestartups/internal/domain"

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
	case "SKIP_PEEK":
		return game.SkipPeek(playerID)
	case "DISCARD_SELECTED_CARD":
		return game.DiscardSelectedCard(playerID)
	case "PASS_HAND":
		return game.PassHands()
	case "TAKE_LOAN":
		return game.TakeLoan(playerID)
	case "REPAY_LOAN":
		return game.RepayLoan(playerID, command.Count)
	case "RESOLVE_LEARNING":
		return game.ResolveLearning()
	case "DRAW_MARKET":
		return game.DrawMarket()
	case "PREPARE_MARKET_BAG":
		return game.PrepareMarketBag()
	default:
		return domain.ErrInvalidAction
	}
}

func PeekDemand(game *domain.Game, playerID, kind string, position int) (domain.DemandCard, error) {
	return game.PeekDemand(playerID, kind, position)
}

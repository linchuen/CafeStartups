package domain

import "fmt"

func (g *Game) SelectCard(playerID, cardID string) error {
	if g.Phase != PhaseExperiment {
		return ErrInvalidPhase
	}
	p, err := g.player(playerID)
	if err != nil {
		return err
	}
	if g.acted[playerID] {
		return ErrAlreadySelected
	}
	for _, c := range p.Hand {
		if c.ID == cardID {
			g.selected[playerID] = c
			g.acted[playerID] = false
			return nil
		}
	}
	return ErrCardNotFound
}

func (g *Game) HasSelected(playerID string) bool {
	_, ok := g.selected[playerID]
	return ok
}

func (g *Game) HasActed(playerID string) bool { return g.acted[playerID] }

func (g *Game) PassHandsIfReady() (bool, error) {
	if g.Phase != PhaseExperiment {
		return false, ErrInvalidPhase
	}
	if len(g.selected) != len(g.Players) {
		return false, nil
	}
	for _, p := range g.Players {
		if !g.acted[p.ID] {
			return false, nil
		}
	}
	if err := g.PassHands(); err != nil {
		return false, err
	}
	return true, nil
}

func (g *Game) PassHands() error {
	if g.Phase != PhaseExperiment {
		return ErrInvalidPhase
	}
	if len(g.selected) != len(g.Players) {
		return ErrInvalidAction
	}
	for _, p := range g.Players {
		if !g.acted[p.ID] {
			return ErrInvalidAction
		}
	}
	remaining := make([][]Card, len(g.Players))
	for i, p := range g.Players {
		remaining[i] = append([]Card(nil), p.Hand...)
	}
	direction := 1
	if g.Period == PeriodTwo {
		direction = -1
	}
	for i, p := range g.Players {
		target := (i + direction + len(g.Players)) % len(g.Players)
		g.Players[target].Hand = append(g.Players[target].Hand[:0], remaining[i]...)
		_ = p
	}
	g.Round++
	g.selected, g.acted = map[string]Card{}, map[string]bool{}
	g.recordCashFlowRound()
	if g.Round == ExperimentRounds {
		for _, p := range g.Players {
			if len(p.Hand) != 1 {
				return fmt.Errorf("%w: hand size", ErrInvalidAction)
			}
		}
		g.Center = g.Players[0].Hand[0]
		g.Phase = PhaseLearning
	}
	return nil
}

func (g *Game) PlaySelectedCard(playerID string) error {
	p, err := g.player(playerID)
	if err != nil {
		return err
	}
	if g.Phase != PhaseExperiment {
		return ErrInvalidPhase
	}
	c, ok := g.selected[playerID]
	if !ok {
		return ErrInvalidAction
	}
	if err := g.pay(p, c); err != nil {
		return err
	}
	p.Tableau = append(p.Tableau, c)
	g.applyCardEffects(p, c)
	removeCard(&p.Hand, c.ID)
	g.acted[playerID] = true
	return nil
}

func (g *Game) DiscardSelectedCard(playerID string) error {
	p, err := g.player(playerID)
	if err != nil {
		return err
	}
	if g.Phase != PhaseExperiment {
		return ErrInvalidPhase
	}
	c, ok := g.selected[playerID]
	if !ok {
		return ErrInvalidAction
	}
	p.Cash += DiscardRefund
	p.cashFlowOtherIncome += DiscardRefund
	p.Discard = append(p.Discard, c)
	removeCard(&p.Hand, c.ID)
	g.acted[playerID] = true
	return nil
}

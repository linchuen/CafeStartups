package main

import (
	"sync"

	"cafestartups/internal/domain"
)

type game struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type gameRoom struct {
	game
	Seed      string
	Token     string
	PlayerID  string
	Domain    *domain.Game
	Version   int64
	Processed map[string]commandResult
}

type commandResult struct {
	Status int
	Body   []byte
}

type gameStore struct {
	sync.RWMutex
	games map[string]*gameRoom
}

// gameLocked resolves a local game by its stable id. Callers must hold the store lock.
func (s *gameStore) gameLocked(id string) (*gameRoom, bool) {
	game, ok := s.games[id]
	return game, ok
}

type setupRequest struct {
	Token         string `json:"token"`
	PartnerID     string `json:"partnerId"`
	StarterShopID string `json:"starterShopId"`
}

type commandRequest struct {
	Token       string   `json:"token"`
	GameVersion int64    `json:"gameVersion"`
	Type        string   `json:"type"`
	CardID      string   `json:"cardId"`
	KPIs        []string `json:"kpis"`
	Count       int      `json:"count"`
	CommandID   string   `json:"commandId"`
}

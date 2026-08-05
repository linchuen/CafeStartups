package server

import (
	"encoding/json"
	"net/http"
	"sync"

	fixturedata "cafestartups/data"
	"cafestartups/internal/application"
	"cafestartups/internal/catalog"
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

func applyCommand(room *gameRoom, playerID string, input commandRequest) error {
	return application.ExecuteCommand(room.Domain, playerID, application.Command{Type: input.Type, Card: input.CardID, KPIs: input.KPIs, Count: input.Count})
}

func beginNextPeriod(room *gameRoom) error {
	return application.BeginNextPeriod(room.Domain, room.Seed, &room.Version)
}

func allKPIsSelected(room *gameRoom) bool {
	return application.AllKPIsSelected(room.Domain)
}

func setRandomBotKPIs(room *gameRoom, player *domain.Player) error {
	return application.BeginNextPeriodForBot(room.Domain, room.Seed, room.Version, player)
}

func addBots(room *gameRoom) error {
	for len(room.Domain.Players) < 4 {
		index := len(room.Domain.Players) + 1
		if err := room.Domain.AddPlayer("bot-"+itoa(int64(index)), "?擗?? "+itoa(int64(index))); err != nil {
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

func passHandsIfReady(room *gameRoom) error {
	if room.Domain.Phase != domain.PhaseExperiment {
		return nil
	}
	passed, err := room.Domain.PassHandsIfReady()
	if passed {
		room.Version++
	}
	return err
}

func (s *gameStore) gamesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var input struct {
			Seed        string `json:"seed"`
			DisplayName string `json:"displayName"`
		}
		_ = json.NewDecoder(r.Body).Decode(&input)
		id := randomID(8)
		seed := input.Seed
		if seed == "" {
			seed = id
		}
		playerID := "player-" + randomID(6)
		g, err := domain.NewGame(seed, []string{playerID})
		if err != nil {
			writeDomainError(w, err)
			return
		}
		cards, err := catalog.Load(fixturedata.MVPFixture)
		if err != nil {
			writeCode(w, http.StatusInternalServerError, "INVALID_CARD_CATALOG")
			return
		}
		g.SetCatalog(cards)
		if input.DisplayName != "" {
			g.Players[0].DisplayName = input.DisplayName
		}
		token := randomID(16)
		localGame := &gameRoom{game: game{ID: id, Status: "lobby"}, Seed: seed, Token: token, PlayerID: playerID, Domain: g, Processed: map[string]commandResult{}}
		s.Lock()
		s.games[id] = localGame
		s.Unlock()
		writeJSON(w, http.StatusCreated, map[string]any{"id": id, "token": token, "playerId": playerID, "state": localGame.view(token)})
		return
	}
	if r.Method == http.MethodGet {
		s.RLock()
		games := make([]game, 0, len(s.games))
		for _, room := range s.games {
			games = append(games, room.game)
		}
		s.RUnlock()
		writeJSON(w, http.StatusOK, games)
		return
	}
	w.Header().Set("Allow", "GET, POST")
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

func (s *gameStore) setupHandler(w http.ResponseWriter, r *http.Request) {
	s.Lock()
	defer s.Unlock()
	room, ok := s.gameLocked(r.PathValue("id"))
	if !ok {
		writeCode(w, http.StatusNotFound, "GAME_NOT_FOUND")
		return
	}
	initialSetup := room.Status == "playing" && room.Domain.Period == domain.PeriodZero && room.Domain.Phase == domain.PhaseHypothesis
	if room.Status != "lobby" && !initialSetup {
		writeCode(w, http.StatusConflict, "GAME_ALREADY_STARTED")
		return
	}
	var input setupRequest
	if err := decodeBody(r, &input); err != nil {
		writeCode(w, http.StatusBadRequest, "INVALID_REQUEST")
		return
	}
	if input.Token != room.Token {
		writeCode(w, http.StatusUnauthorized, "INVALID_SESSION")
		return
	}
	if err := room.Domain.SetInitialCards(room.PlayerID, input.PartnerID, input.StarterShopID); err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, room.view(input.Token))
}

func (s *gameStore) startHandler(w http.ResponseWriter, r *http.Request) {
	s.Lock()
	defer s.Unlock()
	room, ok := s.gameLocked(r.PathValue("id"))
	if !ok {
		writeCode(w, http.StatusNotFound, "GAME_NOT_FOUND")
		return
	}
	token := tokenFrom(r)
	if token != room.Token {
		writeCode(w, http.StatusUnauthorized, "INVALID_SESSION")
		return
	}
	if len(room.Domain.Players) < 4 {
		if err := addBots(room); err != nil {
			writeDomainError(w, err)
			return
		}
	}
	if err := room.Domain.Start(); err != nil {
		writeDomainError(w, err)
		return
	}
	room.Status = "playing"
	room.Version++
	writeJSON(w, http.StatusOK, room.view(token))
}

func (s *gameStore) stateHandler(w http.ResponseWriter, r *http.Request) {
	s.RLock()
	defer s.RUnlock()
	room, ok := s.gameLocked(r.PathValue("id"))
	if !ok {
		writeCode(w, http.StatusNotFound, "GAME_NOT_FOUND")
		return
	}
	writeJSON(w, http.StatusOK, room.view(tokenFrom(r)))
}

func (s *gameStore) commandHandler(w http.ResponseWriter, r *http.Request) {
	s.Lock()
	defer s.Unlock()
	room, ok := s.gameLocked(r.PathValue("id"))
	if !ok {
		writeCode(w, http.StatusNotFound, "GAME_NOT_FOUND")
		return
	}
	if room.Status != "playing" {
		writeCode(w, http.StatusConflict, "GAME_NOT_STARTED")
		return
	}
	var input commandRequest
	if err := decodeBody(r, &input); err != nil {
		writeCode(w, http.StatusBadRequest, "INVALID_REQUEST")
		return
	}
	if input.Token != room.Token {
		writeCode(w, http.StatusUnauthorized, "INVALID_SESSION")
		return
	}
	key := commandKey(input)
	if cached, ok := room.Processed[key]; ok {
		writeRawJSON(w, cached.Status, cached.Body)
		return
	}
	if input.GameVersion != room.Version {
		writeJSON(w, http.StatusConflict, map[string]any{"code": "VERSION_CONFLICT", "gameVersion": room.Version, "state": room.view(input.Token)})
		return
	}
	if err := applyCommand(room, room.PlayerID, input); err != nil {
		writeDomainError(w, err)
		return
	}
	room.Version++
	runBots(room)
	if err := passHandsIfReady(room); err != nil {
		writeDomainError(w, err)
		return
	}
	if allKPIsSelected(room) {
		if err := beginNextPeriod(room); err != nil {
			writeDomainError(w, err)
			return
		}
	}
	if room.Domain.Phase == domain.PhaseFinished {
		room.Status = "finished"
	}
	result := map[string]any{"gameVersion": room.Version, "state": room.view(input.Token)}
	body, _ := json.Marshal(result)
	room.Processed[key] = commandResult{Status: http.StatusOK, Body: body}
	writeRawJSON(w, http.StatusOK, body)
}

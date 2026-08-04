package main

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	fixturedata "cafestartups/data"
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

func main() {
	store := &gameStore{games: make(map[string]*gameRoom)}
	server := &http.Server{Addr: ":8080", Handler: withCORS(newHandler(store))}
	log.Printf("Cafe Startups API listening on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func newHandler(store *gameStore) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("/api/games", store.gamesHandler)
	mux.HandleFunc("POST /api/games/{id}/setup", store.setupHandler)
	mux.HandleFunc("POST /api/games/{id}/start", store.startHandler)
	mux.HandleFunc("GET /api/games/{id}", store.stateHandler)
	mux.HandleFunc("POST /api/games/{id}/commands", store.commandHandler)
	return withCORS(mux)
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
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
		catalog, err := domain.LoadCatalog(fixturedata.MVPFixture)
		if err != nil {
			writeCode(w, http.StatusInternalServerError, "INVALID_CARD_CATALOG")
			return
		}
		g.SetCatalog(catalog)
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
	if room.Status != "lobby" {
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
	if err := room.Domain.BeginExperiment(); err != nil {
		writeDomainError(w, err)
		return
	}
	room.Status = "playing"
	room.Version++
	runBots(room)
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
	err := applyCommand(room, room.PlayerID, input)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	room.Version++
	runBots(room)
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

func applyCommand(room *gameRoom, playerID string, input commandRequest) error {
	switch input.Type {
	case "BEGIN_EXPERIMENT":
		return room.Domain.BeginExperiment()
	case "SET_KPI":
		return room.Domain.SetKPIs(playerID, input.KPIs...)
	case "SELECT_CARD":
		return room.Domain.SelectCard(playerID, input.CardID)
	case "PLAY_SELECTED_CARD":
		return room.Domain.PlaySelectedCard(playerID)
	case "DISCARD_SELECTED_CARD":
		return room.Domain.DiscardSelectedCard(playerID)
	case "PASS_HAND":
		return room.Domain.PassHands()
	case "TAKE_LOAN":
		return room.Domain.TakeLoan(playerID)
	case "REPAY_LOAN":
		return room.Domain.RepayLoan(playerID, input.Count)
	case "CONFIRM_INTEREST":
		return room.Domain.SettleInterest()
	case "CONFIRM_REVENUE":
		room.Domain.SettleRevenue()
		return nil
	case "RESOLVE_LEARNING":
		return room.Domain.ResolveLearning()
	default:
		return domain.ErrInvalidAction
	}
}

func (room *gameRoom) view(token string) map[string]any {
	players := make([]map[string]any, 0, len(room.Domain.Players))
	for _, p := range room.Domain.Players {
		players = append(players, map[string]any{"id": p.ID, "displayName": p.DisplayName, "bot": p.IsBot, "cash": p.Cash, "loans": p.Loans, "revenue": p.Revenue, "score": p.Score, "selectedKPIs": p.SelectedKPIs, "brandAwareness": p.BrandAwareness, "products": p.Products, "values": p.Values, "resources": p.Resources, "handCount": len(p.Hand)})
	}
	result := map[string]any{"id": room.ID, "status": room.Status, "seed": room.Seed, "gameVersion": room.Version, "period": room.Domain.Period, "phase": room.Domain.Phase, "round": room.Domain.Round, "demandBoard": room.Domain.DemandBoard, "center": room.Domain.Center, "partnerOptions": room.Domain.PartnerOptions, "starterShopOptions": room.Domain.StarterShopOptions, "players": players}
	if token == room.Token {
		for _, p := range room.Domain.Players {
			if p.ID == room.PlayerID {
				result["me"] = map[string]any{"id": p.ID, "hand": p.Hand, "tableau": p.Tableau, "discardCount": len(p.Discard), "partner": p.Partner, "starterShop": p.StarterShop, "cash": p.Cash, "loans": p.Loans, "customers": p.Customers, "revenue": p.Revenue, "score": p.Score, "selectedKPIs": p.SelectedKPIs, "cashFlow": p.CashFlow, "cashFlowRounds": p.CashFlowRounds, "brandAwareness": p.BrandAwareness, "products": p.Products, "values": p.Values, "resources": p.Resources}
			}
		}
	}
	return result
}

func beginNextPeriod(room *gameRoom) error {
	for _, player := range room.Domain.Players {
		if player.IsBot && player.KPISelectionPeriod != room.Domain.Period {
			if err := setRandomBotKPIs(room, player); err != nil {
				return err
			}
		}
	}
	if err := room.Domain.BeginExperiment(); err != nil {
		return err
	}
	room.Version++
	return nil
}

func allKPIsSelected(room *gameRoom) bool {
	if room.Domain.Phase != domain.PhaseHypothesis || room.Domain.Period == domain.PeriodOne {
		return false
	}
	for _, player := range room.Domain.Players {
		if len(player.SelectedKPIs) != 2 || player.KPISelectionPeriod != room.Domain.Period {
			return false
		}
	}
	return len(room.Domain.Players) > 0
}

func setRandomBotKPIs(room *gameRoom, player *domain.Player) error {
	kpis := []string{"brand_awareness", "products", "values", "resources"}
	i := botChoice(room, player.ID+"-kpi", len(kpis))
	j := botChoice(room, player.ID+"-kpi-second", len(kpis)-1)
	if j >= i {
		j++
	}
	return room.Domain.SetKPIs(player.ID, kpis[i], kpis[j])
}

func addBots(room *gameRoom) error {
	for len(room.Domain.Players) < 4 {
		index := len(room.Domain.Players) + 1
		id := fmt.Sprintf("bot-%d", index)
		if err := room.Domain.AddPlayer(id, fmt.Sprintf("電腦玩家 %d", index)); err != nil {
			return err
		}
		room.Domain.Players[len(room.Domain.Players)-1].IsBot = true
		room.Version++
	}
	return nil
}

// runBots performs intentionally naive, deterministic random actions. It is
// only a pacing helper for solo MVP games; all actual legality checks remain in
// the domain layer.
func runBots(room *gameRoom) {
	if room.Domain.Phase == domain.PhaseHypothesis {
		for _, player := range room.Domain.Players {
			if player.IsBot && player.KPISelectionPeriod != room.Domain.Period {
				kpis := []string{"brand_awareness", "products", "values", "resources"}
				i := botChoice(room, player.ID, len(kpis))
				j := botChoice(room, player.ID+"-second", len(kpis)-1)
				if j >= i {
					j++
				}
				if err := room.Domain.SetKPIs(player.ID, kpis[i], kpis[j]); err == nil {
					room.Version++
				}
			}
		}
	}
	if room.Domain.Phase != domain.PhaseExperiment {
		return
	}
	for _, player := range room.Domain.Players {
		if !player.IsBot {
			continue
		}
		if !room.Domain.HasSelected(player.ID) {
			if len(player.Hand) == 0 {
				continue
			}
			card := player.Hand[botChoice(room, player.ID, len(player.Hand))]
			if err := room.Domain.SelectCard(player.ID, card.ID); err == nil {
				room.Version++
			}
		}
		if room.Domain.HasActed(player.ID) {
			continue
		}
		if botChoice(room, player.ID+"-action", 2) == 0 {
			if err := room.Domain.PlaySelectedCard(player.ID); err == nil {
				room.Version++
				continue
			}
		}
		if err := room.Domain.DiscardSelectedCard(player.ID); err == nil {
			room.Version++
		}
	}
}

func botChoice(room *gameRoom, key string, size int) int {
	if size <= 1 {
		return 0
	}
	sum := sha1.Sum([]byte(room.Seed + "|" + key + "|" + itoa(room.Version)))
	return int(sum[0]) % size
}
func tokenFrom(r *http.Request) string {
	if t := r.Header.Get("X-Session-Token"); t != "" {
		return t
	}
	return r.URL.Query().Get("token")
}
func decodeBody(r *http.Request, v any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}
func writeCode(w http.ResponseWriter, status int, code string) {
	writeJSON(w, status, map[string]string{"code": code})
}
func writeDomainError(w http.ResponseWriter, err error) {
	status := http.StatusBadRequest
	if errors.Is(err, domain.ErrInsufficientCash) || errors.Is(err, domain.ErrLoanLimit) {
		status = http.StatusConflict
	}
	writeCode(w, status, err.Error())
}
func randomID(size int) string {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeRawJSON(w http.ResponseWriter, status int, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(body)
	_, _ = w.Write([]byte("\n"))
}

func commandKey(input commandRequest) string {
	if input.CommandID != "" {
		return input.Token + "|" + input.CommandID
	}
	b, _ := json.Marshal(input)
	checksum := sha1.Sum(b)
	return hex.EncodeToString(checksum[:])
}
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Session-Token")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func itoa(n int64) string {
	const digits = "0123456789"
	if n == 0 {
		return "0"
	}
	out := ""
	for n > 0 {
		out = string(digits[n%10]) + out
		n /= 10
	}
	return out
}

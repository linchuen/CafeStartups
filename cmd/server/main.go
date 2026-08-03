package main

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	fixturedata "cafestartups/data"
	"cafestartups/internal/domain"
)

type game struct {
	ID       string `json:"id"`
	RoomCode string `json:"roomCode"`
	Status   string `json:"status"`
}
type session struct {
	Token, PlayerID string
	Ready           bool
}
type gameRoom struct {
	game
	Seed      string
	Domain    *domain.Game
	HostToken string
	Sessions  map[string]*session
	Version   int64
	Events    []event
	Processed map[string]commandResult
}
type commandResult struct {
	Status int
	Body   []byte
}
type event struct {
	Version  int64  `json:"version"`
	Type     string `json:"type"`
	PlayerID string `json:"playerId,omitempty"`
}
type gameStore struct {
	sync.RWMutex
	games map[string]*gameRoom
}

// roomLocked resolves either the stable game id or the six-character room code.
// Callers must hold the store lock.
func (s *gameStore) roomLocked(key string) (*gameRoom, bool) {
	if room, ok := s.games[key]; ok {
		return room, true
	}
	for _, room := range s.games {
		if strings.EqualFold(room.RoomCode, key) {
			return room, true
		}
	}
	return nil, false
}

type joinRequest struct {
	DisplayName string `json:"displayName"`
}
type readyRequest struct {
	Token string `json:"token"`
	Ready *bool  `json:"ready"`
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
	mux.HandleFunc("POST /api/games/{id}/join", store.joinHandler)
	mux.HandleFunc("POST /api/games/{id}/ready", store.readyHandler)
	mux.HandleFunc("POST /api/games/{id}/setup", store.setupHandler)
	mux.HandleFunc("POST /api/games/{id}/start", store.startHandler)
	mux.HandleFunc("GET /api/games/{id}", store.stateHandler)
	mux.HandleFunc("POST /api/games/{id}/commands", store.commandHandler)
	mux.HandleFunc("GET /api/games/{id}/events", store.eventsHandler)
	return withCORS(mux)
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *gameStore) gamesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var input struct {
			Seed string `json:"seed"`
		}
		_ = json.NewDecoder(r.Body).Decode(&input)
		id := randomID(8)
		seed := input.Seed
		if seed == "" {
			seed = id
		}
		g, err := domain.NewGame(seed, nil)
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
		room := &gameRoom{game: game{ID: id, RoomCode: strings.ToUpper(id[:6]), Status: "lobby"}, Seed: seed, Domain: g, Sessions: map[string]*session{}, Processed: map[string]commandResult{}}
		s.Lock()
		s.games[id] = room
		s.Unlock()
		writeJSON(w, http.StatusCreated, roomSummary(room))
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

func (s *gameStore) joinHandler(w http.ResponseWriter, r *http.Request) {
	s.Lock()
	defer s.Unlock()
	room, ok := s.roomLocked(r.PathValue("id"))
	if !ok {
		writeCode(w, http.StatusNotFound, "GAME_NOT_FOUND")
		return
	}
	if room.Status != "lobby" {
		writeCode(w, http.StatusConflict, "GAME_ALREADY_STARTED")
		return
	}
	if len(room.Sessions) >= 1 {
		writeCode(w, http.StatusConflict, "SINGLE_PLAYER_MVP")
		return
	}
	if len(room.Domain.Players) >= 4 {
		writeCode(w, http.StatusConflict, "ROOM_FULL")
		return
	}
	var input joinRequest
	if err := decodeBody(r, &input); err != nil || strings.TrimSpace(input.DisplayName) == "" {
		writeCode(w, http.StatusBadRequest, "INVALID_REQUEST")
		return
	}
	token := randomID(16)
	playerID := "player-" + randomID(6)
	if err := room.Domain.AddPlayer(playerID, input.DisplayName); err != nil {
		writeDomainError(w, err)
		return
	}
	room.Sessions[token] = &session{Token: token, PlayerID: playerID}
	if len(room.Sessions) == 1 {
		room.HostToken = token
	}
	room.Version++
	room.Events = append(room.Events, event{room.Version, "PLAYER_JOINED", playerID})
	writeJSON(w, http.StatusCreated, map[string]any{"token": token, "playerId": playerID, "state": room.view(token)})
}

func (s *gameStore) readyHandler(w http.ResponseWriter, r *http.Request) {
	s.Lock()
	defer s.Unlock()
	room, ok := s.roomLocked(r.PathValue("id"))
	if !ok {
		writeCode(w, http.StatusNotFound, "GAME_NOT_FOUND")
		return
	}
	if room.Status != "lobby" {
		writeCode(w, http.StatusConflict, "GAME_ALREADY_STARTED")
		return
	}
	var input readyRequest
	if err := decodeBody(r, &input); err != nil {
		writeCode(w, http.StatusBadRequest, "INVALID_REQUEST")
		return
	}
	sess, ok := room.Sessions[input.Token]
	if !ok {
		writeCode(w, http.StatusUnauthorized, "INVALID_SESSION")
		return
	}
	ready := true
	if input.Ready != nil {
		ready = *input.Ready
	}
	sess.Ready = ready
	room.Version++
	room.Events = append(room.Events, event{room.Version, "PLAYER_READY", sess.PlayerID})
	writeJSON(w, http.StatusOK, room.view(input.Token))
}

func (s *gameStore) setupHandler(w http.ResponseWriter, r *http.Request) {
	s.Lock()
	defer s.Unlock()
	room, ok := s.roomLocked(r.PathValue("id"))
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
	sess, ok := room.Sessions[input.Token]
	if !ok {
		writeCode(w, http.StatusUnauthorized, "INVALID_SESSION")
		return
	}
	if err := room.Domain.SetInitialCards(sess.PlayerID, input.PartnerID, input.StarterShopID); err != nil {
		writeDomainError(w, err)
		return
	}
	room.Version++
	room.Events = append(room.Events, event{room.Version, "INITIAL_CARDS_SELECTED", sess.PlayerID})
	writeJSON(w, http.StatusOK, room.view(input.Token))
}

func (s *gameStore) startHandler(w http.ResponseWriter, r *http.Request) {
	s.Lock()
	defer s.Unlock()
	room, ok := s.roomLocked(r.PathValue("id"))
	if !ok {
		writeCode(w, http.StatusNotFound, "GAME_NOT_FOUND")
		return
	}
	token := tokenFrom(r)
	if token != room.HostToken {
		writeCode(w, http.StatusForbidden, "ONLY_HOST_CAN_START")
		return
	}
	if len(room.Domain.Players) < 4 {
		if err := addBots(room); err != nil {
			writeDomainError(w, err)
			return
		}
	}
	if len(room.Sessions) < 1 {
		writeCode(w, http.StatusConflict, "NOT_ENOUGH_PLAYERS")
		return
	}
	for _, sess := range room.Sessions {
		if !sess.Ready {
			writeCode(w, http.StatusConflict, "PLAYERS_NOT_READY")
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
	room.Events = append(room.Events, event{Version: room.Version, Type: "GAME_STARTED"})
	runBots(room)
	writeJSON(w, http.StatusOK, room.view(token))
}

func (s *gameStore) stateHandler(w http.ResponseWriter, r *http.Request) {
	s.RLock()
	defer s.RUnlock()
	room, ok := s.roomLocked(r.PathValue("id"))
	if !ok {
		writeCode(w, http.StatusNotFound, "GAME_NOT_FOUND")
		return
	}
	writeJSON(w, http.StatusOK, room.view(tokenFrom(r)))
}

func (s *gameStore) commandHandler(w http.ResponseWriter, r *http.Request) {
	s.Lock()
	defer s.Unlock()
	room, ok := s.roomLocked(r.PathValue("id"))
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
	sess, ok := room.Sessions[input.Token]
	if !ok {
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
	err := applyCommand(room, sess.PlayerID, input)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	room.Version++
	room.Events = append(room.Events, event{room.Version, input.Type, sess.PlayerID})
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
	result := map[string]any{"gameVersion": room.Version, "event": room.Events[len(room.Events)-1], "state": room.view(input.Token)}
	body, _ := json.Marshal(result)
	room.Processed[key] = commandResult{Status: http.StatusOK, Body: body}
	writeRawJSON(w, http.StatusOK, body)
}

func applyCommand(room *gameRoom, playerID string, input commandRequest) error {
	switch input.Type {
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

// events is an SSE stream: it keeps the Phase 2 event contract usable without
// adding a third-party WebSocket dependency. A WebSocket adapter can consume the
// same versioned events in a later transport-only change.
func (s *gameStore) eventsHandler(w http.ResponseWriter, r *http.Request) {
	if strings.EqualFold(r.Header.Get("Upgrade"), "websocket") {
		s.websocketEvents(w, r)
		return
	}
	s.RLock()
	room, ok := s.roomLocked(r.PathValue("id"))
	if !ok {
		s.RUnlock()
		writeCode(w, http.StatusNotFound, "GAME_NOT_FOUND")
		return
	}
	since := int64(0)
	if v := r.URL.Query().Get("since"); v != "" {
		_, _ = fmt.Sscan(v, &since)
	}
	events := make([]event, 0)
	for _, e := range room.Events {
		if e.Version > since {
			events = append(events, e)
		}
	}
	s.RUnlock()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	for _, e := range events {
		b, _ := json.Marshal(e)
		_, _ = w.Write([]byte("id: " + itoa(e.Version) + "\ndata: " + string(b) + "\n\n"))
	}
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

func (s *gameStore) websocketEvents(w http.ResponseWriter, r *http.Request) {
	s.RLock()
	_, exists := s.roomLocked(r.PathValue("id"))
	s.RUnlock()
	if !exists {
		writeCode(w, http.StatusNotFound, "GAME_NOT_FOUND")
		return
	}
	key := r.Header.Get("Sec-WebSocket-Key")
	if key == "" {
		writeCode(w, http.StatusBadRequest, "INVALID_WEBSOCKET_HANDSHAKE")
		return
	}
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		writeCode(w, http.StatusInternalServerError, "WEBSOCKET_UNSUPPORTED")
		return
	}
	conn, rw, err := hijacker.Hijack()
	if err != nil {
		return
	}
	acceptHash := sha1.Sum([]byte(key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
	accept := base64.StdEncoding.EncodeToString(acceptHash[:])
	_, _ = rw.WriteString("HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: " + accept + "\r\n\r\n")
	if err := rw.Flush(); err != nil {
		_ = conn.Close()
		return
	}
	defer conn.Close()
	since := int64(0)
	for i := 0; i < 120; i++ { // up to two minutes; clients reconnect with since=<version>
		events := s.eventsSince(r.PathValue("id"), since)
		for _, e := range events {
			b, _ := json.Marshal(e)
			if err := writeWebSocketText(conn, b); err != nil {
				return
			}
			since = e.Version
		}
		time.Sleep(time.Second)
	}
}

func (s *gameStore) eventsSince(id string, since int64) []event {
	s.RLock()
	defer s.RUnlock()
	room, ok := s.roomLocked(id)
	if !ok {
		return nil
	}
	result := make([]event, 0)
	for _, e := range room.Events {
		if e.Version > since {
			result = append(result, e)
		}
	}
	return result
}

func writeWebSocketText(conn interface{ Write([]byte) (int, error) }, payload []byte) error {
	if len(payload) >= 126 {
		return errors.New("websocket event too large")
	}
	frame := append([]byte{0x81, byte(len(payload))}, payload...)
	_, err := conn.Write(frame)
	return err
}

func (room *gameRoom) view(token string) map[string]any {
	players := make([]map[string]any, 0, len(room.Domain.Players))
	for _, p := range room.Domain.Players {
		ready := p.IsBot
		for _, sess := range room.Sessions {
			if sess.PlayerID == p.ID {
				ready = sess.Ready
			}
		}
		players = append(players, map[string]any{"id": p.ID, "displayName": p.DisplayName, "bot": p.IsBot, "ready": ready, "cash": p.Cash, "loans": p.Loans, "revenue": p.Revenue, "score": p.Score, "selectedKPIs": p.SelectedKPIs, "brandAwareness": p.BrandAwareness, "products": p.Products, "values": p.Values, "resources": p.Resources, "handCount": len(p.Hand)})
	}
	result := map[string]any{"id": room.ID, "roomCode": room.RoomCode, "status": room.Status, "seed": room.Seed, "gameVersion": room.Version, "period": room.Domain.Period, "phase": room.Domain.Phase, "round": room.Domain.Round, "demandBoard": room.Domain.DemandBoard, "center": room.Domain.Center, "partnerOptions": room.Domain.PartnerOptions, "starterShopOptions": room.Domain.StarterShopOptions, "players": players}
	if sess := room.Sessions[token]; sess != nil {
		for _, p := range room.Domain.Players {
			if p.ID == sess.PlayerID {
				result["me"] = map[string]any{"id": p.ID, "hand": p.Hand, "tableau": p.Tableau, "discardCount": len(p.Discard), "partner": p.Partner, "starterShop": p.StarterShop, "cash": p.Cash, "loans": p.Loans, "customers": p.Customers, "revenue": p.Revenue, "score": p.Score, "selectedKPIs": p.SelectedKPIs, "brandAwareness": p.BrandAwareness, "products": p.Products, "values": p.Values, "resources": p.Resources}
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
	room.Events = append(room.Events, event{Version: room.Version, Type: "PERIOD_STARTED"})
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
		room.Events = append(room.Events, event{Version: room.Version, Type: "BOT_JOINED", PlayerID: id})
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
					room.Events = append(room.Events, event{room.Version, "BOT_SET_KPI", player.ID})
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
				room.Events = append(room.Events, event{room.Version, "BOT_SELECT_CARD", player.ID})
			}
		}
		if room.Domain.HasActed(player.ID) {
			continue
		}
		if botChoice(room, player.ID+"-action", 2) == 0 {
			if err := room.Domain.PlaySelectedCard(player.ID); err == nil {
				room.Version++
				room.Events = append(room.Events, event{room.Version, "BOT_PLAY_SELECTED_CARD", player.ID})
				continue
			}
		}
		if err := room.Domain.DiscardSelectedCard(player.ID); err == nil {
			room.Version++
			room.Events = append(room.Events, event{room.Version, "BOT_DISCARD_SELECTED_CARD", player.ID})
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
func roomSummary(room *gameRoom) game { return room.game }
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
	return input.Token + "|" + hex.EncodeToString(checksum[:])
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

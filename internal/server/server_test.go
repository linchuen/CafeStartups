package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cafestartups/internal/domain"
)

func TestHealth(t *testing.T) {
	recorder := httptest.NewRecorder()
	healthHandler(recorder, httptest.NewRequest(http.MethodGet, "/health", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
	var body map[string]string
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatalf("decode health response: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("expected status ok, got %q", body["status"])
	}
}

func TestLocalGameLifecycleAndVersionConflict(t *testing.T) {
	store := &gameStore{games: make(map[string]*gameRoom)}
	handler := newHandler(store)
	create := httptest.NewRecorder()
	handler.ServeHTTP(create, httptest.NewRequest(http.MethodPost, "/api/games", strings.NewReader(`{"seed":"test-seed","displayName":"Alice"}`)))
	var created struct {
		ID, Token, PlayerID string
	}
	var response map[string]any
	if err := json.NewDecoder(create.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	created.ID = response["id"].(string)
	created.Token = response["token"].(string)
	created.PlayerID = response["playerId"].(string)
	if created.ID == "" || created.Token == "" || created.PlayerID == "" {
		t.Fatalf("unexpected local game response: %+v", response)
	}

	start := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/games/"+created.ID+"/start", nil)
	req.Header.Set("X-Session-Token", created.Token)
	handler.ServeHTTP(start, req)
	if start.Code != http.StatusOK {
		t.Fatalf("start status=%d body=%s", start.Code, start.Body.String())
	}
	var started map[string]any
	_ = json.NewDecoder(start.Body).Decode(&started)
	version := int64(started["gameVersion"].(float64))
	if started["period"].(float64) != 0 {
		t.Fatalf("start should leave game in period zero: %+v", started)
	}
	begin := httptest.NewRecorder()
	beginBody := `{"token":"` + created.Token + `","gameVersion":` + itoa(version) + `,"commandId":"begin-experiment","type":"BEGIN_EXPERIMENT"}`
	handler.ServeHTTP(begin, httptest.NewRequest(http.MethodPost, "/api/games/"+created.ID+"/commands", strings.NewReader(beginBody)))
	if begin.Code != http.StatusOK {
		t.Fatalf("begin experiment status=%d body=%s", begin.Code, begin.Body.String())
	}
	var begun map[string]any
	_ = json.NewDecoder(begin.Body).Decode(&begun)
	version = int64(begun["gameVersion"].(float64))
	state := httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/games/"+created.ID+"?token="+created.Token, nil)
	handler.ServeHTTP(state, req)
	var visible map[string]any
	_ = json.NewDecoder(state.Body).Decode(&visible)
	if _, ok := visible["me"]; !ok {
		t.Fatal("player should receive private view")
	}
	conflict := httptest.NewRecorder()
	body := `{"token":"` + created.Token + `","gameVersion":` + itoa(version-1) + `,"type":"TAKE_LOAN"}`
	handler.ServeHTTP(conflict, httptest.NewRequest(http.MethodPost, "/api/games/"+created.ID+"/commands", strings.NewReader(body)))
	if conflict.Code != http.StatusConflict {
		t.Fatalf("conflict status=%d", conflict.Code)
	}
	loanBody := `{"token":"` + created.Token + `","gameVersion":` + itoa(version) + `,"commandId":"loan-once","type":"TAKE_LOAN"}`
	loan := httptest.NewRecorder()
	handler.ServeHTTP(loan, httptest.NewRequest(http.MethodPost, "/api/games/"+created.ID+"/commands", strings.NewReader(loanBody)))
	if loan.Code != http.StatusOK {
		t.Fatalf("loan status=%d body=%s", loan.Code, loan.Body.String())
	}
	duplicate := httptest.NewRecorder()
	handler.ServeHTTP(duplicate, httptest.NewRequest(http.MethodPost, "/api/games/"+created.ID+"/commands", strings.NewReader(loanBody)))
	if duplicate.Code != http.StatusOK {
		t.Fatalf("duplicate status=%d", duplicate.Code)
	}
	if store.games[created.ID].Domain.Players[0].Loans != 1 {
		t.Fatalf("duplicate command changed loans to %d", store.games[created.ID].Domain.Players[0].Loans)
	}
}

func TestCreateGame(t *testing.T) {
	store := &gameStore{games: make(map[string]*gameRoom)}
	recorder := httptest.NewRecorder()
	store.gamesHandler(recorder, httptest.NewRequest(http.MethodPost, "/api/games", nil))

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", recorder.Code)
	}
	var created struct {
		ID, Token string
		State     map[string]any `json:"state"`
	}
	if err := json.NewDecoder(recorder.Body).Decode(&created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if created.ID == "" || created.Token == "" || created.State["status"] != "lobby" {
		t.Fatalf("unexpected game response: %+v", created)
	}
}

func TestSoloStartAddsRandomBots(t *testing.T) {
	store := &gameStore{games: make(map[string]*gameRoom)}
	handler := newHandler(store)
	create := httptest.NewRecorder()
	handler.ServeHTTP(create, httptest.NewRequest(http.MethodPost, "/api/games", strings.NewReader(`{"seed":"solo-seed"}`)))
	var created struct {
		ID, Token string
	}
	if err := json.NewDecoder(create.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}
	room := game{ID: created.ID}
	token := created.Token
	start := httptest.NewRecorder()
	startReq := httptest.NewRequest(http.MethodPost, "/api/games/"+room.ID+"/start", nil)
	startReq.Header.Set("X-Session-Token", token)
	handler.ServeHTTP(start, startReq)
	if start.Code != http.StatusOK {
		t.Fatalf("start status=%d body=%s", start.Code, start.Body.String())
	}
	if err := store.games[room.ID].Domain.BeginExperiment(); err != nil {
		t.Fatal(err)
	}
	runBots(store.games[room.ID])
	if got := store.games[room.ID].Domain.Players[0].SelectedKPIs; len(got) != 0 {
		t.Fatalf("KPIs must not be selected before period one ends: %v", got)
	}
	if len(store.games[room.ID].Domain.Players) != 4 {
		t.Fatalf("expected 4 players, got %d", len(store.games[room.ID].Domain.Players))
	}
	for _, player := range store.games[room.ID].Domain.Players[1:] {
		if !player.IsBot || len(player.Hand) < 6 || len(player.Hand) > 7 || !store.games[room.ID].Domain.HasActed(player.ID) {
			t.Fatalf("bot was not initialized and acted: %+v", player)
		}
	}
}

func TestSoloGameCompletesThroughHTTP(t *testing.T) {
	store := &gameStore{games: make(map[string]*gameRoom)}
	handler := newHandler(store)
	create := httptest.NewRecorder()
	handler.ServeHTTP(create, httptest.NewRequest(http.MethodPost, "/api/games", strings.NewReader(`{"seed":"full-solo-seed"}`)))
	var created struct {
		ID, Token string
	}
	if err := json.NewDecoder(create.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}
	room := game{ID: created.ID}
	token := created.Token
	start := httptest.NewRecorder()
	startRequest := httptest.NewRequest(http.MethodPost, "/api/games/"+room.ID+"/start", nil)
	startRequest.Header.Set("X-Session-Token", token)
	handler.ServeHTTP(start, startRequest)
	if start.Code != http.StatusOK {
		t.Fatalf("start status=%d body=%s", start.Code, start.Body.String())
	}

	commandNumber := 0
	command := func(commandType, cardID string) {
		commandNumber++
		payload := map[string]any{
			"token":       token,
			"gameVersion": store.games[room.ID].Version,
			"commandId":   fmt.Sprintf("full-solo-%d", commandNumber),
			"type":        commandType,
		}
		if cardID != "" {
			payload["cardId"] = cardID
		}
		if commandType == "SET_KPI" {
			payload["kpis"] = []string{"values", "resources"}
		}
		body, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/api/games/"+room.ID+"/commands", strings.NewReader(string(body))))
		if recorder.Code != http.StatusOK {
			t.Fatalf("command %s status=%d body=%s", commandType, recorder.Code, recorder.Body.String())
		}
	}

	command("BEGIN_EXPERIMENT", "")

	for period := 1; period <= 3; period++ {
		for round := domain.InitialRound; round < domain.ExperimentRounds; round++ {
			human := store.games[room.ID].PlayerID
			player := store.games[room.ID].Domain.Players[0]
			for _, candidate := range store.games[room.ID].Domain.Players {
				if candidate.ID == human {
					player = candidate
					break
				}
			}
			if len(player.Hand) == 0 {
				t.Fatalf("period %d round %d human has no cards", period, round)
			}
			command("SELECT_CARD", player.Hand[0].ID)
			command("DISCARD_SELECTED_CARD", "")
			command("PASS_HAND", "")
		}
		if store.games[room.ID].Domain.Phase != domain.PhaseLearning {
			t.Fatalf("period %d did not reach learning: %s", period, store.games[room.ID].Domain.Phase)
		}
		command("DRAW_MARKET", "")
		command("RESOLVE_LEARNING", "")
		if got := len(store.games[room.ID].Domain.Players[0].CashFlow); got != period {
			t.Fatalf("period %d cash-flow statements=%d", period, got)
		}
		if period < 3 {
			command("SET_KPI", "")
		}
	}
	if store.games[room.ID].Status != "finished" || store.games[room.ID].Domain.Phase != domain.PhaseFinished {
		t.Fatalf("game did not finish: status=%s phase=%s", store.games[room.ID].Status, store.games[room.ID].Domain.Phase)
	}
}

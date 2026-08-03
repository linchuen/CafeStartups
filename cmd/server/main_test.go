package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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

func TestRoomLifecycleAndVersionConflict(t *testing.T) {
	store := &gameStore{games: make(map[string]*gameRoom)}
	handler := newHandler(store)
	create := httptest.NewRecorder()
	handler.ServeHTTP(create, httptest.NewRequest(http.MethodPost, "/api/games", strings.NewReader(`{"seed":"test-seed"}`)))
	var room game
	if err := json.NewDecoder(create.Body).Decode(&room); err != nil {
		t.Fatal(err)
	}

	join := func(name string) (string, map[string]any) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/games/"+room.ID+"/join", strings.NewReader(`{"displayName":"`+name+`"}`))
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusCreated {
			t.Fatalf("join status=%d body=%s", rec.Code, rec.Body.String())
		}
		var body map[string]any
		if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		return body["token"].(string), body["state"].(map[string]any)
	}
	aToken, _ := join("Alice")
	bToken, _ := join("Bob")
	ready := func(token string) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/games/"+room.ID+"/ready", strings.NewReader(`{"token":"`+token+`"}`))
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("ready status=%d", rec.Code)
		}
	}
	ready(aToken)
	ready(bToken)
	start := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/games/"+room.ID+"/start", nil)
	req.Header.Set("X-Session-Token", aToken)
	handler.ServeHTTP(start, req)
	if start.Code != http.StatusOK {
		t.Fatalf("start status=%d body=%s", start.Code, start.Body.String())
	}
	var started map[string]any
	_ = json.NewDecoder(start.Body).Decode(&started)
	version := int64(started["gameVersion"].(float64))
	state := httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/games/"+room.ID+"?token="+bToken, nil)
	handler.ServeHTTP(state, req)
	var visible map[string]any
	_ = json.NewDecoder(state.Body).Decode(&visible)
	if _, ok := visible["me"]; !ok {
		t.Fatal("player should receive private view")
	}
	conflict := httptest.NewRecorder()
	body := `{"token":"` + aToken + `","gameVersion":` + itoa(version-1) + `,"type":"TAKE_LOAN"}`
	handler.ServeHTTP(conflict, httptest.NewRequest(http.MethodPost, "/api/games/"+room.ID+"/commands", strings.NewReader(body)))
	if conflict.Code != http.StatusConflict {
		t.Fatalf("conflict status=%d", conflict.Code)
	}
	loanBody := `{"token":"` + aToken + `","gameVersion":` + itoa(version) + `,"commandId":"loan-once","type":"TAKE_LOAN"}`
	loan := httptest.NewRecorder()
	handler.ServeHTTP(loan, httptest.NewRequest(http.MethodPost, "/api/games/"+room.ID+"/commands", strings.NewReader(loanBody)))
	if loan.Code != http.StatusOK {
		t.Fatalf("loan status=%d body=%s", loan.Code, loan.Body.String())
	}
	duplicate := httptest.NewRecorder()
	handler.ServeHTTP(duplicate, httptest.NewRequest(http.MethodPost, "/api/games/"+room.ID+"/commands", strings.NewReader(loanBody)))
	if duplicate.Code != http.StatusOK {
		t.Fatalf("duplicate status=%d", duplicate.Code)
	}
	playerID := store.games[room.ID].Sessions[aToken].PlayerID
	for _, player := range store.games[room.ID].Domain.Players {
		if player.ID == playerID && player.Loans != 1 {
			t.Fatalf("duplicate command changed loans to %d", player.Loans)
		}
	}
}

func TestCreateGame(t *testing.T) {
	store := &gameStore{games: make(map[string]*gameRoom)}
	recorder := httptest.NewRecorder()
	store.gamesHandler(recorder, httptest.NewRequest(http.MethodPost, "/api/games", nil))

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", recorder.Code)
	}
	var created game
	if err := json.NewDecoder(recorder.Body).Decode(&created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if created.ID == "" || len(created.RoomCode) != 6 || created.Status != "lobby" {
		t.Fatalf("unexpected game response: %+v", created)
	}
}

func TestSoloStartAddsRandomBots(t *testing.T) {
	store := &gameStore{games: make(map[string]*gameRoom)}
	handler := newHandler(store)
	create := httptest.NewRecorder()
	handler.ServeHTTP(create, httptest.NewRequest(http.MethodPost, "/api/games", strings.NewReader(`{"seed":"solo-seed"}`)))
	var room game
	if err := json.NewDecoder(create.Body).Decode(&room); err != nil {
		t.Fatal(err)
	}

	join := httptest.NewRecorder()
	handler.ServeHTTP(join, httptest.NewRequest(http.MethodPost, "/api/games/"+room.ID+"/join", strings.NewReader(`{"displayName":"Solo"}`)))
	var joined map[string]any
	if err := json.NewDecoder(join.Body).Decode(&joined); err != nil {
		t.Fatal(err)
	}
	token := joined["token"].(string)

	ready := httptest.NewRecorder()
	handler.ServeHTTP(ready, httptest.NewRequest(http.MethodPost, "/api/games/"+room.ID+"/ready", strings.NewReader(`{"token":"`+token+`"}`)))
	if ready.Code != http.StatusOK {
		t.Fatalf("ready status=%d body=%s", ready.Code, ready.Body.String())
	}
	start := httptest.NewRecorder()
	startReq := httptest.NewRequest(http.MethodPost, "/api/games/"+room.ID+"/start", nil)
	startReq.Header.Set("X-Session-Token", token)
	handler.ServeHTTP(start, startReq)
	if start.Code != http.StatusOK {
		t.Fatalf("start status=%d body=%s", start.Code, start.Body.String())
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

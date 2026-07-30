package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
)

type game struct {
	ID       string `json:"id"`
	RoomCode string `json:"roomCode"`
	Status   string `json:"status"`
}

type gameStore struct {
	sync.RWMutex
	games map[string]game
}

func main() {
	store := &gameStore{games: make(map[string]game)}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("/api/games", store.gamesHandler)

	server := &http.Server{Addr: ":8080", Handler: withCORS(mux)}
	log.Printf("Café Startups API listening on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *gameStore) gamesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := randomID(8)
		g := game{ID: id, RoomCode: strings.ToUpper(id[:6]), Status: "lobby"}
		s.Lock()
		s.games[g.ID] = g
		s.Unlock()
		writeJSON(w, http.StatusCreated, g)
		return
	}

	if r.Method == http.MethodGet {
		s.RLock()
		games := make([]game, 0, len(s.games))
		for _, g := range s.games {
			games = append(games, g)
		}
		s.RUnlock()
		writeJSON(w, http.StatusOK, games)
		return
	}

	w.Header().Set("Allow", "GET, POST")
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
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

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

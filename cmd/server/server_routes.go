package main

import "net/http"

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

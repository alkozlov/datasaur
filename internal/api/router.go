package api

import (
	"net/http"

	"block-flow/internal/api/handlers"
	"block-flow/internal/api/middleware"
	"block-flow/internal/engine"
	"block-flow/internal/storage"

	"github.com/gorilla/mux"
)

// NewRouter creates a new HTTP router with all routes configured
func NewRouter(engine *engine.Engine, storage storage.Storage) http.Handler {
	r := mux.NewRouter()

	// Apply middleware
	r.Use(middleware.CORS())
	r.Use(middleware.Logging())
	r.Use(middleware.Recovery())

	// Create handlers
	flowHandler := handlers.NewFlowHandler(engine, storage)
	blockHandler := handlers.NewBlockHandler(engine)
	wsHandler := handlers.NewWebSocketHandler(engine)

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()

	// Flow routes
	api.HandleFunc("/flows", flowHandler.ListFlows).Methods("GET")
	api.HandleFunc("/flows", flowHandler.CreateFlow).Methods("POST")
	api.HandleFunc("/flows/{id}", flowHandler.GetFlow).Methods("GET")
	api.HandleFunc("/flows/{id}", flowHandler.UpdateFlow).Methods("PUT")
	api.HandleFunc("/flows/{id}", flowHandler.DeleteFlow).Methods("DELETE")
	api.HandleFunc("/flows/{id}/start", flowHandler.StartFlow).Methods("POST")
	api.HandleFunc("/flows/{id}/run", flowHandler.StartFlow).Methods("POST") // Alias for start
	api.HandleFunc("/flows/{id}/stop", flowHandler.StopFlow).Methods("POST")
	api.HandleFunc("/flows/{id}/trigger", flowHandler.TriggerFlow).Methods("POST")
	api.HandleFunc("/flows/{id}/status", flowHandler.GetFlowStatus).Methods("GET")

	// Block routes
	api.HandleFunc("/blocks", blockHandler.ListBlocks).Methods("GET")
	api.HandleFunc("/blocks/{type}", blockHandler.GetBlockInfo).Methods("GET")

	// WebSocket route
	api.HandleFunc("/ws", wsHandler.HandleWebSocket).Methods("GET")

	// Health check
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}).Methods("GET")

	// Static files (for future frontend)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/public/")))

	return r
}

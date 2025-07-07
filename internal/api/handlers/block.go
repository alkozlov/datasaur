package handlers

import (
	"encoding/json"
	"net/http"

	"block-flow/internal/engine"

	"github.com/gorilla/mux"
)

// BlockHandler handles block-related HTTP requests
type BlockHandler struct {
	engine *engine.Engine
}

// NewBlockHandler creates a new block handler
func NewBlockHandler(engine *engine.Engine) *BlockHandler {
	return &BlockHandler{
		engine: engine,
	}
}

// ListBlocks handles GET /api/v1/blocks
func (h *BlockHandler) ListBlocks(w http.ResponseWriter, r *http.Request) {
	registry := h.engine.GetRegistry()
	blockInfo := registry.GetBlockInfo()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blockInfo)
}

// GetBlockInfo handles GET /api/v1/blocks/{type}
func (h *BlockHandler) GetBlockInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blockType := vars["type"]

	registry := h.engine.GetRegistry()
	blockInfo, err := registry.GetBlockInfoByType(blockType)
	if err != nil {
		http.Error(w, "Block type not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blockInfo)
}

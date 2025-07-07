package handlers

import (
	"encoding/json"
	"net/http"

	"block-flow/internal/engine"
	"block-flow/internal/models"
	"block-flow/internal/storage"

	"github.com/gorilla/mux"
)

// FlowHandler handles flow-related HTTP requests
type FlowHandler struct {
	engine  *engine.Engine
	storage storage.Storage
}

// NewFlowHandler creates a new flow handler
func NewFlowHandler(engine *engine.Engine, storage storage.Storage) *FlowHandler {
	return &FlowHandler{
		engine:  engine,
		storage: storage,
	}
}

// ListFlows handles GET /api/v1/flows
func (h *FlowHandler) ListFlows(w http.ResponseWriter, r *http.Request) {
	flows, err := h.storage.LoadAllFlows(r.Context())
	if err != nil {
		http.Error(w, "Failed to load flows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(flows)
}

// CreateFlow handles POST /api/v1/flows
func (h *FlowHandler) CreateFlow(w http.ResponseWriter, r *http.Request) {
	var flow models.Flow
	if err := json.NewDecoder(r.Body).Decode(&flow); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Generate ID if not provided
	if flow.ID == "" {
		flow = *models.NewFlow(flow.Name)
	}

	// Validate flow
	if err := flow.Validate(); err != nil {
		http.Error(w, "Flow validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Save flow
	if err := h.storage.SaveFlow(r.Context(), &flow); err != nil {
		http.Error(w, "Failed to save flow", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&flow)
}

// GetFlow handles GET /api/v1/flows/{id}
func (h *FlowHandler) GetFlow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	flowID := vars["id"]

	flow, err := h.storage.LoadFlow(r.Context(), flowID)
	if err != nil {
		http.Error(w, "Flow not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(flow)
}

// UpdateFlow handles PUT /api/v1/flows/{id}
func (h *FlowHandler) UpdateFlow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	flowID := vars["id"]

	var flow models.Flow
	if err := json.NewDecoder(r.Body).Decode(&flow); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Ensure ID matches
	flow.ID = flowID

	// Validate flow
	if err := flow.Validate(); err != nil {
		http.Error(w, "Flow validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Save flow
	if err := h.storage.SaveFlow(r.Context(), &flow); err != nil {
		http.Error(w, "Failed to save flow", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&flow)
}

// DeleteFlow handles DELETE /api/v1/flows/{id}
func (h *FlowHandler) DeleteFlow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	flowID := vars["id"]

	// Stop flow if running
	h.engine.StopFlow(r.Context(), flowID)

	// Delete flow
	if err := h.storage.DeleteFlow(r.Context(), flowID); err != nil {
		http.Error(w, "Failed to delete flow", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// StartFlow handles POST /api/v1/flows/{id}/run
func (h *FlowHandler) StartFlow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	flowID := vars["id"]

	if err := h.engine.StartFlow(r.Context(), flowID); err != nil {
		http.Error(w, "Failed to start flow: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "started"})
}

// StopFlow handles POST /api/v1/flows/{id}/stop
func (h *FlowHandler) StopFlow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	flowID := vars["id"]

	if err := h.engine.StopFlow(r.Context(), flowID); err != nil {
		http.Error(w, "Failed to stop flow: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "stopped"})
}

// TriggerFlow handles POST /api/v1/flows/{id}/trigger
func (h *FlowHandler) TriggerFlow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	flowID := vars["id"]

	var input models.Message
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		// If no input provided, create a default message
		input = *models.NewMessage(map[string]interface{}{"trigger": true})
	}

	if err := h.engine.TriggerFlow(r.Context(), flowID, &input); err != nil {
		http.Error(w, "Failed to trigger flow: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "triggered"})
}

// GetFlowStatus handles GET /api/v1/flows/{id}/status
func (h *FlowHandler) GetFlowStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	flowID := vars["id"]

	status, err := h.engine.GetFlowStatus(flowID)
	if err != nil {
		http.Error(w, "Flow not running", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

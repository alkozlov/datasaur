package engine

import (
	"context"
	"fmt"
	"sync"

	"block-flow/internal/blocks"
	"block-flow/internal/blocks/builtin"
	"block-flow/internal/models"
	"block-flow/internal/storage"
)

// Logger interface for engine logging
type Logger interface {
	Debug(message string, fields map[string]interface{})
	Info(message string, fields map[string]interface{})
	Warn(message string, fields map[string]interface{})
	Error(message string, fields map[string]interface{})
}

// SimpleLogger is a basic implementation of the Logger interface
type SimpleLogger struct{}

// Debug logs debug messages
func (l *SimpleLogger) Debug(message string, fields map[string]interface{}) {
	fmt.Printf("[DEBUG] %s: %v\n", message, fields)
}

// Info logs info messages
func (l *SimpleLogger) Info(message string, fields map[string]interface{}) {
	fmt.Printf("[INFO] %s: %v\n", message, fields)
}

// Warn logs warning messages
func (l *SimpleLogger) Warn(message string, fields map[string]interface{}) {
	fmt.Printf("[WARN] %s: %v\n", message, fields)
}

// Error logs error messages
func (l *SimpleLogger) Error(message string, fields map[string]interface{}) {
	fmt.Printf("[ERROR] %s: %v\n", message, fields)
}

// Engine manages flow execution using the new FlowExecutor
type Engine struct {
	storage  storage.Storage
	registry *blocks.Registry
	executor *FlowExecutor
	logger   Logger
	mu       sync.RWMutex
}

// New creates a new flow engine
func New(storage storage.Storage, logger Logger) *Engine {
	registry := blocks.NewRegistry()

	// Register built-in blocks
	builtin.RegisterBuiltinBlocks(registry)

	engine := &Engine{
		storage:  storage,
		registry: registry,
		executor: NewFlowExecutor(registry, logger),
		logger:   logger,
	}

	return engine
}

// LoadAndStartFlows loads all flows from storage and starts active ones
func (e *Engine) LoadAndStartFlows(ctx context.Context) error {
	flows, err := e.storage.LoadAllFlows(ctx)
	if err != nil {
		return fmt.Errorf("failed to load flows: %w", err)
	}

	for _, flow := range flows {
		if flow.Active {
			if err := e.StartFlow(ctx, flow.ID); err != nil {
				e.logger.Error("Failed to start flow", map[string]interface{}{
					"flow_id": flow.ID,
					"error":   err.Error(),
				})
				// Continue with other flows
				continue
			}
		}
	}

	return nil
}

// StartFlow starts execution of a flow
func (e *Engine) StartFlow(ctx context.Context, flowID string) error {
	// Load flow from storage
	flow, err := e.storage.LoadFlow(ctx, flowID)
	if err != nil {
		return fmt.Errorf("failed to load flow: %w", err)
	}

	// Validate flow
	if err := flow.Validate(); err != nil {
		return fmt.Errorf("flow validation failed: %w", err)
	}

	// Use the new executor to prepare and start the flow
	return e.executor.PrepareAndStartFlow(flow)
}

// StopFlow stops execution of a flow
func (e *Engine) StopFlow(ctx context.Context, flowID string) error {
	return e.executor.StopFlow(flowID)
}

// TriggerFlow triggers a flow with an input message (manual trigger)
func (e *Engine) TriggerFlow(ctx context.Context, flowID string, input *models.Message) error {
	// For now, we'll implement this as a start operation
	// In the future, this could send a message to a running flow
	return e.StartFlow(ctx, flowID)
}

// GetFlowStatus returns the status of a flow
func (e *Engine) GetFlowStatus(flowID string) (map[string]interface{}, error) {
	running, err := e.executor.GetFlowStatus(flowID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"running": running,
		"flow_id": flowID,
	}, nil
}

// GetRegistry returns the block registry
func (e *Engine) GetRegistry() *blocks.Registry {
	return e.registry
}

// Shutdown gracefully shuts down the engine
func (e *Engine) Shutdown(ctx context.Context) {
	// For now, just log the shutdown
	e.logger.Info("Engine shutting down", map[string]interface{}{})
}

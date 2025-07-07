package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"block-flow/internal/models"
)

// FileStorage implements Storage interface using the file system
type FileStorage struct {
	dataDir string
	mu      sync.RWMutex
}

// NewFileStorage creates a new file-based storage
func NewFileStorage(dataDir string) *FileStorage {
	return &FileStorage{
		dataDir: dataDir,
	}
}

// SaveFlow saves a flow to a JSON file
func (fs *FileStorage) SaveFlow(ctx context.Context, flow *models.Flow) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	// Ensure flows directory exists
	flowsDir := filepath.Join(fs.dataDir, "flows")
	if err := os.MkdirAll(flowsDir, 0o755); err != nil {
		return fmt.Errorf("failed to create flows directory: %w", err)
	}

	// Save flow to file
	filename := filepath.Join(flowsDir, flow.ID+".json")
	data, err := flow.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal flow: %w", err)
	}

	if err := os.WriteFile(filename, data, 0o644); err != nil {
		return fmt.Errorf("failed to write flow file: %w", err)
	}

	return nil
}

// LoadFlow loads a flow from a JSON file
func (fs *FileStorage) LoadFlow(ctx context.Context, flowID string) (*models.Flow, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	filename := filepath.Join(fs.dataDir, "flows", flowID+".json")
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, NewStorageError("flow not found", flowID, err)
		}
		return nil, fmt.Errorf("failed to read flow file: %w", err)
	}

	flow, err := models.FromJSON(data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal flow: %w", err)
	}

	return flow, nil
}

// LoadAllFlows loads all flows from the flows directory
func (fs *FileStorage) LoadAllFlows(ctx context.Context) ([]*models.Flow, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	flowsDir := filepath.Join(fs.dataDir, "flows")

	// Check if flows directory exists
	if _, err := os.Stat(flowsDir); os.IsNotExist(err) {
		return []*models.Flow{}, nil // Return empty slice if directory doesn't exist
	}

	entries, err := os.ReadDir(flowsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read flows directory: %w", err)
	}

	flows := make([]*models.Flow, 0)
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		flowID := entry.Name()[:len(entry.Name())-5] // Remove .json extension
		flow, err := fs.LoadFlow(ctx, flowID)
		if err != nil {
			// Log error but continue loading other flows
			continue
		}
		flows = append(flows, flow)
	}

	return flows, nil
}

// DeleteFlow deletes a flow file
func (fs *FileStorage) DeleteFlow(ctx context.Context, flowID string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	filename := filepath.Join(fs.dataDir, "flows", flowID+".json")
	if err := os.Remove(filename); err != nil {
		if os.IsNotExist(err) {
			return NewStorageError("flow not found", flowID, err)
		}
		return fmt.Errorf("failed to delete flow file: %w", err)
	}

	return nil
}

// FlowExists checks if a flow file exists
func (fs *FileStorage) FlowExists(ctx context.Context, flowID string) bool {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	filename := filepath.Join(fs.dataDir, "flows", flowID+".json")
	_, err := os.Stat(filename)
	return err == nil
}

// SaveFlowExecution saves a flow execution to a JSON file
func (fs *FileStorage) SaveFlowExecution(ctx context.Context, execution *models.FlowExecution) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	// Ensure executions directory exists
	execDir := filepath.Join(fs.dataDir, "executions")
	if err := os.MkdirAll(execDir, 0o755); err != nil {
		return fmt.Errorf("failed to create executions directory: %w", err)
	}

	// Save execution to file
	filename := filepath.Join(execDir, execution.ID+".json")
	data, err := json.MarshalIndent(execution, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal execution: %w", err)
	}

	if err := os.WriteFile(filename, data, 0o644); err != nil {
		return fmt.Errorf("failed to write execution file: %w", err)
	}

	return nil
}

// LoadFlowExecution loads a flow execution from a JSON file
func (fs *FileStorage) LoadFlowExecution(ctx context.Context, executionID string) (*models.FlowExecution, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	filename := filepath.Join(fs.dataDir, "executions", executionID+".json")
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, NewStorageError("execution not found", executionID, err)
		}
		return nil, fmt.Errorf("failed to read execution file: %w", err)
	}

	var execution models.FlowExecution
	if err := json.Unmarshal(data, &execution); err != nil {
		return nil, fmt.Errorf("failed to unmarshal execution: %w", err)
	}

	return &execution, nil
}

// LoadFlowExecutions loads all executions for a specific flow
func (fs *FileStorage) LoadFlowExecutions(ctx context.Context, flowID string) ([]*models.FlowExecution, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	execDir := filepath.Join(fs.dataDir, "executions")

	// Check if executions directory exists
	if _, err := os.Stat(execDir); os.IsNotExist(err) {
		return []*models.FlowExecution{}, nil
	}

	entries, err := os.ReadDir(execDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read executions directory: %w", err)
	}

	executions := make([]*models.FlowExecution, 0)
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		executionID := entry.Name()[:len(entry.Name())-5] // Remove .json extension
		execution, err := fs.LoadFlowExecution(ctx, executionID)
		if err != nil {
			continue // Skip invalid executions
		}

		if execution.FlowID == flowID {
			executions = append(executions, execution)
		}
	}

	return executions, nil
}

// DeleteFlowExecution deletes an execution file
func (fs *FileStorage) DeleteFlowExecution(ctx context.Context, executionID string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	filename := filepath.Join(fs.dataDir, "executions", executionID+".json")
	if err := os.Remove(filename); err != nil {
		if os.IsNotExist(err) {
			return NewStorageError("execution not found", executionID, err)
		}
		return fmt.Errorf("failed to delete execution file: %w", err)
	}

	return nil
}

// SaveConfig saves configuration data
func (fs *FileStorage) SaveConfig(ctx context.Context, key string, value interface{}) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	// Ensure config directory exists
	configDir := filepath.Join(fs.dataDir, "config")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Save config to file
	filename := filepath.Join(configDir, key+".json")
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(filename, data, 0o644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// LoadConfig loads configuration data
func (fs *FileStorage) LoadConfig(ctx context.Context, key string, target interface{}) error {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	filename := filepath.Join(fs.dataDir, "config", key+".json")
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return NewStorageError("config not found", key, err)
		}
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

// DeleteConfig deletes configuration data
func (fs *FileStorage) DeleteConfig(ctx context.Context, key string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	filename := filepath.Join(fs.dataDir, "config", key+".json")
	if err := os.Remove(filename); err != nil {
		if os.IsNotExist(err) {
			return NewStorageError("config not found", key, err)
		}
		return fmt.Errorf("failed to delete config file: %w", err)
	}

	return nil
}

// Health checks if the storage is healthy
func (fs *FileStorage) Health(ctx context.Context) error {
	// Check if data directory is accessible
	if err := os.MkdirAll(fs.dataDir, 0o755); err != nil {
		return fmt.Errorf("data directory not accessible: %w", err)
	}
	return nil
}

// Close closes the storage (no-op for file storage)
func (fs *FileStorage) Close() error {
	return nil
}

// StorageError represents a storage-related error
type StorageError struct {
	Message string
	Key     string
	Cause   error
}

func (e *StorageError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("storage error [%s]: %s: %v", e.Key, e.Message, e.Cause)
	}
	return fmt.Sprintf("storage error [%s]: %s", e.Key, e.Message)
}

func (e *StorageError) Unwrap() error {
	return e.Cause
}

// NewStorageError creates a new storage error
func NewStorageError(message, key string, cause error) *StorageError {
	return &StorageError{
		Message: message,
		Key:     key,
		Cause:   cause,
	}
}

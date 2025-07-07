package models

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// generateID generates a unique identifier
func generateID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based ID if crypto/rand fails
		return fmt.Sprintf("id_%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

// ValidationError represents a validation error
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// NewValidationError creates a new validation error
func NewValidationError(message string) *ValidationError {
	return &ValidationError{Message: message}
}

// ExecutionError represents an execution error
type ExecutionError struct {
	NodeID  string
	Message string
	Cause   error
}

func (e *ExecutionError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("execution error in node %s: %s: %v", e.NodeID, e.Message, e.Cause)
	}
	return fmt.Sprintf("execution error in node %s: %s", e.NodeID, e.Message)
}

func (e *ExecutionError) Unwrap() error {
	return e.Cause
}

// NewExecutionError creates a new execution error
func NewExecutionError(nodeID, message string, cause error) *ExecutionError {
	return &ExecutionError{
		NodeID:  nodeID,
		Message: message,
		Cause:   cause,
	}
}

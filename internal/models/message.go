package models

import (
	"context"
	"time"
)

// Message represents a message passed between blocks in a flow
type Message struct {
	ID        string                 `json:"id"`
	Payload   interface{}            `json:"payload"` // Main message payload
	Topic     string                 `json:"topic,omitempty"`
	Headers   map[string]string      `json:"headers,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Source    string                 `json:"source"`            // Source node ID
	Target    string                 `json:"target"`            // Target node ID
	Context   map[string]interface{} `json:"context,omitempty"` // Execution context
}

// NewMessage creates a new message
func NewMessage(payload interface{}) *Message {
	return &Message{
		ID:        generateID(),
		Payload:   payload,
		Headers:   make(map[string]string),
		Timestamp: time.Now(),
		Context:   make(map[string]interface{}),
	}
}

// Clone creates a deep copy of the message
func (m *Message) Clone() *Message {
	clone := &Message{
		ID:        generateID(), // New ID for cloned message
		Payload:   m.Payload,    // Note: shallow copy of payload
		Topic:     m.Topic,
		Timestamp: time.Now(),
		Source:    m.Source,
		Target:    m.Target,
	}

	// Deep copy headers
	if m.Headers != nil {
		clone.Headers = make(map[string]string, len(m.Headers))
		for k, v := range m.Headers {
			clone.Headers[k] = v
		}
	}

	// Deep copy context
	if m.Context != nil {
		clone.Context = make(map[string]interface{}, len(m.Context))
		for k, v := range m.Context {
			clone.Context[k] = v
		}
	}

	return clone
}

// SetHeader sets a header value
func (m *Message) SetHeader(key, value string) {
	if m.Headers == nil {
		m.Headers = make(map[string]string)
	}
	m.Headers[key] = value
}

// GetHeader gets a header value
func (m *Message) GetHeader(key string) (string, bool) {
	if m.Headers == nil {
		return "", false
	}
	value, exists := m.Headers[key]
	return value, exists
}

// SetContext sets a context value
func (m *Message) SetContext(key string, value interface{}) {
	if m.Context == nil {
		m.Context = make(map[string]interface{})
	}
	m.Context[key] = value
}

// GetContext gets a context value
func (m *Message) GetContext(key string) (interface{}, bool) {
	if m.Context == nil {
		return nil, false
	}
	value, exists := m.Context[key]
	return value, exists
}

// BlockExecutionContext provides context for block execution
type BlockExecutionContext struct {
	Context   context.Context
	NodeID    string
	FlowID    string
	Message   *Message
	Logger    BlockLogger
	State     map[string]interface{} // Block-specific state storage
	Debug     bool
	Timestamp time.Time
}

// NewBlockExecutionContext creates a new block execution context
func NewBlockExecutionContext(ctx context.Context, nodeID, flowID string, msg *Message, logger BlockLogger) *BlockExecutionContext {
	return &BlockExecutionContext{
		Context:   ctx,
		NodeID:    nodeID,
		FlowID:    flowID,
		Message:   msg,
		Logger:    logger,
		State:     make(map[string]interface{}),
		Debug:     true, // Default to debug mode
		Timestamp: time.Now(),
	}
}

// SetState sets a state value for the block
func (ctx *BlockExecutionContext) SetState(key string, value interface{}) {
	ctx.State[key] = value
}

// GetState gets a state value for the block
func (ctx *BlockExecutionContext) GetState(key string) (interface{}, bool) {
	value, exists := ctx.State[key]
	return value, exists
}

// BlockLogger provides logging interface for blocks
type BlockLogger interface {
	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Error(msg string, err error, fields map[string]interface{})
}

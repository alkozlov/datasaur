package models

import (
	"encoding/json"
	"time"
)

// Flow represents a complete data flow configuration
type Flow struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Nodes       []Node            `json:"nodes"`
	Connections []Connection      `json:"connections"`
	Properties  map[string]string `json:"properties,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Version     string            `json:"version"`
	Active      bool              `json:"active"`
}

// Node represents a single block/node in the flow
type Node struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`       // Block type (e.g., "inject", "debug", "function")
	Name       string                 `json:"name"`       // User-defined name
	X          float64                `json:"x"`          // X position in UI
	Y          float64                `json:"y"`          // Y position in UI
	Properties map[string]interface{} `json:"properties"` // Block-specific configuration
	Inputs     int                    `json:"inputs"`     // Number of input ports
	Outputs    int                    `json:"outputs"`    // Number of output ports
	Wires      [][]string             `json:"wires"`      // Output connections [output_port][connected_node_ids]
}

// Connection represents a wire between two nodes
type Connection struct {
	ID         string `json:"id"`
	Source     string `json:"source"`      // Source node ID
	SourcePort int    `json:"source_port"` // Source output port (0-based)
	Target     string `json:"target"`      // Target node ID
	TargetPort int    `json:"target_port"` // Target input port (0-based)
	Label      string `json:"label,omitempty"`
}

// FlowExecution represents the runtime state of a flow execution
type FlowExecution struct {
	ID        string                `json:"id"`
	FlowID    string                `json:"flow_id"`
	Status    ExecutionStatus       `json:"status"`
	StartedAt time.Time             `json:"started_at"`
	EndedAt   *time.Time            `json:"ended_at,omitempty"`
	Error     string                `json:"error,omitempty"`
	Nodes     map[string]*NodeState `json:"nodes"`
	Messages  []ExecutionMessage    `json:"messages,omitempty"`
}

// NodeState represents the runtime state of a node during execution
type NodeState struct {
	NodeID      string            `json:"node_id"`
	Status      NodeStatus        `json:"status"`
	ExecutedAt  *time.Time        `json:"executed_at,omitempty"`
	Duration    time.Duration     `json:"duration"`
	InputCount  int               `json:"input_count"`
	OutputCount int               `json:"output_count"`
	Error       string            `json:"error,omitempty"`
	LastMessage *Message          `json:"last_message,omitempty"`
	Properties  map[string]string `json:"properties,omitempty"`
}

// ExecutionStatus represents the status of flow execution
type ExecutionStatus string

const (
	ExecutionStatusPending   ExecutionStatus = "pending"
	ExecutionStatusRunning   ExecutionStatus = "running"
	ExecutionStatusCompleted ExecutionStatus = "completed"
	ExecutionStatusFailed    ExecutionStatus = "failed"
	ExecutionStatusStopped   ExecutionStatus = "stopped"
)

// NodeStatus represents the status of node execution
type NodeStatus string

const (
	NodeStatusIdle    NodeStatus = "idle"
	NodeStatusRunning NodeStatus = "running"
	NodeStatusSuccess NodeStatus = "success"
	NodeStatusError   NodeStatus = "error"
	NodeStatusSkipped NodeStatus = "skipped"
)

// ExecutionMessage represents a message during flow execution for debugging
type ExecutionMessage struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	NodeID    string    `json:"node_id"`
	Type      string    `json:"type"` // "input", "output", "error", "debug"
	Message   *Message  `json:"message,omitempty"`
	Error     string    `json:"error,omitempty"`
	Debug     string    `json:"debug,omitempty"`
}

// NewFlow creates a new flow with default values
func NewFlow(name string) *Flow {
	now := time.Now()
	return &Flow{
		ID:          generateID(),
		Name:        name,
		Nodes:       make([]Node, 0),
		Connections: make([]Connection, 0),
		Properties:  make(map[string]string),
		CreatedAt:   now,
		UpdatedAt:   now,
		Version:     "1.0.0",
		Active:      false,
	}
}

// NewNode creates a new node with default values
func NewNode(nodeType, name string) *Node {
	return &Node{
		ID:         generateID(),
		Type:       nodeType,
		Name:       name,
		Properties: make(map[string]interface{}),
		Inputs:     1,
		Outputs:    1,
		Wires:      make([][]string, 1),
	}
}

// NewFlowExecution creates a new flow execution
func NewFlowExecution(flowID string) *FlowExecution {
	return &FlowExecution{
		ID:        generateID(),
		FlowID:    flowID,
		Status:    ExecutionStatusPending,
		StartedAt: time.Now(),
		Nodes:     make(map[string]*NodeState),
		Messages:  make([]ExecutionMessage, 0),
	}
}

// AddNode adds a node to the flow
func (f *Flow) AddNode(node Node) {
	f.Nodes = append(f.Nodes, node)
	f.UpdatedAt = time.Now()
}

// AddConnection adds a connection to the flow
func (f *Flow) AddConnection(conn Connection) {
	f.Connections = append(f.Connections, conn)
	f.UpdatedAt = time.Now()
}

// GetNode returns a node by ID
func (f *Flow) GetNode(nodeID string) (*Node, bool) {
	for i := range f.Nodes {
		if f.Nodes[i].ID == nodeID {
			return &f.Nodes[i], true
		}
	}
	return nil, false
}

// RemoveNode removes a node and its connections from the flow
func (f *Flow) RemoveNode(nodeID string) bool {
	// Remove the node
	for i, node := range f.Nodes {
		if node.ID == nodeID {
			f.Nodes = append(f.Nodes[:i], f.Nodes[i+1:]...)
			break
		}
	}

	// Remove connections involving this node
	f.Connections = f.filterConnections(func(conn Connection) bool {
		return conn.Source != nodeID && conn.Target != nodeID
	})

	f.UpdatedAt = time.Now()
	return true
}

// filterConnections is a helper function to filter connections
func (f *Flow) filterConnections(predicate func(Connection) bool) []Connection {
	filtered := make([]Connection, 0)
	for _, conn := range f.Connections {
		if predicate(conn) {
			filtered = append(filtered, conn)
		}
	}
	return filtered
}

// Validate checks if the flow configuration is valid
func (f *Flow) Validate() error {
	// Check for duplicate node IDs
	nodeIDs := make(map[string]bool)
	for _, node := range f.Nodes {
		if nodeIDs[node.ID] {
			return NewValidationError("duplicate node ID: " + node.ID)
		}
		nodeIDs[node.ID] = true
	}

	// Check connection validity
	for _, conn := range f.Connections {
		// Check if source and target nodes exist
		if !nodeIDs[conn.Source] {
			return NewValidationError("connection references non-existent source node: " + conn.Source)
		}
		if !nodeIDs[conn.Target] {
			return NewValidationError("connection references non-existent target node: " + conn.Target)
		}
	}

	return nil
}

// ToJSON converts the flow to JSON
func (f *Flow) ToJSON() ([]byte, error) {
	return json.MarshalIndent(f, "", "  ")
}

// FromJSON creates a flow from JSON data
func FromJSON(data []byte) (*Flow, error) {
	var flow Flow
	if err := json.Unmarshal(data, &flow); err != nil {
		return nil, err
	}
	return &flow, nil
}

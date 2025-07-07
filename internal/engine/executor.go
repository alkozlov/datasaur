package engine

import (
	"fmt"
	"sync"
	"time"

	"block-flow/internal/blocks"
	"block-flow/internal/models"
)

// LoggerAdapter adapts our Logger interface to models.BlockLogger
type LoggerAdapter struct {
	logger Logger
}

func (la *LoggerAdapter) Debug(msg string, fields map[string]interface{}) {
	la.logger.Debug(msg, fields)
}

func (la *LoggerAdapter) Info(msg string, fields map[string]interface{}) {
	la.logger.Info(msg, fields)
}

func (la *LoggerAdapter) Warn(msg string, fields map[string]interface{}) {
	la.logger.Warn(msg, fields)
}

func (la *LoggerAdapter) Error(msg string, err error, fields map[string]interface{}) {
	// Convert the error to a field and call our logger
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["error"] = err
	la.logger.Error(msg, fields)
}

// RuntimeNode represents a node during execution with its channels
type RuntimeNode struct {
	ID         string
	Type       string
	Name       string
	Block      blocks.Block
	Properties map[string]interface{}

	// Message channels
	InputChan  chan *models.Message
	OutputChan chan *models.Message

	// Connection management
	OutputConnections []string // Target node IDs

	// Execution control
	StopChan  chan struct{}
	WaitGroup *sync.WaitGroup
}

// RuntimeFlow represents a flow during execution
type RuntimeFlow struct {
	ID          string
	Name        string
	Nodes       map[string]*RuntimeNode
	Connections []models.Connection

	// Flow control
	StopChan  chan struct{}
	WaitGroup sync.WaitGroup
	Running   bool
	mutex     sync.RWMutex
}

// FlowExecutor manages the execution of flows
type FlowExecutor struct {
	registry *blocks.Registry
	logger   Logger
	flows    map[string]*RuntimeFlow
	mutex    sync.RWMutex
}

// NewFlowExecutor creates a new flow executor
func NewFlowExecutor(registry *blocks.Registry, logger Logger) *FlowExecutor {
	return &FlowExecutor{
		registry: registry,
		logger:   logger,
		flows:    make(map[string]*RuntimeFlow),
	}
}

// ValidateFlow validates a flow before execution
func (fe *FlowExecutor) ValidateFlow(flow *models.Flow) error {
	if flow == nil {
		return fmt.Errorf("flow is nil")
	}

	if len(flow.Nodes) == 0 {
		return fmt.Errorf("flow must contain at least one block")
	}

	// Validate all nodes have valid block types
	for _, node := range flow.Nodes {
		_, err := fe.registry.GetBlockInfoByType(node.Type)
		if err != nil {
			return fmt.Errorf("unknown block type '%s' in node '%s'", node.Type, node.ID)
		}
	}

	// Validate connections
	nodeMap := make(map[string]*models.Node)
	for _, node := range flow.Nodes {
		nodeMap[node.ID] = &node
	}

	for _, conn := range flow.Connections {
		sourceNode, exists := nodeMap[conn.Source]
		if !exists {
			return fmt.Errorf("connection references non-existent source node '%s'", conn.Source)
		}

		targetNode, exists := nodeMap[conn.Target]
		if !exists {
			return fmt.Errorf("connection references non-existent target node '%s'", conn.Target)
		}

		// Validate port ranges
		if conn.SourcePort >= sourceNode.Outputs {
			return fmt.Errorf("connection references invalid source port %d (node '%s' has %d outputs)",
				conn.SourcePort, conn.Source, sourceNode.Outputs)
		}

		if conn.TargetPort >= targetNode.Inputs {
			return fmt.Errorf("connection references invalid target port %d (node '%s' has %d inputs)",
				conn.TargetPort, conn.Target, targetNode.Inputs)
		}
	}

	return nil
}

// PrepareFlow prepares a flow for execution by creating runtime structures
func (fe *FlowExecutor) PrepareFlow(flow *models.Flow) (*RuntimeFlow, error) {
	err := fe.ValidateFlow(flow)
	if err != nil {
		return nil, fmt.Errorf("flow validation failed: %w", err)
	}

	runtimeFlow := &RuntimeFlow{
		ID:          flow.ID,
		Name:        flow.Name,
		Nodes:       make(map[string]*RuntimeNode),
		Connections: flow.Connections,
		StopChan:    make(chan struct{}),
		Running:     false,
	}

	// Create runtime nodes
	for _, node := range flow.Nodes {
		block, err := fe.registry.CreateBlock(node.Type)
		if err != nil {
			return nil, fmt.Errorf("failed to create block for node '%s': %w", node.ID, err)
		}

		blockInfo, err := fe.registry.GetBlockInfoByType(node.Type)
		if err != nil {
			return nil, fmt.Errorf("failed to get block info for node '%s': %w", node.ID, err)
		}

		runtimeNode := &RuntimeNode{
			ID:         node.ID,
			Type:       node.Type,
			Name:       node.Name,
			Block:      block,
			Properties: node.Properties,
			InputChan:  make(chan *models.Message, 100), // Buffered channel
			OutputChan: make(chan *models.Message, 100), // Buffered channel
			StopChan:   make(chan struct{}),
			WaitGroup:  &runtimeFlow.WaitGroup,
		}

		// Determine output connections for this node
		for _, conn := range flow.Connections {
			if conn.Source == node.ID {
				runtimeNode.OutputConnections = append(runtimeNode.OutputConnections, conn.Target)
			}
		}

		runtimeFlow.Nodes[node.ID] = runtimeNode

		fe.logger.Debug("Created runtime node", map[string]interface{}{
			"node_id":     node.ID,
			"node_type":   node.Type,
			"block_group": blockInfo.BlockGroup,
			"connections": len(runtimeNode.OutputConnections),
		})
	}

	return runtimeFlow, nil
}

// StartFlow starts the execution of a flow
func (fe *FlowExecutor) StartFlow(flowID string) error {
	fe.mutex.Lock()
	defer fe.mutex.Unlock()

	runtimeFlow, exists := fe.flows[flowID]
	if !exists {
		return fmt.Errorf("flow '%s' is not prepared for execution", flowID)
	}

	if runtimeFlow.Running {
		return fmt.Errorf("flow '%s' is already running", flowID)
	}

	runtimeFlow.mutex.Lock()
	runtimeFlow.Running = true
	runtimeFlow.mutex.Unlock()

	// Start all nodes
	for _, node := range runtimeFlow.Nodes {
		runtimeFlow.WaitGroup.Add(1)
		go fe.runNode(node, runtimeFlow)
	}

	fe.logger.Info("Flow started", map[string]interface{}{
		"flow_id":   flowID,
		"flow_name": runtimeFlow.Name,
		"nodes":     len(runtimeFlow.Nodes),
	})

	return nil
}

// StopFlow stops the execution of a flow
func (fe *FlowExecutor) StopFlow(flowID string) error {
	fe.mutex.Lock()
	defer fe.mutex.Unlock()

	runtimeFlow, exists := fe.flows[flowID]
	if !exists {
		return fmt.Errorf("flow '%s' is not running", flowID)
	}

	if !runtimeFlow.Running {
		return fmt.Errorf("flow '%s' is not running", flowID)
	}

	// Signal all nodes to stop
	close(runtimeFlow.StopChan)

	// Wait for all nodes to finish
	runtimeFlow.WaitGroup.Wait()

	runtimeFlow.mutex.Lock()
	runtimeFlow.Running = false
	runtimeFlow.mutex.Unlock()

	fe.logger.Info("Flow stopped", map[string]interface{}{
		"flow_id": flowID,
	})

	return nil
}

// runNode runs a single node in the flow
func (fe *FlowExecutor) runNode(node *RuntimeNode, flow *RuntimeFlow) {
	defer node.WaitGroup.Done()

	blockInfo := fe.registry.GetBlockInfo()[0] // Get block info for logging
	for _, info := range fe.registry.GetBlockInfo() {
		if info.Type == node.Type {
			blockInfo = info
			break
		}
	}

	fe.logger.Debug("Starting node", map[string]interface{}{
		"node_id":     node.ID,
		"node_type":   node.Type,
		"block_group": blockInfo.BlockGroup,
	})

	// Handle different block groups
	switch blockInfo.BlockGroup {
	case blocks.InputGroup:
		fe.runInputNode(node, flow)
	case blocks.PropagationGroup:
		fe.runPropagationNode(node, flow)
	case blocks.ActionGroup:
		fe.runActionNode(node, flow)
	}

	fe.logger.Debug("Node finished", map[string]interface{}{
		"node_id": node.ID,
	})
}

// runInputNode runs an input group node (generates messages)
func (fe *FlowExecutor) runInputNode(node *RuntimeNode, flow *RuntimeFlow) {
	// For inject blocks, we can implement interval-based message generation
	// For now, we'll implement a simple trigger mechanism

	ticker := time.NewTicker(1 * time.Second) // Default 1 second interval
	defer ticker.Stop()

	for {
		select {
		case <-flow.StopChan:
			return
		case <-node.StopChan:
			return
		case <-ticker.C: // Generate message from input block
			ctx := &models.BlockExecutionContext{
				NodeID:  node.ID,
				Logger:  &LoggerAdapter{logger: fe.logger},
				Message: nil, // Input blocks don't have input messages
			}

			messages, err := node.Block.Execute(ctx, node.Properties)
			if err != nil {
				fe.logger.Error("Error executing input node", map[string]interface{}{
					"node_id": node.ID,
					"error":   err.Error(),
				})
				continue
			}

			// Send messages to output connections
			for _, msg := range messages {
				fe.distributeMessage(node, msg, flow)
			}
		}
	}
}

// runPropagationNode runs a propagation group node (processes messages)
func (fe *FlowExecutor) runPropagationNode(node *RuntimeNode, flow *RuntimeFlow) {
	for {
		select {
		case <-flow.StopChan:
			return
		case <-node.StopChan:
			return
		case msg := <-node.InputChan: // Process message
			ctx := &models.BlockExecutionContext{
				NodeID:  node.ID,
				Logger:  &LoggerAdapter{logger: fe.logger},
				Message: msg,
			}

			messages, err := node.Block.Execute(ctx, node.Properties)
			if err != nil {
				fe.logger.Error("Error executing propagation node", map[string]interface{}{
					"node_id": node.ID,
					"error":   err.Error(),
				})
				continue
			}

			// Send messages to output connections
			for _, outMsg := range messages {
				fe.distributeMessage(node, outMsg, flow)
			}
		}
	}
}

// runActionNode runs an action group node (consumes messages)
func (fe *FlowExecutor) runActionNode(node *RuntimeNode, flow *RuntimeFlow) {
	for {
		select {
		case <-flow.StopChan:
			return
		case <-node.StopChan:
			return
		case msg := <-node.InputChan: // Process message (no output)
			ctx := &models.BlockExecutionContext{
				NodeID:  node.ID,
				Logger:  &LoggerAdapter{logger: fe.logger},
				Message: msg,
			}

			_, err := node.Block.Execute(ctx, node.Properties)
			if err != nil {
				fe.logger.Error("Error executing action node", map[string]interface{}{
					"node_id": node.ID,
					"error":   err.Error(),
				})
			}
			// Action blocks don't generate output messages
		}
	}
}

// distributeMessage sends a message to all connected target nodes
func (fe *FlowExecutor) distributeMessage(sourceNode *RuntimeNode, msg *models.Message, flow *RuntimeFlow) {
	for _, targetNodeID := range sourceNode.OutputConnections {
		targetNode, exists := flow.Nodes[targetNodeID]
		if !exists {
			fe.logger.Error("Target node not found", map[string]interface{}{
				"source_node": sourceNode.ID,
				"target_node": targetNodeID,
			})
			continue
		}

		// Clone message for each target to avoid shared state issues
		clonedMsg := msg.Clone()

		// Non-blocking send (drop message if channel is full)
		select {
		case targetNode.InputChan <- clonedMsg:
			fe.logger.Debug("Message sent", map[string]interface{}{
				"from":    sourceNode.ID,
				"to":      targetNodeID,
				"payload": clonedMsg.Payload,
			})
		default:
			fe.logger.Warn("Target node input channel full, dropping message", map[string]interface{}{
				"source_node": sourceNode.ID,
				"target_node": targetNodeID,
			})
		}
	}
}

// PrepareAndStartFlow is a convenience method to prepare and start a flow
func (fe *FlowExecutor) PrepareAndStartFlow(flow *models.Flow) error {
	runtimeFlow, err := fe.PrepareFlow(flow)
	if err != nil {
		return fmt.Errorf("failed to prepare flow: %w", err)
	}

	fe.mutex.Lock()
	fe.flows[flow.ID] = runtimeFlow
	fe.mutex.Unlock()

	return fe.StartFlow(flow.ID)
}

// GetFlowStatus returns the status of a flow
func (fe *FlowExecutor) GetFlowStatus(flowID string) (bool, error) {
	fe.mutex.RLock()
	defer fe.mutex.RUnlock()

	runtimeFlow, exists := fe.flows[flowID]
	if !exists {
		return false, fmt.Errorf("flow '%s' not found", flowID)
	}

	runtimeFlow.mutex.RLock()
	running := runtimeFlow.Running
	runtimeFlow.mutex.RUnlock()

	return running, nil
}

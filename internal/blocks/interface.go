package blocks

import (
	"fmt"

	"block-flow/internal/models"
)

// BlockGroup represents the functional group of a block
type BlockGroup string

const (
	// InputGroup - blocks that generate data with no input streams
	InputGroup BlockGroup = "input"
	// PropagationGroup - blocks that process data with input and output streams
	PropagationGroup BlockGroup = "propagation"
	// ActionGroup - blocks that consume data with only input streams
	ActionGroup BlockGroup = "action"
)

// Block represents a processing block in the flow
type Block interface {
	// GetType returns the block type identifier
	GetType() string

	// GetName returns the human-readable block name
	GetName() string

	// GetDescription returns the block description
	GetDescription() string

	// GetCategory returns the block category for UI grouping
	GetCategory() string

	// GetBlockGroup returns the block group (Input, Propagation, Action)
	GetBlockGroup() BlockGroup

	// GetInputs returns the number of input ports
	GetInputs() int

	// GetOutputs returns the number of output ports
	GetOutputs() int

	// GetProperties returns the block property definitions
	GetProperties() []PropertyDefinition

	// Validate validates the block configuration
	Validate(properties map[string]interface{}) error

	// Execute processes the input message and returns output messages
	Execute(ctx *models.BlockExecutionContext, properties map[string]interface{}) ([]*models.Message, error)
}

// PropertyDefinition defines a configurable property of a block
type PropertyDefinition struct {
	Name         string      `json:"name"`
	Type         string      `json:"type"` // "string", "number", "boolean", "select", "json"
	DisplayName  string      `json:"display_name"`
	Description  string      `json:"description"`
	Required     bool        `json:"required"`
	DefaultValue interface{} `json:"default_value,omitempty"`
	Options      []Option    `json:"options,omitempty"` // For select type
	Validation   Validation  `json:"validation,omitempty"`
}

// Option represents a select option
type Option struct {
	Label string      `json:"label"`
	Value interface{} `json:"value"`
}

// Validation defines validation rules for properties
type Validation struct {
	Min    *float64 `json:"min,omitempty"`
	Max    *float64 `json:"max,omitempty"`
	Regex  string   `json:"regex,omitempty"`
	Length *int     `json:"length,omitempty"`
}

// BlockFactory creates new instances of blocks
type BlockFactory interface {
	CreateBlock() Block
	GetBlockInfo() BlockInfo
}

// BlockInfo provides metadata about a block type
type BlockInfo struct {
	Type        string     `json:"type"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Category    string     `json:"category"`
	BlockGroup  BlockGroup `json:"block_group"`
	Inputs      int        `json:"inputs"`
	Outputs     int        `json:"outputs"`
	Version     string     `json:"version"`
	Author      string     `json:"author,omitempty"`
	Icon        string     `json:"icon,omitempty"`
	Color       string     `json:"color,omitempty"`
}

// Registry manages available blocks
type Registry struct {
	blocks map[string]BlockFactory
}

// NewRegistry creates a new block registry
func NewRegistry() *Registry {
	return &Registry{
		blocks: make(map[string]BlockFactory),
	}
}

// Register registers a block factory
func (r *Registry) Register(factory BlockFactory) {
	info := factory.GetBlockInfo()
	r.blocks[info.Type] = factory
}

// CreateBlock creates a new block instance by type
func (r *Registry) CreateBlock(blockType string) (Block, error) {
	factory, exists := r.blocks[blockType]
	if !exists {
		return nil, NewBlockError("unknown block type", blockType, nil)
	}
	return factory.CreateBlock(), nil
}

// GetBlockTypes returns all registered block types
func (r *Registry) GetBlockTypes() []string {
	types := make([]string, 0, len(r.blocks))
	for blockType := range r.blocks {
		types = append(types, blockType)
	}
	return types
}

// GetBlockInfo returns information about all registered blocks
func (r *Registry) GetBlockInfo() []BlockInfo {
	info := make([]BlockInfo, 0, len(r.blocks))
	for _, factory := range r.blocks {
		info = append(info, factory.GetBlockInfo())
	}
	return info
}

// GetBlockInfoByType returns information about a specific block type
func (r *Registry) GetBlockInfoByType(blockType string) (BlockInfo, error) {
	factory, exists := r.blocks[blockType]
	if !exists {
		return BlockInfo{}, NewBlockError("unknown block type", blockType, nil)
	}
	return factory.GetBlockInfo(), nil
}

// BlockError represents a block-related error
type BlockError struct {
	Message   string
	BlockType string
	Cause     error
}

func (e *BlockError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("block error [%s]: %s: %v", e.BlockType, e.Message, e.Cause)
	}
	return fmt.Sprintf("block error [%s]: %s", e.BlockType, e.Message)
}

func (e *BlockError) Unwrap() error {
	return e.Cause
}

// NewBlockError creates a new block error
func NewBlockError(message, blockType string, cause error) *BlockError {
	return &BlockError{
		Message:   message,
		BlockType: blockType,
		Cause:     cause,
	}
}

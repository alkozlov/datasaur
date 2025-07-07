package builtin

import (
	"fmt"
	"log"

	"block-flow/internal/blocks"
	"block-flow/internal/models"
)

// DebugBlock outputs debug information about messages
type DebugBlock struct{}

func (b *DebugBlock) GetType() string {
	return "debug"
}

func (b *DebugBlock) GetName() string {
	return "Debug"
}

func (b *DebugBlock) GetDescription() string {
	return "Output debug information to console or log"
}

func (b *DebugBlock) GetCategory() string {
	return "output"
}

func (b *DebugBlock) GetBlockGroup() blocks.BlockGroup {
	return blocks.ActionGroup
}

func (b *DebugBlock) GetInputs() int {
	return 1
}

func (b *DebugBlock) GetOutputs() int {
	return 0
}

func (b *DebugBlock) GetProperties() []blocks.PropertyDefinition {
	return []blocks.PropertyDefinition{
		{
			Name:         "name",
			Type:         "string",
			DisplayName:  "Name",
			Description:  "Block name for identification",
			Required:     false,
			DefaultValue: "Debug",
		},
		{
			Name:         "console",
			Type:         "boolean",
			DisplayName:  "Console Output",
			Description:  "Output to console/log",
			Required:     false,
			DefaultValue: true,
		},
		{
			Name:         "complete",
			Type:         "select",
			DisplayName:  "Output",
			Description:  "What to output",
			Required:     false,
			DefaultValue: "payload",
			Options: []blocks.Option{
				{Label: "Payload only", Value: "payload"},
				{Label: "Complete message", Value: "complete"},
			},
		},
		{
			Name:         "prefix",
			Type:         "string",
			DisplayName:  "Prefix",
			Description:  "Optional prefix for debug output",
			Required:     false,
			DefaultValue: "",
		},
	}
}

func (b *DebugBlock) Validate(properties map[string]interface{}) error {
	// No validation needed for debug block
	return nil
}

func (b *DebugBlock) Execute(ctx *models.BlockExecutionContext, properties map[string]interface{}) ([]*models.Message, error) {
	if ctx.Message == nil {
		return nil, fmt.Errorf("no input message to debug")
	}

	// Get properties
	console, _ := properties["console"].(bool)
	complete, _ := properties["complete"].(string)
	prefix, _ := properties["prefix"].(string)

	if complete == "" {
		complete = "payload"
	}

	// Prepare output
	var output interface{}
	if complete == "complete" {
		output = ctx.Message
	} else {
		output = ctx.Message.Payload
	}

	// Format debug message
	debugMsg := fmt.Sprintf("[%s] Debug", ctx.NodeID)
	if prefix != "" {
		debugMsg = fmt.Sprintf("[%s] %s", ctx.NodeID, prefix)
	}

	// Output to console if enabled
	if console {
		log.Printf("%s: %v", debugMsg, output)
	}

	// Log debug information
	ctx.Logger.Debug("Debug block output", map[string]interface{}{
		"node_id": ctx.NodeID,
		"prefix":  prefix,
		"output":  output,
		"topic":   ctx.Message.Topic,
	})

	// Debug blocks don't pass messages forward
	return []*models.Message{}, nil
}

// DebugBlockFactory creates debug block instances
type DebugBlockFactory struct{}

func (f *DebugBlockFactory) CreateBlock() blocks.Block {
	return &DebugBlock{}
}

func (f *DebugBlockFactory) GetBlockInfo() blocks.BlockInfo {
	block := &DebugBlock{}
	return blocks.BlockInfo{
		Type:        "debug",
		Name:        "Debug",
		Description: "Output debug information to console or log",
		Category:    "output",
		BlockGroup:  blocks.ActionGroup,
		Inputs:      block.GetInputs(),
		Outputs:     block.GetOutputs(),
		Version:     "1.0.0",
		Author:      "Block-Flow",
		Icon:        "bug",
		Color:       "#FF9800",
	}
}

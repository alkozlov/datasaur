package builtin

import (
	"fmt"
	"strconv"

	"block-flow/internal/blocks"
	"block-flow/internal/models"
)

// InjectBlock provides manual input trigger with configurable payload
type InjectBlock struct{}

func (b *InjectBlock) GetType() string {
	return "inject"
}

func (b *InjectBlock) GetName() string {
	return "Inject"
}

func (b *InjectBlock) GetDescription() string {
	return "Manually trigger the flow with a configurable payload"
}

func (b *InjectBlock) GetCategory() string {
	return "input"
}

func (b *InjectBlock) GetBlockGroup() blocks.BlockGroup {
	return blocks.InputGroup
}

func (b *InjectBlock) GetInputs() int {
	return 0
}

func (b *InjectBlock) GetOutputs() int {
	return 1
}

func (b *InjectBlock) GetProperties() []blocks.PropertyDefinition {
	return []blocks.PropertyDefinition{
		{
			Name:         "name",
			Type:         "string",
			DisplayName:  "Name",
			Description:  "Block name for identification",
			Required:     false,
			DefaultValue: "Inject",
		},
		{
			Name:         "payload",
			Type:         "string",
			DisplayName:  "Output Value",
			Description:  "The value to inject (will be converted to number if possible)",
			Required:     true,
			DefaultValue: "0",
		},
		{
			Name:         "interval",
			Type:         "number",
			DisplayName:  "Interval (ms)",
			Description:  "Message publishing interval in milliseconds (0 = manual trigger only)",
			Required:     false,
			DefaultValue: 1000,
			Validation: blocks.Validation{
				Min: &[]float64{0}[0],
			},
		},
		{
			Name:         "topic",
			Type:         "string",
			DisplayName:  "Topic",
			Description:  "Optional topic for the message",
			Required:     false,
			DefaultValue: "",
		},
		{
			Name:         "payloadType",
			Type:         "select",
			DisplayName:  "Payload Type",
			Description:  "The type of the payload",
			Required:     false,
			DefaultValue: "number",
			Options: []blocks.Option{
				{Label: "Number", Value: "number"},
				{Label: "String", Value: "string"},
				{Label: "Boolean", Value: "boolean"},
			},
		},
	}
}

func (b *InjectBlock) Validate(properties map[string]interface{}) error {
	if _, ok := properties["payload"]; !ok {
		return fmt.Errorf("payload property is required")
	}
	return nil
}

func (b *InjectBlock) Execute(ctx *models.BlockExecutionContext, properties map[string]interface{}) ([]*models.Message, error) {
	// Get properties
	payloadStr, _ := properties["payload"].(string)
	topic, _ := properties["topic"].(string)
	payloadType, _ := properties["payloadType"].(string)
	if payloadType == "" {
		payloadType = "number"
	}

	var payload interface{}
	var err error

	// Convert payload based on type
	switch payloadType {
	case "number":
		payload, err = strconv.ParseFloat(payloadStr, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse payload as number: %w", err)
		}
	case "boolean":
		payload, err = strconv.ParseBool(payloadStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse payload as boolean: %w", err)
		}
	default:
		payload = payloadStr
	}

	// Create output message
	outputMsg := models.NewMessage(payload)
	outputMsg.Topic = topic
	outputMsg.Source = ctx.NodeID

	ctx.Logger.Debug("Inject block executed", map[string]interface{}{
		"payload":      payload,
		"payload_type": payloadType,
		"topic":        topic,
	})

	return []*models.Message{outputMsg}, nil
}

// InjectBlockFactory creates inject block instances
type InjectBlockFactory struct{}

func (f *InjectBlockFactory) CreateBlock() blocks.Block {
	return &InjectBlock{}
}

func (f *InjectBlockFactory) GetBlockInfo() blocks.BlockInfo {
	block := &InjectBlock{}
	return blocks.BlockInfo{
		Type:        "inject",
		Name:        "Inject",
		Description: "Manually trigger the flow with configurable payload",
		Category:    "input",
		BlockGroup:  blocks.InputGroup,
		Inputs:      block.GetInputs(),
		Outputs:     block.GetOutputs(),
		Version:     "1.0.0",
		Author:      "Block-Flow",
		Icon:        "play-circle",
		Color:       "#4CAF50",
	}
}

package builtin

import (
	"fmt"

	"block-flow/internal/blocks"
	"block-flow/internal/models"
)

// Helper function to extract number from message or property
func extractNumber(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("value is not a number: %T", value)
	}
}

// AdditionBlock performs addition operation
type AdditionBlock struct{}

func (b *AdditionBlock) GetType() string {
	return "add"
}

func (b *AdditionBlock) GetName() string {
	return "Addition"
}

func (b *AdditionBlock) GetDescription() string {
	return "Add a number to the input payload"
}

func (b *AdditionBlock) GetCategory() string {
	return "math"
}

func (b *AdditionBlock) GetBlockGroup() blocks.BlockGroup {
	return blocks.PropagationGroup
}

func (b *AdditionBlock) GetInputs() int {
	return 1
}

func (b *AdditionBlock) GetOutputs() int {
	return 1
}

func (b *AdditionBlock) GetProperties() []blocks.PropertyDefinition {
	return []blocks.PropertyDefinition{
		{
			Name:         "name",
			Type:         "string",
			DisplayName:  "Name",
			Description:  "Block name for identification",
			Required:     false,
			DefaultValue: "Add",
		},
		{
			Name:         "value",
			Type:         "number",
			DisplayName:  "Second Addend",
			Description:  "The number to add to the input",
			Required:     true,
			DefaultValue: 0.0,
		},
	}
}

func (b *AdditionBlock) Validate(properties map[string]interface{}) error {
	if _, ok := properties["value"]; !ok {
		return fmt.Errorf("value property is required")
	}
	return nil
}

func (b *AdditionBlock) Execute(ctx *models.BlockExecutionContext, properties map[string]interface{}) ([]*models.Message, error) {
	if ctx.Message == nil {
		return nil, fmt.Errorf("no input message")
	}

	// Extract input number
	inputNum, err := extractNumber(ctx.Message.Payload)
	if err != nil {
		return nil, fmt.Errorf("input payload is not a number: %w", err)
	}

	// Extract value to add
	addValue, err := extractNumber(properties["value"])
	if err != nil {
		return nil, fmt.Errorf("add value is not a number: %w", err)
	}

	// Perform addition
	result := inputNum + addValue

	// Create output message
	outputMsg := ctx.Message.Clone()
	outputMsg.Payload = result
	outputMsg.Source = ctx.NodeID

	ctx.Logger.Debug("Addition performed", map[string]interface{}{
		"input":  inputNum,
		"add":    addValue,
		"result": result,
	})

	return []*models.Message{outputMsg}, nil
}

// SubtractionBlock performs subtraction operation
type SubtractionBlock struct{}

func (b *SubtractionBlock) GetType() string {
	return "subtract"
}

func (b *SubtractionBlock) GetName() string {
	return "Subtraction"
}

func (b *SubtractionBlock) GetDescription() string {
	return "Subtract a number from the input payload"
}

func (b *SubtractionBlock) GetCategory() string {
	return "math"
}

func (b *SubtractionBlock) GetBlockGroup() blocks.BlockGroup {
	return blocks.PropagationGroup
}

func (b *SubtractionBlock) GetInputs() int {
	return 1
}

func (b *SubtractionBlock) GetOutputs() int {
	return 1
}

func (b *SubtractionBlock) GetProperties() []blocks.PropertyDefinition {
	return []blocks.PropertyDefinition{
		{
			Name:         "name",
			Type:         "string",
			DisplayName:  "Name",
			Description:  "Block name for identification",
			Required:     false,
			DefaultValue: "Subtract",
		},
		{
			Name:         "value",
			Type:         "number",
			DisplayName:  "Subtracted Value",
			Description:  "The number to subtract from the input",
			Required:     true,
			DefaultValue: 0.0,
		},
	}
}

func (b *SubtractionBlock) Validate(properties map[string]interface{}) error {
	if _, ok := properties["value"]; !ok {
		return fmt.Errorf("value property is required")
	}
	return nil
}

func (b *SubtractionBlock) Execute(ctx *models.BlockExecutionContext, properties map[string]interface{}) ([]*models.Message, error) {
	if ctx.Message == nil {
		return nil, fmt.Errorf("no input message")
	}

	// Extract input number
	inputNum, err := extractNumber(ctx.Message.Payload)
	if err != nil {
		return nil, fmt.Errorf("input payload is not a number: %w", err)
	}

	// Extract value to subtract
	subValue, err := extractNumber(properties["value"])
	if err != nil {
		return nil, fmt.Errorf("subtract value is not a number: %w", err)
	}

	// Perform subtraction
	result := inputNum - subValue

	// Create output message
	outputMsg := ctx.Message.Clone()
	outputMsg.Payload = result
	outputMsg.Source = ctx.NodeID

	ctx.Logger.Debug("Subtraction performed", map[string]interface{}{
		"input":    inputNum,
		"subtract": subValue,
		"result":   result,
	})

	return []*models.Message{outputMsg}, nil
}

// MultiplicationBlock performs multiplication operation
type MultiplicationBlock struct{}

func (b *MultiplicationBlock) GetType() string {
	return "multiply"
}

func (b *MultiplicationBlock) GetName() string {
	return "Multiplication"
}

func (b *MultiplicationBlock) GetDescription() string {
	return "Multiply the input payload by a number"
}

func (b *MultiplicationBlock) GetCategory() string {
	return "math"
}

func (b *MultiplicationBlock) GetBlockGroup() blocks.BlockGroup {
	return blocks.PropagationGroup
}

func (b *MultiplicationBlock) GetInputs() int {
	return 1
}

func (b *MultiplicationBlock) GetOutputs() int {
	return 1
}

func (b *MultiplicationBlock) GetProperties() []blocks.PropertyDefinition {
	return []blocks.PropertyDefinition{
		{
			Name:         "name",
			Type:         "string",
			DisplayName:  "Name",
			Description:  "Block name for identification",
			Required:     false,
			DefaultValue: "Multiply",
		},
		{
			Name:         "value",
			Type:         "number",
			DisplayName:  "Second Multiplier",
			Description:  "The number to multiply the input by",
			Required:     true,
			DefaultValue: 1.0,
		},
	}
}

func (b *MultiplicationBlock) Validate(properties map[string]interface{}) error {
	if _, ok := properties["value"]; !ok {
		return fmt.Errorf("value property is required")
	}
	return nil
}

func (b *MultiplicationBlock) Execute(ctx *models.BlockExecutionContext, properties map[string]interface{}) ([]*models.Message, error) {
	if ctx.Message == nil {
		return nil, fmt.Errorf("no input message")
	}

	// Extract input number
	inputNum, err := extractNumber(ctx.Message.Payload)
	if err != nil {
		return nil, fmt.Errorf("input payload is not a number: %w", err)
	}

	// Extract multiplier
	mulValue, err := extractNumber(properties["value"])
	if err != nil {
		return nil, fmt.Errorf("multiplier is not a number: %w", err)
	}

	// Perform multiplication
	result := inputNum * mulValue

	// Create output message
	outputMsg := ctx.Message.Clone()
	outputMsg.Payload = result
	outputMsg.Source = ctx.NodeID

	ctx.Logger.Debug("Multiplication performed", map[string]interface{}{
		"input":      inputNum,
		"multiplier": mulValue,
		"result":     result,
	})

	return []*models.Message{outputMsg}, nil
}

// DivisionBlock performs division operation
type DivisionBlock struct{}

func (b *DivisionBlock) GetType() string {
	return "divide"
}

func (b *DivisionBlock) GetName() string {
	return "Division"
}

func (b *DivisionBlock) GetDescription() string {
	return "Divide the input payload by a number"
}

func (b *DivisionBlock) GetCategory() string {
	return "math"
}

func (b *DivisionBlock) GetBlockGroup() blocks.BlockGroup {
	return blocks.PropagationGroup
}

func (b *DivisionBlock) GetInputs() int {
	return 1
}

func (b *DivisionBlock) GetOutputs() int {
	return 1
}

func (b *DivisionBlock) GetProperties() []blocks.PropertyDefinition {
	return []blocks.PropertyDefinition{
		{
			Name:         "name",
			Type:         "string",
			DisplayName:  "Name",
			Description:  "Block name for identification",
			Required:     false,
			DefaultValue: "Divide",
		},
		{
			Name:         "value",
			Type:         "number",
			DisplayName:  "Divisor",
			Description:  "The number to divide the input by",
			Required:     true,
			DefaultValue: 1.0,
			Validation: blocks.Validation{
				Min: &[]float64{0.000001}[0], // Prevent division by zero
			},
		},
	}
}

func (b *DivisionBlock) Validate(properties map[string]interface{}) error {
	value, ok := properties["value"]
	if !ok {
		return fmt.Errorf("value property is required")
	}

	divisor, err := extractNumber(value)
	if err != nil {
		return fmt.Errorf("divisor must be a number: %w", err)
	}

	if divisor == 0 {
		return fmt.Errorf("division by zero is not allowed")
	}

	return nil
}

func (b *DivisionBlock) Execute(ctx *models.BlockExecutionContext, properties map[string]interface{}) ([]*models.Message, error) {
	if ctx.Message == nil {
		return nil, fmt.Errorf("no input message")
	}

	// Extract input number
	inputNum, err := extractNumber(ctx.Message.Payload)
	if err != nil {
		return nil, fmt.Errorf("input payload is not a number: %w", err)
	}

	// Extract divisor
	divValue, err := extractNumber(properties["value"])
	if err != nil {
		return nil, fmt.Errorf("divisor is not a number: %w", err)
	}

	// Check for division by zero
	if divValue == 0 {
		return nil, fmt.Errorf("division by zero")
	}

	// Perform division
	result := inputNum / divValue

	// Create output message
	outputMsg := ctx.Message.Clone()
	outputMsg.Payload = result
	outputMsg.Source = ctx.NodeID

	ctx.Logger.Debug("Division performed", map[string]interface{}{
		"input":   inputNum,
		"divisor": divValue,
		"result":  result,
	})

	return []*models.Message{outputMsg}, nil
}

// Block factories

type AdditionBlockFactory struct{}

func (f *AdditionBlockFactory) CreateBlock() blocks.Block { return &AdditionBlock{} }
func (f *AdditionBlockFactory) GetBlockInfo() blocks.BlockInfo {
	block := &AdditionBlock{}
	return blocks.BlockInfo{
		Type: "add", Name: "Addition", Description: "Add a number to the input",
		Category: "math", BlockGroup: blocks.PropagationGroup, Inputs: block.GetInputs(), Outputs: block.GetOutputs(),
		Version: "1.0.0", Author: "Block-Flow", Icon: "plus", Color: "#2196F3",
	}
}

type SubtractionBlockFactory struct{}

func (f *SubtractionBlockFactory) CreateBlock() blocks.Block { return &SubtractionBlock{} }
func (f *SubtractionBlockFactory) GetBlockInfo() blocks.BlockInfo {
	block := &SubtractionBlock{}
	return blocks.BlockInfo{
		Type: "subtract", Name: "Subtraction", Description: "Subtract a number from the input",
		Category: "math", BlockGroup: blocks.PropagationGroup, Inputs: block.GetInputs(), Outputs: block.GetOutputs(),
		Version: "1.0.0", Author: "Block-Flow", Icon: "minus", Color: "#F44336",
	}
}

type MultiplicationBlockFactory struct{}

func (f *MultiplicationBlockFactory) CreateBlock() blocks.Block { return &MultiplicationBlock{} }
func (f *MultiplicationBlockFactory) GetBlockInfo() blocks.BlockInfo {
	block := &MultiplicationBlock{}
	return blocks.BlockInfo{
		Type: "multiply", Name: "Multiplication", Description: "Multiply the input by a number",
		Category: "math", BlockGroup: blocks.PropagationGroup, Inputs: block.GetInputs(), Outputs: block.GetOutputs(),
		Version: "1.0.0", Author: "Block-Flow", Icon: "times", Color: "#9C27B0",
	}
}

type DivisionBlockFactory struct{}

func (f *DivisionBlockFactory) CreateBlock() blocks.Block { return &DivisionBlock{} }
func (f *DivisionBlockFactory) GetBlockInfo() blocks.BlockInfo {
	block := &DivisionBlock{}
	return blocks.BlockInfo{
		Type: "divide", Name: "Division", Description: "Divide the input by a number",
		Category: "math", BlockGroup: blocks.PropagationGroup, Inputs: block.GetInputs(), Outputs: block.GetOutputs(),
		Version: "1.0.0", Author: "Block-Flow", Icon: "divide", Color: "#FF5722",
	}
}

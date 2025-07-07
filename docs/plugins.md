# Plugin Development Guide

## Overview

Block-Flow supports a plugin-based architecture that allows developers to create custom blocks that extend the platform's functionality. Plugins are compiled as Go shared libraries (.so files on Linux, .dll on Windows) and loaded dynamically at runtime.

## Plugin Architecture

### Plugin Interface

All plugins must implement the `Plugin` interface:

```go
type Plugin interface {
    // GetInfo returns metadata about the plugin
    GetInfo() PluginInfo
    
    // GetBlocks returns the block factories provided by this plugin
    GetBlocks() []blocks.BlockFactory
    
    // Initialize is called when the plugin is loaded
    Initialize() error
    
    // Shutdown is called when the plugin is unloaded
    Shutdown() error
}
```

### Plugin Info

```go
type PluginInfo struct {
    Name        string `json:"name"`
    Version     string `json:"version"`
    Description string `json:"description"`
    Author      string `json:"author"`
    License     string `json:"license"`
    Website     string `json:"website,omitempty"`
}
```

## Creating a Plugin

### 1. Plugin Structure

Create a new Go module for your plugin:

```
my-plugin/
├── go.mod
├── go.sum
├── plugin.go          # Main plugin implementation
├── blocks/
│   ├── my_block.go    # Custom block implementation
│   └── another_block.go
└── README.md
```

### 2. Implement the Plugin Interface

```go
package main

import (
    "block-flow/internal/blocks"
    "block-flow/internal/models"
)

// MyPlugin implements the Plugin interface
type MyPlugin struct {
    blocks []blocks.BlockFactory
}

func (p *MyPlugin) GetInfo() PluginInfo {
    return PluginInfo{
        Name:        "My Custom Plugin",
        Version:     "1.0.0",
        Description: "A collection of custom blocks",
        Author:      "Your Name",
        License:     "MIT",
    }
}

func (p *MyPlugin) GetBlocks() []blocks.BlockFactory {
    return p.blocks
}

func (p *MyPlugin) Initialize() error {
    // Initialize your blocks
    p.blocks = []blocks.BlockFactory{
        &MyBlockFactory{},
        &AnotherBlockFactory{},
    }
    return nil
}

func (p *MyPlugin) Shutdown() error {
    // Clean up resources
    return nil
}

// Plugin entry point - required for Go plugins
func NewPlugin() Plugin {
    return &MyPlugin{}
}
```

### 3. Implement Custom Blocks

```go
package main

import (
    "context"
    "fmt"
    
    "block-flow/internal/blocks"
    "block-flow/internal/models"
)

// MyBlock implements the Block interface
type MyBlock struct{}

func (b *MyBlock) GetType() string {
    return "my-custom-block"
}

func (b *MyBlock) GetName() string {
    return "My Custom Block"
}

func (b *MyBlock) GetDescription() string {
    return "This block does something custom"
}

func (b *MyBlock) GetCategory() string {
    return "custom"
}

func (b *MyBlock) GetInputs() int {
    return 1
}

func (b *MyBlock) GetOutputs() int {
    return 1
}

func (b *MyBlock) GetProperties() []blocks.PropertyDefinition {
    return []blocks.PropertyDefinition{
        {
            Name:         "message",
            Type:         "string",
            DisplayName:  "Message",
            Description:  "The message to process",
            Required:     true,
            DefaultValue: "Hello",
        },
        {
            Name:         "multiplier",
            Type:         "number",
            DisplayName:  "Multiplier",
            Description:  "How many times to repeat the message",
            Required:     false,
            DefaultValue: 1,
            Validation: blocks.Validation{
                Min: &[]float64{1}[0],
                Max: &[]float64{10}[0],
            },
        },
    }
}

func (b *MyBlock) Validate(properties map[string]interface{}) error {
    if _, ok := properties["message"]; !ok {
        return fmt.Errorf("message property is required")
    }
    return nil
}

func (b *MyBlock) Execute(ctx *models.BlockExecutionContext, properties map[string]interface{}) ([]*models.Message, error) {
    // Get the input message
    inputMsg := ctx.Message
    
    // Get properties
    message, _ := properties["message"].(string)
    multiplier, _ := properties["multiplier"].(float64)
    if multiplier == 0 {
        multiplier = 1
    }
    
    // Process the message
    result := ""
    for i := 0; i < int(multiplier); i++ {
        result += message + " "
    }
    
    // Create output message
    outputMsg := models.NewMessage(result)
    outputMsg.Topic = inputMsg.Topic
    outputMsg.Source = ctx.NodeID
    
    // Log debug information
    ctx.Logger.Debug("Custom block executed", map[string]interface{}{
        "input":      inputMsg.Payload,
        "output":     result,
        "multiplier": multiplier,
    })
    
    return []*models.Message{outputMsg}, nil
}

// MyBlockFactory creates instances of MyBlock
type MyBlockFactory struct{}

func (f *MyBlockFactory) CreateBlock() blocks.Block {
    return &MyBlock{}
}

func (f *MyBlockFactory) GetBlockInfo() blocks.BlockInfo {
    return blocks.BlockInfo{
        Type:        "my-custom-block",
        Name:        "My Custom Block",
        Description: "A custom block that repeats messages",
        Category:    "custom",
        Version:     "1.0.0",
        Author:      "Your Name",
        Color:       "#9C27B0",
    }
}
```

## Building Plugins

### 1. Go Module Setup

```bash
# Initialize the module
go mod init my-plugin

# Add Block-Flow as a dependency
go get block-flow@latest

# Tidy dependencies
go mod tidy
```

### 2. Build as Shared Library

```bash
# Build for Linux
go build -buildmode=plugin -o my-plugin.so .

# Build for Windows (cross-compilation)
GOOS=windows go build -buildmode=plugin -o my-plugin.dll .
```

### 3. Plugin Deployment

1. Copy the compiled plugin file to the `data/plugins/` directory
2. Restart Block-Flow to load the plugin
3. The new blocks will be available in the block registry

## Plugin Configuration

### Plugin Manifest (optional)

Create a `plugin.json` file alongside your plugin binary:

```json
{
  "name": "my-custom-plugin",
  "version": "1.0.0",
  "description": "A collection of custom blocks",
  "author": "Your Name",
  "license": "MIT",
  "website": "https://github.com/yourname/my-plugin",
  "binary": "my-plugin.so",
  "blocks": [
    {
      "type": "my-custom-block",
      "name": "My Custom Block",
      "category": "custom"
    }
  ],
  "dependencies": [
    "some-other-plugin"
  ]
}
```

## Best Practices

### 1. Error Handling

Always handle errors gracefully and provide meaningful error messages:

```go
func (b *MyBlock) Execute(ctx *models.BlockExecutionContext, properties map[string]interface{}) ([]*models.Message, error) {
    // Validate inputs
    if ctx.Message == nil {
        return nil, fmt.Errorf("no input message provided")
    }
    
    // Process with error handling
    result, err := b.processMessage(ctx.Message)
    if err != nil {
        return nil, fmt.Errorf("failed to process message: %w", err)
    }
    
    return []*models.Message{result}, nil
}
```

### 2. Context Handling

Respect the context for cancellation and timeouts:

```go
func (b *MyBlock) Execute(ctx *models.BlockExecutionContext, properties map[string]interface{}) ([]*models.Message, error) {
    // Check for cancellation
    select {
    case <-ctx.Context.Done():
        return nil, ctx.Context.Err()
    default:
    }
    
    // Long-running operation with context
    result := make(chan string, 1)
    go func() {
        // Do work
        result <- "processed"
    }()
    
    select {
    case <-ctx.Context.Done():
        return nil, ctx.Context.Err()
    case r := <-result:
        return []*models.Message{models.NewMessage(r)}, nil
    }
}
```

### 3. Logging

Use the provided logger for debugging and monitoring:

```go
func (b *MyBlock) Execute(ctx *models.BlockExecutionContext, properties map[string]interface{}) ([]*models.Message, error) {
    ctx.Logger.Info("Starting block execution", map[string]interface{}{
        "node_id": ctx.NodeID,
        "flow_id": ctx.FlowID,
    })
    
    // Process...
    
    ctx.Logger.Debug("Block execution completed", map[string]interface{}{
        "output_count": len(outputs),
    })
    
    return outputs, nil
}
```

### 4. State Management

Use the block execution context for temporary state:

```go
func (b *MyBlock) Execute(ctx *models.BlockExecutionContext, properties map[string]interface{}) ([]*models.Message, error) {
    // Get previous state
    counter, _ := ctx.GetState("counter")
    if counter == nil {
        counter = 0
    }
    
    // Update state
    newCounter := counter.(int) + 1
    ctx.SetState("counter", newCounter)
    
    // Process...
    
    return outputs, nil
}
```

### 5. Property Validation

Provide comprehensive property validation:

```go
func (b *MyBlock) Validate(properties map[string]interface{}) error {
    // Check required properties
    url, ok := properties["url"].(string)
    if !ok || url == "" {
        return fmt.Errorf("url property is required and must be a non-empty string")
    }
    
    // Validate URL format
    if _, err := url.Parse(url); err != nil {
        return fmt.Errorf("invalid URL format: %w", err)
    }
    
    // Check optional properties with defaults
    timeout, ok := properties["timeout"].(float64)
    if ok && timeout < 0 {
        return fmt.Errorf("timeout must be a positive number")
    }
    
    return nil
}
```

## Testing Plugins

### Unit Testing

```go
package main

import (
    "context"
    "testing"
    
    "block-flow/internal/models"
)

func TestMyBlock_Execute(t *testing.T) {
    block := &MyBlock{}
    
    // Create test context
    ctx := &models.BlockExecutionContext{
        Context: context.Background(),
        NodeID:  "test-node",
        FlowID:  "test-flow",
        Message: models.NewMessage("test input"),
        Logger:  &TestLogger{},
    }
    
    // Test properties
    properties := map[string]interface{}{
        "message":    "Hello",
        "multiplier": 3.0,
    }
    
    // Execute block
    outputs, err := block.Execute(ctx, properties)
    
    // Assertions
    if err != nil {
        t.Fatalf("Expected no error, got: %v", err)
    }
    
    if len(outputs) != 1 {
        t.Fatalf("Expected 1 output, got: %d", len(outputs))
    }
    
    expected := "Hello Hello Hello "
    if outputs[0].Payload != expected {
        t.Errorf("Expected %q, got %q", expected, outputs[0].Payload)
    }
}

type TestLogger struct{}

func (l *TestLogger) Debug(msg string, fields map[string]interface{}) {}
func (l *TestLogger) Info(msg string, fields map[string]interface{})  {}
func (l *TestLogger) Warn(msg string, fields map[string]interface{})  {}
func (l *TestLogger) Error(msg string, err error, fields map[string]interface{}) {}
```

## Example Plugins

### HTTP Request Block

```go
type HTTPRequestBlock struct{}

func (b *HTTPRequestBlock) Execute(ctx *models.BlockExecutionContext, properties map[string]interface{}) ([]*models.Message, error) {
    url := properties["url"].(string)
    method := properties["method"].(string)
    
    // Create HTTP request
    req, err := http.NewRequestWithContext(ctx.Context, method, url, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    // Send request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()
    
    // Read response
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }
    
    // Create output message
    output := models.NewMessage(map[string]interface{}{
        "statusCode": resp.StatusCode,
        "headers":    resp.Header,
        "body":       string(body),
    })
    
    return []*models.Message{output}, nil
}
```

### Database Query Block

```go
type DatabaseQueryBlock struct{}

func (b *DatabaseQueryBlock) Execute(ctx *models.BlockExecutionContext, properties map[string]interface{}) ([]*models.Message, error) {
    query := properties["query"].(string)
    connectionString := properties["connection"].(string)
    
    // Connect to database
    db, err := sql.Open("postgres", connectionString)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }
    defer db.Close()
    
    // Execute query
    rows, err := db.QueryContext(ctx.Context, query)
    if err != nil {
        return nil, fmt.Errorf("query failed: %w", err)
    }
    defer rows.Close()
    
    // Process results
    var results []map[string]interface{}
    columns, _ := rows.Columns()
    
    for rows.Next() {
        values := make([]interface{}, len(columns))
        valuePtrs := make([]interface{}, len(columns))
        for i := range values {
            valuePtrs[i] = &values[i]
        }
        
        rows.Scan(valuePtrs...)
        
        row := make(map[string]interface{})
        for i, col := range columns {
            row[col] = values[i]
        }
        results = append(results, row)
    }
    
    // Create output message
    output := models.NewMessage(results)
    output.SetHeader("row-count", fmt.Sprintf("%d", len(results)))
    
    return []*models.Message{output}, nil
}
```

## Plugin Distribution

### Publishing

1. Create a GitHub repository for your plugin
2. Tag releases with semantic versioning
3. Provide pre-compiled binaries for different platforms
4. Include comprehensive documentation and examples

### Plugin Registry (Future)

A centralized plugin registry will be available for discovering and installing plugins:

```bash
# Install a plugin from registry
block-flow plugin install username/plugin-name

# List installed plugins
block-flow plugin list

# Update plugins
block-flow plugin update
```

This plugin system provides a powerful way to extend Block-Flow's functionality while maintaining clean separation between core functionality and custom implementations.

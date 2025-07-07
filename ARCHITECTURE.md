# Block-Flow Application Architecture

## Overview
Block-Flow is a simplified Node-RED-like visual programming platform built in Go. It enables users to create data processing workflows using a visual block-based interface with a React frontend and Go backend.

## System Architecture

### High-Level Architecture
```
┌─────────────────┐    HTTP/WebSocket    ┌─────────────────┐
│  React Frontend │ ◄─────────────────► │   Go Backend    │
│   (Web UI)      │                     │  (REST + WS)    │
└─────────────────┘                     └─────────────────┘
                                                │
                                                ▼
                                        ┌─────────────────┐
                                        │ Flow Engine     │
                                        │ (Execution)     │
                                        └─────────────────┘
                                                │
                                                ▼
                                        ┌─────────────────┐
                                        │ Block Registry  │
                                        │ (Built-in +     │
                                        │  Plugins)       │
                                        └─────────────────┘
                                                │
                                                ▼
                                        ┌─────────────────┐
                                        │ Persistence     │
                                        │ (JSON Files)    │
                                        └─────────────────┘
```

### Core Components

#### 1. Flow Engine
- **Purpose**: Executes data flows based on block configurations
- **Responsibilities**:
  - Parse flow definitions from JSON
  - Execute blocks in the correct order
  - Handle data passing between blocks
  - Manage flow state and debugging information
  - Provide real-time execution feedback

#### 2. Block System
- **Block Interface**: Common interface for all blocks
- **Built-in Blocks**: Standard set of processing blocks
- **Plugin System**: Dynamic loading of custom block libraries
- **Block Registry**: Central registry for all available blocks

#### 3. API Layer
- **REST API**: CRUD operations for flows and blocks
- **WebSocket**: Real-time communication for flow execution and debugging
- **Authentication**: Basic security (future enhancement)

#### 4. Persistence Layer
- **Flow Storage**: JSON-based flow configuration storage
- **Plugin Management**: Dynamic plugin discovery and loading
- **Configuration**: Application settings and preferences

## Project Structure

```
block-flow/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   │   ├── flow.go            # Flow CRUD handlers
│   │   │   ├── block.go           # Block management handlers
│   │   │   └── websocket.go       # WebSocket handlers
│   │   ├── middleware/
│   │   │   ├── cors.go            # CORS middleware
│   │   │   ├── logging.go         # Request logging
│   │   │   └── recovery.go        # Panic recovery
│   │   └── router.go              # HTTP router setup
│   ├── engine/
│   │   ├── flow.go                # Flow execution engine
│   │   ├── executor.go            # Block execution logic
│   │   ├── context.go             # Execution context
│   │   └── debugger.go            # Flow debugging support
│   ├── blocks/
│   │   ├── interface.go           # Block interface definition
│   │   ├── registry.go            # Block registry
│   │   ├── builtin/
│   │   │   ├── input.go           # Input blocks (inject, http in, etc.)
│   │   │   ├── output.go          # Output blocks (debug, http out, etc.)
│   │   │   ├── processing.go      # Processing blocks (function, switch, etc.)
│   │   │   └── utility.go         # Utility blocks (delay, change, etc.)
│   │   └── plugin/
│   │       ├── loader.go          # Plugin loading mechanism
│   │       └── manager.go         # Plugin lifecycle management
│   ├── models/
│   │   ├── flow.go                # Flow data structures
│   │   ├── block.go               # Block data structures
│   │   ├── message.go             # Message/payload structures
│   │   └── plugin.go              # Plugin metadata structures
│   ├── storage/
│   │   ├── interface.go           # Storage interface
│   │   ├── file.go                # File-based storage implementation
│   │   └── memory.go              # In-memory storage (for testing)
│   └── config/
│       └── config.go              # Application configuration
├── pkg/
│   ├── logger/
│   │   └── logger.go              # Structured logging
│   └── utils/
│       ├── json.go                # JSON utilities
│       └── validation.go          # Input validation
├── plugins/
│   ├── example/
│   │   ├── plugin.go              # Example plugin implementation
│   │   └── blocks.go              # Custom blocks in plugin
│   └── README.md                  # Plugin development guide
├── web/
│   ├── public/
│   │   └── index.html             # Static files
│   ├── src/                       # React frontend (future)
│   └── package.json               # Frontend dependencies (future)
├── data/
│   ├── flows/
│   │   └── default.json           # Default flow configuration
│   └── plugins/                   # Plugin binaries directory
├── docs/
│   ├── api.md                     # API documentation
│   ├── blocks.md                  # Block development guide
│   └── plugins.md                 # Plugin development guide
├── scripts/
│   ├── build.sh                   # Build scripts
│   └── dev.sh                     # Development scripts
├── tests/
│   ├── integration/
│   └── testdata/
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── ARCHITECTURE.md                # This file
```

## Functional Requirements

### Core Features

#### 1. Flow Management
- **FR-1.1**: Create, read, update, delete flow configurations
- **FR-1.2**: Load flow configuration from JSON file on startup
- **FR-1.3**: Save flow configuration to JSON file
- **FR-1.4**: Support single data flow (simplified version)

#### 2. Block System
- **FR-2.1**: Provide standard set of built-in blocks:
  - **Input blocks**: Inject (manual trigger), Timer, HTTP In
  - **Processing blocks**: Function (JavaScript-like), Switch, Change
  - **Output blocks**: Debug, HTTP Response, File Write
  - **Utility blocks**: Delay, Template
- **FR-2.2**: Execute blocks based on flow configuration
- **FR-2.3**: Pass messages between connected blocks
- **FR-2.4**: Support block configuration and parameters

#### 3. Plugin System
- **FR-3.1**: Load custom block libraries as plugins
- **FR-3.2**: Discover plugins in designated directory
- **FR-3.3**: Register custom blocks from plugins
- **FR-3.4**: Hot-reload plugins during development

#### 4. API Interface
- **FR-4.1**: REST API for flow and block management
- **FR-4.2**: WebSocket for real-time flow execution updates
- **FR-4.3**: Get available blocks and their configurations
- **FR-4.4**: Trigger flow execution manually

#### 5. Debugging and Monitoring
- **FR-5.1**: Debug output for message flow between blocks
- **FR-5.2**: Flow execution status and error reporting
- **FR-5.3**: Block-level execution timing and statistics
- **FR-5.4**: Real-time execution visualization data

## Non-Functional Requirements

### Performance
- **NFR-1.1**: Support at least 100 blocks in a single flow
- **NFR-1.2**: Execute simple flows in under 100ms
- **NFR-1.3**: Handle concurrent flow executions efficiently
- **NFR-1.4**: Memory usage should not exceed 100MB for basic operations

### Scalability
- **NFR-2.1**: Plugin system should support 50+ custom blocks
- **NFR-2.2**: Support for multiple concurrent WebSocket connections
- **NFR-2.3**: Horizontal scaling preparation (stateless design)

### Reliability
- **NFR-3.1**: Graceful error handling and recovery
- **NFR-3.2**: Flow execution should not crash the application
- **NFR-3.3**: Plugin failures should be isolated
- **NFR-3.4**: Automatic flow state persistence

### Maintainability
- **NFR-4.1**: Clean architecture with separation of concerns
- **NFR-4.2**: Comprehensive logging and debugging support
- **NFR-4.3**: Unit test coverage > 80%
- **NFR-4.4**: Clear documentation and API specifications

### Security
- **NFR-5.1**: Input validation for all API endpoints
- **NFR-5.2**: Safe plugin loading (sandbox future enhancement)
- **NFR-5.3**: Protection against code injection in function blocks
- **NFR-5.4**: CORS configuration for frontend integration

### Usability
- **NFR-6.1**: RESTful API design following OpenAPI standards
- **NFR-6.2**: Real-time feedback for flow execution
- **NFR-6.3**: Clear error messages and debugging information
- **NFR-6.4**: Easy plugin development workflow

## Technology Stack

### Backend
- **Language**: Go 1.24.3
- **HTTP Framework**: Gorilla Mux or Gin
- **WebSocket**: Gorilla WebSocket
- **JSON Processing**: Standard library encoding/json
- **Logging**: Structured logging (logrus or zap)
- **Testing**: Standard testing package + testify

### Frontend (Future)
- **Framework**: React 18+
- **State Management**: Redux Toolkit or Zustand
- **UI Components**: Material-UI or Ant Design
- **Flow Visualization**: React Flow or custom canvas
- **HTTP Client**: Axios
- **WebSocket**: Native WebSocket API

### Development Tools
- **Build**: Make + Go build
- **Linting**: golangci-lint
- **Formatting**: gofmt, goimports
- **Documentation**: godoc
- **Testing**: Go test + testify

## Implementation Phases

### Phase 1: Core Backend (Week 1-2)
1. Project setup and structure
2. Basic models and interfaces
3. File-based storage implementation
4. Simple flow execution engine
5. Built-in blocks (inject, debug, function)

### Phase 2: API and Plugin System (Week 3-4)
1. REST API implementation
2. WebSocket for real-time updates
3. Plugin loading mechanism
4. Example plugin development
5. Extended built-in block library

### Phase 3: Advanced Features (Week 5-6)
1. Flow debugging and monitoring
2. Error handling and recovery
3. Performance optimization
4. Comprehensive testing
5. Documentation and examples

### Phase 4: Frontend Integration (Week 7-8)
1. React frontend setup
2. Flow visualization component
3. Block palette and properties
4. Real-time execution monitoring
5. Full integration testing

## Getting Started

### Prerequisites
- Go 1.24.3 or later
- Node.js 18+ (for frontend)
- Make (optional, for build scripts)

### Development Setup
1. Clone the repository
2. Run `go mod init block-flow`
3. Install dependencies
4. Set up development environment
5. Run tests to verify setup

### Building and Running
```bash
# Build the application
make build

# Run in development mode
make dev

# Run tests
make test

# Build for production
make build-prod
```

This architecture provides a solid foundation for your block-flow application with clear separation of concerns, extensibility through plugins, and a path for future enhancements.

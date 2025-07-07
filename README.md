# Block-Flow

A simplified Node-RED-like visual programming platform built in Go. Block-Flow enables users to create data processing workflows using a visual block-based interface.

## 🎯 Project Status

**✅ Phase 1 - Core Architecture & Backend (COMPLETED)**
- [x] Project structure and architecture design
- [x] Core models (Flow, Node, Message, etc.)
- [x] File-based storage implementation
- [x] Flow execution engine (basic)
- [x] REST API with all endpoints
- [x] Block registry and interface system
- [x] Plugin architecture foundation
- [x] Configuration management
- [x] Error handling and logging
- [x] Comprehensive documentation

**✅ Phase 2 - Built-in Blocks & Execution (COMPLETED)**
- [x] Built-in math blocks (add, subtract, multiply, divide)
- [x] Input blocks (inject with number/string/boolean support)
- [x] Output blocks (debug with console logging)
- [x] Complete flow execution logic with message passing
- [x] Block registry with all built-in blocks
- [x] Comprehensive error handling and validation
- [x] Example flows demonstrating math operations

**✅ Phase 3 - React Frontend (COMPLETED)**
- [x] React frontend with TypeScript
- [x] Material-UI design system
- [x] Visual flow editor with drag-and-drop
- [x] Block palette with categorized blocks
- [x] Real-time console for debugging output
- [x] WebSocket integration for live updates
- [x] Flow persistence and execution controls

**🚧 Phase 4 - Enhancement (NEXT STEPS)**
- [ ] Block property editing panels
- [ ] More advanced block types (conditionals, loops, etc.)
- [ ] Flow templates and examples
- [ ] Plugin development tools
- [ ] Performance optimization
- [ ] Comprehensive testing

## 🚀 Features

- 🎯 **Clean Architecture**: Modular design with separation of concerns
- 🧩 **Plugin System**: Extensible architecture for custom blocks
- 💾 **JSON Storage**: File-based flow configuration persistence
- 🌐 **REST API**: Comprehensive RESTful interface
- ⚡ **Real-time Updates**: WebSocket-based live debugging and execution feedback
- 🎨 **Visual Editor**: React-based drag-and-drop flow designer
- 🧮 **Built-in Blocks**: Math operations, input/output, and debugging blocks
- 📱 **Responsive UI**: Modern Material-UI interface
- 📚 **Well Documented**: Extensive documentation and examples

## 🏃 Quick Start

### Prerequisites
- Go 1.24.3 or later
- Node.js 18+ and npm

### 1. Clone and Setup
```bash
git clone <repository-url>
cd block-flow
make setup  # Installs all dependencies (Go + npm)
```

### 2. Start the Backend
```bash
make run
# Or for development with hot reload:
make dev
```

### 3. Start the Frontend (in a new terminal)
```bash
make dev-frontend
# Or manually:
cd web && npm start
```

### 4. Open the Application
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080/api/v1

## 🎮 Using the Visual Editor

1. **Adding Blocks**: Drag blocks from the left palette onto the canvas
2. **Connecting Blocks**: Click output ports (right side) then input ports (left side)
3. **Moving Blocks**: Drag blocks around by their header
4. **Deleting Blocks**: Right-click for context menu
5. **Saving Flows**: Click "Save" in the toolbar
6. **Running Flows**: Click "Start" then "Trigger" to execute

## 🧮 Available Blocks

### Input Blocks
- **Inject**: Manual input with configurable payload (numbers, strings, booleans)

### Math Blocks
- **Addition**: Add a number to the input
- **Subtraction**: Subtract a number from the input
- **Multiplication**: Multiply input by a number
- **Division**: Divide input by a number (with zero-division protection)

### Output Blocks
- **Debug**: Output values to console for debugging

### Example Flow
Create a flow: `Inject(5) → Add(3) → Multiply(2) → Debug` = Result: 16

## 📁 Project Structure

```
block-flow/
├── cmd/server/              # Application entry point
├── internal/                # Private application code
│   ├── api/                # HTTP handlers and middleware
│   │   ├── handlers/       # Route handlers
│   │   └── middleware/     # HTTP middleware
│   ├── engine/             # Flow execution engine
│   ├── blocks/             # Block system and interfaces
│   ├── models/             # Data structures
│   ├── storage/            # Persistence layer
│   └── config/             # Configuration management
├── pkg/                    # Public library code
│   └── logger/             # Logging utilities
├── plugins/                # Plugin examples (future)
├── web/                    # Frontend code (future)
│   └── public/             # Static files
├── data/                   # Runtime data
│   └── flows/              # Flow configurations
├── docs/                   # Documentation
│   ├── api.md              # API documentation
│   └── plugins.md          # Plugin development guide
└── ARCHITECTURE.md         # System architecture
```

## 🛠️ Installation & Setup

### Prerequisites
- Go 1.24.3 or later
- Make (optional)

### Quick Start

```bash
# Clone the repository
git clone <repository-url>
cd block-flow

# Download dependencies
go mod tidy

# Build the application
go build -o bin/block-flow-server.exe ./cmd/server

# Run the server
./bin/block-flow-server.exe
```

### Using Make

```bash
# Build
make build

# Run
make run

# Development mode (with hot reload, requires 'air' tool)
make dev

# Run tests
make test

# Format code
make fmt

# Lint code
make lint
```

## 🔌 API Endpoints

The server runs on `http://localhost:8080` by default.

### Core Endpoints

- **Health Check**: `GET /api/v1/health`
- **Flows**: 
  - `GET /api/v1/flows` - List all flows
  - `POST /api/v1/flows` - Create new flow
  - `GET /api/v1/flows/{id}` - Get flow details
  - `PUT /api/v1/flows/{id}` - Update flow
  - `DELETE /api/v1/flows/{id}` - Delete flow
- **Flow Control**:
  - `POST /api/v1/flows/{id}/start` - Start flow execution
  - `POST /api/v1/flows/{id}/stop` - Stop flow execution
  - `POST /api/v1/flows/{id}/trigger` - Manually trigger flow
  - `GET /api/v1/flows/{id}/status` - Get execution status
- **Blocks**:
  - `GET /api/v1/blocks` - List available block types
  - `GET /api/v1/blocks/{type}` - Get block type info

### Testing the API

```bash
# Check server health
curl http://localhost:8080/api/v1/health

# List flows
curl http://localhost:8080/api/v1/flows

# Create a simple flow
curl -X POST http://localhost:8080/api/v1/flows \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Flow",
    "description": "A simple test flow",
    "nodes": [
      {
        "id": "inject-1",
        "type": "inject",
        "name": "Start",
        "x": 100,
        "y": 100,
        "properties": {"payload": "Hello World"},
        "outputs": 1,
        "wires": [["debug-1"]]
      },
      {
        "id": "debug-1", 
        "type": "debug",
        "name": "Output",
        "x": 300,
        "y": 100,
        "properties": {"console": true},
        "inputs": 1
      }
    ]
  }'
```

## 📚 Documentation

- **[Architecture Guide](ARCHITECTURE.md)** - System design and architecture
- **[API Documentation](docs/api.md)** - REST API reference  
- **[Plugin Development](docs/plugins.md)** - Creating custom blocks

## 🔧 Configuration

Configure the application using environment variables:

```bash
# Server configuration
SERVER_ADDRESS=:8080
SERVER_READ_TIMEOUT=15s
SERVER_WRITE_TIMEOUT=15s

# Storage configuration  
DATA_DIR=./data
PLUGINS_DIR=./data/plugins

# Engine configuration
MAX_CONCURRENT_FLOWS=10
DEFAULT_TIMEOUT=30s
DEBUG_MODE=true

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

## 🧪 Example Flow

An example flow is included in `data/flows/example.json`:

```json
{
  "id": "example-flow-001",
  "name": "Example Flow", 
  "nodes": [
    {
      "id": "inject-1",
      "type": "inject",
      "name": "Manual Trigger",
      "properties": {
        "topic": "test",
        "payload": "Hello, World!"
      }
    },
    {
      "id": "debug-1",
      "type": "debug", 
      "name": "Debug Output",
      "properties": {
        "console": true,
        "complete": "payload"
      }
    }
  ],
  "connections": [...]
}
```

## 🔌 Plugin System

Block-Flow supports a plugin-based architecture for custom blocks:

```go
// Example custom block
type MyBlock struct{}

func (b *MyBlock) Execute(ctx *models.BlockExecutionContext, properties map[string]interface{}) ([]*models.Message, error) {
    // Custom block logic
    input := ctx.Message.Payload
    output := fmt.Sprintf("Processed: %v", input)
    
    return []*models.Message{models.NewMessage(output)}, nil
}
```

See [Plugin Development Guide](docs/plugins.md) for details.

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Test the running server
go run test-server.go
```

## 🔄 Development Workflow

1. **Start development server**:
   ```bash
   make dev  # Uses 'air' for hot reloading
   ```

2. **Make changes** to the code

3. **Test endpoints**:
   ```bash
   curl http://localhost:8080/api/v1/health
   ```

4. **Run tests**:
   ```bash
   make test
   ```

5. **Build for production**:
   ```bash
   make build-prod
   ```

## 🏗️ Architecture Highlights

### Clean Architecture
- **Separation of Concerns**: Clear boundaries between API, business logic, and storage
- **Dependency Inversion**: Core business logic doesn't depend on external frameworks
- **Testability**: Each layer can be tested independently

### Extensibility
- **Plugin System**: Dynamic loading of custom blocks
- **Block Registry**: Centralized management of available blocks
- **Message Pipeline**: Flexible message passing between blocks

### Scalability
- **Stateless Design**: Ready for horizontal scaling
- **Concurrent Execution**: Multiple flows can run simultaneously
- **Resource Management**: Proper context handling and cancellation

## 📝 TODO / Roadmap

### Immediate (Phase 2)
- [ ] Implement built-in blocks (inject, debug, function, switch, etc.)
- [ ] Complete flow execution engine with proper message routing
- [ ] Add WebSocket handlers for real-time updates
- [ ] Implement plugin loading mechanism
- [ ] Add flow validation and error handling

### Short-term
- [ ] React frontend for visual flow editing
- [ ] Block property editor UI
- [ ] Real-time execution visualization
- [ ] Flow import/export functionality
- [ ] Basic authentication and authorization

### Long-term
- [ ] Advanced debugging features (breakpoints, step-through)
- [ ] Flow versioning and history
- [ ] Performance monitoring and metrics
- [ ] Cloud deployment options
- [ ] Plugin marketplace

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes and add tests
4. Commit: `git commit -m 'Add amazing feature'`
5. Push: `git push origin feature/amazing-feature`
6. Submit a pull request

### Development Guidelines
- Follow Go conventions and best practices
- Write tests for new functionality
- Update documentation as needed
- Use meaningful commit messages
- Ensure code passes linting and formatting checks

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- Inspired by [Node-RED](https://nodered.org/) visual programming platform
- Built with [Gorilla Mux](https://github.com/gorilla/mux) for HTTP routing
- Uses [Gorilla WebSocket](https://github.com/gorilla/websocket) for real-time communication

---

**Happy Flow Building! 🚀**

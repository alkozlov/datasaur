#!/bin/bash

# Block-Flow Development Startup Script

echo "🚀 Starting Block-Flow Development Environment"
echo "=============================================="

# Check if required tools are installed
check_tool() {
    if ! command -v $1 &> /dev/null; then
        echo "❌ $1 is not installed. Please install $1 first."
        exit 1
    fi
}

echo "📋 Checking prerequisites..."
check_tool "go"
check_tool "node"
check_tool "npm"

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo "❌ Please run this script from the block-flow project root directory"
    exit 1
fi

# Install dependencies if needed
if [ ! -d "web/node_modules" ]; then
    echo "📦 Installing frontend dependencies..."
    cd web && npm install && cd ..
fi

# Build backend first to check for errors
echo "🔨 Building backend..."
if ! go build -o bin/block-flow-server ./cmd/server; then
    echo "❌ Backend build failed. Please fix the errors and try again."
    exit 1
fi

echo "✅ Backend build successful"

# Start backend in background
echo "🖥️  Starting backend server..."
./bin/block-flow-server &
BACKEND_PID=$!

# Wait a moment for backend to start
sleep 2

# Check if backend is running
if ! kill -0 $BACKEND_PID 2>/dev/null; then
    echo "❌ Backend failed to start"
    exit 1
fi

echo "✅ Backend running on http://localhost:8080"
echo "📊 API available at http://localhost:8080/api/v1"

# Start frontend
echo "🎨 Starting frontend development server..."
cd web
npm start &
FRONTEND_PID=$!
cd ..

echo ""
echo "🎉 Block-Flow is now running!"
echo ""
echo "📱 Frontend: http://localhost:3000"
echo "🖥️  Backend:  http://localhost:8080"
echo "📊 API:      http://localhost:8080/api/v1"
echo ""
echo "Press Ctrl+C to stop all services"

# Function to cleanup background processes
cleanup() {
    echo ""
    echo "🛑 Shutting down services..."
    
    if kill -0 $BACKEND_PID 2>/dev/null; then
        echo "   Stopping backend server..."
        kill $BACKEND_PID
    fi
    
    if kill -0 $FRONTEND_PID 2>/dev/null; then
        echo "   Stopping frontend server..."
        kill $FRONTEND_PID
    fi
    
    echo "✅ All services stopped"
    exit 0
}

# Set up trap to cleanup on exit
trap cleanup INT TERM

# Wait for either process to exit
wait

@echo off
REM Block-Flow Development Startup Script for Windows

echo 🚀 Starting Block-Flow Development Environment
echo ==============================================

REM Check if required tools are installed
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo ❌ Go is not installed. Please install Go first.
    exit /b 1
)

where node >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo ❌ Node.js is not installed. Please install Node.js first.
    exit /b 1
)

where npm >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo ❌ npm is not installed. Please install npm first.
    exit /b 1
)

REM Check if we're in the right directory
if not exist "go.mod" (
    echo ❌ Please run this script from the block-flow project root directory
    exit /b 1
)

echo 📋 Prerequisites check complete

REM Install frontend dependencies if needed
if not exist "web\node_modules" (
    echo 📦 Installing frontend dependencies...
    cd web
    npm install
    if %ERRORLEVEL% NEQ 0 (
        echo ❌ Failed to install frontend dependencies
        exit /b 1
    )
    cd ..
)

REM Build backend first to check for errors
echo 🔨 Building backend...
if not exist "bin" mkdir bin
go build -o bin\block-flow-server.exe .\cmd\server
if %ERRORLEVEL% NEQ 0 (
    echo ❌ Backend build failed. Please fix the errors and try again.
    exit /b 1
)

echo ✅ Backend build successful

REM Start backend in background
echo 🖥️  Starting backend server...
start "Block-Flow Backend" bin\block-flow-server.exe

REM Wait a moment for backend to start
timeout /t 3 /nobreak > nul

echo ✅ Backend started on http://localhost:8080
echo 📊 API available at http://localhost:8080/api/v1

REM Start frontend
echo 🎨 Starting frontend development server...
cd web
start "Block-Flow Frontend" npm start
cd ..

echo.
echo 🎉 Block-Flow is now starting!
echo.
echo 📱 Frontend: http://localhost:3000
echo 🖥️  Backend:  http://localhost:8080
echo 📊 API:      http://localhost:8080/api/v1
echo.
echo Services are starting in separate windows.
echo Close the console windows to stop the services.

pause

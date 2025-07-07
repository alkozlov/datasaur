# Frontend Development Guide

## Quick Start

1. **Install dependencies**:
```bash
cd web
npm install
```

2. **Start development server**:
```bash
npm start
```

3. **Build for production**:
```bash
npm run build
```

## Development Notes

### Fixed Issues:
- ✅ TypeScript configuration updated for ES2017+ features
- ✅ All TypeScript type errors resolved
- ✅ Added missing type dependencies (@types/react-draggable)
- ✅ Proper type annotations throughout codebase
- ✅ Material-UI color prop type issues fixed

### Architecture:
- **App.tsx**: Main application container with state management
- **components/BlockPalette.tsx**: Left sidebar with draggable blocks
- **components/FlowCanvas.tsx**: Central canvas with nodes and connections
- **components/Console.tsx**: Bottom console for debugging output
- **api.ts**: Backend API integration layer
- **types.ts**: TypeScript interface definitions

### Key Features:
1. **Drag & Drop**: Implemented with react-draggable and native HTML5 drag/drop
2. **Visual Connections**: SVG-based wire rendering between block ports
3. **Real-time Updates**: WebSocket integration for live console output
4. **Material Design**: Professional UI with Material-UI components
5. **Type Safety**: Full TypeScript coverage with proper error handling

### API Integration:
- REST endpoints for flow CRUD operations
- WebSocket for real-time debugging messages
- Automatic proxy configuration in development
- Error handling with user-friendly notifications

### Usage:
1. Drag blocks from left palette to center canvas
2. Click output ports (right side) then input ports (left side) to connect
3. Move blocks by dragging their headers
4. Right-click blocks for context menu (delete, etc.)
5. Use toolbar buttons to save, start, and trigger flows
6. Watch console for real-time debug output

The frontend is now fully functional and matches all specified requirements!

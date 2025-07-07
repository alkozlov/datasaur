import React, { useState, useRef, useCallback, useEffect } from 'react';
import { 
  Box, 
  Typography, 
  Menu, 
  MenuItem, 
  ListItemIcon, 
  ListItemText, 
  Button,
  ButtonGroup,
  Paper,
  Chip
} from '@mui/material';
import { 
  PlayArrow, 
  Stop, 
  Save, 
  Clear
} from '@mui/icons-material';
import { CanvasNode, CanvasConnection, ConsoleMessage } from '../types/canvas';
import { saveFlow, runFlow, stopFlow, getFlowStatus, FlowData } from '../api';
import Block from './Block';
import Connections from './Connection';
import { useConnectionLogic } from '../hooks/useConnectionLogic';
import { useBlockHoverLogic } from '../hooks/useBlockHoverLogic';

interface CanvasProps {
  onConsoleMessage: (message: ConsoleMessage) => void;
}

function Canvas({ onConsoleMessage }: CanvasProps) {
  const [nodes, setNodes] = useState<CanvasNode[]>([]);
  const [connections, setConnections] = useState<CanvasConnection[]>([]);
  const [selectedNodeId, setSelectedNodeId] = useState<string | null>(null);
  const [isDragging, setIsDragging] = useState(false);
  const [dragOffset, setDragOffset] = useState({ x: 0, y: 0 });
  const [contextMenu, setContextMenu] = useState<{ x: number; y: number; nodeId: string } | null>(null);
  const [flowStatus, setFlowStatus] = useState<{ running: boolean; flowId?: string }>({ running: false });
  const [currentFlowId, setCurrentFlowId] = useState<string | null>(null);
  const canvasRef = useRef<HTMLDivElement>(null);
  const dragNodeId = useRef<string | null>(null);

  const addConsoleMessage = useCallback((type: ConsoleMessage['type'], message: string) => {
    const consoleMessage: ConsoleMessage = {
      id: Date.now().toString(),
      timestamp: new Date(),
      type,
      message,
    };
    onConsoleMessage(consoleMessage);
  }, [onConsoleMessage]);

  // Connection logic hook
  const connectionLogic = useConnectionLogic({
    nodes,
    connections,
    setConnections,
    addConsoleMessage,
  });

  // Block hover logic hook
  const blockHoverLogic = useBlockHoverLogic({
    nodes,
    canBeConnectionSource: connectionLogic.canBeConnectionSource,
    isConnecting: connectionLogic.isConnecting,
    handleConnectionNodeMouseEnter: connectionLogic.handleConnectionNodeMouseEnter,
    handleConnectionNodeMouseLeave: connectionLogic.handleConnectionNodeMouseLeave,
  });

  // Flow management functions
  const handleSaveFlow = useCallback(async () => {
    if (nodes.length === 0) {
      addConsoleMessage('warning', 'Cannot save empty flow');
      return;
    }

    try {
      const flowData: FlowData = {
        id: currentFlowId || `flow-${Date.now()}`,
        name: `Flow ${new Date().toLocaleTimeString()}`,
        description: 'Created from canvas',
        nodes: nodes.map(node => ({
          id: node.id,
          type: node.type,
          name: node.name,
          x: node.x,
          y: node.y,
          properties: {}, // Default properties
          inputs: node.inputs || 1,
          outputs: node.outputs || 1,
        })),
        connections: connections.map(conn => ({
          id: conn.id,
          source: conn.sourceNodeId,
          source_port: conn.sourcePort,
          target: conn.targetNodeId,
          target_port: conn.targetPort,
        })),
      };

      const flowId = await saveFlow(flowData);
      setCurrentFlowId(flowId);
      addConsoleMessage('success', `Flow saved with ID: ${flowId}`);
    } catch (error) {
      addConsoleMessage('error', `Failed to save flow: ${error}`);
    }
  }, [nodes, connections, currentFlowId, addConsoleMessage]);

  const handleRunFlow = useCallback(async () => {
    if (!currentFlowId) {
      addConsoleMessage('warning', 'Please save the flow first');
      return;
    }

    try {
      await runFlow(currentFlowId);
      setFlowStatus({ running: true, flowId: currentFlowId });
      addConsoleMessage('success', 'Flow started successfully');
    } catch (error) {
      addConsoleMessage('error', `Failed to start flow: ${error}`);
    }
  }, [currentFlowId, addConsoleMessage]);

  const handleStopFlow = useCallback(async () => {
    if (!currentFlowId) {
      addConsoleMessage('warning', 'No flow is running');
      return;
    }

    try {
      await stopFlow(currentFlowId);
      setFlowStatus({ running: false });
      addConsoleMessage('success', 'Flow stopped');
    } catch (error) {
      addConsoleMessage('error', `Failed to stop flow: ${error}`);
    }
  }, [currentFlowId, addConsoleMessage]);

  const handleClearCanvas = useCallback(() => {
    setNodes([]);
    setConnections([]);
    setSelectedNodeId(null);
    setCurrentFlowId(null);
    setFlowStatus({ running: false });
    addConsoleMessage('info', 'Canvas cleared');
  }, [addConsoleMessage]);

  // Poll flow status if flow is running
  useEffect(() => {
    if (!currentFlowId || !flowStatus.running) return;

    const interval = setInterval(async () => {
      try {
        const status = await getFlowStatus(currentFlowId);
        setFlowStatus(prev => ({ ...prev, running: status.running }));
      } catch (error) {
        // Flow might have stopped or error occurred
        setFlowStatus(prev => ({ ...prev, running: false }));
      }
    }, 2000);

    return () => clearInterval(interval);
  }, [currentFlowId, flowStatus.running]);

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    const nodeData = e.dataTransfer.getData('application/json');
    if (!nodeData) return;

    const droppedNode = JSON.parse(nodeData);
    const canvasRect = canvasRef.current?.getBoundingClientRect();
    if (!canvasRect) return;

    const newNode: CanvasNode = {
      id: `${droppedNode.type}-${Date.now()}`,
      type: droppedNode.type,
      name: droppedNode.name,
      description: droppedNode.description,
      category: droppedNode.category,
      x: e.clientX - canvasRect.left - 60, // Center the node
      y: e.clientY - canvasRect.top - 16,
      selected: false,
      inputs: droppedNode.inputs || 1,
      outputs: droppedNode.outputs || 1,
    };

    setNodes(prev => [...prev, newNode]);
    addConsoleMessage('success', `Added ${newNode.name} block to canvas`);
  }, [addConsoleMessage]);

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault();
  }, []);

  const handleNodeMouseDown = useCallback((e: React.MouseEvent, nodeId: string) => {
    e.stopPropagation();
    if (e.button === 0 && !connectionLogic.isConnecting) { // Left click and not connecting
      setSelectedNodeId(nodeId);
      setNodes(prev => prev.map(node => ({ ...node, selected: node.id === nodeId })));
      
      const node = nodes.find(n => n.id === nodeId);
      if (node) {
        setIsDragging(true);
        dragNodeId.current = nodeId;
        setDragOffset({
          x: e.clientX - node.x,
          y: e.clientY - node.y,
        });
      }
    }
  }, [nodes, connectionLogic.isConnecting]);

  const handleNodeContextMenu = useCallback((e: React.MouseEvent, nodeId: string) => {
    e.preventDefault();
    e.stopPropagation();
    setSelectedNodeId(nodeId);
    setNodes(prev => prev.map(node => ({ ...node, selected: node.id === nodeId })));
    setContextMenu({ x: e.clientX, y: e.clientY, nodeId });
  }, []);

  const handleMouseMove = useCallback((e: MouseEvent) => {
    if (isDragging && dragNodeId.current) {
      setNodes(prev => prev.map(node => 
        node.id === dragNodeId.current
          ? { ...node, x: e.clientX - dragOffset.x, y: e.clientY - dragOffset.y }
          : node
      ));
    } else if (connectionLogic.isConnecting && connectionLogic.tempConnection) {
      // Update temporary connection line
      const canvasRect = canvasRef.current?.getBoundingClientRect();
      if (canvasRect) {
        connectionLogic.updateTempConnection(e.clientX, e.clientY, canvasRect);
      }
    }
  }, [isDragging, dragOffset, connectionLogic]);

  const handleMouseUp = useCallback(() => {
    setIsDragging(false);
    dragNodeId.current = null;
  }, []);

  const handleCanvasClick = useCallback((e: React.MouseEvent) => {
    if (e.target === e.currentTarget) {
      setSelectedNodeId(null);
      connectionLogic.setSelectedConnectionId(null);
      blockHoverLogic.clearHoverStates();
      setNodes(prev => prev.map(node => ({ ...node, selected: false })));
      
      // Cancel connection if clicking on empty canvas
      if (connectionLogic.isConnecting) {
        connectionLogic.cancelConnection();
      }
    }
    setContextMenu(null);
  }, [connectionLogic, blockHoverLogic]);

  const handleDeleteNode = useCallback((nodeId: string) => {
    const node = nodes.find(n => n.id === nodeId);
    if (node) {
      // Remove node
      setNodes(prev => prev.filter(n => n.id !== nodeId));
      
      // Remove all connections involving this node
      const removedConnections = connections.filter(
        c => c.sourceNodeId === nodeId || c.targetNodeId === nodeId
      );
      setConnections(prev => prev.filter(
        c => c.sourceNodeId !== nodeId && c.targetNodeId !== nodeId
      ));
      
      addConsoleMessage('info', `Deleted ${node.name} block and ${removedConnections.length} connections`);
      setSelectedNodeId(null);
    }
    setContextMenu(null);
  }, [nodes, connections, addConsoleMessage]);

  const handleDeleteConnection = useCallback((connectionId: string) => {
    setConnections(prev => prev.filter(c => c.id !== connectionId));
    addConsoleMessage('info', 'Connection deleted');
  }, [addConsoleMessage]);

  const handleKeyDown = useCallback((e: KeyboardEvent) => {
    if (e.key === 'Delete') {
      if (selectedNodeId) {
        handleDeleteNode(selectedNodeId);
      } else if (connectionLogic.selectedConnectionId) {
        connectionLogic.handleDeleteConnection(connectionLogic.selectedConnectionId);
        connectionLogic.setSelectedConnectionId(null);
      }
    } else if (e.key === 'Escape') {
      if (connectionLogic.isConnecting) {
        connectionLogic.cancelConnection();
      } else {
        setSelectedNodeId(null);
        connectionLogic.setSelectedConnectionId(null);
        blockHoverLogic.clearHoverStates();
        setNodes(prev => prev.map(node => ({ ...node, selected: false })));
      }
    }
  }, [selectedNodeId, connectionLogic, blockHoverLogic, handleDeleteNode]);

  useEffect(() => {
    document.addEventListener('mousemove', handleMouseMove);
    document.addEventListener('mouseup', handleMouseUp);
    document.addEventListener('keydown', handleKeyDown);

    return () => {
      document.removeEventListener('mousemove', handleMouseMove);
      document.removeEventListener('mouseup', handleMouseUp);
      document.removeEventListener('keydown', handleKeyDown);
    };
  }, [handleMouseMove, handleMouseUp, handleKeyDown]);

  return (
    <Box
      ref={canvasRef}
      sx={{
        width: '80%',
        height: '85vh',
        backgroundColor: '#ffffff',
        border: '1px solid #ddd',
        padding: 2,
        boxSizing: 'border-box',
        position: 'relative',
        overflow: 'hidden',
        cursor: isDragging ? 'grabbing' : connectionLogic.isConnecting ? 'crosshair' : 'default',
      }}
      onDrop={handleDrop}
      onDragOver={handleDragOver}
      onClick={handleCanvasClick}
    >
      {/* Canvas Header with Controls */}
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
          <Typography variant="h6">
            Canvas {connectionLogic.isConnecting && '(Connecting...)'}
          </Typography>
          {flowStatus.running && (
            <Chip 
              size="small" 
              color="success" 
              label="Running" 
              sx={{ fontSize: '10px' }}
            />
          )}
        </Box>
        
        <Paper elevation={1} sx={{ p: 0.5 }}>
          <ButtonGroup size="small" variant="outlined">
            <Button
              startIcon={<Save />}
              onClick={handleSaveFlow}
              disabled={nodes.length === 0}
            >
              Save
            </Button>
            <Button
              startIcon={<PlayArrow />}
              onClick={handleRunFlow}
              disabled={!currentFlowId || flowStatus.running}
              color="success"
            >
              Run
            </Button>
            <Button
              startIcon={<Stop />}
              onClick={handleStopFlow}
              disabled={!flowStatus.running}
              color="error"
            >
              Stop
            </Button>
            <Button
              startIcon={<Clear />}
              onClick={handleClearCanvas}
              color="warning"
            >
              Clear
            </Button>
          </ButtonGroup>
        </Paper>
      </Box>
      
      {/* Connections component */}
      <Connections
        connections={connections}
        nodes={nodes}
        selectedConnectionId={connectionLogic.selectedConnectionId}
        tempConnection={connectionLogic.tempConnection}
        onConnectionClick={connectionLogic.handleConnectionClick}
      />
      
      {/* Render nodes using Block component */}
      {nodes.map((node) => (
        <Block
          key={node.id}
          node={node}
          isSelected={node.selected}
          isHovered={blockHoverLogic.hoveredNodeId === node.id}
          isConnectionTarget={connectionLogic.hoveredConnectionTarget === node.id}
          isConnecting={connectionLogic.isConnecting}
          showConnectionButton={blockHoverLogic.showConnectionButton === node.id}
          canBeConnectionSource={connectionLogic.canBeConnectionSource(node)}
          onMouseDown={handleNodeMouseDown}
          onMouseEnter={blockHoverLogic.handleNodeMouseEnter}
          onMouseLeave={blockHoverLogic.handleNodeMouseLeave}
          onMouseUp={connectionLogic.handleNodeMouseUp}
          onContextMenu={handleNodeContextMenu}
          onConnectionButtonMouseDown={connectionLogic.handleConnectionButtonMouseDown}
          onConnectionButtonMouseEnter={blockHoverLogic.handleConnectionButtonMouseEnter}
          onConnectionButtonMouseLeave={blockHoverLogic.handleConnectionButtonMouseLeave}
          isDragging={isDragging}
        />
      ))}

      {contextMenu && (
        <Menu
          open={Boolean(contextMenu)}
          onClose={() => setContextMenu(null)}
          anchorReference="anchorPosition"
          anchorPosition={{ top: contextMenu.y, left: contextMenu.x }}
        >
          <MenuItem onClick={() => handleDeleteNode(contextMenu.nodeId)}>
            <ListItemIcon>
              <Typography>🗑️</Typography>
            </ListItemIcon>
            <ListItemText>Delete</ListItemText>
          </MenuItem>
        </Menu>
      )}
    </Box>
  );
}

export default Canvas;

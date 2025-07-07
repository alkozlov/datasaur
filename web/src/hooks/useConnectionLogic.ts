import { useState, useCallback } from 'react';
import { CanvasNode, CanvasConnection, ConsoleMessage } from '../types/canvas';

interface UseConnectionLogicProps {
  nodes: CanvasNode[];
  connections: CanvasConnection[];
  setConnections: React.Dispatch<React.SetStateAction<CanvasConnection[]>>;
  addConsoleMessage: (type: ConsoleMessage['type'], message: string) => void;
}

export function useConnectionLogic({
  nodes,
  connections,
  setConnections,
  addConsoleMessage,
}: UseConnectionLogicProps) {
  const [isConnecting, setIsConnecting] = useState(false);
  const [connectionStart, setConnectionStart] = useState<{ nodeId: string; x: number; y: number } | null>(null);
  const [tempConnection, setTempConnection] = useState<{ startX: number; startY: number; endX: number; endY: number } | null>(null);
  const [hoveredConnectionTarget, setHoveredConnectionTarget] = useState<string | null>(null);
  const [selectedConnectionId, setSelectedConnectionId] = useState<string | null>(null);

  // Get block group from category (mapping frontend categories to backend groups)
  const getBlockGroup = useCallback((category: string): string => {
    switch (category) {
      case 'input': return 'Input';
      case 'math': return 'Propagation';
      case 'output': return 'Action';
      default: return 'Unknown';
    }
  }, []);

  // Check if a node can be a connection source (Input or Propagation groups)
  const canBeConnectionSource = useCallback((node: CanvasNode): boolean => {
    const blockGroup = getBlockGroup(node.category);
    return blockGroup === 'Input' || blockGroup === 'Propagation';
  }, [getBlockGroup]);

  // Check if a node can be a connection target (Propagation or Action groups, but not Input)
  const canBeConnectionTarget = useCallback((node: CanvasNode): boolean => {
    const blockGroup = getBlockGroup(node.category);
    return blockGroup === 'Propagation' || blockGroup === 'Action';
  }, [getBlockGroup]);

  // Check if a connection between two nodes is valid
  const isValidConnection = useCallback((sourceNode: CanvasNode, targetNode: CanvasNode): boolean => {
    return canBeConnectionSource(sourceNode) && 
           canBeConnectionTarget(targetNode) && 
           sourceNode.id !== targetNode.id;
  }, [canBeConnectionSource, canBeConnectionTarget]);

  // Get the connection anchor points (right side of source, left side of target)
  const getConnectionAnchor = useCallback((node: CanvasNode) => {
    return {
      x: node.x + 120, // Right side of block
      y: node.y + 16,  // Center height of block
    };
  }, []);

  // Handle connection button mouse down (start connection)
  const handleConnectionButtonMouseDown = useCallback((e: React.MouseEvent, nodeId: string) => {
    e.stopPropagation();
    if (e.button === 0) { // Left click only
      const node = nodes.find(n => n.id === nodeId);
      if (!node || !canBeConnectionSource(node)) return;

      const anchor = getConnectionAnchor(node);
      setIsConnecting(true);
      setConnectionStart({
        nodeId,
        x: anchor.x,
        y: anchor.y,
      });
      setTempConnection({
        startX: anchor.x,
        startY: anchor.y,
        endX: anchor.x,
        endY: anchor.y,
      });
      addConsoleMessage('info', 'Creating connection...');
    }
  }, [nodes, canBeConnectionSource, getConnectionAnchor, addConsoleMessage]);

  // Handle node hover for connection target highlighting
  const handleConnectionNodeMouseEnter = useCallback((nodeId: string) => {
    if (isConnecting && connectionStart) {
      const node = nodes.find(n => n.id === nodeId);
      if (node && canBeConnectionTarget(node) && node.id !== connectionStart.nodeId) {
        setHoveredConnectionTarget(nodeId);
      }
    }
  }, [isConnecting, connectionStart, nodes, canBeConnectionTarget]);

  const handleConnectionNodeMouseLeave = useCallback(() => {
    setHoveredConnectionTarget(null);
  }, []);

  // Complete connection when dropping on a valid target
  const handleNodeMouseUp = useCallback((nodeId: string) => {
    if (isConnecting && connectionStart && hoveredConnectionTarget === nodeId) {
      const sourceNode = nodes.find(n => n.id === connectionStart.nodeId);
      const targetNode = nodes.find(n => n.id === nodeId);
      
      if (sourceNode && targetNode && isValidConnection(sourceNode, targetNode)) {
        // Check if connection already exists
        const existingConnection = connections.find(
          c => c.sourceNodeId === connectionStart.nodeId && c.targetNodeId === nodeId
        );
        
        if (!existingConnection) {
          const newConnection: CanvasConnection = {
            id: `conn-${Date.now()}`,
            sourceNodeId: connectionStart.nodeId,
            sourcePort: 0, // Single port for now
            targetNodeId: nodeId,
            targetPort: 0, // Single port for now
          };
          
          setConnections(prev => [...prev, newConnection]);
          addConsoleMessage('success', `Connected ${sourceNode.name} to ${targetNode.name}`);
        } else {
          addConsoleMessage('warning', 'Connection already exists');
        }
      }
    }
    
    // Reset connection state
    setIsConnecting(false);
    setConnectionStart(null);
    setTempConnection(null);
    setHoveredConnectionTarget(null);
  }, [isConnecting, connectionStart, hoveredConnectionTarget, nodes, connections, isValidConnection, setConnections, addConsoleMessage]);

  // Update temporary connection line during mouse move
  const updateTempConnection = useCallback((clientX: number, clientY: number, canvasRect: DOMRect) => {
    if (isConnecting && tempConnection) {
      setTempConnection(prev => prev ? {
        ...prev,
        endX: clientX - canvasRect.left,
        endY: clientY - canvasRect.top,
      } : null);
    }
  }, [isConnecting, tempConnection]);

  // Cancel connection
  const cancelConnection = useCallback(() => {
    setIsConnecting(false);
    setConnectionStart(null);
    setTempConnection(null);
    setHoveredConnectionTarget(null);
    addConsoleMessage('info', 'Connection cancelled');
  }, [addConsoleMessage]);

  // Delete connection
  const handleDeleteConnection = useCallback((connectionId: string) => {
    setConnections(prev => prev.filter(c => c.id !== connectionId));
    addConsoleMessage('info', 'Connection deleted');
  }, [setConnections, addConsoleMessage]);

  // Handle connection click
  const handleConnectionClick = useCallback((connectionId: string) => {
    setSelectedConnectionId(connectionId);
  }, []);

  return {
    // State
    isConnecting,
    connectionStart,
    tempConnection,
    hoveredConnectionTarget,
    selectedConnectionId,
    
    // Functions
    canBeConnectionSource,
    canBeConnectionTarget,
    isValidConnection,
    getConnectionAnchor,
    handleConnectionButtonMouseDown,
    handleConnectionNodeMouseEnter,
    handleConnectionNodeMouseLeave,
    handleNodeMouseUp,
    updateTempConnection,
    cancelConnection,
    handleDeleteConnection,
    handleConnectionClick,
    setSelectedConnectionId,
  };
}

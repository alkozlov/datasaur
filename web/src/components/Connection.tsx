import React, { useCallback } from 'react';
import { CanvasNode, CanvasConnection } from '../types/canvas';

interface ConnectionProps {
  connection: CanvasConnection;
  sourceNode: CanvasNode;
  targetNode: CanvasNode;
  isSelected: boolean;
  onConnectionClick: (connectionId: string) => void;
}

interface ConnectionsProps {
  connections: CanvasConnection[];
  nodes: CanvasNode[];
  selectedConnectionId: string | null;
  tempConnection: { startX: number; startY: number; endX: number; endY: number } | null;
  onConnectionClick: (connectionId: string) => void;
}

// Get the connection anchor points (right side of source, left side of target)
const getConnectionAnchor = (node: CanvasNode) => {
  return {
    x: node.x + 120, // Right side of block
    y: node.y + 16,  // Center height of block
  };
};

function Connection({ 
  connection, 
  sourceNode, 
  targetNode, 
  isSelected, 
  onConnectionClick 
}: ConnectionProps) {
  
  const handleClick = useCallback((e: React.MouseEvent) => {
    e.stopPropagation();
    onConnectionClick(connection.id);
  }, [onConnectionClick, connection.id]);

  // Use connection anchors (right side of source, left side of target)
  const startAnchor = getConnectionAnchor(sourceNode);
  const endAnchor = {
    x: targetNode.x, // Left side of target node
    y: targetNode.y + 16, // Center height
  };
  
  const startX = startAnchor.x;
  const startY = startAnchor.y;
  const endX = endAnchor.x;
  const endY = endAnchor.y;
  
  // Create a smooth Bezier curve
  const controlX1 = startX + Math.abs(endX - startX) * 0.5;
  const controlY1 = startY;
  const controlX2 = endX - Math.abs(endX - startX) * 0.5;
  const controlY2 = endY;
  
  const path = `M ${startX} ${startY} C ${controlX1} ${controlY1}, ${controlX2} ${controlY2}, ${endX} ${endY}`;
  
  return (
    <g key={connection.id}>
      <path
        d={path}
        stroke={isSelected ? "#1976d2" : "#666"}
        strokeWidth={isSelected ? "3" : "2"}
        fill="none"
        style={{ cursor: 'pointer' }}
        onClick={handleClick}
      />
      {/* Connection end marker */}
      <circle
        cx={endX}
        cy={endY}
        r="4"
        fill={isSelected ? "#1976d2" : "#666"}
      />
    </g>
  );
}

function Connections({ 
  connections, 
  nodes, 
  selectedConnectionId, 
  tempConnection, 
  onConnectionClick 
}: ConnectionsProps) {
  
  const renderConnection = useCallback((connection: CanvasConnection) => {
    const sourceNode = nodes.find(n => n.id === connection.sourceNodeId);
    const targetNode = nodes.find(n => n.id === connection.targetNodeId);
    
    if (!sourceNode || !targetNode) return null;
    
    return (
      <Connection
        key={connection.id}
        connection={connection}
        sourceNode={sourceNode}
        targetNode={targetNode}
        isSelected={selectedConnectionId === connection.id}
        onConnectionClick={onConnectionClick}
      />
    );
  }, [nodes, selectedConnectionId, onConnectionClick]);

  return (
    <svg
      style={{
        position: 'absolute',
        top: 0,
        left: 0,
        width: '100%',
        height: '100%',
        pointerEvents: 'none',
        zIndex: 1,
      }}
    >
      {/* Render existing connections */}
      <g style={{ pointerEvents: 'all' }}>
        {connections.map(renderConnection)}
      </g>
      
      {/* Render temporary connection line */}
      {tempConnection && (
        <path
          d={`M ${tempConnection.startX} ${tempConnection.startY} C ${tempConnection.startX + 50} ${tempConnection.startY}, ${tempConnection.endX - 50} ${tempConnection.endY}, ${tempConnection.endX} ${tempConnection.endY}`}
          stroke="#1976d2"
          strokeWidth="2"
          strokeDasharray="5,5"
          fill="none"
        />
      )}
    </svg>
  );
}

export default Connections;
export { Connection };

import React, { useCallback } from 'react';
import { Box, Typography } from '@mui/material';
import { CanvasNode } from '../types/canvas';

interface BlockProps {
  node: CanvasNode;
  isSelected: boolean;
  isHovered: boolean;
  isConnectionTarget: boolean;
  isConnecting: boolean;
  showConnectionButton: boolean;
  canBeConnectionSource: boolean;
  onMouseDown: (e: React.MouseEvent, nodeId: string) => void;
  onMouseEnter: (nodeId: string) => void;
  onMouseLeave: (nodeId: string) => void;
  onMouseUp: (nodeId: string) => void;
  onContextMenu: (e: React.MouseEvent, nodeId: string) => void;
  onConnectionButtonMouseDown: (e: React.MouseEvent, nodeId: string) => void;
  onConnectionButtonMouseEnter: (nodeId: string) => void;
  onConnectionButtonMouseLeave: (nodeId: string) => void;
  isDragging: boolean;
}

const getNodeColor = (category: string) => {
  switch (category) {
    case 'input': return '#8BBEE8';
    case 'output': return '#F9C74F';
    case 'math': return '#90E0EF';
    default: return '#B3B3B3';
  }
};

function Block({
  node,
  isSelected,
  isHovered,
  isConnectionTarget,
  isConnecting,
  showConnectionButton,
  canBeConnectionSource,
  onMouseDown,
  onMouseEnter,
  onMouseLeave,
  onMouseUp,
  onContextMenu,
  onConnectionButtonMouseDown,
  onConnectionButtonMouseEnter,
  onConnectionButtonMouseLeave,
  isDragging,
}: BlockProps) {
  
  const handleMouseDown = useCallback((e: React.MouseEvent) => {
    onMouseDown(e, node.id);
  }, [onMouseDown, node.id]);

  const handleMouseEnter = useCallback(() => {
    onMouseEnter(node.id);
  }, [onMouseEnter, node.id]);

  const handleMouseLeave = useCallback(() => {
    onMouseLeave(node.id);
  }, [onMouseLeave, node.id]);

  const handleMouseUp = useCallback(() => {
    onMouseUp(node.id);
  }, [onMouseUp, node.id]);

  const handleContextMenu = useCallback((e: React.MouseEvent) => {
    onContextMenu(e, node.id);
  }, [onContextMenu, node.id]);

  const handleConnectionButtonMouseDown = useCallback((e: React.MouseEvent) => {
    onConnectionButtonMouseDown(e, node.id);
  }, [onConnectionButtonMouseDown, node.id]);

  const handleConnectionButtonMouseEnter = useCallback(() => {
    onConnectionButtonMouseEnter(node.id);
  }, [onConnectionButtonMouseEnter, node.id]);

  const handleConnectionButtonMouseLeave = useCallback(() => {
    onConnectionButtonMouseLeave(node.id);
  }, [onConnectionButtonMouseLeave, node.id]);

  return (
    <Box key={node.id}>
      {/* Main node */}
      <Box
        sx={{
          position: 'absolute',
          left: node.x,
          top: node.y,
          width: '120px',
          height: '32px',
          backgroundColor: getNodeColor(node.category),
          border: isSelected 
            ? '2px solid #1976d2' 
            : isConnectionTarget 
              ? '2px solid #4caf50'
              : '1px solid #999',
          borderRadius: '4px',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          cursor: isDragging ? 'grabbing' : 'grab',
          userSelect: 'none',
          zIndex: 2,
          '&:hover': {
            borderColor: '#666',
          },
          '&:active': {
            cursor: 'grabbing',
          },
        }}
        onMouseDown={handleMouseDown}
        onMouseEnter={handleMouseEnter}
        onMouseLeave={handleMouseLeave}
        onMouseUp={handleMouseUp}
        onContextMenu={handleContextMenu}
      >
        <Typography
          variant="caption"
          sx={{
            color: '#fff',
            fontSize: '11px',
            fontWeight: 500,
            textAlign: 'center',
          }}
        >
          {node.name}
        </Typography>
      </Box>
      
      {/* Connection button - only show for source nodes on hover when not connecting */}
      {canBeConnectionSource && !isConnecting && showConnectionButton && (
        <Box
          sx={{
            position: 'absolute',
            left: node.x + 60 - 12, // Center of node minus half button width
            top: node.y + 16 - 12,  // Center of node minus half button height
            width: '24px',
            height: '24px',
            backgroundColor: '#1976d2',
            border: '2px solid #fff',
            borderRadius: '50%',
            cursor: 'crosshair',
            zIndex: 4,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            boxShadow: '0 2px 4px rgba(0,0,0,0.2)',
            transition: 'transform 0.1s ease',
            '&:hover': {
              backgroundColor: '#1565c0',
              transform: 'scale(1.1)',
            },
          }}
          onMouseDown={handleConnectionButtonMouseDown}
          onMouseEnter={handleConnectionButtonMouseEnter}
          onMouseLeave={handleConnectionButtonMouseLeave}
        >
          <Typography 
            sx={{ 
              color: '#fff', 
              fontSize: '14px', 
              fontWeight: 'bold',
              lineHeight: 1,
              pointerEvents: 'none', // Prevent text from interfering with mouse events
            }}
          >
            ⚡
          </Typography>
        </Box>
      )}
    </Box>
  );
}

export default Block;

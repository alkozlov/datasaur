import { useState, useCallback } from 'react';
import { CanvasNode } from '../types/canvas';

interface UseBlockHoverLogicProps {
  nodes: CanvasNode[];
  canBeConnectionSource: (node: CanvasNode) => boolean;
  isConnecting: boolean;
  handleConnectionNodeMouseEnter: (nodeId: string) => void;
  handleConnectionNodeMouseLeave: () => void;
}

export function useBlockHoverLogic({
  nodes,
  canBeConnectionSource,
  isConnecting,
  handleConnectionNodeMouseEnter,
  handleConnectionNodeMouseLeave,
}: UseBlockHoverLogicProps) {
  const [hoveredNodeId, setHoveredNodeId] = useState<string | null>(null);
  const [showConnectionButton, setShowConnectionButton] = useState<string | null>(null);

  // Handle node hover for connection target highlighting
  const handleNodeMouseEnter = useCallback((nodeId: string) => {
    setHoveredNodeId(nodeId);
    
    // Show connection button for source nodes
    const node = nodes.find(n => n.id === nodeId);
    if (node && canBeConnectionSource(node) && !isConnecting) {
      setShowConnectionButton(nodeId);
    }
    
    // Handle connection target highlighting
    handleConnectionNodeMouseEnter(nodeId);
  }, [nodes, canBeConnectionSource, isConnecting, handleConnectionNodeMouseEnter]);

  const handleNodeMouseLeave = useCallback((nodeId: string) => {
    // Use a timeout to prevent flickering when moving between node and connection button
    setTimeout(() => {
      setHoveredNodeId(prev => prev === nodeId ? null : prev);
      setShowConnectionButton(prev => prev === nodeId ? null : prev);
    }, 100);
    
    handleConnectionNodeMouseLeave();
  }, [handleConnectionNodeMouseLeave]);

  // Handle connection button hover to maintain the hover state
  const handleConnectionButtonMouseEnter = useCallback((nodeId: string) => {
    setHoveredNodeId(nodeId);
    setShowConnectionButton(nodeId);
  }, []);

  const handleConnectionButtonMouseLeave = useCallback((nodeId: string) => {
    // Small delay to allow moving back to the node
    setTimeout(() => {
      setShowConnectionButton(prev => prev === nodeId ? null : prev);
    }, 100);
  }, []);

  // Clear all hover states
  const clearHoverStates = useCallback(() => {
    setHoveredNodeId(null);
    setShowConnectionButton(null);
  }, []);

  return {
    hoveredNodeId,
    showConnectionButton,
    handleNodeMouseEnter,
    handleNodeMouseLeave,
    handleConnectionButtonMouseEnter,
    handleConnectionButtonMouseLeave,
    clearHoverStates,
  };
}

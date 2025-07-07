// Types for the Block-Flow canvas functionality

export interface CanvasNode {
  id: string;
  type: string;
  name: string;
  description: string;
  category: string;
  x: number;
  y: number;
  selected: boolean;
  inputs?: number;  // Number of input ports
  outputs?: number; // Number of output ports
}

export interface CanvasConnection {
  id: string;
  sourceNodeId: string;
  sourcePort: number;
  targetNodeId: string;
  targetPort: number;
}

export interface ConnectionPoint {
  nodeId: string;
  port: number;
  isOutput: boolean;
  x: number;
  y: number;
}

export interface ConsoleMessage {
  id: string;
  timestamp: Date;
  type: 'info' | 'success' | 'warning' | 'error';
  message: string;
}

export interface CanvasState {
  nodes: CanvasNode[];
  connections: CanvasConnection[];
  selectedNodeId: string | null;
  isDragging: boolean;
  dragOffset: { x: number; y: number };
  isConnecting: boolean;
  connectionStart: ConnectionPoint | null;
  tempConnection: { startX: number; startY: number; endX: number; endY: number } | null;
}

export interface Flow {
  id: string;
  name: string;
  description?: string;
  nodes: FlowNode[];
  connections: FlowConnection[];
  created_at: string;
  updated_at: string;
  version: string;
  active: boolean;
}

export interface FlowNode {
  id: string;
  type: string;
  name: string;
  x: number;
  y: number;
  properties: Record<string, any>;
  inputs: number;
  outputs: number;
}

export interface FlowConnection {
  id: string;
  source: string;
  source_port: number;
  target: string;
  target_port: number;
}

// API utility functions for Block-Flow backend communication

const API_BASE_URL = '/api/v1';

export interface NodeType {
  type: string;
  name: string;
  description: string;
  category: string;
  version: string;
  author?: string;
  icon?: string;
  color?: string;
  inputs?: number;
  outputs?: number;
  block_group?: string;
}

export interface FlowData {
  id: string;
  name: string;
  description?: string;
  nodes: FlowNodeData[];
  connections: FlowConnectionData[];
}

export interface FlowNodeData {
  id: string;
  type: string;
  name: string;
  x: number;
  y: number;
  properties: Record<string, any>;
  inputs: number;
  outputs: number;
}

export interface FlowConnectionData {
  id: string;
  source: string;
  source_port: number;
  target: string;
  target_port: number;
}

export async function fetchAvailableNodes(): Promise<NodeType[]> {
  try {
    const response = await fetch(`${API_BASE_URL}/blocks`);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    const nodes = await response.json();
    return nodes;
  } catch (error) {
    console.error('Error fetching nodes:', error);
    // Return mock data if backend is not available
    return [
      { type: 'inject', name: 'Inject', description: 'Input injection node', category: 'input', version: '1.0.0', inputs: 0, outputs: 1 },
      { type: 'debug', name: 'Debug', description: 'Debug output node', category: 'output', version: '1.0.0', inputs: 1, outputs: 0 },
      { type: 'add', name: 'Add', description: 'Mathematical addition', category: 'math', version: '1.0.0', inputs: 2, outputs: 1 },
      { type: 'subtract', name: 'Subtract', description: 'Mathematical subtraction', category: 'math', version: '1.0.0', inputs: 2, outputs: 1 },
      { type: 'multiply', name: 'Multiply', description: 'Mathematical multiplication', category: 'math', version: '1.0.0', inputs: 2, outputs: 1 },
      { type: 'divide', name: 'Divide', description: 'Mathematical division', category: 'math', version: '1.0.0', inputs: 2, outputs: 1 },
    ];
  }
}

export async function saveFlow(flowData: FlowData): Promise<string> {
  const response = await fetch(`${API_BASE_URL}/flows`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(flowData),
  });
  
  if (!response.ok) {
    throw new Error(`Failed to save flow: ${response.status}`);
  }
  
  const result = await response.json();
  return result.id;
}

export async function runFlow(flowId: string): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/flows/${flowId}/run`, {
    method: 'POST',
  });
  
  if (!response.ok) {
    throw new Error(`Failed to run flow: ${response.status}`);
  }
}

export async function stopFlow(flowId: string): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/flows/${flowId}/stop`, {
    method: 'POST',
  });
  
  if (!response.ok) {
    throw new Error(`Failed to stop flow: ${response.status}`);
  }
}

export async function getFlowStatus(flowId: string): Promise<{ running: boolean }> {
  const response = await fetch(`${API_BASE_URL}/flows/${flowId}/status`);
  
  if (!response.ok) {
    throw new Error(`Failed to get flow status: ${response.status}`);
  }
  
  return response.json();
}

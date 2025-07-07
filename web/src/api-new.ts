// API utility functions for Block-Flow backend communication

const API_BASE_URL = 'http://localhost:8080/api';

export interface NodeType {
  id: string;
  name: string;
  description: string;
  category: string;
  inputs: number;
  outputs: number;
}

export async function fetchAvailableNodes(): Promise<NodeType[]> {
  try {
    const response = await fetch(`${API_BASE_URL}/nodes`);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    const nodes = await response.json();
    return nodes;
  } catch (error) {
    console.error('Error fetching nodes:', error);
    // Return mock data if backend is not available
    return [
      { id: 'input', name: 'Input', description: 'Input node', category: 'io', inputs: 0, outputs: 1 },
      { id: 'output', name: 'Output', description: 'Output node', category: 'io', inputs: 1, outputs: 0 },
      { id: 'math-add', name: 'Add', description: 'Mathematical addition', category: 'math', inputs: 2, outputs: 1 },
      { id: 'math-sub', name: 'Subtract', description: 'Mathematical subtraction', category: 'math', inputs: 2, outputs: 1 },
      { id: 'logic-and', name: 'AND', description: 'Logical AND operation', category: 'logic', inputs: 2, outputs: 1 },
      { id: 'logic-or', name: 'OR', description: 'Logical OR operation', category: 'logic', inputs: 2, outputs: 1 },
    ];
  }
}

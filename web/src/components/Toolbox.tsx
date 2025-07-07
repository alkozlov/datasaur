import React, { useState, useEffect } from 'react';
import { 
  Box, 
  Typography, 
  List, 
  ListItem, 
  ListItemText, 
  CircularProgress,
  Alert
} from '@mui/material';
import { fetchAvailableNodes, NodeType } from '../api';

function Toolbox() {
  const [nodes, setNodes] = useState<NodeType[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadNodes = async () => {
      try {
        setLoading(true);
        const availableNodes = await fetchAvailableNodes();
        setNodes(availableNodes);
        setError(null);
      } catch (err) {
        setError('Failed to load nodes');
        console.error('Error loading nodes:', err);
      } finally {
        setLoading(false);
      }
    };

    loadNodes();
  }, []);

  const getChipColor = (category: string) => {
    switch (category) {
      case 'io': return 'primary';
      case 'math': return 'secondary';
      case 'logic': return 'success';
      default: return 'default';
    }
  };

  const getNodeColor = (category: string) => {
    switch (category) {
      case 'input': return '#8BBEE8'; // Light blue for input nodes
      case 'output': return '#F9C74F'; // Yellow for output nodes  
      case 'math': return '#90E0EF'; // Light cyan for math operations
      default: return '#B3B3B3'; // Gray for unknown categories
    }
  };

  const getNodeHoverColor = (category: string) => {
    switch (category) {
      case 'input': return '#7BB3E0';
      case 'output': return '#F1BF41';
      case 'math': return '#82D8E8';
      default: return '#A5A5A5';
    }
  };

  const handleDragStart = (e: React.DragEvent, node: NodeType) => {
    e.dataTransfer.setData('application/json', JSON.stringify(node));
    e.dataTransfer.effectAllowed = 'copy';
  };

  return (
    <Box
      sx={{
        width: '20%',
        height: '85vh',
        backgroundColor: '#f8f8f8',
        border: '1px solid #ccc',
        padding: 1,
        boxSizing: 'border-box',
        overflow: 'auto',
      }}
    >
      <Typography variant="subtitle2" gutterBottom sx={{ fontWeight: 600, fontSize: '12px', mb: 1 }}>
        Nodes
      </Typography>
      
      {loading && (
        <Box display="flex" justifyContent="center" mt={2}>
          <CircularProgress size={24} />
        </Box>
      )}
      
      {error && (
        <Alert severity="error" sx={{ mt: 1, mb: 1 }}>
          {error}
        </Alert>
      )}
      
      {!loading && !error && (
        <List dense sx={{ p: 0 }}>
          {nodes.map((node) => (
            <ListItem
              key={node.type}
              draggable
              onDragStart={(e) => handleDragStart(e, node)}
              sx={{
                border: '1px solid #999',
                borderRadius: '4px',
                mb: 0.5,
                backgroundColor: getNodeColor(node.category),
                cursor: 'pointer',
                minHeight: '32px',
                maxWidth: '120px',
                width: '120px',
                padding: '4px 8px',
                '&:hover': {
                  backgroundColor: getNodeHoverColor(node.category),
                  borderColor: '#666',
                },
                transition: 'all 0.2s ease',
              }}
            >
              <ListItemText
                primary={
                  <Typography 
                    variant="caption" 
                    fontWeight="500"
                    sx={{ 
                      color: '#fff',
                      fontSize: '11px',
                      lineHeight: '1.2',
                    }}
                  >
                    {node.name}
                  </Typography>
                }
                secondary={null}
                sx={{ 
                  margin: 0,
                  '& .MuiListItemText-primary': {
                    margin: 0,
                  }
                }}
              />
            </ListItem>
          ))}
        </List>
      )}
    </Box>
  );
}

export default Toolbox;

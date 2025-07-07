import React from 'react';
import { Box, Typography, List, ListItem, ListItemText, Chip } from '@mui/material';
import { ConsoleMessage } from '../types/canvas';

interface ConsoleProps {
  messages: ConsoleMessage[];
}

function Console({ messages }: ConsoleProps) {
  const getChipColor = (type: ConsoleMessage['type']) => {
    switch (type) {
      case 'error': return 'error';
      case 'warning': return 'warning';
      case 'success': return 'success';
      case 'info': return 'info';
      default: return 'default';
    }
  };

  return (
    <Box
      sx={{
        width: '100%',
        height: '15vh',
        backgroundColor: '#2e2e2e',
        color: '#ffffff',
        border: '1px solid #ddd',
        padding: 2,
        boxSizing: 'border-box',
        overflow: 'auto',
      }}
    >
      <Typography variant="h6" gutterBottom sx={{ color: '#ffffff' }}>
        Console
      </Typography>
      <List dense sx={{ pt: 0 }}>
        {messages.slice(-10).map((message) => (
          <ListItem key={message.id} sx={{ py: 0.5, px: 0 }}>
            <ListItemText
              primary={
                <Box display="flex" alignItems="center" gap={1}>
                  <Typography variant="caption" sx={{ color: '#888', minWidth: '60px' }}>
                    {message.timestamp.toLocaleTimeString()}
                  </Typography>
                  <Chip 
                    label={message.type.toUpperCase()} 
                    size="small" 
                    color={getChipColor(message.type) as any}
                    sx={{ minWidth: '60px' }}
                  />
                  <Typography variant="body2" sx={{ color: '#fff' }}>
                    {message.message}
                  </Typography>
                </Box>
              }
            />
          </ListItem>
        ))}
      </List>
    </Box>
  );
}

export default Console;

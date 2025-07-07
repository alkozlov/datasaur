import React, { useState, useCallback } from 'react';
import { Box } from '@mui/material';
import Toolbox from './components/Toolbox';
import Canvas from './components/Canvas';
import Console from './components/Console';
import { ConsoleMessage } from './types/canvas';

function App() {
  const [consoleMessages, setConsoleMessages] = useState<ConsoleMessage[]>([]);

  const handleConsoleMessage = useCallback((message: ConsoleMessage) => {
    setConsoleMessages(prev => [...prev, message]);
  }, []);

  return (
    <Box
      sx={{
        display: 'flex',
        flexDirection: 'column',
        height: '100vh',
        width: '100vw',
        margin: 0,
        padding: 0,
      }}
    >
      {/* Top section with Toolbox and Canvas */}
      <Box
        sx={{
          display: 'flex',
          flexDirection: 'row',
          height: '85vh',
          width: '100%',
        }}
      >
        <Toolbox />
        <Canvas onConsoleMessage={handleConsoleMessage} />
      </Box>
      
      {/* Bottom section with Console */}
      <Console messages={consoleMessages} />
    </Box>
  );
}

export default App;

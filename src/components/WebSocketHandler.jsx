export const connectWebSocket = (setGameState) => {
  const socket = new WebSocket("ws://localhost:8080/ws");

  socket.onopen = () => console.log("Connected to server");

  socket.onmessage = (event) => {
      try {
          const gameState = JSON.parse(event.data);
          
          console.log("Received gameState from WebSocket:", gameState); // Debugging

          // Ensure `paddles` is an array
          const normalizedGameState = {
              ...gameState,
              paddles: Array.isArray(gameState.paddles) 
                  ? gameState.paddles 
                  : Object.values(gameState.paddles || {})
          };

          setGameState(normalizedGameState);
      } catch (error) {
          console.error("Error parsing gameState:", error, event.data);
      }
  };

  socket.onerror = (error) => console.error("WebSocket error:", error);

  return socket;
};

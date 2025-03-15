export const connectWebSocket = (setGameState) => {
    const socket = new WebSocket("ws://localhost:8080/ws");
  
    socket.onopen = () => console.log("Connected to server");
  
    socket.onmessage = (event) => {
      const gameState = JSON.parse(event.data);
      setGameState(gameState);
    };
  
    socket.onerror = (error) => console.error("WebSocket error:", error);
  
    return socket;
  };
  
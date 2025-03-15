import React, { useState, useEffect } from "react";
import GameCanvas from "./components/GameCanvas";
import { connectWebSocket } from "./components/WebSocketHandler";

const App = () => {
  const [gameState, setGameState] = useState(null);
  const [ws, setWs] = useState(null);

  useEffect(() => {
    const socket = connectWebSocket(setGameState);
    setWs(socket);

    return () => {
      socket.close();
    };
  }, []);

  return (
    <div>
      <h1>Multiplayer Pong</h1>
      {gameState ? <GameCanvas gameState={gameState} /> : <p>Connecting...</p>}
    </div>
  );
};

export default App;

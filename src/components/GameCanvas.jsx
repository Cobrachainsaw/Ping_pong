import React, { useEffect, useRef } from "react";

const GameCanvas = ({ gameState }) => {
  const canvasRef = useRef(null);

  useEffect(() => {
    const canvas = canvasRef.current;
    const ctx = canvas.getContext("2d");

    // Debug: Log gameState structure
    console.log("Received gameState:", gameState);

    // Ensure gameState exists and has correct structure
    if (!gameState || typeof gameState !== "object" || !gameState.ball || !gameState.paddles) {
      console.warn("Invalid gameState received:", gameState);
      return;
    }

    // Ensure paddles is an array
    const paddles = Array.isArray(gameState.paddles) ? gameState.paddles : Object.values(gameState.paddles);

    const draw = () => {
      ctx.clearRect(0, 0, canvas.width, canvas.height);

      // Draw ball
      ctx.fillStyle = "red";
      ctx.beginPath();
      ctx.arc(gameState.ball?.x || 300, gameState.ball?.y || 300, 10, 0, Math.PI * 2);
      ctx.fill();

      // Draw paddles
      paddles.forEach((paddle) => {
        if (!paddle) return;
        ctx.fillStyle = "blue";
        ctx.fillRect(paddle.x || 0, paddle.y || 0, paddle.width || 10, paddle.height || 50);
      });
    };

    draw();
  }, [gameState]);

  return <canvas ref={canvasRef} width={600} height={600} />;
};

export default GameCanvas;

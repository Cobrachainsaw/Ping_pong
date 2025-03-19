package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type GameState struct {
	Ball    Ball     `json:"ball"`
	Paddles []Paddle `json:"paddles"` // Now an array of paddles
}

type Ball struct {
	X, Y   int
	VX, VY int
}

type Paddle struct {
	X      int `json:"x"`
	Y      int `json:"y"` // ✅ Correct: "y" instead of "x"
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Client struct {
	conn *websocket.Conn
}

var (
	clients      = make(map[*Client]bool)
	clientsMutex sync.Mutex
	broadcast    = make(chan GameState)
	upgrader     = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	gameState = GameState{
		Ball: Ball{X: 300, Y: 300, VX: 3, VY: 3},
		Paddles: []Paddle{
			{X: 50, Y: 250, Width: 10, Height: 50},  // Left Paddle
			{X: 540, Y: 250, Width: 10, Height: 50}, // Right Paddle
		},
	}
	mutex sync.Mutex
)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return
	}
	client := &Client{conn: conn}

	clientsMutex.Lock()
	clients[client] = true
	clientsMutex.Unlock()

	defer func() {
		clientsMutex.Lock()
		delete(clients, client)
		clientsMutex.Unlock()
		conn.Close()
	}()

	for {
		var msg map[string]string
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}

		mutex.Lock()
		switch msg["key"] {
		case "w":
			gameState.Paddles[0].Y -= 10 // Move left paddle up
		case "s":
			gameState.Paddles[0].Y += 10 // Move left paddle down
		case "ArrowUp":
			gameState.Paddles[1].Y -= 10 // Move right paddle up
		case "ArrowDown":
			gameState.Paddles[1].Y += 10 // Move right paddle down
		}
		mutex.Unlock()
	}
}

func gameLoop() {
	ticker := time.NewTicker(16 * time.Millisecond) // ~60 FPS
	for range ticker.C {
		mutex.Lock()
		gameState.Ball.X += gameState.Ball.VX
		gameState.Ball.Y += gameState.Ball.VY

		// Ball collision with walls
		if gameState.Ball.X <= 0 || gameState.Ball.X >= 580 {
			gameState.Ball.VX *= -1
		}
		if gameState.Ball.Y <= 0 || gameState.Ball.Y >= 580 {
			gameState.Ball.VY *= -1
		}
		mutex.Unlock()

		broadcast <- gameState
	}
}

func handleBroadcast() {
	for {
		state := <-broadcast
		jsonState, _ := json.Marshal(state)
		clientsMutex.Lock()
		for client := range clients {
			client.conn.WriteMessage(websocket.TextMessage, jsonState)
		}
		clientsMutex.Unlock()
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)
	go handleBroadcast()
	go gameLoop()

	fmt.Println("Game server running on :8080")
	http.ListenAndServe(":8080", nil)
}

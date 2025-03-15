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
	Ball    Ball    `json:"ball"`
	Paddles Paddles `json:"paddles"`
}

type Ball struct {
	X, Y   int
	VX, VY int
}

type Paddles struct {
	Top, Bottom, Left, Right int
}

type Client struct {
	conn *websocket.Conn
}

var (
	clients   = make(map[*Client]bool)
	broadcast = make(chan GameState)
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	gameState = GameState{
		Ball:    Ball{X: 300, Y: 300, VX: 3, VY: 3},
		Paddles: Paddles{Top: 250, Bottom: 250, Left: 250, Right: 250},
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
	clients[client] = true

	defer func() {
		delete(clients, client)
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
		// Handle paddle movement
		switch msg["key"] {
		case "a":
			gameState.Paddles.Top -= 10
		case "d":
			gameState.Paddles.Top += 10
		case "j":
			gameState.Paddles.Bottom -= 10
		case "l":
			gameState.Paddles.Bottom += 10
		case "w":
			gameState.Paddles.Left -= 10
		case "s":
			gameState.Paddles.Left += 10
		case "i":
			gameState.Paddles.Right -= 10
		case "k":
			gameState.Paddles.Right += 10
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
		for client := range clients {
			client.conn.WriteMessage(websocket.TextMessage, jsonState)
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)
	go handleBroadcast()
	go gameLoop()

	fmt.Println("Game server running on :8080")
	http.ListenAndServe(":8080", nil)
}

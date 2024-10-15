// Package websocket provides functionality for handling WebSocket connections.
package websocket

import (
	"fmt"
	"log"
	"net/http"

	"quiz-app/internal/quiz"

	"github.com/gorilla/websocket"
)

var ListConnections []*websocket.Conn

// Handler represents a WebSocket handler.
type Handler struct {
	quizSession *quiz.QuizSession
}

// NewHandler returns a new WebSocket handler.
func NewHandler(quizSession *quiz.QuizSession) *Handler {
	return &Handler{quizSession: quizSession}
}

// HandleConnection handles a new WebSocket connection.
func (h *Handler) HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	ListConnections = append(ListConnections, conn)

	// Handle incoming messages from the client
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// Handle incoming message
		switch messageType {
		case websocket.BinaryMessage:
			// Handle binary message
			fmt.Println("Received binary message:", message)
			h.handleBinaryMessage(conn, message)
		default:
			// Handle unknown message type
			fmt.Println("Received unknown message type:", messageType)
			h.handleUnknownMessage(conn, messageType)
		}
	}
}

// handleBinaryMessage handles a binary message from the client.
func (h *Handler) handleBinaryMessage(conn *websocket.Conn, message []byte) {
	// Handle binary message
	log.Println("Received binary message:", message)
}

// handleUnknownMessage handles an unknown message type from the client.
func (h *Handler) handleUnknownMessage(conn *websocket.Conn, messageType int) {
	// Handle unknown message type
	log.Println("Received unknown message type:", messageType)
}

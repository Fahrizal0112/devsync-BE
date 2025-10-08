package websocket

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // Allow all origins in development
    },
}

type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
}

type Client struct {
    hub      *Hub
    conn     *websocket.Conn
    send     chan []byte
    userID   uint
    projectID uint
}

type Message struct {
    Type      string      `json:"type"`
    ProjectID uint        `json:"project_id"`
    UserID    uint        `json:"user_id"`
    Data      interface{} `json:"data"`
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan []byte),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true
            log.Printf("Client connected: User %d, Project %d", client.userID, client.projectID)

        case client := <-h.unregister:
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
                log.Printf("Client disconnected: User %d, Project %d", client.userID, client.projectID)
            }

        case message := <-h.broadcast:
            var msg Message
            if err := json.Unmarshal(message, &msg); err != nil {
                continue
            }

            for client := range h.clients {
                if client.projectID == msg.ProjectID {
                    select {
                    case client.send <- message:
                    default:
                        close(client.send)
                        delete(h.clients, client)
                    }
                }
            }
        }
    }
}

// Broadcast sends a message to all connected clients
func (h *Hub) Broadcast(message []byte) {
    h.broadcast <- message
}

func (h *Hub) HandleWebSocket(c *gin.Context) {
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Println("WebSocket upgrade error:", err)
        return
    }

    // Get user and project from query params or JWT
    userID := uint(1)    // TODO: Extract from JWT
    projectID := uint(1) // TODO: Extract from query params

    client := &Client{
        hub:       h,
        conn:      conn,
        send:      make(chan []byte, 256),
        userID:    userID,
        projectID: projectID,
    }

    client.hub.register <- client

    go client.writePump()
    go client.readPump()
}

func (c *Client) readPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()

    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            break
        }
        c.hub.broadcast <- message
    }
}

func (c *Client) writePump() {
    defer c.conn.Close()

    for message := range c.send {
        if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
            return
        }
    }
}
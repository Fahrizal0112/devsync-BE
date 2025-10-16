package websocket

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"

    "devsync-be/internal/auth"
    "devsync-be/internal/config"

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
    config     *config.Config
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

func NewHub(cfg *config.Config) *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan []byte),
        register:   make(chan *Client),
        unregister: make(chan *Client),
        config:     cfg,
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

    // Extract and validate JWT token from query params
    tokenString := c.Query("token")
    if tokenString == "" {
        log.Println("WebSocket: token query parameter required")
        conn.Close()
        return
    }

    // Validate JWT token and extract userID
    claims, err := auth.ValidateToken(tokenString, h.config.JWTSecret)
    if err != nil {
        log.Println("WebSocket: Invalid token:", err)
        conn.Close()
        return
    }

    userID := claims.UserID

    // Extract projectID from query params
    projectIDStr := c.Query("project_id")
    if projectIDStr == "" {
        log.Println("WebSocket: project_id query parameter required")
        conn.Close()
        return
    }
    
    var projectID uint
    if _, err := fmt.Sscanf(projectIDStr, "%d", &projectID); err != nil {
        log.Println("WebSocket: Invalid project_id format")
        conn.Close()
        return
    }

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
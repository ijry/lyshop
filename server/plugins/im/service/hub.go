package service

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Frame is the WebSocket message envelope.
type Frame struct {
	Type      string         `json:"type"`       // msg|ack|typing|assign|close|ping
	SessionID uint64         `json:"session_id"`
	Payload   map[string]any `json:"payload"`
}

// Client represents one WebSocket connection.
type Client struct {
	ID   string // "user_{id}" or "staff_{id}"
	Conn *websocket.Conn
	Send chan []byte
}

// Hub manages all active WebSocket connections.
type Hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan *delivery
	mu         sync.RWMutex
}

type delivery struct {
	targetID string
	data     []byte
}

var GlobalHub = NewHub()

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client, 16),
		unregister: make(chan *Client, 16),
		broadcast:  make(chan *delivery, 256),
	}
}

// Run starts the hub event loop. Call once in a goroutine.
func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			h.clients[c.ID] = c
			h.mu.Unlock()

		case c := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[c.ID]; ok {
				delete(h.clients, c.ID)
				close(c.Send)
			}
			h.mu.Unlock()

		case d := <-h.broadcast:
			h.mu.RLock()
			target, ok := h.clients[d.targetID]
			h.mu.RUnlock()
			if ok {
				select {
				case target.Send <- d.data:
				default:
					// Slow client — drop and unregister
					h.mu.Lock()
					delete(h.clients, d.targetID)
					close(target.Send)
					h.mu.Unlock()
				}
			}
		}
	}
}

// Send delivers bytes to a specific client by ID.
func (h *Hub) Send(targetID string, data []byte) {
	h.broadcast <- &delivery{targetID: targetID, data: data}
}

// Register adds a client to the hub.
func (h *Hub) Register(c *Client) { h.register <- c }

// Unregister removes a client from the hub.
func (h *Hub) Unregister(c *Client) { h.unregister <- c }

// Online returns true if clientID is connected.
func (h *Hub) Online(id string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[id]
	return ok
}

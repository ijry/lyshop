package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ijry/lyshop/core/cache"
)

// Frame is the WebSocket message envelope.
type Frame struct {
	Type      string         `json:"type"` // msg|ack|typing|assign|close|ping
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
	clients      map[string]*Client
	register     chan *Client
	unregister   chan *Client
	broadcast    chan *delivery
	nodeID       string
	redisEnabled bool
	mu           sync.RWMutex
}

type delivery struct {
	targetID string
	data     []byte
}

const imWSChannel = "lyshop:im:ws"

type hubEnvelope struct {
	NodeID    string          `json:"node_id"`
	TargetID  string          `json:"target_id"`
	Data      json.RawMessage `json:"data"`
	CreatedAt int64           `json:"created_at"`
}

var GlobalHub = NewHub()

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client, 16),
		unregister: make(chan *Client, 16),
		broadcast:  make(chan *delivery, 256),
		nodeID:     newNodeID(),
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
	h.sendLocal(targetID, data)
	h.publishRemote(context.Background(), targetID, data)
}

func (h *Hub) sendLocal(targetID string, data []byte) {
	h.broadcast <- &delivery{targetID: targetID, data: data}
}

func (h *Hub) InitRedisBus(ctx context.Context) {
	if cache.Client == nil {
		return
	}
	h.redisEnabled = true
	pubsub := cache.Client.Subscribe(ctx, imWSChannel)
	go func() {
		defer pubsub.Close()
		ch := pubsub.Channel()
		for msg := range ch {
			if !h.shouldDeliverRemote([]byte(msg.Payload)) {
				continue
			}
			var env hubEnvelope
			if err := json.Unmarshal([]byte(msg.Payload), &env); err != nil {
				continue
			}
			h.sendLocal(env.TargetID, []byte(env.Data))
		}
	}()
}

func (h *Hub) publishRemote(ctx context.Context, targetID string, data []byte) {
	if !h.redisEnabled || cache.Client == nil {
		return
	}
	env := hubEnvelope{
		NodeID:    h.nodeID,
		TargetID:  targetID,
		Data:      json.RawMessage(data),
		CreatedAt: time.Now().UnixMilli(),
	}
	raw, err := json.Marshal(env)
	if err != nil {
		return
	}
	if err := cache.Client.Publish(ctx, imWSChannel, raw).Err(); err != nil {
		slog.Warn("im websocket redis publish failed", "error", err)
	}
}

func (h *Hub) shouldDeliverRemote(raw []byte) bool {
	var env hubEnvelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return false
	}
	return env.NodeID != "" && env.NodeID != h.nodeID && env.TargetID != "" && len(env.Data) > 0
}

func newNodeID() string {
	var buf [4]byte
	_, _ = rand.Read(buf[:])
	host, _ := os.Hostname()
	if host == "" {
		host = "node"
	}
	return host + "-" + hex.EncodeToString(buf[:])
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

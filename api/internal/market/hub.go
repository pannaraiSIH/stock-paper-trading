package market

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type HubManager interface {
	Connect(ctx *gin.Context)
	Broadcast(price PriceEvent)
	SubscriptionEvent() <-chan SubscribeEvent
}

type Client struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

type Hub struct {
	clients           map[*Client]map[string]bool
	subscriptions     map[string]map[*Client]bool
	subscriptionEvent chan SubscribeEvent
	mu                sync.RWMutex
}

func NewHub() HubManager {
	return &Hub{
		clients:           make(map[*Client]map[string]bool),
		subscriptions:     make(map[string]map[*Client]bool),
		subscriptionEvent: make(chan SubscribeEvent, 100),
	}
}

func (c *Client) WriteJSON(v any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.conn.WriteJSON(v)
}

func (c *Client) WritePing() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.conn.WriteMessage(websocket.PingMessage, nil)
}

func (h *Hub) Connect(c *gin.Context) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("failed to upgrade connection: %v", err)
		return
	}

	defer conn.Close()

	client := &Client{conn: conn}

	defer h.removeClient(client)

	h.mu.Lock()
	h.clients[client] = make(map[string]bool)
	h.mu.Unlock()

	done := make(chan struct{})
	defer close(done)

	h.heartbeat(client, done)

	for {
		var msg SubscribeEvent

		if err := conn.ReadJSON(&msg); err != nil {
			log.Printf("connection closed or read error encountered: %v", err)
			break
		}

		switch msg.Action {
		case "subscribe":
			if h.Subscribe(msg.Symbol, client) {
				h.subscriptionEvent <- msg
			}

		case "unsubscribe":
			if h.Unsubscribe(msg.Symbol, client) {
				h.subscriptionEvent <- msg
			}
		}
	}
}

func (h *Hub) Subscribe(symbol string, client *Client) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client][symbol] = true

	if h.subscriptions[symbol] == nil {
		h.subscriptions[symbol] = make(map[*Client]bool)
	}

	isFirst := len(h.subscriptions[symbol]) == 0
	h.subscriptions[symbol][client] = true

	return isFirst
}

func (h *Hub) Unsubscribe(symbol string, client *Client) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	if symbols, ok := h.clients[client]; ok {
		delete(symbols, symbol)
	}

	conns, ok := h.subscriptions[symbol]
	if !ok {
		return false
	}

	delete(conns, client)

	if len(conns) == 0 {
		delete(h.subscriptions, symbol)
		return true
	}

	return false
}

func (h *Hub) SubscriptionEvent() <-chan SubscribeEvent {
	return h.subscriptionEvent
}

func (h *Hub) Broadcast(price PriceEvent) {
	h.mu.RLock()
	clients := h.subscriptions[price.Symbol]
	h.mu.RUnlock()

	for client := range clients {
		if err := client.WriteJSON(price); err != nil {
			log.Printf("failed to send price: %v", err)
		}
	}
}

func (h *Hub) removeClient(client *Client) {
	h.mu.Lock()

	symbols := h.clients[client]

	var subscriptionEvents []SubscribeEvent

	for symbol := range symbols {
		delete(h.subscriptions[symbol], client)

		if len(h.subscriptions[symbol]) == 0 {
			delete(h.subscriptions, symbol)

			subscriptionEvents = append(subscriptionEvents, SubscribeEvent{
				Action: "unsubscribe",
				Symbol: symbol,
			})

		}

	}

	delete(h.clients, client)

	h.mu.Unlock()

	for _, event := range subscriptionEvents {
		h.subscriptionEvent <- event
	}
}

func (h *Hub) heartbeat(client *Client, done chan struct{}) {
	pongWait := 30 * time.Second

	client.conn.SetReadDeadline(time.Now().Add(pongWait))

	client.conn.SetPongHandler(func(appData string) error {
		client.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return

			case <-ticker.C:
				if err := client.WritePing(); err != nil {
					return
				}
			}
		}
	}()
}

package market

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type HubManager interface {
	Connect(ctx *gin.Context)
	Broadcast(price PriceEvent)
	SubscriptionEvent() <-chan SubscribeEvent
}

type Hub struct {
	clients           map[*websocket.Conn]bool
	subscriptions     map[string]map[*websocket.Conn]bool
	subscriptionEvent chan SubscribeEvent
	mu                sync.RWMutex
}

func NewHub() HubManager {
	return &Hub{
		clients:           map[*websocket.Conn]bool{},
		subscriptions:     map[string]map[*websocket.Conn]bool{},
		subscriptionEvent: make(chan SubscribeEvent, 100),
	}
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

	h.mu.Lock()
	h.clients[conn] = true
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.clients, conn)
		h.mu.Unlock()
	}()

	for {
		var msg SubscribeEvent

		if err := conn.ReadJSON(&msg); err != nil {
			log.Printf("connection closed or read error encountered: %v", err)
			break
		}

		switch msg.Action {
		case "subscribe":
			if h.Subscribe(msg.Symbol, conn) {
				h.subscriptionEvent <- msg
			}

		case "unsubscribe":
			if h.Unsubscribe(msg.Symbol, conn) {
				h.subscriptionEvent <- msg
			}
		}
	}
}

func (h *Hub) Subscribe(symbol string, conn *websocket.Conn) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.subscriptions[symbol] == nil {
		h.subscriptions[symbol] = make(map[*websocket.Conn]bool)
	}

	isFirst := len(h.subscriptions[symbol]) == 0
	h.subscriptions[symbol][conn] = true

	return isFirst
}

func (h *Hub) Unsubscribe(symbol string, conn *websocket.Conn) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	conns, ok := h.subscriptions[symbol]
	if !ok {
		return false
	}

	delete(conns, conn)

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
	defer h.mu.RUnlock()

	for conn := range h.subscriptions[price.Symbol] {
		if err := conn.WriteJSON(price); err != nil {
			log.Printf("failed to send price: %v", err)
		}
	}
}

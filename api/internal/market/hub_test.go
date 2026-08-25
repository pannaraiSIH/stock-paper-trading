package market

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHub_Connect(t *testing.T) {
	hub := NewHub()

	r := gin.New()
	r.GET("market/ws", hub.Connect)

	server := httptest.NewServer(r)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/market/ws"

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()
}

func TestHub_Connect_UpgradeFailed(t *testing.T) {
	hub := NewHub()

	r := gin.New()
	r.GET("market/ws", hub.Connect)

	req := httptest.NewRequest(http.MethodGet, "/market/ws", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)

	assert.Equal(t, http.StatusBadRequest, w.Code)

}

func TestHub_Subscribe(t *testing.T) {
	tests := []struct {
		name        string
		body        SubscribeEvent
		expectEvent bool
	}{
		{
			name: "does not emit event for invalid action",
			body: SubscribeEvent{
				Action: "invalid-action",
				Symbol: "AAPL",
			},
			expectEvent: false,
		},
		{
			name: "emits subscribe event",
			body: SubscribeEvent{
				Action: "subscribe",
				Symbol: "AAPL",
			},
			expectEvent: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hub := NewHub()

			r := gin.New()
			r.GET("/market/ws", hub.Connect)

			server := httptest.NewServer(r)
			defer server.Close()

			wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/market/ws"

			conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
			require.NoError(t, err)
			defer conn.Close()

			err = conn.WriteJSON(tt.body)
			require.NoError(t, err)

			if tt.expectEvent {
				select {
				case event := <-hub.SubscriptionEvent():
					assert.Equal(t, tt.body.Action, event.Action)
					assert.Equal(t, tt.body.Symbol, event.Symbol)

				case <-time.After(time.Second):
					t.Fatal("expected subscription event")
				}

				return
			}

			select {
			case event := <-hub.SubscriptionEvent():
				t.Fatalf("unexpected event: %+v", event)

			case <-time.After(100 * time.Millisecond):
			}
		})
	}
}

func TestHub_Unsubscribe(t *testing.T) {
	tests := []struct {
		name        string
		body        SubscribeEvent
		expectEvent bool
	}{
		{
			name: "does not emit event for invalid action",
			body: SubscribeEvent{
				Action: "invalid-action",
				Symbol: "AAPL",
			},
			expectEvent: false,
		},
		{
			name: "emits unsubscribe event",
			body: SubscribeEvent{
				Action: "unsubscribe",
				Symbol: "AAPL",
			},
			expectEvent: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hub := NewHub()

			r := gin.New()
			r.GET("/market/ws", hub.Connect)

			server := httptest.NewServer(r)
			defer server.Close()

			wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/market/ws"

			conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
			require.NoError(t, err)
			defer conn.Close()

			err = conn.WriteJSON(SubscribeEvent{
				Action: "subscribe",
				Symbol: "AAPL",
			})
			require.NoError(t, err)

			select {
			case <-hub.SubscriptionEvent():
			case <-time.After(time.Second):
				t.Fatal("expected subscribe event")
			}

			err = conn.WriteJSON(tt.body)
			require.NoError(t, err)

			if tt.expectEvent {
				select {
				case event := <-hub.SubscriptionEvent():
					assert.Equal(t, tt.body.Action, event.Action)
					assert.Equal(t, tt.body.Symbol, event.Symbol)

				case <-time.After(time.Second):
					log.Fatal("expected subscription event")
				}

				return
			}

			select {
			case <-hub.SubscriptionEvent():

			case <-time.After(100 * time.Millisecond):
			}
		})
	}
}

func TestHub_Broadcast(t *testing.T) {
	tests := []struct {
		name            string
		subscribeSymbol string
		broadcastSymbol string
		expectMessage   bool
	}{
		{
			name:            "does not broadcast to different symbol subscriber",
			subscribeSymbol: "A",
			broadcastSymbol: "AAPL",
			expectMessage:   false,
		},
		{
			name:            "broadcasts to subscribed symbol",
			subscribeSymbol: "AAPL",
			broadcastSymbol: "AAPL",
			expectMessage:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hub := NewHub()

			r := gin.New()
			r.GET("/market/ws", hub.Connect)

			server := httptest.NewServer(r)
			defer server.Close()

			wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/market/ws"

			conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
			require.NoError(t, err)
			defer conn.Close()

			err = conn.WriteJSON(SubscribeEvent{
				Action: "subscribe",
				Symbol: tt.subscribeSymbol,
			})
			require.NoError(t, err)

			select {
			case <-hub.SubscriptionEvent():
			case <-time.After(time.Second):
				log.Fatal("expected subscription event")
			}

			price := PriceEvent{
				Symbol: tt.broadcastSymbol,
			}

			hub.Broadcast(price)

			if tt.expectMessage {
				var msg PriceEvent

				err := conn.ReadJSON(&msg)
				require.NoError(t, err)

				assert.Equal(t, tt.broadcastSymbol, msg.Symbol)
				return
			}

			err = conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
			require.NoError(t, err)

			var msg PriceEvent
			err = conn.ReadJSON(&msg)

			assert.Error(t, err)
		})
	}
}

package market

import (
	"context"
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockRealtimeProvider struct {
	ConnectFunc         func(ctx context.Context) error
	DisconnectFunc      func()
	SubscribeStatusFunc func() <-chan SubscribeStatusEvent
	SubscribeFunc       func(symbol string) error
	UnsubscribeFunc     func(symbol string) error
	PricesFunc          func() <-chan PriceEvent
	ErrorsFunc          func() <-chan error
}

func (m *MockRealtimeProvider) Connect(ctx context.Context) error {
	return m.ConnectFunc(ctx)
}

func (m *MockRealtimeProvider) Disconnect() {
	m.DisconnectFunc()
}

func (m *MockRealtimeProvider) SubscribeStatus() <-chan SubscribeStatusEvent {
	return m.SubscribeStatusFunc()
}

func (m *MockRealtimeProvider) Subscribe(symbol string) error {
	return m.SubscribeFunc(symbol)
}

func (m *MockRealtimeProvider) Unsubscribe(symbol string) error {
	return m.UnsubscribeFunc(symbol)
}

func (m *MockRealtimeProvider) Prices() <-chan PriceEvent {
	return m.PricesFunc()
}

func (m *MockRealtimeProvider) Errors() <-chan error {
	return m.ErrorsFunc()
}

type MockHubManager struct {
	ConnectFunc           func(ctx *gin.Context)
	BroadcastFunc         func(price PriceEvent)
	SubscriptionEventFunc func() <-chan SubscribeEvent
}

func (m *MockHubManager) Connect(ctx *gin.Context) {
	if m.ConnectFunc != nil {
		m.ConnectFunc(ctx)
	}
}

func (m *MockHubManager) Broadcast(price PriceEvent) {
	m.BroadcastFunc(price)
}

func (m *MockHubManager) SubscriptionEvent() <-chan SubscribeEvent {
	return m.SubscriptionEventFunc()
}

type MockMarketCache struct {
	GetLatestPriceFunc func(ctx context.Context, symbol string) (PriceEvent, error)
	SetLatestPriceFunc func(ctx context.Context, symbol string, price PriceEvent) error
}

func (m *MockMarketCache) GetLatestPrice(
	ctx context.Context,
	symbol string,
) (PriceEvent, error) {
	return m.GetLatestPriceFunc(ctx, symbol)
}

func (m *MockMarketCache) SetLatestPrice(
	ctx context.Context,
	symbol string,
	price PriceEvent,
) error {
	return m.SetLatestPriceFunc(ctx, symbol, price)
}

func TestRunWorker(t *testing.T) {
	tests := []struct {
		name                      string
		event                     *SubscribeEvent
		priceEvent                *PriceEvent
		connectErr                error
		cacheErr                  error
		expectedSubscribeCalled   bool
		expectedUnsubscribeCalled bool
		expectedCacheCalled       bool
		expectedBroadcastCalled   bool
	}{
		{
			name:       "provider connect error",
			connectErr: errors.New("connect error"),
		},
		{
			name: "invalid action",
			event: &SubscribeEvent{
				Action: "invalid-action",
				Symbol: "AAPL",
			},
			expectedSubscribeCalled: false,
		},
		{
			name: "subscribe event success",
			event: &SubscribeEvent{
				Action: "subscribe",
				Symbol: "AAPL",
			},
			expectedSubscribeCalled: true,
		},
		{
			name: "unsubscribe event success",
			event: &SubscribeEvent{
				Action: "unsubscribe",
				Symbol: "AAPL",
			},
			expectedUnsubscribeCalled: true,
		},
		{
			name: "cache error",
			priceEvent: &PriceEvent{
				Event:  "price",
				Symbol: "AAPL",
			},
			cacheErr:                errors.New("cache error"),
			expectedCacheCalled:     true,
			expectedBroadcastCalled: true,
		},
		{
			name: "price received",
			priceEvent: &PriceEvent{
				Event:  "price",
				Symbol: "AAPL",
			},
			expectedCacheCalled:     true,
			expectedBroadcastCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var broadcastCalled bool
			var cacheCalled bool
			var subscribeCalled bool
			var unsubscribeCalled bool

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			provider := &MockRealtimeProvider{
				ConnectFunc: func(ctx context.Context) error {
					return tt.connectErr
				},
				SubscribeFunc: func(symbol string) error {
					subscribeCalled = true
					cancel()
					return nil
				},
				UnsubscribeFunc: func(symbol string) error {
					unsubscribeCalled = true
					cancel()
					return nil
				},
				PricesFunc: func() <-chan PriceEvent {
					ch := make(chan PriceEvent, 1)

					if tt.priceEvent != nil {
						ch <- *tt.priceEvent
					}

					return ch
				},
				SubscribeStatusFunc: func() <-chan SubscribeStatusEvent {
					ch := make(chan SubscribeStatusEvent)
					close(ch)
					return ch
				},
				ErrorsFunc: func() <-chan error {
					ch := make(chan error)
					close(ch)
					return ch
				},
			}
			hub := &MockHubManager{
				BroadcastFunc: func(price PriceEvent) {
					broadcastCalled = true
					cancel()
				},
				SubscriptionEventFunc: func() <-chan SubscribeEvent {
					ch := make(chan SubscribeEvent, 1)

					if tt.event != nil {
						ch <- *tt.event

						if tt.event.Action == "invalid-action" {
							cancel()
						}
					}

					close(ch)
					return ch
				},
			}
			marketCache := &MockMarketCache{
				SetLatestPriceFunc: func(ctx context.Context, symbol string, price PriceEvent) error {
					cacheCalled = true
					return tt.cacheErr
				},
			}

			worker := NewMarketWorker(provider, hub, marketCache)

			worker.Run(ctx)

			assert.Equal(t, tt.expectedBroadcastCalled, broadcastCalled)
			assert.Equal(t, tt.expectedCacheCalled, cacheCalled)
			assert.Equal(t, tt.expectedSubscribeCalled, subscribeCalled)
			assert.Equal(t, tt.expectedUnsubscribeCalled, unsubscribeCalled)
		})
	}
}

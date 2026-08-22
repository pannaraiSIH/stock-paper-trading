package market

import (
	"context"
)

type SubscribeStatusItem struct {
	Symbol   string `json:"symbol"`
	Exchange string `json:"exchange,omitempty"`
	MicCode  string `json:"mic_code,omitempty"`
}

type SubscribeStatusEvent struct {
	Event   string                `json:"event"`
	Status  string                `json:"status"`
	Success []SubscribeStatusItem `json:"success"`
	Fails   []SubscribeStatusItem `json:"fails"`
}

type PriceEvent struct {
	Event         string   `json:"event"`
	Symbol        string   `json:"symbol"`
	Exchange      string   `json:"exchange,omitempty"`
	MicCode       string   `json:"mic_code,omitempty"`
	Type          string   `json:"type,omitempty"`
	Currency      string   `json:"currency,omitempty"`
	CurrencyBase  string   `json:"currency_base,omitempty"`
	CurrencyQuote string   `json:"currency_quote,omitempty"`
	Timestamp     int64    `json:"timestamp"`
	Price         float64  `json:"price"`
	DayVolume     *int64   `json:"day_volume,omitempty"`
	Bid           *float64 `json:"bid,omitempty"`
	Ask           *float64 `json:"ask,omitempty"`
}

type RealtimeProvider interface {
	Connect(ctx context.Context) error
	Disconnect()
	SubscribeStatus() <-chan SubscribeStatusEvent
	Subscribe(symbol string) error
	Unsubscribe(symbol string) error
	Prices() <-chan PriceEvent
	Errors() <-chan error
}

type SubscribeEvent struct {
	Action string `json:"action"`
	Symbol string `json:"symbol"`
}

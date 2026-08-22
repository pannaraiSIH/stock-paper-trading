package twelvedata

import (
	"context"

	"github.com/pannaraiSIH/stock-paper-trading/internal/market"
	"github.com/twelvedata/twelvedata-go/twelvedata/ws"
)

type TwelveDataWebsocket struct {
	websocket *ws.Client
}

func NewTwelveDataWebsocket(apiKey string) (market.RealtimeProvider, error) {
	client, err := ws.NewClient(ws.Options{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, err
	}

	return &TwelveDataWebsocket{
		websocket: client,
	}, nil
}

func (tw *TwelveDataWebsocket) Connect(ctx context.Context) error {
	return tw.websocket.Connect(ctx)
}

func (tw *TwelveDataWebsocket) Disconnect() {
	tw.websocket.Disconnect()
}

func (tw *TwelveDataWebsocket) SubscribeStatus() <-chan market.SubscribeStatusEvent {
	out := make(chan market.SubscribeStatusEvent)

	go func() {
		defer close(out)

		for event := range tw.websocket.SubscribeStatuses() {
			success := make([]market.SubscribeStatusItem, 0, len(event.Success))
			for _, item := range event.Success {
				success = append(success, market.SubscribeStatusItem{
					Symbol:   item.Symbol,
					Exchange: item.Exchange,
					MicCode:  item.MicCode,
				})
			}

			fails := make([]market.SubscribeStatusItem, 0, len(event.Fails))
			for _, item := range event.Fails {
				fails = append(fails, market.SubscribeStatusItem{
					Symbol:   item.Symbol,
					Exchange: item.Exchange,
					MicCode:  item.MicCode,
				})
			}

			out <- market.SubscribeStatusEvent{
				Event:   event.Event,
				Status:  event.Status,
				Success: success,
				Fails:   fails,
			}
		}
	}()

	return out
}

func (tw *TwelveDataWebsocket) Subscribe(symbol string) error {
	return tw.websocket.Subscribe(symbol)
}

func (tw *TwelveDataWebsocket) Unsubscribe(symbol string) error {
	return tw.websocket.Unsubscribe(symbol)
}

func (tw *TwelveDataWebsocket) Prices() <-chan market.PriceEvent {
	out := make(chan market.PriceEvent)

	go func() {
		defer close(out)

		for event := range tw.websocket.Prices() {
			out <- market.PriceEvent{
				Event:         event.Event,
				Symbol:        event.Symbol,
				Exchange:      event.Exchange,
				MicCode:       event.MicCode,
				Type:          event.Type,
				Currency:      event.Currency,
				CurrencyBase:  event.CurrencyBase,
				CurrencyQuote: event.CurrencyQuote,
				Timestamp:     event.Timestamp,
				Price:         event.Price,
				DayVolume:     event.DayVolume,
				Bid:           event.Bid,
				Ask:           event.Ask,
			}
		}
	}()

	return out
}

func (tw *TwelveDataWebsocket) Errors() <-chan error {
	return tw.websocket.Errors()
}

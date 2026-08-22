package market

import (
	"context"
	"fmt"
	"log"
)

type marketWorker struct {
	provider    RealtimeProvider
	hubManager  HubManager
	marketCache MarketCache
}

func NewMarketWorker(
	provider RealtimeProvider,
	hubManager HubManager,
	priceCache MarketCache,
) *marketWorker {
	return &marketWorker{
		provider:    provider,
		hubManager:  hubManager,
		marketCache: priceCache,
	}
}

func (w *marketWorker) Run(ctx context.Context) {
	if err := w.provider.Connect(ctx); err != nil {
		log.Printf("failed to connect market provider: %v", err)
	}

	go func() {
		for event := range w.hubManager.SubscriptionEvent() {
			switch event.Action {
			case "subscribe":
				log.Printf("subscribe event: %v\n", event.Symbol)
				if err := w.provider.Subscribe(event.Symbol); err != nil {
					log.Printf("failed to subscribe %s: %v", event.Symbol, err)
				}

			case "unsubscribe":
				if err := w.provider.Unsubscribe(event.Symbol); err != nil {
					log.Printf("failed to unsubscribe %s: %v", event.Symbol, err)
				}
			}
		}
	}()

	subscribeStatusCh := w.provider.SubscribeStatus()

	go func() {
		for s := range subscribeStatusCh {
			log.Printf(
				"subscribe status: %s success=%v fails=%v\n",
				s.Status,
				s.Success,
				s.Fails,
			)
		}
	}()

	go func() {
		for err := range w.provider.Errors() {
			log.Printf("websocket error: %v\n", err)
		}
	}()

	prices := w.provider.Prices()

	for {
		select {
		case <-ctx.Done():
			return
		case price, ok := <-prices:
			if !ok {
				return
			}

			fmt.Printf("%s @ %v\n", price.Symbol, price.Price)

			if err := w.marketCache.SetLatestPrice(ctx, price.Symbol, price); err != nil {
				log.Printf("failed to cache latest price: %v", err)
			}

			w.hubManager.Broadcast(price)
		}
	}
}

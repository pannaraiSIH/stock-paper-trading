package twelvedata

import (
	"context"

	"github.com/pannaraiSIH/stock-paper-trading/internal/market"
	"github.com/twelvedata/twelvedata-go/twelvedata"
)

type TwelveDataClient struct {
	client *twelvedata.APIClient
}

func NewTwelveDataClient(apiKey string) (*TwelveDataClient, error) {
	cfg, err := twelvedata.NewConfig(apiKey)
	if err != nil {
		return nil, err
	}

	return &TwelveDataClient{
		client: twelvedata.NewAPIClient(cfg),
	}, nil
}

func (tw *TwelveDataClient) SearchStocks(ctx context.Context, query string, outputSize int64) ([]market.StockSearchResponse, error) {
	resp, _, err := tw.client.ReferenceDataAPI.
		GetSymbolSearch(ctx).
		Symbol(query).
		Outputsize(outputSize).
		ShowPlan(true).
		Execute()
	if err != nil {
		return nil, err
	}

	var stocks []market.StockSearchResponse

	for _, r := range resp.Data {
		stocks = append(stocks, market.StockSearchResponse{
			Symbol:   r.Symbol,
			Name:     r.InstrumentName,
			Exchange: r.Exchange,
			Currency: r.Currency,
		})
	}

	return stocks, nil
}

func (tw *TwelveDataClient) GetCandles(
	ctx context.Context,
	symbol string,
	interval market.Interval,
	outputSize int64) ([]market.GetCandleResponse, error) {

	resp, _, err := tw.client.MarketDataAPI.
		GetTimeSeries(ctx).
		Symbol(symbol).
		Interval(twelvedata.IntervalEnum(interval)).
		Outputsize(outputSize).
		Execute()
	if err != nil {
		return nil, err
	}

	var candles []market.GetCandleResponse

	for _, r := range resp.Values {
		candles = append(candles, market.GetCandleResponse{
			Datetime: r.Datetime,
			Open:     r.Open,
			High:     r.High,
			Low:      r.Low,
			Close:    r.Close,
			Volume:   r.GetVolume(),
		})
	}

	return candles, nil
}

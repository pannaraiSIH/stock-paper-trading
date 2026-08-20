package twelvedata

import (
	"context"
	"errors"
	"net/http"

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

func (tw *TwelveDataClient) SearchStocks(ctx context.Context, query market.SearchStocksQuery) ([]market.SearchStocksResponse, error) {
	resp, _, err := tw.client.ReferenceDataAPI.
		GetSymbolSearch(ctx).
		Symbol(query.Query).
		Outputsize(query.OutputSize).
		ShowPlan(true).
		Execute()
	if err != nil {
		return nil, handleTwelveDataError(err)
	}

	var stocks []market.SearchStocksResponse

	for _, r := range resp.Data {
		stocks = append(stocks, market.SearchStocksResponse{
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
	query market.GetCandlesQuery,
) ([]market.GetCandleResponse, error) {

	resp, _, err := tw.client.MarketDataAPI.
		GetTimeSeries(ctx).
		Symbol(symbol).
		Interval(twelvedata.IntervalEnum(query.Interval)).
		Outputsize(query.OutputSize).
		Execute()
	if err != nil {
		return nil, handleTwelveDataError(err)
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

func (tw *TwelveDataClient) GetStockDetails(ctx context.Context, symbol string) (market.GetStockDetailsResponse, error) {
	resp, _, err := tw.client.FundamentalsAPI.GetProfile(ctx).Symbol(symbol).Execute()
	if err != nil {
		return market.GetStockDetailsResponse{}, handleTwelveDataError(err)
	}

	return market.GetStockDetailsResponse{
		Symbol:      resp.Symbol,
		Name:        resp.Name,
		Exchange:    resp.Exchange,
		MicCode:     resp.MicCode,
		Sector:      resp.Sector,
		Industry:    resp.Industry,
		Website:     resp.Website,
		Description: resp.Description,
		Type:        resp.Type,
		CEO:         resp.CEO,
		Address:     resp.Address,
		City:        resp.City,
		State:       resp.State,
		Country:     resp.Country,
		Phone:       resp.Phone,
	}, nil
}

func handleTwelveDataError(err error) error {
	var apiErr twelvedata.TwelvedataApiError

	if errors.As(err, &apiErr) {
		statusCode := apiErr.GetStatusCode()

		switch {
		case statusCode == http.StatusNotFound:
			return market.ErrStockNotFound
		case statusCode == http.StatusForbidden:
			return market.ErrMarketProviderUnavailable
		}
	}

	return err
}

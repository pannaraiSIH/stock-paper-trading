package market

type SearchStocksQuery struct {
	Query      string `form:"query" binding:"required"`
	OutputSize int64  `form:"outputSize" binding:"required"`
}

type StockSearchResponse struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Exchange string `json:"exchange"`
	Currency string `json:"currency"`
}

type GetCandlesQuery struct {
	Interval   Interval `form:"interval" binding:"required"`
	OutputSize int64    `form:"outputSize" binding:"required"`
}

type GetCandleResponse struct {
	Datetime string `json:"datetime"`
	Open     string `json:"open"`
	High     string `json:"high"`
	Low      string `json:"low"`
	Close    string `json:"close"`
	Volume   string `json:"volume"`
}

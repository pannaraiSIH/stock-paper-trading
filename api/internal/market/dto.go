package market

type SearchStocksQuery struct {
	Query      string `form:"query" binding:"required"`
	OutputSize int64  `form:"outputSize" binding:"required"`
}

type SearchStocksResponse struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Exchange string `json:"exchange"`
	Currency string `json:"currency"`
}

type GetCandlesQuery struct {
	Symbol     string   `form:"symbol" binding:"required"`
	Interval   Interval `form:"interval" binding:"required,oneof=1min 5min 1h 1day"`
	OutputSize int64    `form:"outputSize" binding:"required"`
}

type GetStockDetailsQuery struct {
	Symbol string `form:"symbol" binding:"required"`
}

type GetCandleResponse struct {
	Datetime string `json:"datetime"`
	Open     string `json:"open"`
	High     string `json:"high"`
	Low      string `json:"low"`
	Close    string `json:"close"`
	Volume   string `json:"volume"`
}

type GetStockDetailsResponse struct {
	Symbol      string  `json:"symbol"`
	Name        string  `json:"name"`
	Exchange    string  `json:"exchange"`
	MicCode     string  `json:"micCode"`
	Sector      *string `json:"sector"`
	Industry    *string `json:"industry"`
	Website     *string `json:"website"`
	Description *string `json:"description"`
	Type        *string `json:"type"`
	CEO         *string `json:"ceo"`
	Address     *string `json:"address"`
	City        *string `json:"city"`
	State       *string `json:"state"`
	Country     *string `json:"country"`
	Phone       *string `json:"phone"`
}

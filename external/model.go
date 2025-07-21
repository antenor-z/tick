package external

type Stock []struct {
	Ticker       string  `json:"ticker"`
	AveragePrice float64 `json:"averagePrice"`
}

type YahooResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				RegularMarketPrice float64 `json:"regularMarketPrice"`
			} `json:"meta"`
		} `json:"result"`
	} `json:"chart"`
}

type BinanceResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

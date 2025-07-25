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
			Indicators struct {
				Quote []struct {
					Close []float64 `json:"close"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
	} `json:"chart"`
}

type PriceData struct {
	Price  float64
	Change float64
	Color  string
}

type Prices struct {
	Now          PriceData
	SevenDaysAgo PriceData
}

type BinanceResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

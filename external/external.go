package external

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func LoadStocks(filename string) (Stock, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config Stock
	err = json.Unmarshal(data, &config)
	return config, err
}

func getColor(change float64) string {
	if change < 0 {
		return "red"
	}
	return "green"
}

func fetchPriceBinance(symbol string, averagePrice float64) (Prices, error) {
	nowResp, err := http.Get("https://api.binance.com/api/v3/ticker/price?symbol=" + symbol)
	if err != nil {
		return Prices{}, err
	}
	defer nowResp.Body.Close()

	var nowResult struct {
		Price string `json:"price"`
	}
	if err := json.NewDecoder(nowResp.Body).Decode(&nowResult); err != nil {
		return Prices{}, err
	}
	nowPrice, err := strconv.ParseFloat(nowResult.Price, 64)
	if err != nil {
		return Prices{}, err
	}

	sevenDaysAgo := time.Now().AddDate(0, 0, -7).UnixMilli()
	klineURL := fmt.Sprintf("https://api.binance.com/api/v3/klines?symbol=%s&interval=1d&limit=1&startTime=%d", symbol, sevenDaysAgo)

	klineResp, err := http.Get(klineURL)
	if err != nil {
		return Prices{}, err
	}
	defer klineResp.Body.Close()

	var klineData [][]interface{}
	if err := json.NewDecoder(klineResp.Body).Decode(&klineData); err != nil {
		return Prices{}, err
	}
	if len(klineData) == 0 || len(klineData[0]) < 5 {
		return Prices{}, fmt.Errorf("invalid kline data for %s", symbol)
	}

	closeStr, ok := klineData[0][4].(string)
	if !ok {
		return Prices{}, fmt.Errorf("unexpected close price format")
	}
	sevenDaysAgoPrice, err := strconv.ParseFloat(closeStr, 64)
	if err != nil {
		return Prices{}, err
	}

	nowChange := ((nowPrice - averagePrice) / averagePrice) * 100
	oldChange := ((sevenDaysAgoPrice - averagePrice) / averagePrice) * 100

	return Prices{
		Now: PriceData{
			Price:  nowPrice,
			Change: nowChange,
			Color:  getColor(nowChange),
		},
		SevenDaysAgo: PriceData{
			Price:  sevenDaysAgoPrice,
			Change: oldChange,
			Color:  getColor(oldChange),
		},
	}, nil
}

func fetchPriceYahoo(ticker string, averagePrice float64) (Prices, error) {
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?range=7d&interval=1d", ticker)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	if err != nil {
		return Prices{}, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Prices{}, err
	}
	defer resp.Body.Close()

	var result YahooResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return Prices{}, err
	}

	if len(result.Chart.Result) == 0 {
		return Prices{}, fmt.Errorf("no result for ticker %s", ticker)
	}

	chart := result.Chart.Result[0]
	closes := chart.Indicators.Quote[0].Close
	if len(closes) < 2 {
		return Prices{}, fmt.Errorf("not enough historical data for ticker %s", ticker)
	}

	return Prices{
		Now: PriceData{
			Price:  chart.Meta.RegularMarketPrice,
			Change: ((chart.Meta.RegularMarketPrice - averagePrice) / averagePrice) * 100,
			Color:  getColor(((chart.Meta.RegularMarketPrice - averagePrice) / averagePrice) * 100),
		},
		SevenDaysAgo: PriceData{
			Price:  closes[0],
			Change: ((chart.Meta.RegularMarketPrice - closes[0]) / closes[0]) * 100,
			Color:  getColor(((chart.Meta.RegularMarketPrice - closes[0]) / closes[0]) * 100),
		},
	}, nil
}

func FetchPrice(ticker string, averagePrice float64) (Prices, error) {
	if strings.HasPrefix(ticker, "BTC") {
		return fetchPriceBinance(ticker, averagePrice)
	} else {
		return fetchPriceYahoo(ticker, averagePrice)
	}
}

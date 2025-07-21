package external

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
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

func fetchPriceBTC(ticker string) (float64, error) {
	url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s", ticker)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
	if err != nil {
		return 0, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result BinanceResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return 0, err
	}

	price, err := strconv.ParseFloat(result.Price, 64)
	if err != nil {
		return 0, err
	}

	return price, nil
}

func fetchPriceYahoo(ticker string) (float64, error) {
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s", ticker)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
	if err != nil {
		return 0, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result YahooResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return 0, err
	}

	if len(result.Chart.Result) == 0 {
		return 0, fmt.Errorf("no result for ticker %s", ticker)
	}

	return result.Chart.Result[0].Meta.RegularMarketPrice, nil
}

func FetchPrice(ticker string) (float64, error) {
	if strings.HasPrefix(ticker, "BTC") {
		return fetchPriceBTC(ticker)
	} else {
		return fetchPriceYahoo(ticker)
	}
}

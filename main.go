package main

import (
	"fmt"
	"os"
	"strings"
	"tick/external"
	"tick/sendMail"
	"time"
)

func main() {
	stocks, err := external.LoadStocks("stocks.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	mailContent := "<html><body><span style='color: black'>"
	cleanContent := "################\n"
	cleanContent += fmt.Sprintf("Tick Report %s\n", time.Now().Format("2006-01-02 15:04"))
	for _, stock := range stocks {
		price, err := external.FetchPrice(stock.Ticker, stock.AveragePrice)
		if err != nil {
			mailContent += fmt.Sprintf("%s: (Error fetching current price: %v)<br />\n", stock.Ticker, err)
			continue
		}

		tickerSafe := strings.ReplaceAll(stock.Ticker, ".SA", "")
		mailContent += fmt.Sprintf("<strong>%s</strong><br />\n", tickerSafe)
		mailContent += fmt.Sprintf("&nbsp;&nbsp;AVG: R$ %.2f<br />\n", stock.AveragePrice)
		mailContent += fmt.Sprintf("&nbsp;&nbsp;NOW: R$ %.2f<br />\n", price.Now.Price)
		mailContent += fmt.Sprintf("&nbsp;&nbsp;7DAYSAGO: R$ %.2f<br />\n", price.SevenDaysAgo.Price)
		mailContent += fmt.Sprintf("&nbsp;&nbsp;AVG x NOW: <span style='color:%s'>%.2f%%</span><br />\n", price.Now.Color, price.Now.Change)
		mailContent += fmt.Sprintf("&nbsp;&nbsp;7DAYSAGO x NOW: <span style='color:%s'>%.2f%%</span><br />\n", price.SevenDaysAgo.Color, price.SevenDaysAgo.Change)
		if !strings.HasPrefix(stock.Ticker, "BTC") {
			mailContent += fmt.Sprintf("&nbsp;&nbsp;https://finance.yahoo.com/quote/%s<br />\n", stock.Ticker)
		}
		mailContent += "<br />\n"

		cleanContent += fmt.Sprintf("-> %s\n", stock.Ticker)
		cleanContent += fmt.Sprintf("  AVG: R$ %.2f\n", stock.AveragePrice)
		cleanContent += fmt.Sprintf("  NOW: R$ %.2f\n", price.Now.Price)
		cleanContent += fmt.Sprintf("  7DAYSAGO: R$ %.2f\n", price.SevenDaysAgo.Price)
		cleanContent += fmt.Sprintf("  AVG x NOW: %.2f%%\n", price.Now.Change)
		cleanContent += fmt.Sprintf("  7DAYSAGO x NOW: %.2f%%\n", price.SevenDaysAgo.Change)
	}

	mailContent += "</span></body></html>"
	cleanContent += "end\n\n\n\n"
	fmt.Print(cleanContent)

	sendMail.SendMail(
		sendMail.Email{
			To:      sendMail.GetConfig().Mailgun.Receiver,
			Subject: fmt.Sprintf("Tick Report %s", time.Now().Format("2006-01-02 15:04")),
			Text:    mailContent,
		},
		sendMail.GetConfig())
}

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

	mailContent := "<span style='color: black'>"

	for _, stock := range stocks {
		currentPrice, err := external.FetchPrice(stock.Ticker)
		if err != nil {
			mailContent += fmt.Sprintf("%s: (Error fetching current price: %v)<br />\n", stock.Ticker, err)
			continue
		}

		change := ((currentPrice - stock.AveragePrice) / stock.AveragePrice) * 100
		var color string
		if change > 0 {
			color = "green"
		} else if change < 0 {
			color = "red"
		} else {
			color = "gray"
		}

		tickerSafe := strings.ReplaceAll(stock.Ticker, ".SA", "")
		mailContent += fmt.Sprintf("<strong>%s</strong><br />\n", tickerSafe)
		mailContent += fmt.Sprintf("&nbsp;&nbsp;Avg Price: R$ %.2f<br />\n", stock.AveragePrice)
		mailContent += fmt.Sprintf("&nbsp;&nbsp;Now Price: R$ %.2f<br />\n", currentPrice)
		mailContent += fmt.Sprintf("&nbsp;&nbsp;Change: <span style='color:%s'>%.2f%%</span><br /><br />\n\n", color, change)
	}

	mailContent += "</span>"

	sendMail.SendMail(
		sendMail.Email{
			To:      sendMail.GetConfig().Mailgun.Receiver,
			Subject: fmt.Sprintf("Tick Report %s", time.Now().Format("2006-01-02 15:04")),
			Text:    mailContent,
		},
		sendMail.GetConfig())
}

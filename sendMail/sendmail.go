package sendMail

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Email struct {
	To      string
	Subject string
	Text    string
}

func SendMail(email Email, config Config) {
	form := url.Values{}
	form.Add("from", config.Mailgun.Sender)
	form.Add("to", email.To)
	form.Add("subject", email.Subject)
	form.Add("html", email.Text)
	req, err := http.NewRequest("POST", "https://api.mailgun.net/v3/"+config.Mailgun.Domain+"/messages", strings.NewReader(form.Encode()))
	if err != nil {
		panic(err)
	}

	req.SetBasicAuth("api", config.Mailgun.ApiKey)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if res.StatusCode != 200 {
		panic("Failed to send email. " + res.Status)
	}
}

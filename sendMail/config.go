package sendMail

import (
	"encoding/json"
	"os"
)

type Config struct {
	Mailgun struct {
		ApiKey   string `json:"apiKey"`
		Sender   string `json:"sender"`
		Domain   string `json:"domain"`
		Receiver string `json:"receiver"`
	} `json:"mailgun"`
}

func GetConfig() Config {
	var config Config
	dat, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(dat, &config)
	if err != nil {
		panic(err)
	}
	return config
}

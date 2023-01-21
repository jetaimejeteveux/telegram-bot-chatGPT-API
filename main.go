package main

import (
	// "context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	// "github.com/PullRequestInc/go-gpt3"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func getToken(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}
func getAPIkey(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func main() {
	bot, err := tgbotapi.NewBotAPI(getToken("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	//picking model
	model := "text-davinci-003"

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			prompt := update.Message.Text
			//create request body
			body := strings.NewReader(fmt.Sprintf(`{
				"model": "%s",
				"prompt": "%s"
			}`, model, prompt))

			req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", body)
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", getAPIkey("API_KEY")))
			if err != nil {
				panic(err)
			}

			// Send the request
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			// Read the response
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}
			log.Printf(" response dari API = %v", data)
			// Print the response
			log.Printf(" response dari API bentuk string = %s", string(data))

			log.Printf("\n ini printf 1 [%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, string(data))
			log.Printf("\n ini printf 2[%v] %v", update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			log.Printf("ini printf 3 %v %v", msg.ReplyToMessageID, update.Message.MessageID)

			bot.Send(msg)
		}
	}
}

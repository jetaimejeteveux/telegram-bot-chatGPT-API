package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	// "github.com/PullRequestInc/go-gpt3"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jetaimejeteveux/tele-bot-spreadsheetAPI/models"
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

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			//create request body
			reqBody := models.Request{
				ModelRequest: "text-davinci-003",
				Prompt:       update.Message.Text,
				Temperature:  1,
				MaxTokens:    100,
			}
			jsonBody, _ := json.Marshal(reqBody)
			req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewBuffer(jsonBody))
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
			jsonData, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}

			var data models.TextCompletionResponse
			json.Unmarshal([]byte(jsonData), &data)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, data.Choices[0].Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}

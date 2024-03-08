package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	db  *sql.DB
	bot *tgbotapi.BotAPI
)

func main() {
	d, err := ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}

	db = d

	b, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	log.Print("bot successfully authenticated")

	bot = b
	bot.Debug = false
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log.Printf("Running bot service...")

	for u := range updates {
		if u.Message != nil {
			if u.Message.IsCommand() {
				go handleCommand(ctx, u)
			} else {
				go handleMessage(ctx, u)
			}
		}
	}
}

func handleCommand(ctx context.Context, u tgbotapi.Update) {
	select {
	case <-ctx.Done():
		return
	default:
		switch u.Message.Command() {
		case "start":
			if err := DeleteHistory(db, u.Message.Chat.ID); err != nil {
				failedProcessUpdate(u, err)
				return
			}
			githubBtn := tgbotapi.NewInlineKeyboardButtonURL("Github of Creator", "https://github.com/dnridwn")
			sendMessageToUser(u, "Hello! I am an AI Chatbot with Gemini model.\nI was created by Den Ridwan Saputra\n\nPlease write your message.", tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{githubBtn}))
		}
	}
}

func handleMessage(ctx context.Context, u tgbotapi.Update) {
	select {
	case <-ctx.Done():
		return
	default:
		log.Printf("Received new message from %s", u.Message.From.UserName)

		c, err := GetLatestHistory(db, u.Message.Chat.ID)
		if err != nil {
			failedProcessUpdate(u, err)
			return
		}

		history := []Content{}
		if c.ID > 0 {
			if err := json.Unmarshal([]byte(c.History), &history); err != nil {
				failedProcessUpdate(u, err)
				return
			}
		}

		g, err := NewGemini(os.Getenv("GEMINI_API_KEY"), geminiProModel)
		if err != nil {
			failedProcessUpdate(u, err)
			return
		}

		g.StartChat(history)
		resp, err := g.SendMessage(u.Message.Text)
		if err != nil {
			failedProcessUpdate(u, err)
			return
		}

		sendMessageToUser(u, resp.String(), nil)
		SaveHistory(db, u.Message.Chat.ID, g.GetHistory())
	}
}

func failedProcessUpdate(u tgbotapi.Update, err error) {
	log.Println(err)
	sendMessageToUser(u, "Sorry, something went wrong. Please re-send your message", nil)
}

func sendMessageToUser(u tgbotapi.Update, message string, replyMarkup interface{}) {
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, message)
	msg.ReplyToMessageID = u.Message.MessageID
	msg.ReplyMarkup = replyMarkup
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

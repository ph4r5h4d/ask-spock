package main

import (
	"database/sql"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ph4r5h4d/ask-spock/pkg/gtp35turbo"
	"github.com/ph4r5h4d/ask-spock/repository"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"os"
)

type Dependencies struct {
	db           *sql.DB
	bot          *tgbotapi.BotAPI
	openaiClient *openai.Client
}

func main() {
	d, err := setup()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
	}(d.db)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := d.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			user, _ := repository.GetOrCreateUser(d.db, update.Message.From.UserName)
			if !user.Active {
				d.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "You are not allowed to use the service"))
				continue
			}
			log.Info().Msgf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			resp, err := gtp35turbo.Response(d.openaiClient, update.Message.Text)
			if err != nil {
				tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())

			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, resp)
			msg.ReplyToMessageID = update.Message.MessageID

			d.bot.Send(msg)
		}
	}
}

func setup() (Dependencies, error) {
	telegramBotAPIKey := os.Getenv("TGBOT_API_KEY")
	openaiAPIKey := os.Getenv("OPENAI_API_KEY")
	if telegramBotAPIKey == "" || openaiAPIKey == "" {
		return Dependencies{}, errors.New("you need to set the API keys environment variables. [TGBOT_API_KEY, OPENAI_API_KEY]")
	}

	d := Dependencies{}
	db, err := sql.Open("sqlite3", "./app.db")
	if err != nil {
		return Dependencies{}, err
	}
	d.db = db

	bot, err := tgbotapi.NewBotAPI(telegramBotAPIKey)
	if err != nil {
		return Dependencies{}, err
	}
	bot.Debug = false
	log.Info().Msgf("Authorized on account %s", bot.Self.UserName)
	d.bot = bot

	client := openai.NewClient(openaiAPIKey)
	d.openaiClient = client

	return d, nil
}

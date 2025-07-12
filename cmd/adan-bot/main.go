package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	ErrMissingToken = errors.New("missing Telegram API token")
)

func Main() error {
	apiToken := strings.Trim(Getenv("TELEGRAM_API_TOKEN", ""), `"`)
	if apiToken == "" {
		return ErrMissingToken
	}

	bot, errNBA := tgbotapi.NewBotAPI(apiToken)
	if errNBA != nil {
		return fmt.Errorf("cannot connect to Telegram API: %w", errNBA)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.Text = handleCommand(update.Message.Command())

		_, err := bot.Send(msg)
		if err != nil {
			log.Printf("Cannot send message to '%v'", update.Message.Chat.UserName)
		}
	}

	return nil
}

func handleCommand(command string) string {
	switch command {
	case "help":
		return "/hola and /status."
	case "hola":
		return "Hola mi nombre es Adan el Bot 🤖 de la comunidad de Golang" +
			" Venezuela. Y como la cancion: <<naci en esta ribera del " +
			"arauca vibrador, soy hermano de la espuma de las garzas de " +
			"las rosas y del sol.>> "
	case "status":
		//nolint:misspell
		return "De momento todo esta bien"
	default:
		return "I don't know that command"
	}
}

func run() {
	err := Profile(Main)
	if err != nil {
		log.Println(err)
		return
	}
}

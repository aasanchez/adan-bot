package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Errors.
var (
	ErrMissingToken = errors.New("missing Telegram API token")
)

func Main() error {
	APITOKEN := strings.Trim(Getenv("TELEGRAM_API_TOKEN", ""), `"`)
	if APITOKEN == "" {
		return ErrMissingToken
	}

	bot, errNBA := tgbotapi.NewBotAPI(APITOKEN)
	if errNBA != nil {
		return fmt.Errorf("cannot connect to Telegram API: %w", errNBA)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			continue
		}

		chat := update.Message.Chat
		msg := tgbotapi.NewMessage(chat.ID, "")

		switch update.Message.Command() {
		case "help":
			msg.Text = "/hola and /status."
		case "hola":
			msg.Text = "Hola mi nombre es Adan el Bot 🤖 de la comunidad de Golang"
			msg.Text += " Venezuela. Y como la cancion: <<naci en esta ribera del "
			msg.Text += "arauca vibrador, soy hermano de la espuma de las garzas de "
			msg.Text += "las rosas y del sol.>> "
		case "status":
			msg.Text = "De momento todo esta bien"
		default:
			msg.Text = "I don't know that command"
		}

		if _, err := bot.Send(msg); err != nil {
			log.Printf("Cannot send message to '%v'", chat.UserName)
		}
	}

	return nil
}

func main() {
	if err := Profile(Main); err != nil {
		log.Fatalln(err)
	}
}

package main

import (
	"fabertoolbox_bot/config"
	"fabertoolbox_bot/internal/bot"
	"gopkg.in/telebot.v4"
	"log"
	"time"
)

func main() {
	cfg := config.LoadConfig()

	pref := telebot.Settings{
		Token:  cfg.Telegram.Token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	telegramBot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatalf("Ошибка создания Telegram бота: %v", err)
	}

	service := &bot.Service{Bot: telegramBot}

	service.RegisterHandlers()

	telegramBot.Start()
}

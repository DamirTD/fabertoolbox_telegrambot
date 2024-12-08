package bot

import "gopkg.in/telebot.v4"

func (s *Service) handleStart(c telebot.Context) error {
	return c.Send("Привет! Я Telegram-бот.")
}

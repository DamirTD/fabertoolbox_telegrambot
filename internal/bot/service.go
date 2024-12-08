package bot

import (
	"gopkg.in/telebot.v4"
	"sync"
)

type Service struct {
	Bot      *telebot.Bot
	userInfo sync.Map
}

func (s *Service) RegisterHandlers() {
	s.Bot.Handle("/start", func(c telebot.Context) error {
		return s.handleStart(c)
	})
}

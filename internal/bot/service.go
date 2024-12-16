package bot

import (
	"fabertoolbox_bot/internal/bot/handlers"
	"fabertoolbox_bot/internal/ui"
	"gopkg.in/telebot.v4"
	"sync"
)

type Service struct {
	Bot      *telebot.Bot
	userInfo sync.Map
	Menu     *telebot.ReplyMarkup
}

func (s *Service) RegisterHandlers() {
	s.Menu = ui.CreateMainMenu()

	s.Bot.Handle("/start", handlers.HandleStartCommand(s.Menu))
	s.Bot.Handle("О каналах", handlers.HandleAboutChannelsButton)
	s.Bot.Handle("Подать заявку на вступление", handlers.HandleJoinRequestButton)
	s.Bot.Handle("Оплатить подписку", handlers.HandleSubscriptionPaymentButton)
	s.Bot.Handle("Задать вопрос", handlers.HandleAskQuestionButton)
}

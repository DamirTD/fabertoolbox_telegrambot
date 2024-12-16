package handlers

import "gopkg.in/telebot.v4"

func HandleStartCommand(menu *telebot.ReplyMarkup) func(c telebot.Context) error {
	return func(c telebot.Context) error {
		return c.Send("Добро пожаловать! Выберите действие из меню ниже:", menu)
	}
}

func HandleAboutChannelsButton(c telebot.Context) error {
	return c.Send("Информация о каналах.")
}

func HandleJoinRequestButton(c telebot.Context) error {
	return c.Send("Чтобы подать заявку, заполните форму.")
}

func HandleSubscriptionPaymentButton(c telebot.Context) error {
	return c.Send("Оплата подписки доступна по следующей ссылке.")
}

func HandleAskQuestionButton(c telebot.Context) error {
	return c.Send("Задайте ваш вопрос.")
}

package ui

import "gopkg.in/telebot.v4"

func CreateMainMenu() *telebot.ReplyMarkup {
	menu := &telebot.ReplyMarkup{}

	btnInfo := menu.Text("О каналах")
	btnJoin := menu.Text("Подать заявку на вступление")
	btnPay := menu.Text("Оплатить подписку")
	btnQuestion := menu.Text("Задать вопрос")

	menu.Reply(
		menu.Row(btnInfo, btnJoin),
		menu.Row(btnPay, btnQuestion),
	)

	return menu
}

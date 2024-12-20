package handlers

import (
	"fabertoolbox_bot/internal/utils"
	"fmt"
	"gopkg.in/telebot.v4"
	"log"
	"path/filepath"
	"regexp"
	"strings"
)

var userSessions = make(map[int64]string)

func HandleStartCommand(menu *telebot.ReplyMarkup) func(c telebot.Context) error {
	return func(c telebot.Context) error {
		return c.Send("Добро пожаловать! Выберите действие из меню ниже:", menu)
	}
}

func HandleAboutChannelsButton(c telebot.Context) error {
	projectRoot := filepath.Join("..", ".")

	textFilePath := filepath.Join(projectRoot, "assets", "text", "about_channels.md")

	text, err := utils.ReadFileContent(textFilePath)
	if err != nil {
		return c.Send("Не удалось загрузить информацию о каналах.")
	}

	if err := c.Send(text); err != nil {
		return err
	}

	videoFilePath := filepath.Join(projectRoot, "assets", "videos", "about_channels_video.mp4")

	video := &telebot.Video{
		File: telebot.FromDisk(videoFilePath),
	}

	return c.Send(video)
}

func HandleJoinRequestButton(c telebot.Context) error {
	projectRoot := filepath.Join("..", ".")

	textFilePath := filepath.Join(projectRoot, "assets", "text", "join_text.md")

	text, err := utils.ReadFileContent(textFilePath)
	if err != nil {
		return c.Send("Не удалось загрузить информацию о каналах.")
	}

	username := c.Sender().FirstName
	if username == "" {
		username = "Пользователь"
	}

	message := strings.ReplaceAll(text, "ИМЯ", username)

	return c.Send(message, &telebot.SendOptions{ParseMode: telebot.ModeMarkdown})
}

func HandleSubscriptionPaymentButton(c telebot.Context) error {
	if userData, exists := userSessions[c.Sender().ID]; !exists || userData == "" {
		return c.Send("Введите свои данные: Фамилия Имя, Никнейк в телеграмме, номер телефона. Пример: Нарегеева Айгуль, aigul, 87078627777")
	}

	err := c.Send("Пожалуйста, отправьте скриншот оплаты.")
	if err != nil {
		return err
	}

	return nil
}

func HandleUserData(c telebot.Context) error {
	userData := c.Message().Text

	const dataFormat = `^[А-Яа-я\s]+, [\w]+, \d{11}$`

	matched, err := regexp.MatchString(dataFormat, userData)
	if err != nil || !matched {
		return c.Send("Пожалуйста, введите данные в правильном формате:\nФамилия Имя, Никнейк в телеграмме, номер телефона. Пример: Нарегеева Айгуль, aigul, 87078627777")
	}

	userSessions[c.Sender().ID] = userData
	return c.Send("Пожалуйста, отправьте скриншот оплаты.")
}

func HandlePaymentScreenshot(c telebot.Context) error {
	// Проверка на наличие фото в сообщении
	if c.Message().Photo == nil {
		return c.Send("Пожалуйста, отправьте скриншот оплаты.")
	}

	// Проверка существования данных пользователя в сессии
	userData, exists := userSessions[c.Sender().ID]
	if !exists {
		return c.Send("Произошла ошибка. Пожалуйста, начните сначала.")
	}

	// ID группы, куда отправляются сообщения
	groupID := int64(-4686187812)
	message := fmt.Sprintf("Новые данные:\n%s\nОт пользователя: %s", userData, c.Sender().FirstName)

	// Генерация уникального идентификатора для сообщения (сокращаем данные для callback)
	messageID := fmt.Sprintf("%d_%d", c.Sender().ID, c.Message().ID) // Используем ID пользователя и ID сообщения

	// Кнопки для одобрения и отклонения
	approveButton := telebot.InlineButton{
		Text: "Одобрить",
		Data: "approve_" + messageID, // Привязка уникального идентификатора к кнопке
	}
	declineButton := telebot.InlineButton{
		Text: "Отклонить",
		Data: "decline_" + messageID, // Привязка уникального идентификатора к кнопке
	}

	// Создание клавиатуры с кнопками
	inlineKeyboard := &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{approveButton, declineButton},
		},
	}

	// Отправка сообщения в группу с кнопками
	_, err := c.Bot().Send(&telebot.Chat{ID: groupID}, message, inlineKeyboard)
	if err != nil {
		return err
	}

	// Сохранение уникального идентификатора сообщения в сессии
	userSessions[c.Sender().ID] = messageID

	// Пересылка фото в группу
	_, err = c.Bot().Forward(&telebot.Chat{ID: groupID}, c.Message())
	if err != nil {
		return err
	}

	// Информирование пользователя о том, что его запрос был отправлен
	return c.Send("Ваш запрос был отправлен в группу для проверки.")
}

func HandleApprovalButtons(c telebot.Context) error {
	action := c.Data()
	log.Printf("Received callback data: %s", action)

	parts := strings.Split(action, "_")
	if len(parts) < 2 {
		return c.Send("Произошла ошибка, невозможно обработать действие.")
	}

	actionType := parts[0]
	messageID := strings.Join(parts[1:], "_")

	log.Printf("Action type: %s, Message ID: %s", actionType, messageID)

	userData, exists := userSessions[c.Sender().ID]
	if !exists || userData != messageID {
		return c.Send("Произошла ошибка, данные не найдены. Пожалуйста, начните сначала.")
	}

	var message string
	if actionType == "approve" {
		message = fmt.Sprintf(`Привет, %s! Ваша оплата за подписку зачислена!
Дата оплаты: 04.11.2024
Кол-во оплаченных периодов: 1
Дата следующего платежа: 04.12.2024
Продлить подписку вы можете в течение 5 дней с даты следующего платежа.

Спасибо, что вы с нами!`, c.Sender().FirstName)
	} else if actionType == "decline" {
		message = "Отклонено. Пожалуйста, свяжитесь с администратором."
	}

	_, err := c.Bot().Send(&telebot.User{ID: c.Sender().ID}, message)

	if err != nil {
		log.Printf("Error sending response message: %v", err)
		return err
	}

	return nil
}

func HandleAskQuestionButton(c telebot.Context) error {
	return c.Send("Задайте ваш вопрос.")
}

package handlers

import (
	"fabertoolbox_bot/internal/utils"
	"fabertoolbox_bot/internal/variables"
	"fmt"
	"gopkg.in/telebot.v4"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

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
	// Проверка существования данных в сессии
	if userData, exists := variables.UserSessions[c.Sender().ID]; !exists || userData == "" {
		// Если данные не введены, попросить ввести их
		return c.Send("Введите свои данные: Фамилия Имя, Никнейк в телеграмме, номер телефона и периоды через запятую. Пример: Нарегеева Айгуль, aigul, 87078627777 и периоды: 3, 4.")
	} else {
		// Если данные существуют, валидируем их
		const dataAndPeriodsFormat = `^[А-Яа-яёЁ\s]+, [\w]+, \d{11} и периоды: \d{1,2}(?:,\s*\d{1,2}){0,18}$`
		matched, err := regexp.MatchString(dataAndPeriodsFormat, userData)
		if err != nil || !matched {
			return c.Send("Пожалуйста, введите данные в правильном формате:\nФамилия Имя, Никнейк в телеграмме, номер телефона и периоды через запятую. Пример: Нарегеева Айгуль, aigul, 87078627777 и периоды: 3, 4.")
		}
	}

	// Если все валидно, просим отправить скриншот
	return c.Send("Пожалуйста, отправьте скриншот оплаты.")
}

func HandleUserData(c telebot.Context) error {
	userDataAndPeriods := strings.TrimSpace(c.Message().Text)

	// Регулярное выражение для проверки формата (периоды от 1 до 19)
	const dataAndPeriodsFormat = `^[А-Яа-яёЁ\s]+, [\w]+, \d{11} и периоды: \d{1,2}(?:,\s*\d{1,2}){0,18}$`
	matched, err := regexp.MatchString(dataAndPeriodsFormat, userDataAndPeriods)
	if err != nil || !matched {
		return c.Send("Пожалуйста, введите данные в правильном формате:\nФамилия Имя, Никнейк в телеграмме, номер телефона и периоды через запятую. Пример: Нарегеева Айгуль, aigul, 87078627777 и периоды: 3, 4.")
	}

	parts := strings.Split(userDataAndPeriods, " и периоды: ")
	if len(parts) != 2 {
		return c.Send("Ошибка: Неверный формат ввода. Пожалуйста, укажите данные и периоды в указанном формате.")
	}

	userData := parts[0]
	periodsText := parts[1]

	variables.UserSessionsMutex.Lock()
	variables.UserSessions[c.Sender().ID] = fmt.Sprintf("%s|%s", userData, periodsText)
	variables.UserSessionsMutex.Unlock()

	log.Printf("Данные пользователя и периоды сохранены: User ID: %d, Data: %s", c.Sender().ID, userData)

	// Подтверждение успешного ввода
	return c.Send("Данные и периоды успешно сохранены. Пожалуйста, отправьте скриншот оплаты.")
}
func HandlePaymentScreenshot(c telebot.Context) error {
	if c.Message().Photo == nil {
		return c.Send("Пожалуйста, отправьте скриншот оплаты.")
	}

	userID := c.Sender().ID
	variables.UserSessionsMutex.Lock()
	defer variables.UserSessionsMutex.Unlock()

	userData, exists := variables.UserSessions[userID]
	if !exists || !strings.Contains(userData, "|") {
		return c.Send("Произошла ошибка с данными. Пожалуйста, начните сначала.")
	}

	parts := strings.Split(userData, "|")
	if len(parts) != 2 {
		return c.Send("Произошла ошибка с данными. Пожалуйста, начните сначала.")
	}

	userInfo := parts[0]
	periods := parts[1]
	message := fmt.Sprintf("Новые данные:\n%s\nОплаченные периоды: %s\nОт пользователя: %s", userInfo, periods, c.Sender().FirstName)

	groupID := int64(-4686187812)
	approveButton := telebot.InlineButton{Text: "Одобрить", Data: fmt.Sprintf("approve_%d", userID)}
	declineButton := telebot.InlineButton{Text: "Отклонить", Data: fmt.Sprintf("decline_%d", userID)}
	inlineKeyboard := &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{approveButton, declineButton},
		},
	}

	_, err := c.Bot().Forward(&telebot.Chat{ID: groupID}, c.Message())
	if err != nil {
		return c.Send("Произошла ошибка при пересылке скриншота.")
	}

	// После этого отправляем сообщение с данными
	_, err = c.Bot().Send(&telebot.Chat{ID: groupID}, message, inlineKeyboard)
	if err != nil {
		return c.Send("Произошла ошибка при отправке данных в группу.")
	}

	return c.Send("Ваш запрос был отправлен в группу для проверки.")
}
func HandleApprovalButtons(c telebot.Context) error {
	action := c.Data()
	parts := strings.Split(action, "_")
	if len(parts) != 2 {
		return c.Send("Произошла ошибка, невозможно обработать действие.")
	}

	actionType := parts[0]
	userIDStr := parts[1]
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.Send("Произошла ошибка, некорректный идентификатор пользователя.")
	}

	userData, exists := variables.UserSessions[userID]
	if !exists || !strings.Contains(userData, "|") {
		return c.Send("Произошла ошибка, данные не найдены. Пожалуйста, начните сначала.")
	}

	parts = strings.Split(userData, "|")
	if len(parts) != 2 {
		return c.Send("Произошла ошибка с данными. Пожалуйста, начните сначала.")
	}
	periods := parts[1]

	// Получаем массив чисел из строки периодов
	periodsList := parsePeriods(periods)

	// Находим максимальный период
	maxPeriod := getMaxPeriod(periodsList)

	// Находим дату следующего платежа для максимального периода
	nextPaymentDate := getNextPaymentDate(maxPeriod)

	loc, err := time.LoadLocation("Asia/Almaty")
	currentDate := time.Now().In(loc).Format("02.01.2006 15:04:05")

	var message string
	if actionType == "approve" {
		message = fmt.Sprintf(`Привет, %s! Ваша оплата за подписку зачислена!
Дата оплаты: %s
Оплаченные периоды: %s
Дата следующего платежа: %s
Продлить подписку вы можете в течение 5 дней с даты следующего платежа.

Спасибо, что вы с нами!`, c.Sender().FirstName, currentDate, periods, nextPaymentDate)
	} else if actionType == "decline" {
		message = "Отклонено. Пожалуйста, свяжитесь с администратором."
	} else {
		return c.Send("Произошла ошибка, неизвестное действие.")
	}

	_, err = c.Bot().Send(&telebot.User{ID: userID}, message)
	if err != nil {
		log.Printf("Ошибка при отправке сообщения: %v", err)
		return err
	}

	return nil
}
func parsePeriods(periods string) []int {
	var periodsList []int
	periodsArray := strings.Split(periods, ", ")
	for _, p := range periodsArray {
		period, err := strconv.Atoi(p)
		if err == nil {
			periodsList = append(periodsList, period)
		}
	}
	return periodsList
}
func getMaxPeriod(periods []int) int {
	maxPeriod := 0
	for _, period := range periods {
		if period > maxPeriod {
			maxPeriod = period
		}
	}
	return maxPeriod
}
func getNextPaymentDate(period int) string {
	if period < 19 {
		nextPeriod := period + 1
		nextPeriodRange := variables.Periods[nextPeriod]
		dateRangeParts := strings.Split(nextPeriodRange, " - ")
		nextPaymentDate := dateRangeParts[0]
		return nextPaymentDate
	}
	return variables.Periods[1][:10]
}

func HandleAskQuestionButton(c telebot.Context) error {
	return c.Send("Задайте ваш вопрос.")
}

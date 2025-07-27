package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type TelegramBot struct {
	bot      *tgbotapi.BotAPI
	database *Database
	email    *EmailService
	logger   *logrus.Logger
	ctx      context.Context
	cancel   context.CancelFunc
}

type UserState struct {
	UserID   int64
	State    string
	Data     map[string]string
	LastSeen time.Time
}

var userStates = make(map[int64]*UserState)

func NewTelegramBot(db *Database, email *EmailService, logger *logrus.Logger) (*TelegramBot, error) {
	token := getEnv("TELEGRAM_BOT_TOKEN", "")
	if token == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &TelegramBot{
		bot:      bot,
		database: db,
		email:    email,
		logger:   logger,
		ctx:      ctx,
		cancel:   cancel,
	}, nil
}

func (t *TelegramBot) Start() error {
	t.logger.Info("Starting Telegram bot...")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.bot.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			go t.handleUpdate(update)
		case <-t.ctx.Done():
			return nil
		}
	}
}

func (t *TelegramBot) Stop() error {
	t.cancel()
	return nil
}

func (t *TelegramBot) handleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	message := update.Message
	userID := message.From.ID

	// Получаем или создаем состояние пользователя
	state, exists := userStates[userID]
	if !exists {
		state = &UserState{
			UserID: userID,
			State:  "start",
			Data:   make(map[string]string),
		}
		userStates[userID] = state
	}

	state.LastSeen = time.Now()

	// Обрабатываем команды
	if message.IsCommand() {
		t.handleCommand(message, state)
		return
	}

	// Обрабатываем сообщения в зависимости от состояния
	switch state.State {
	case "start":
		t.handleStart(message, state)
	case "waiting_for_type":
		t.handleTypeSelection(message, state)
	case "waiting_for_message":
		t.handleMessageInput(message, state)
	default:
		t.sendMessage(message.Chat.ID, "Пожалуйста, используйте /start для начала работы с ботом.")
	}
}

func (t *TelegramBot) handleCommand(message *tgbotapi.Message, state *UserState) {
	switch message.Command() {
	case "start":
		t.handleStart(message, state)
	case "help":
		t.sendHelp(message.Chat.ID)
	case "stats":
		t.handleStats(message.Chat.ID)
	default:
		t.sendMessage(message.Chat.ID, "Неизвестная команда. Используйте /help для получения справки.")
	}
}

func (t *TelegramBot) handleStart(message *tgbotapi.Message, state *UserState) {
	state.State = "waiting_for_type"
	state.Data = make(map[string]string)

	text := `🏥 Добро пожаловать в систему обратной связи больницы!

Пожалуйста, выберите тип вашего обращения:

1️⃣ Жалоба - если у вас есть претензии к качеству обслуживания
2️⃣ Отзыв - если вы хотите поделиться положительными впечатлениями

Выберите номер (1 или 2):`

	t.sendMessage(message.Chat.ID, text)
}

func (t *TelegramBot) handleTypeSelection(message *tgbotapi.Message, state *UserState) {
	text := strings.ToLower(strings.TrimSpace(message.Text))

	var feedbackType string
	switch text {
	case "1", "жалоба", "complaint":
		feedbackType = "complaint"
		state.Data["type"] = feedbackType
		state.State = "waiting_for_message"
		t.sendMessage(message.Chat.ID, "📝 Пожалуйста, опишите вашу жалобу подробно. Мы рассмотрим её в кратчайшие сроки.")
	case "2", "отзыв", "review":
		feedbackType = "review"
		state.Data["type"] = feedbackType
		state.State = "waiting_for_message"
		t.sendMessage(message.Chat.ID, "📝 Пожалуйста, поделитесь вашими впечатлениями о работе больницы.")
	default:
		t.sendMessage(message.Chat.ID, "Пожалуйста, выберите 1 (жалоба) или 2 (отзыв).")
	}
}

func (t *TelegramBot) handleMessageInput(message *tgbotapi.Message, state *UserState) {
	feedbackType := state.Data["type"]

	feedback := &Feedback{
		UserID:    message.From.ID,
		Username:  message.From.UserName,
		FirstName: message.From.FirstName,
		LastName:  message.From.LastName,
		Message:   message.Text,
		Type:      feedbackType,
		Status:    "new",
	}

	// Сохраняем в базу данных
	if err := t.database.SaveFeedback(feedback); err != nil {
		t.logger.Error("Failed to save feedback: ", err)
		t.sendMessage(message.Chat.ID, "❌ Произошла ошибка при сохранении вашего обращения. Пожалуйста, попробуйте позже.")
		return
	}

	// Отправляем на email
	go func() {
		if err := t.email.SendFeedbackEmail(feedback); err != nil {
			t.logger.Error("Failed to send email: ", err)
		}
	}()

	// Обновляем статус
	t.database.UpdateFeedbackStatus(feedback.ID, "sent")

	// Отправляем подтверждение пользователю
	responseText := "✅ Спасибо! Ваше обращение успешно отправлено."
	if feedbackType == "complaint" {
		responseText += "\n\nМы рассмотрим вашу жалобу и примем необходимые меры."
	} else {
		responseText += "\n\nВаш отзыв очень важен для нас!"
	}

	t.sendMessage(message.Chat.ID, responseText)

	// Сбрасываем состояние
	state.State = "start"
	state.Data = make(map[string]string)
}

func (t *TelegramBot) sendHelp(chatID int64) {
	helpText := `🤖 Справка по использованию бота:

/start - Начать работу с ботом
/help - Показать эту справку
/stats - Статистика обращений (только для администраторов)

Бот предназначен для сбора жалоб и отзывов о работе больницы.
Все обращения обрабатываются в кратчайшие сроки.`

	t.sendMessage(chatID, helpText)
}

func (t *TelegramBot) handleStats(chatID int64) {
	// Проверяем, является ли пользователь администратором
	adminUserID := getEnv("ADMIN_USER_ID", "")
	if adminUserID == "" || fmt.Sprintf("%d", chatID) != adminUserID {
		t.sendMessage(chatID, "⛔ У вас нет прав для просмотра статистики.")
		return
	}

	stats, err := t.database.GetFeedbackStats()
	if err != nil {
		t.logger.Error("Failed to get stats: ", err)
		t.sendMessage(chatID, "❌ Ошибка при получении статистики.")
		return
	}

	statsText := "📊 Статистика обращений:\n\n"
	statsText += fmt.Sprintf("Жалобы: %d\n", stats["complaint"])
	statsText += fmt.Sprintf("Отзывы: %d\n", stats["review"])
	statsText += fmt.Sprintf("Всего: %d", stats["complaint"]+stats["review"])

	t.sendMessage(chatID, statsText)
}

func (t *TelegramBot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"

	if _, err := t.bot.Send(msg); err != nil {
		t.logger.Error("Failed to send message: ", err)
	}
}

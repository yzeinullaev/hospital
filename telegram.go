package main

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type UserState struct {
	State string
	Data  map[string]string
}

type TelegramBot struct {
	bot      *tgbotapi.BotAPI
	database *Database
	email    *EmailService
	logger   *logrus.Logger
	users    map[int64]*UserState
}

func NewTelegramBot(database *Database, email *EmailService, logger *logrus.Logger) (*TelegramBot, error) {
	token := getEnv("TELEGRAM_BOT_TOKEN", "")
	if token == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is not set")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	return &TelegramBot{
		bot:      bot,
		database: database,
		email:    email,
		logger:   logger,
		users:    make(map[int64]*UserState),
	}, nil
}

func (t *TelegramBot) Start() error {
	t.logger.Info("Bot started: @", t.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			t.handleMessage(update.Message)
		} else if update.CallbackQuery != nil {
			t.handleCallbackQuery(update.CallbackQuery)
		}
	}

	return nil
}

func (t *TelegramBot) Stop() error {
	return nil
}

func (t *TelegramBot) handleMessage(message *tgbotapi.Message) {
	userID := message.From.ID
	state, exists := t.users[userID]
	if !exists {
		state = &UserState{
			State: "start",
			Data:  make(map[string]string),
		}
		t.users[userID] = state
	}

	// Обрабатываем команды
	if message.IsCommand() {
		t.handleCommand(message, state)
		return
	}

	// Обрабатываем текст в зависимости от состояния
	switch state.State {
	case "waiting_for_type":
		t.handleTypeSelection(message, state)
	case "waiting_for_message":
		t.handleMessageInput(message, state)
	default:
		t.sendMainMenu(message.Chat.ID, "Выберите действие:")
	}
}

func (t *TelegramBot) handleCommand(message *tgbotapi.Message, state *UserState) {
	switch message.Command() {
	case "start":
		state.State = "start"
		state.Data = make(map[string]string)
		t.sendMainMenu(message.Chat.ID, "Добро пожаловать! Выберите действие:")
	case "menu":
		t.sendMainMenu(message.Chat.ID, "Главное меню:")
	default:
		t.sendMainMenu(message.Chat.ID, "Используйте /start для начала работы")
	}
}

func (t *TelegramBot) handleCallbackQuery(callback *tgbotapi.CallbackQuery) {
	userID := callback.From.ID
	state, exists := t.users[userID]
	if !exists {
		state = &UserState{
			State: "start",
			Data:  make(map[string]string),
		}
		t.users[userID] = state
	}

	data := callback.Data
	switch data {
	case "complaint":
		state.State = "waiting_for_message"
		state.Data["type"] = "complaint"
		t.sendMessage(callback.Message.Chat.ID, "Пожалуйста, напишите вашу жалобу:")
	case "review":
		state.State = "waiting_for_message"
		state.Data["type"] = "review"
		t.sendMessage(callback.Message.Chat.ID, "Пожалуйста, напишите ваш отзыв:")
	default:
		t.sendMainMenu(callback.Message.Chat.ID, "Выберите действие:")
	}
}

func (t *TelegramBot) handleTypeSelection(message *tgbotapi.Message, state *UserState) {
	text := strings.ToLower(strings.TrimSpace(message.Text))

	switch text {
	case "жалоба", "complaint":
		state.State = "waiting_for_message"
		state.Data["type"] = "complaint"
		t.sendMessage(message.Chat.ID, "Пожалуйста, напишите вашу жалобу:")
	case "отзыв", "review":
		state.State = "waiting_for_message"
		state.Data["type"] = "review"
		t.sendMessage(message.Chat.ID, "Пожалуйста, напишите ваш отзыв:")
	default:
		t.sendMessage(message.Chat.ID, "Пожалуйста, выберите 'жалоба' или 'отзыв'")
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
		t.sendMessage(message.Chat.ID, "Произошла ошибка при сохранении. Попробуйте позже.")
		return
	}

	// Отправляем email
	if err := t.email.SendFeedbackEmail(feedback); err != nil {
		t.logger.Error("Failed to send email: ", err)
	}

	// Отправляем подтверждение пользователю
	responseText := fmt.Sprintf("✅ Ваш %s успешно отправлен!\n\nТекст: %s",
		getTypeDisplayName(feedbackType), message.Text)

	t.sendMessage(message.Chat.ID, responseText)

	// Сбрасываем состояние
	state.State = "start"
	state.Data = make(map[string]string)

	// Показываем главное меню
	t.sendMainMenu(message.Chat.ID, "Хотите отправить еще один отзыв или жалобу?")
}

func (t *TelegramBot) sendMainMenu(chatID int64, text string) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Жалоба", "complaint"),
			tgbotapi.NewInlineKeyboardButtonData("⭐ Отзыв", "review"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	t.bot.Send(msg)
}

func (t *TelegramBot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	t.bot.Send(msg)
}

func getTypeDisplayName(feedbackType string) string {
	switch feedbackType {
	case "complaint":
		return "жалоба"
	case "review":
		return "отзыв"
	default:
		return feedbackType
	}
}

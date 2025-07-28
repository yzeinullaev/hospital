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
	if update.Message == nil && update.CallbackQuery == nil {
		return
	}

	// Обрабатываем callback query (нажатие на кнопки)
	if update.CallbackQuery != nil {
		t.handleCallbackQuery(update.CallbackQuery)
		return
	}

	// Обрабатываем обычные сообщения
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
		t.sendMessageWithMenu(message.Chat.ID, "Пожалуйста, используйте /start для начала работы с ботом.")
	}
}

func (t *TelegramBot) handleCallbackQuery(callback *tgbotapi.CallbackQuery) {
	userID := callback.From.ID
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

	switch callback.Data {
	case "complaint":
		state.Data["type"] = "complaint"
		state.State = "waiting_for_message"
		t.sendMessage(callback.Message.Chat.ID, "📝 Пожалуйста, опишите вашу жалобу подробно. Мы рассмотрим её в кратчайшие сроки.")
	case "review":
		state.Data["type"] = "review"
		state.State = "waiting_for_message"
		t.sendMessage(callback.Message.Chat.ID, "📝 Пожалуйста, поделитесь вашими впечатлениями о работе больницы.")
	case "help":
		t.sendHelp(callback.Message.Chat.ID)
	case "stats":
		t.handleStats(callback.Message.Chat.ID)
	case "main_menu":
		t.sendMainMenu(callback.Message.Chat.ID)
	default:
		t.logger.Warnf("Unknown callback data: %s", callback.Data)
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
	case "menu":
		t.sendMainMenu(message.Chat.ID)
	default:
		t.sendMessage(message.Chat.ID, "Неизвестная команда. Используйте /help для получения справки.")
	}
}

func (t *TelegramBot) handleStart(message *tgbotapi.Message, state *UserState) {
	state.State = "waiting_for_type"
	state.Data = make(map[string]string)

	text := `🏥 Добро пожаловать в систему обратной связи больницы!

Пожалуйста, выберите тип вашего обращения:`

	t.sendMessageWithTypeButtons(message.Chat.ID, text)
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
		t.sendMessageWithTypeButtons(message.Chat.ID, "Пожалуйста, выберите тип обращения, используя кнопки ниже:")
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

	// Обрабатываем медиафайлы
	if message.Photo != nil && len(message.Photo) > 0 {
		// Берем последнее (самое большое) фото
		photo := message.Photo[len(message.Photo)-1]
		mediaFile := MediaFile{
			FileID:   photo.FileID,
			FileType: "photo",
			FileSize: int64(photo.FileSize),
		}
		feedback.MediaFiles = append(feedback.MediaFiles, mediaFile)
	} else if message.Video != nil {
		mediaFile := MediaFile{
			FileID:   message.Video.FileID,
			FileType: "video",
			FileName: message.Video.FileName,
			FileSize: int64(message.Video.FileSize),
			MimeType: message.Video.MimeType,
		}
		feedback.MediaFiles = append(feedback.MediaFiles, mediaFile)
	} else if message.Document != nil {
		mediaFile := MediaFile{
			FileID:   message.Document.FileID,
			FileType: "document",
			FileName: message.Document.FileName,
			FileSize: int64(message.Document.FileSize),
			MimeType: message.Document.MimeType,
		}
		feedback.MediaFiles = append(feedback.MediaFiles, mediaFile)
	} else if message.Audio != nil {
		mediaFile := MediaFile{
			FileID:   message.Audio.FileID,
			FileType: "audio",
			FileName: message.Audio.FileName,
			FileSize: int64(message.Audio.FileSize),
			MimeType: message.Audio.MimeType,
		}
		feedback.MediaFiles = append(feedback.MediaFiles, mediaFile)
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

	// Добавляем информацию о медиафайлах
	if len(feedback.MediaFiles) > 0 {
		responseText += fmt.Sprintf("\n\n📎 Прикреплено файлов: %d", len(feedback.MediaFiles))
	}

	responseText += "\n\nХотите отправить еще одно обращение?"

	t.sendMessageWithMainMenu(message.Chat.ID, responseText)

	// Сбрасываем состояние
	state.State = "start"
	state.Data = make(map[string]string)
}

func (t *TelegramBot) sendMessageWithTypeButtons(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"

	// Создаем inline кнопки
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Жалоба", "complaint"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⭐ Отзыв", "review"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❓ Помощь", "help"),
		),
	)

	msg.ReplyMarkup = keyboard

	if _, err := t.bot.Send(msg); err != nil {
		t.logger.Error("Failed to send message with buttons: ", err)
	}
}

func (t *TelegramBot) sendMessageWithMainMenu(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"

	// Создаем inline кнопки для главного меню
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏥 Новое обращение", "main_menu"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❓ Помощь", "help"),
		),
	)

	msg.ReplyMarkup = keyboard

	if _, err := t.bot.Send(msg); err != nil {
		t.logger.Error("Failed to send message with main menu: ", err)
	}
}

func (t *TelegramBot) sendMainMenu(chatID int64) {
	text := `🏥 Главное меню системы обратной связи больницы

Выберите действие:`

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"

	// Создаем inline кнопки для главного меню
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Отправить жалобу", "complaint"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⭐ Оставить отзыв", "review"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❓ Помощь", "help"),
		),
	)

	msg.ReplyMarkup = keyboard

	if _, err := t.bot.Send(msg); err != nil {
		t.logger.Error("Failed to send main menu: ", err)
	}
}

func (t *TelegramBot) sendHelp(chatID int64) {
	helpText := `🤖 Справка по использованию бота:

📝 <b>Отправка обращений:</b>
• Жалоба - для претензий к качеству обслуживания
• Отзыв - для положительных впечатлений

📋 <b>Команды:</b>
/start - Начать работу с ботом
/help - Показать эту справку
/menu - Главное меню
/stats - Статистика обращений (только для администраторов)

💡 <b>Как использовать:</b>
1. Выберите тип обращения
2. Опишите вашу проблему или впечатление
3. Отправьте сообщение

✅ Все обращения обрабатываются в кратчайшие сроки.`

	msg := tgbotapi.NewMessage(chatID, helpText)
	msg.ParseMode = "HTML"

	// Добавляем кнопку возврата в главное меню
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 Главное меню", "main_menu"),
		),
	)

	msg.ReplyMarkup = keyboard

	if _, err := t.bot.Send(msg); err != nil {
		t.logger.Error("Failed to send help: ", err)
	}
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

	statsText := "📊 <b>Статистика обращений:</b>\n\n"
	statsText += fmt.Sprintf("📝 Жалобы: %d\n", stats["complaint"])
	statsText += fmt.Sprintf("⭐ Отзывы: %d\n", stats["review"])
	statsText += fmt.Sprintf("📈 Всего: %d", stats["complaint"]+stats["review"])

	msg := tgbotapi.NewMessage(chatID, statsText)
	msg.ParseMode = "HTML"

	// Добавляем кнопку возврата в главное меню
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 Главное меню", "main_menu"),
		),
	)

	msg.ReplyMarkup = keyboard

	if _, err := t.bot.Send(msg); err != nil {
		t.logger.Error("Failed to send stats: ", err)
	}
}

func (t *TelegramBot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"

	if _, err := t.bot.Send(msg); err != nil {
		t.logger.Error("Failed to send message: ", err)
	}
}

func (t *TelegramBot) sendMessageWithMenu(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"

	// Добавляем кнопку возврата в главное меню
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 Главное меню", "main_menu"),
		),
	)

	msg.ReplyMarkup = keyboard

	if _, err := t.bot.Send(msg); err != nil {
		t.logger.Error("Failed to send message with menu: ", err)
	}
}

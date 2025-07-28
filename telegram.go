package main

import (
	"fmt"
	"strings"
	"time"

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
		t.sendMainMenu(message.Chat.ID, "Әрекетті таңдаңыз:")
	}
}

func (t *TelegramBot) handleCommand(message *tgbotapi.Message, state *UserState) {
	switch message.Command() {
	case "start":
		state.State = "start"
		state.Data = make(map[string]string)
		t.sendMainMenu(message.Chat.ID, "Қош келдіңіз! Әрекетті таңдаңыз:")
	case "menu":
		t.sendMainMenu(message.Chat.ID, "Басты мәзір:")
	case "stats":
		if t.isAdmin(message.From.ID) {
			t.handleStats(message.Chat.ID)
		} else {
			t.sendMessage(message.Chat.ID, "❌ Сізде статистикаға қолжетімділік жоқ")
		}
	default:
		t.sendMainMenu(message.Chat.ID, "Жұмысты бастау үшін /start пәрменін пайдаланыңыз")
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
		t.sendMessage(callback.Message.Chat.ID, "📝 Өтініш, шағымыңызды толық сипаттаңыз. Біз оны мүмкіндігінше қысқа мерзімде қарастырамыз.")
	case "review":
		state.State = "waiting_for_message"
		state.Data["type"] = "review"
		t.sendMessage(callback.Message.Chat.ID, "⭐ Өтініш, пікіріңізді толық сипаттаңыз. Біз сіздің пікіріңізді жоғары бағалаймыз.")
	case "stats":
		if t.isAdmin(userID) {
			t.handleStats(callback.Message.Chat.ID)
		} else {
			t.sendMessage(callback.Message.Chat.ID, "❌Сізде статистикаға қолжетімділік жоқ")
		}
	case "new_request":
		state.State = "start"
		state.Data = make(map[string]string)
		t.sendMainMenu(callback.Message.Chat.ID, "")
	case "help":
		t.sendHelp(callback.Message.Chat.ID)
	case "back_to_menu":
		state.State = "start"
		state.Data = make(map[string]string)
		t.sendMainMenu(callback.Message.Chat.ID, "")
	default:
		t.sendMainMenu(callback.Message.Chat.ID, "")
	}
}

func (t *TelegramBot) handleTypeSelection(message *tgbotapi.Message, state *UserState) {
	text := strings.ToLower(strings.TrimSpace(message.Text))

	switch text {
	case "жалоба", "complaint":
		state.State = "waiting_for_message"
		state.Data["type"] = "complaint"
		t.sendMessage(message.Chat.ID, "📝 Өтініш, шағымыңызды толық сипаттаңыз. Біз оны мүмкіндігінше қысқа мерзімде қарастырамыз.")
	case "отзыв", "review":
		state.State = "waiting_for_message"
		state.Data["type"] = "review"
		t.sendMessage(message.Chat.ID, "⭐ Өтініш, пікіріңізді толық сипаттаңыз. Біз сіздің пікіріңізді жоғары бағалаймыз.")
	default:
		t.sendMessage(message.Chat.ID, "Өтініш, ‘шағым’ немесе ‘пікір’ таңдаңыз")
	}
}

func (t *TelegramBot) handleMessageInput(message *tgbotapi.Message, state *UserState) {
	feedbackType := state.Data["type"]

	// Используем правильный часовой пояс
	timezone := getEnv("TIMEZONE", "Asia/Almaty")

	// Используем фиксированное смещение для Asia/Almaty (UTC+5)
	var currentTime time.Time
	if timezone == "Asia/Almaty" {
		// Создаем фиксированное смещение UTC+5
		loc := time.FixedZone("Asia/Almaty", 5*60*60) // +5 часов в секундах
		currentTime = time.Now().In(loc)
	} else {
		// Пытаемся загрузить часовой пояс
		loc, err := time.LoadLocation(timezone)
		if err != nil {
			// Если не удалось загрузить часовой пояс, используем UTC
			loc = time.UTC
		}
		currentTime = time.Now().In(loc)
	}

	feedback := &Feedback{
		UserID:    message.From.ID,
		Username:  message.From.UserName,
		FirstName: message.From.FirstName,
		LastName:  message.From.LastName,
		Message:   message.Text,
		Type:      feedbackType,
		Status:    "new",
		CreatedAt: currentTime, // Устанавливаем время в правильном часовом поясе
	}

	// Сохраняем в базу данных
	if err := t.database.SaveFeedback(feedback); err != nil {
		t.logger.Error("Failed to save feedback: ", err)
		t.sendMessage(message.Chat.ID, " Сақтау кезінде қате орын алды. Кейінірек қайталап көріңіз.")
		return
	}

	// Отправляем email
	if err := t.email.SendFeedbackEmail(feedback); err != nil {
		t.logger.Error("Failed to send email: ", err)
	}

	// Отправляем подтверждение пользователю с кнопками
	responseText := fmt.Sprintf("✅ Сіздің %s сәтті жіберілді!\n\nБіз сіздің %s қарап, қажетті шараларды қабылдаймыз.\n\nТағы бір өтініш жібергіңіз келе ме?",
		getTypeDisplayName(feedbackType), getTypeDisplayName(feedbackType))

	t.sendConfirmationMenu(message.Chat.ID, responseText)

	// Сбрасываем состояние
	state.State = "start"
	state.Data = make(map[string]string)
}

func (t *TelegramBot) sendMainMenu(chatID int64, text string) {
	// Проверяем, является ли пользователь администратором
	isAdmin := t.isAdmin(chatID)

	// Заголовок меню
	menuText := "🏥 Аурухананың кері байланыс жүйесінің басты мәзірі\n\nӘрекетті таңдаңыз:"

	var keyboard tgbotapi.InlineKeyboardMarkup

	if isAdmin {
		// Меню для администратора с кнопкой статистики
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("📝 Шағым жіберу", "complaint"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("⭐ Пікір қалдыру", "review"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("❓ Көмек", "help"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("📊 Статистика", "stats"),
			),
		)
	} else {
		// Меню для обычных пользователей без статистики
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("📝 Шағым жіберу", "complaint"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("⭐ Пікір қалдыру", "review"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("❓ Көмек", "help"),
			),
		)
	}

	msg := tgbotapi.NewMessage(chatID, menuText)
	msg.ReplyMarkup = keyboard
	t.bot.Send(msg)
}

func (t *TelegramBot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	t.bot.Send(msg)
}

func (t *TelegramBot) handleStats(chatID int64) {
	stats, err := t.database.GetFeedbackStats()
	if err != nil {
		t.logger.Error("Failed to get stats: ", err)
		t.sendMessage(chatID, "❌ Статистиканы алу кезінде қате орын алды")
		return
	}

	complaints := stats["complaint"]
	reviews := stats["review"]
	total := complaints + reviews

	statsText := fmt.Sprintf("📊 Өтініштер статистикасы\n\n"+
		"📝 Жалобы: %d\n"+
		"⭐ Отзывы: %d\n"+
		"📈 Всего: %d", complaints, reviews, total)

	// Отправляем статистику с кнопкой возврата в главное меню
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 Басты мәзір", "back_to_menu"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, statsText)
	msg.ReplyMarkup = keyboard
	t.bot.Send(msg)
}

func (t *TelegramBot) sendConfirmationMenu(chatID int64, text string) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏥 Жаңа өтініш", "new_request"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❓ Көмек", "help"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	t.bot.Send(msg)
}

func (t *TelegramBot) sendHelp(chatID int64) {
	helpText := `ℹ️ Көмек

📝 Шағым жіберу үшін:
1. "📝 Шағым жіберу" батырмасын шертіңіз
2. Шағымыңызды толық сипаттаңыз
3. Хабарламаны жіберіңіз

⭐ Пікірді қалай қалдыруға болады:
1. ⭐ Пікір қалдыру” батырмасын басыңыз
2. Пікіріңізді толық сипаттаңыз
3. Хабарламаны жіберіңіз

📧 Сіздің өтінішіңіз әкімшілікке email арқылы жіберіледі..

🔙 Басты мәзірге оралу үшін /start немесе /menu пәрменін пайдаланыңыз`

	// Отправляем помощь с кнопкой возврата в главное меню
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 Басты мәзір", "back_to_menu"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, helpText)
	msg.ReplyMarkup = keyboard
	t.bot.Send(msg)
}

func getTypeDisplayName(feedbackType string) string {
	switch feedbackType {
	case "complaint":
		return "шағым жіберу"
	case "review":
		return "Пікір қалдыру"
	default:
		return feedbackType
	}
}

func (t *TelegramBot) isAdmin(userID int64) bool {
	adminID := getEnvAsInt("ADMIN_USER_ID", 0)
	return userID == int64(adminID)
}

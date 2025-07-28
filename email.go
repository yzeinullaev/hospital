package main

import (
	"fmt"
	"time"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	fromEmail    string
	fromPassword string
	toEmail      string
	smtpHost     string
	smtpPort     int
}

func NewEmailService() *EmailService {
	return &EmailService{
		fromEmail:    getEnv("EMAIL_FROM", ""),
		fromPassword: getEnv("EMAIL_PASSWORD", ""),
		toEmail:      getEnv("EMAIL_TO", ""),
		smtpHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		smtpPort:     getEnvAsInt("SMTP_PORT", 587),
	}
}

func (e *EmailService) SendFeedbackEmail(feedback *Feedback) error {
	if e.fromEmail == "" || e.fromPassword == "" || e.toEmail == "" {
		return fmt.Errorf("email configuration is incomplete")
	}

	// Используем текущее время в правильном часовом поясе
	timezone := getEnv("TIMEZONE", "Asia/Almaty")
	fmt.Printf("🔧 DEBUG: Переменная TIMEZONE = '%s'\n", timezone)

	// Используем фиксированное смещение для Asia/Almaty (UTC+5)
	var currentTime time.Time
	if timezone == "Asia/Almaty" {
		// Создаем фиксированное смещение UTC+5
		loc := time.FixedZone("Asia/Almaty", 5*60*60) // +5 часов в секундах
		currentTime = time.Now().In(loc)
		fmt.Printf("✅ DEBUG: Используем фиксированное смещение UTC+5\n")
	} else {
		// Пытаемся загрузить часовой пояс
		loc, err := time.LoadLocation(timezone)
		if err != nil {
			// Если не удалось загрузить часовой пояс, используем UTC
			loc = time.UTC
			fmt.Printf("⚠️ DEBUG: Не удалось загрузить часовой пояс '%s', используем UTC\n", timezone)
		} else {
			fmt.Printf("✅ DEBUG: Используем часовой пояс: %s\n", timezone)
		}
		currentTime = time.Now().In(loc)
	}

	fmt.Printf("🕐 DEBUG: Время для email: %s\n", currentTime.Format("02.01.2006 15:04:05"))
	fmt.Printf("🕐 DEBUG: UTC время: %s\n", time.Now().UTC().Format("02.01.2006 15:04:05"))

	// Формируем тему письма
	subject := fmt.Sprintf("Новое обращение: %s", getTypeDisplayName(feedback.Type))

	// Формируем тело письма
	body := fmt.Sprintf(`🏥 Новое обращение в системе обратной связи

👤 Отправитель:
• Имя: %s %s
• Username: @%s
• ID: %d

📝 Тип обращения: %s
📅 Дата: %s

💬 Сообщение:
%s

---
Это автоматическое уведомление от системы обратной связи больницы.`,
		feedback.FirstName,
		feedback.LastName,
		feedback.Username,
		feedback.UserID,
		getTypeDisplayName(feedback.Type),
		currentTime.Format("02.01.2006 15:04:05"),
		feedback.Message,
	)

	m := gomail.NewMessage()
	m.SetHeader("From", e.fromEmail)
	m.SetHeader("To", e.toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(e.smtpHost, e.smtpPort, e.fromEmail, e.fromPassword)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

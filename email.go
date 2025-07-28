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

	// Используем текущее время вместо feedback.CreatedAt
	currentTime := time.Now()

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

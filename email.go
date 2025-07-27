package main

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	fromEmail    string
	fromPassword string
	smtpHost     string
	smtpPort     int
	toEmail      string
}

func NewEmailService() *EmailService {
	return &EmailService{
		fromEmail:    getEnv("EMAIL_FROM", ""),
		fromPassword: getEnv("EMAIL_PASSWORD", ""),
		smtpHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		smtpPort:     getEnvAsInt("SMTP_PORT", 587),
		toEmail:      getEnv("EMAIL_TO", ""),
	}
}

func (e *EmailService) SendFeedbackEmail(feedback *Feedback) error {
	if e.fromEmail == "" || e.fromPassword == "" || e.toEmail == "" {
		return fmt.Errorf("email configuration is incomplete")
	}

	var subject, typeText string
	if feedback.Type == "complaint" {
		subject = "🚨 НОВАЯ ЖАЛОБА - Больница"
		typeText = "ЖАЛОБА"
	} else {
		subject = "⭐ НОВЫЙ ОТЗЫВ - Больница"
		typeText = "ОТЗЫВ"
	}

	body := fmt.Sprintf(`
🏥 Система обратной связи больницы

Тип: %s
ID обращения: %d
Пользователь: %s %s (@%s)
ID пользователя: %d

Сообщение:
%s

Дата: %s
`,
		typeText,
		feedback.ID,
		feedback.FirstName, feedback.LastName, feedback.Username,
		feedback.UserID,
		feedback.Message,
		feedback.CreatedAt.Format("02.01.2006 15:04:05"),
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

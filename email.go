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

	// Формируем тело письма
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

	// Добавляем информацию о медиафайлах
	if len(feedback.MediaFiles) > 0 {
		body += fmt.Sprintf("\n📎 Прикрепленные файлы (%d):\n", len(feedback.MediaFiles))
		for i, media := range feedback.MediaFiles {
			body += fmt.Sprintf("%d. %s", i+1, media.FileType)
			if media.FileName != "" {
				body += fmt.Sprintf(" - %s", media.FileName)
			}
			if media.FileSize > 0 {
				body += fmt.Sprintf(" (%d байт)", media.FileSize)
			}
			body += "\n"
		}
	}

	m := gomail.NewMessage()
	m.SetHeader("From", e.fromEmail)
	m.SetHeader("To", e.toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	// Добавляем информацию о медиафайлах в тело письма
	if len(feedback.MediaFiles) > 0 {
		body += "\n\n📎 Медиафайлы:\n"
		for i, media := range feedback.MediaFiles {
			body += fmt.Sprintf("%d. Тип: %s", i+1, media.FileType)
			if media.FileName != "" {
				body += fmt.Sprintf(", Файл: %s", media.FileName)
			}
			if media.FileSize > 0 {
				body += fmt.Sprintf(", Размер: %d байт", media.FileSize)
			}
			body += "\n"
		}
	}

	d := gomail.NewDialer(e.smtpHost, e.smtpPort, e.fromEmail, e.fromPassword)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

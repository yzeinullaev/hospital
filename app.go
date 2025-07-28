package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

type App struct {
	logger   *logrus.Logger
	bot      *TelegramBot
	database *Database
	email    *EmailService
	server   *http.Server
}

func NewApp(logger *logrus.Logger) *App {
	return &App{
		logger: logger,
	}
}

func (a *App) Run() error {
	// Инициализируем базу данных с повторными попытками
	var db *Database
	var err error

	for i := 0; i < 30; i++ {
		db, err = NewDatabase()
		if err == nil {
			break
		}
		a.logger.Warnf("Failed to connect to database (attempt %d/30): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return fmt.Errorf("failed to initialize database after 30 attempts: %w", err)
	}

	a.database = db
	a.logger.Info("Database connection established")

	// Инициализируем email сервис
	emailService := NewEmailService()
	a.email = emailService

	// Инициализируем Telegram бота
	bot, err := NewTelegramBot(a.database, a.email, a.logger)
	if err != nil {
		return fmt.Errorf("failed to initialize telegram bot: %w", err)
	}
	a.bot = bot

	// Запускаем Telegram бота
	go func() {
		if err := a.bot.Start(); err != nil {
			a.logger.Error("Bot stopped with error: ", err)
		}
	}()

	// Настраиваем HTTP сервер
	mux := http.NewServeMux()
	mux.HandleFunc("/health", a.healthHandler)
	mux.HandleFunc("/feedback", a.feedbackHandler)

	a.server = &http.Server{
		Addr:         ":" + getEnv("PORT", "8080"),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запускаем HTTP сервер
	go func() {
		a.logger.Info("Starting HTTP server on port ", a.server.Addr)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Error("HTTP server error: ", err)
		}
	}()

	// Ожидаем сигнал для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.logger.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error("Server forced to shutdown: ", err)
	}

	if err := a.bot.Stop(); err != nil {
		a.logger.Error("Bot shutdown error: ", err)
	}

	a.logger.Info("Server exited")
	return nil
}

func (a *App) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func (a *App) feedbackHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Feedback endpoint"}`))
}

// getEnv function moved to utils.go

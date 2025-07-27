package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// Загружаем переменные окружения
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Инициализируем логгер
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	// Создаем экземпляр приложения
	app := NewApp(logger)

	// Запускаем приложение
	if err := app.Run(); err != nil {
		logger.Fatal("Failed to run application: ", err)
	}
}

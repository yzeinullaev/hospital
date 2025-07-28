#!/bin/bash
# Скрипт для локального деплоя Hospital Feedback Bot
# Использование: ./deploy-local.sh

set -e

echo "🏥 Hospital Feedback Bot - Local Deploy Script"
echo "=============================================="

# Проверяем наличие Docker и Docker Compose
if ! command -v docker &> /dev/null; then
    echo "❌ Docker не установлен"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose не установлен"
    exit 1
fi

echo "✅ Docker и Docker Compose найдены"

# Останавливаем и удаляем старые контейнеры
echo "🔄 Остановка старых контейнеров..."
docker-compose -f docker-compose.local.yml down

# Очищаем неиспользуемые ресурсы
echo "🧹 Очистка неиспользуемых ресурсов..."
docker system prune -f

# Создаем директории для MySQL если их нет
echo "📁 Создание директорий для данных..."
mkdir -p mysql/data mysql/init mysql/backups

# Собираем и запускаем контейнеры
echo "🚀 Сборка и запуск контейнеров..."
docker-compose -f docker-compose.local.yml up --build -d

# Ждем инициализации MySQL
echo "⏳ Ожидание инициализации MySQL..."
sleep 30

# Проверяем статус
echo "🔍 Проверка статуса приложения..."
docker-compose -f docker-compose.local.yml ps

# Проверяем состояние данных
echo "📊 Проверка состояния данных..."
if [ -f "./mysql-status.sh" ]; then
    ./mysql-status.sh
fi

echo ""
echo "📊 Полезные команды:"
echo "• docker-compose -f docker-compose.local.yml logs app - логи приложения"
echo "• docker-compose -f docker-compose.local.yml logs mysql - логи MySQL"
echo "• docker-compose -f docker-compose.local.yml exec mysql mysql -u root -ppassword hospital_feedback - подключение к БД"
echo "• ./mysql-backup.sh - создание резервной копии"
echo "• ./mysql-status.sh - проверка состояния данных"
echo "• ./mysql-restore.sh <file.sql> - восстановление из резервной копии"
echo ""
echo "✅ Локальный деплой завершен!" 
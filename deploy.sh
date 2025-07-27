#!/bin/bash

# Скрипт для деплоя Hospital Feedback Bot на VPS
# Использование: ./deploy.sh

set -e

echo "🏥 Hospital Feedback Bot - Deploy Script"
echo "========================================"

# Проверяем наличие Docker
if ! command -v docker &> /dev/null; then
    echo "❌ Docker не установлен. Установите Docker сначала."
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose не установлен. Установите Docker Compose сначала."
    exit 1
fi

echo "✅ Docker и Docker Compose установлены"

# Останавливаем существующие контейнеры
echo "🛑 Останавливаем существующие контейнеры..."
docker-compose down

# Удаляем старые образы
echo "🧹 Очищаем старые образы..."
docker system prune -f

# Собираем и запускаем новые контейнеры
echo "🔨 Собираем и запускаем приложение..."
docker-compose up --build -d

# Ждем немного для инициализации
echo "⏳ Ждем инициализации сервисов..."
sleep 10

# Проверяем статус контейнеров
echo "📊 Статус контейнеров:"
docker-compose ps

# Проверяем health endpoint
echo "🏥 Проверяем health endpoint..."
if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo "✅ Приложение успешно запущено!"
    echo "🌐 Доступно по адресу: http://localhost:8080"
    echo "📊 Health check: http://localhost:8080/health"
else
    echo "❌ Приложение не отвечает. Проверьте логи:"
    echo "docker-compose logs app"
fi

echo ""
echo "📝 Полезные команды:"
echo "  Просмотр логов: docker-compose logs -f"
echo "  Остановка: docker-compose down"
echo "  Перезапуск: docker-compose restart"
echo "  Обновление: ./deploy.sh" 
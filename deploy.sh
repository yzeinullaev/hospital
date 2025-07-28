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

# Создаем директории для MySQL если их нет
echo "📁 Создание директорий для данных..."
mkdir -p mysql/init mysql/backups
sudo mkdir -p /opt/mysql_data
sudo chown -R 999:999 /opt/mysql_data

# Собираем и запускаем новые контейнеры
echo "🔨 Собираем и запускаем приложение..."
docker-compose up --build -d

# Ждем инициализации MySQL
echo "⏳ Ждем инициализации MySQL..."
sleep 30

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

# Проверяем состояние данных
echo "📊 Проверка состояния данных..."
if [ -f "./mysql-status.sh" ]; then
    ./mysql-status.sh
fi

echo ""
echo "📝 Полезные команды:"
echo "  Просмотр логов: docker-compose logs -f"
echo "  Остановка: docker-compose down"
echo "  Перезапуск: docker-compose restart"
echo "  Обновление: ./deploy.sh"
echo "  Резервная копия: ./mysql-backup.sh"
echo "  Проверка данных: ./mysql-status.sh"
echo "  Восстановление: ./mysql-restore.sh <file.sql>" 
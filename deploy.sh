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

# Ждем инициализации MySQL с проверкой
echo "⏳ Ждем инициализации MySQL..."
for i in {1..60}; do
    if docker-compose exec -T mysql mysqladmin ping -h localhost -u hospital_user -phospital_password > /dev/null 2>&1; then
        echo "✅ MySQL готов к работе"
        break
    fi
    echo "⏳ Ожидание MySQL... ($i/60)"
    sleep 5
done

# Проверяем статус контейнеров
echo "📊 Статус контейнеров:"
docker-compose ps

# Ждем еще немного для полной инициализации приложения
echo "⏳ Ждем запуска приложения..."
sleep 10

# Проверяем health endpoint
echo "🏥 Проверяем health endpoint..."
if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo "✅ Приложение успешно запущено!"
    echo "🌐 Доступно по адресу: http://localhost:8080"
    echo "📊 Health check: http://localhost:8080/health"
else
    echo "❌ Приложение не отвечает. Проверьте логи:"
    echo "docker-compose logs app"
    echo ""
    echo "🔍 Последние логи приложения:"
    docker-compose logs app --tail=20
fi

# Проверяем состояние данных
echo "📊 Проверка состояния данных..."
if [ -f "./mysql-status.sh" ]; then
    chmod +x ./mysql-status.sh
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
#!/bin/bash

# Скрипт для исправления структуры базы данных
# Использование: ./fix-database.sh

set -e

echo "🔧 Исправление структуры базы данных"
echo "===================================="

# Проверяем, что MySQL контейнер запущен
if ! docker-compose ps mysql | grep -q "Up"; then
    echo "❌ MySQL контейнер не запущен. Запустите сначала: docker-compose up -d mysql"
    exit 1
fi

echo "✅ MySQL контейнер запущен"

# Ждем готовности MySQL
echo "⏳ Ждем готовности MySQL..."
for i in {1..30}; do
    if docker-compose exec -T mysql mysqladmin ping -h localhost -u hospital_user -phospital_password > /dev/null 2>&1; then
        echo "✅ MySQL готов к работе"
        break
    fi
    echo "⏳ Ожидание MySQL... ($i/30)"
    sleep 2
done

# Применяем исправление структуры таблиц
echo "🔧 Применяем исправление структуры таблиц..."
docker-compose exec -T mysql mysql -u hospital_user -phospital_password hospital_feedback < mysql/init/02-fix-tables.sql

echo "✅ Структура таблиц исправлена!"

# Проверяем структуру таблиц
echo "📊 Проверяем структуру таблиц..."
docker-compose exec -T mysql mysql -u hospital_user -phospital_password hospital_feedback -e "
SHOW TABLES;
DESCRIBE feedback;
DESCRIBE media_files;
"

echo ""
echo "🎯 Теперь можно перезапустить приложение:"
echo "docker-compose restart app"
echo ""
echo "📝 Или полный перезапуск:"
echo "./deploy.sh" 
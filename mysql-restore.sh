#!/bin/bash
# Скрипт для восстановления данных MySQL

set -e

echo "🔄 MySQL Restore Script"
echo "======================"

# Проверяем наличие аргумента с именем файла
if [ $# -eq 0 ]; then
    echo "❌ Укажите файл для восстановления"
    echo "Использование: $0 <backup_file.sql>"
    echo ""
    echo "Доступные резервные копии:"
    ls -la ./mysql/backups/*.sql 2>/dev/null || echo "Нет доступных резервных копий"
    exit 1
fi

BACKUP_FILE="$1"

# Проверяем существование файла
if [ ! -f "$BACKUP_FILE" ]; then
    echo "❌ Файл $BACKUP_FILE не найден"
    exit 1
fi

echo "📦 Восстановление из файла: $BACKUP_FILE"
echo "⚠️  ВНИМАНИЕ: Это перезапишет все существующие данные!"
echo ""

# Запрашиваем подтверждение
read -p "Продолжить восстановление? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ Восстановление отменено"
    exit 1
fi

echo "🔄 Восстановление данных..."
docker-compose -f docker-compose.local.yml exec -T mysql mysql -u root -ppassword hospital_feedback < "$BACKUP_FILE"

if [ $? -eq 0 ]; then
    echo "✅ Данные восстановлены успешно!"
    echo "📊 Проверка данных..."
    docker-compose -f docker-compose.local.yml exec -T mysql mysql -u root -ppassword -e "USE hospital_feedback; SELECT COUNT(*) as total_feedback FROM feedback;"
else
    echo "❌ Ошибка при восстановлении данных"
    exit 1
fi

echo "🎉 Восстановление завершено!" 
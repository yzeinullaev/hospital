#!/bin/bash
# Скрипт для настройки автоматического резервного копирования

set -e

echo "⏰ Настройка автоматического резервного копирования"
echo "=================================================="

# Проверяем, что мы в корне проекта
if [ ! -f "docker-compose.yml" ]; then
    echo "❌ Запустите скрипт из корня проекта"
    exit 1
fi

# Получаем абсолютный путь к проекту
PROJECT_DIR=$(pwd)
BACKUP_SCRIPT="$PROJECT_DIR/mysql-backup.sh"

echo "📁 Проект: $PROJECT_DIR"
echo "📦 Скрипт бэкапа: $BACKUP_SCRIPT"

# Делаем скрипт исполняемым
chmod +x "$BACKUP_SCRIPT"

# Создаем cron задачу
CRON_JOB="0 2 * * * $BACKUP_SCRIPT >> $PROJECT_DIR/mysql/backups/backup.log 2>&1"

echo ""
echo "📋 Настройка cron задачи..."
echo "⏰ Время: каждый день в 2:00 утра"
echo "📝 Задача: $CRON_JOB"

# Проверяем, есть ли уже такая задача
if crontab -l 2>/dev/null | grep -q "$BACKUP_SCRIPT"; then
    echo "⚠️  Задача уже существует в cron"
    echo "📋 Текущие задачи:"
    crontab -l | grep -E "(backup|mysql)" || echo "Нет задач резервного копирования"
else
    # Добавляем задачу в cron
    (crontab -l 2>/dev/null; echo "$CRON_JOB") | crontab -
    echo "✅ Задача добавлена в cron"
fi

echo ""
echo "📊 Настройка завершена!"
echo ""
echo "🔧 Полезные команды:"
echo "• crontab -l - посмотреть все задачи"
echo "• crontab -r - удалить все задачи (осторожно!)"
echo "• tail -f $PROJECT_DIR/mysql/backups/backup.log - логи бэкапов"
echo "• $BACKUP_SCRIPT - запустить бэкап вручную"
echo ""
echo "💾 Политика очистки:"
echo "• Файлы старше 7 дней удаляются автоматически"
echo "• Максимум 10 файлов бэкапа"
echo "• Логи сохраняются в backup.log" 
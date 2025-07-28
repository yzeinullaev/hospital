#!/bin/bash
# Скрипт для мониторинга места на диске

set -e

echo "💾 Мониторинг места на диске"
echo "============================"

# Проверяем свободное место
FREE_SPACE=$(df -h . | awk 'NR==2 {print $4}')
USED_SPACE=$(df -h . | awk 'NR==2 {print $3}')
TOTAL_SPACE=$(df -h . | awk 'NR==2 {print $2}')
USAGE_PERCENT=$(df . | awk 'NR==2 {print $5}' | sed 's/%//')

echo "📊 Информация о диске:"
echo "• Всего места: $TOTAL_SPACE"
echo "• Использовано: $USED_SPACE"
echo "• Свободно: $FREE_SPACE"
echo "• Заполнено: $USAGE_PERCENT%"

# Проверяем размер директории с бэкапами
BACKUP_DIR="./mysql/backups"
if [ -d "$BACKUP_DIR" ]; then
    BACKUP_SIZE=$(du -sh "$BACKUP_DIR" 2>/dev/null | cut -f1)
    BACKUP_COUNT=$(find "$BACKUP_DIR" -name "*.sql" -type f 2>/dev/null | wc -l)
    echo ""
    echo "📦 Резервные копии:"
    echo "• Размер: $BACKUP_SIZE"
    echo "• Количество файлов: $BACKUP_COUNT"
else
    echo ""
    echo "📦 Резервные копии: директория не найдена"
fi

# Проверяем размер данных MySQL
MYSQL_DATA_DIR="./mysql/data"
if [ -d "$MYSQL_DATA_DIR" ]; then
    MYSQL_SIZE=$(du -sh "$MYSQL_DATA_DIR" 2>/dev/null | cut -f1)
    echo "🗄️ Данные MySQL: $MYSQL_SIZE"
else
    echo "🗄️ Данные MySQL: директория не найдена"
fi

# Предупреждения
echo ""
if [ "$USAGE_PERCENT" -gt 90 ]; then
    echo "⚠️  ВНИМАНИЕ: Диск заполнен на $USAGE_PERCENT%!"
    echo "🔧 Рекомендуется очистка:"
    echo "• ./mysql-backup.sh - создать бэкап и очистить старые"
    echo "• docker system prune -f - очистить неиспользуемые Docker образы"
    echo "• find ./mysql/backups -name '*.sql' -mtime +7 -delete - удалить старые бэкапы"
elif [ "$USAGE_PERCENT" -gt 80 ]; then
    echo "⚠️  Предупреждение: Диск заполнен на $USAGE_PERCENT%"
    echo "💡 Рекомендуется мониторинг места"
else
    echo "✅ Места достаточно: $USAGE_PERCENT%"
fi

echo ""
echo "🎯 Полезные команды для очистки:"
echo "• docker system prune -f - очистить Docker"
echo "• ./mysql-backup.sh - создать бэкап с очисткой"
echo "• find ./mysql/backups -name '*.sql' -mtime +7 -delete - удалить старые бэкапы"
echo "• du -sh ./mysql/* - показать размер директорий" 
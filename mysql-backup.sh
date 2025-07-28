#!/bin/bash
# Скрипт для резервного копирования данных MySQL

set -e

echo "🗄️ MySQL Backup Script"
echo "======================"

# Создаем директорию для бэкапов если её нет
BACKUP_DIR="./mysql/backups"
mkdir -p "$BACKUP_DIR"

# Генерируем имя файла с текущей датой
BACKUP_FILE="$BACKUP_DIR/hospital_feedback_$(date +%Y%m%d_%H%M%S).sql"

echo "📦 Создание резервной копии..."
echo "Файл: $BACKUP_FILE"

# Создаем резервную копию
docker-compose -f docker-compose.local.yml exec -T mysql mysqldump -u root -ppassword hospital_feedback > "$BACKUP_FILE"

if [ $? -eq 0 ]; then
    echo "✅ Резервная копия создана успешно!"
    echo "📁 Файл: $BACKUP_FILE"
    echo "📊 Размер: $(du -h "$BACKUP_FILE" | cut -f1)"
else
    echo "❌ Ошибка при создании резервной копии"
    exit 1
fi

# Автоматическая очистка старых бэкапов
echo "🧹 Очистка старых резервных копий..."

# Удаляем файлы старше 7 дней (оставляем последние 10)
cd "$BACKUP_DIR"

# Удаляем файлы старше 7 дней
find . -name "*.sql" -type f -mtime +7 -delete

# Если файлов больше 10, удаляем самые старые
BACKUP_COUNT=$(ls -t *.sql 2>/dev/null | wc -l)
if [ "$BACKUP_COUNT" -gt 10 ]; then
    echo "📊 Найдено $BACKUP_COUNT файлов, удаляем старые..."
    ls -t *.sql | tail -n +11 | xargs -r rm
    echo "✅ Удалено старых файлов: $((BACKUP_COUNT - 10))"
else
    echo "✅ Файлов в пределах нормы: $BACKUP_COUNT"
fi

# Показываем текущие бэкапы
echo ""
echo "📋 Текущие резервные копии:"
ls -la *.sql 2>/dev/null | head -5 || echo "Резервные копии не найдены"

echo ""
echo "🎉 Резервное копирование завершено!"
echo "💾 Очистка: файлы старше 7 дней + максимум 10 файлов" 
#!/bin/bash
# Скрипт для проверки состояния данных MySQL

set -e

echo "📊 MySQL Status Check"
echo "===================="

echo "🔍 Проверка подключения к MySQL..."
if docker-compose -f docker-compose.local.yml exec -T mysql mysqladmin ping -u root -ppassword --silent; then
    echo "✅ MySQL работает"
else
    echo "❌ MySQL не отвечает"
    exit 1
fi

echo ""
echo "📋 Информация о базе данных:"
docker-compose -f docker-compose.local.yml exec -T mysql mysql -u root -ppassword -e "
SHOW DATABASES;
USE hospital_feedback;
SHOW TABLES;
SELECT 
    'Всего записей' as metric,
    COUNT(*) as value 
FROM feedback
UNION ALL
SELECT 
    'Жалобы' as metric,
    COUNT(*) as value 
FROM feedback 
WHERE type = 'complaint'
UNION ALL
SELECT 
    'Отзывы' as metric,
    COUNT(*) as value 
FROM feedback 
WHERE type = 'review'
UNION ALL
SELECT 
    'Новые' as metric,
    COUNT(*) as value 
FROM feedback 
WHERE status = 'new'
UNION ALL
SELECT 
    'Отправленные' as metric,
    COUNT(*) as value 
FROM feedback 
WHERE status = 'sent';
"

echo ""
echo "📁 Информация о томах:"
docker volume ls | grep mysql_data || echo "Том mysql_data не найден"

echo ""
echo "💾 Размер данных:"
if [ -d "./mysql/data" ]; then
    du -sh ./mysql/data
else
    echo "Директория данных не найдена"
fi

echo ""
echo "🔄 Последние резервные копии:"
ls -la ./mysql/backups/*.sql 2>/dev/null | head -5 || echo "Резервные копии не найдены"

echo ""
echo "✅ Проверка завершена!" 
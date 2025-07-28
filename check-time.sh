#!/bin/bash

echo "🕐 Проверка времени на сервере:"
echo ""

echo "📅 Системное время (UTC):"
date -u
echo ""

echo "📅 Системное время (локальное):"
date
echo ""

echo "📅 Время в часовом поясе Asia/Almaty:"
TZ=Asia/Almaty date
echo ""

echo "📅 Время в часовом поясе Europe/Moscow:"
TZ=Europe/Moscow date
echo ""

echo "📅 Доступные часовые пояса:"
timedatectl list-timezones | grep -E "(Asia|Europe)" | head -10
echo "..." 
#!/bin/bash

# Скрипт для настройки SSL сертификата с Let's Encrypt
# Использование: ./setup-ssl.sh your-domain.com

set -e

if [ $# -eq 0 ]; then
    echo "❌ Укажите домен в качестве аргумента"
    echo "Использование: ./setup-ssl.sh your-domain.com"
    exit 1
fi

DOMAIN=$1

echo "🔒 SSL Setup Script"
echo "=================="
echo "Домен: $DOMAIN"

# Проверяем, что скрипт запущен от root
if [ "$EUID" -ne 0 ]; then
    echo "❌ Этот скрипт должен быть запущен от имени root"
    echo "Используйте: sudo ./setup-ssl.sh $DOMAIN"
    exit 1
fi

# Устанавливаем Certbot
echo "📦 Устанавливаем Certbot..."
apt update
apt install -y certbot python3-certbot-nginx

# Получаем SSL сертификат
echo "🔐 Получаем SSL сертификат для $DOMAIN..."
certbot --nginx -d $DOMAIN --non-interactive --agree-tos --email admin@$DOMAIN

# Настраиваем автообновление сертификата
echo "⏰ Настраиваем автообновление сертификата..."
(crontab -l 2>/dev/null; echo "0 12 * * * /usr/bin/certbot renew --quiet") | crontab -

# Проверяем конфигурацию Nginx
echo "✅ Проверяем конфигурацию Nginx..."
nginx -t

# Перезапускаем Nginx
echo "🔄 Перезапускаем Nginx..."
systemctl restart nginx

echo ""
echo "✅ SSL сертификат успешно настроен!"
echo "🌐 Ваш сайт доступен по адресу: https://$DOMAIN"
echo ""
echo "📝 Полезные команды:"
echo "  Проверка сертификата: sudo certbot certificates"
echo "  Обновление сертификата: sudo certbot renew"
echo "  Проверка статуса Nginx: sudo systemctl status nginx" 
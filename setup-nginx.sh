#!/bin/bash

# Скрипт для настройки Nginx на VPS
# Использование: ./setup-nginx.sh

set -e

echo "🌐 Nginx Setup Script"
echo "===================="

# Проверяем, что скрипт запущен от root
if [ "$EUID" -ne 0 ]; then
    echo "❌ Этот скрипт должен быть запущен от имени root"
    echo "Используйте: sudo ./setup-nginx.sh"
    exit 1
fi

# Устанавливаем Nginx
echo "📦 Устанавливаем Nginx..."
apt update
apt install -y nginx

# Создаем конфигурацию для приложения
echo "⚙️ Создаем конфигурацию Nginx..."
cat > /etc/nginx/sites-available/hospital-bot << 'EOF'
server {
    listen 80;
    server_name your-domain.com;  # Замените на ваш домен

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebSocket support (если понадобится)
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }

    # Health check endpoint
    location /health {
        proxy_pass http://localhost:8080/health;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Gzip compression
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_proxied any;
    gzip_comp_level 6;
    gzip_types
        text/plain
        text/css
        text/xml
        text/javascript
        application/json
        application/javascript
        application/xml+rss
        application/atom+xml
        image/svg+xml;
}
EOF

# Активируем сайт
echo "🔗 Активируем конфигурацию..."
ln -sf /etc/nginx/sites-available/hospital-bot /etc/nginx/sites-enabled/

# Удаляем дефолтную конфигурацию
rm -f /etc/nginx/sites-enabled/default

# Проверяем конфигурацию
echo "✅ Проверяем конфигурацию Nginx..."
nginx -t

# Перезапускаем Nginx
echo "🔄 Перезапускаем Nginx..."
systemctl restart nginx
systemctl enable nginx

echo ""
echo "✅ Nginx успешно настроен!"
echo "📝 Не забудьте:"
echo "  1. Заменить 'your-domain.com' на ваш домен в конфигурации"
echo "  2. Настроить SSL сертификат (Let's Encrypt)"
echo "  3. Открыть порт 80 в файрволе"
echo ""
echo "🔧 Полезные команды:"
echo "  Перезапуск Nginx: sudo systemctl restart nginx"
echo "  Просмотр логов: sudo tail -f /var/log/nginx/access.log"
echo "  Проверка статуса: sudo systemctl status nginx" 
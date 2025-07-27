#!/bin/bash

# Полный скрипт для настройки VPS
# Использование: ./vps-setup.sh

set -e

echo "🚀 VPS Setup Script - Hospital Feedback Bot"
echo "==========================================="

# Проверяем, что скрипт запущен от root
if [ "$EUID" -ne 0 ]; then
    echo "❌ Этот скрипт должен быть запущен от имени root"
    echo "Используйте: sudo ./vps-setup.sh"
    exit 1
fi

# Обновляем систему
echo "📦 Обновляем систему..."
apt update && apt upgrade -y

# Устанавливаем необходимые пакеты
echo "📦 Устанавливаем необходимые пакеты..."
apt install -y curl wget git unzip software-properties-common apt-transport-https ca-certificates gnupg lsb-release

# Устанавливаем Docker
echo "🐳 Устанавливаем Docker..."
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
apt update
apt install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

# Добавляем пользователя в группу docker
echo "👤 Настраиваем права пользователя..."
usermod -aG docker $SUDO_USER

# Устанавливаем Docker Compose
echo "📦 Устанавливаем Docker Compose..."
curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# Устанавливаем Nginx
echo "🌐 Устанавливаем Nginx..."
apt install -y nginx

# Настраиваем файрвол
echo "🔥 Настраиваем файрвол..."
ufw allow ssh
ufw allow 80/tcp
ufw allow 443/tcp
ufw --force enable

# Создаем директорию для проекта
echo "📁 Создаем директорию проекта..."
mkdir -p /opt/hospital-bot
cd /opt/hospital-bot

# Клонируем репозиторий (если нужно)
echo "📥 Клонируем репозиторий..."
# git clone https://github.com/your-username/hospital-feedback-bot.git .

# Делаем скрипты исполняемыми
echo "🔧 Настраиваем скрипты..."
chmod +x deploy.sh
chmod +x setup-nginx.sh
chmod +x setup-ssl.sh

# Настраиваем права на файлы
echo "🔐 Настраиваем права доступа..."
chown -R $SUDO_USER:$SUDO_USER /opt/hospital-bot

echo ""
echo "✅ VPS успешно настроен!"
echo ""
echo "📝 Следующие шаги:"
echo "  1. Настройте .env файл с вашими данными"
echo "  2. Запустите: ./deploy.sh"
echo "  3. Настройте Nginx: sudo ./setup-nginx.sh"
echo "  4. Настройте SSL: sudo ./setup-ssl.sh your-domain.com"
echo ""
echo "🔧 Полезные команды:"
echo "  Перезагрузка: sudo reboot"
echo "  Проверка Docker: docker --version"
echo "  Проверка Nginx: sudo systemctl status nginx" 
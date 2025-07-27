# ⚡ Быстрый деплой на VPS

## 🚀 Команды для быстрого деплоя

### 1. Подключение к серверу
```bash
ssh root@your-server-ip
```

### 2. Загрузка и настройка
```bash
# Создаем директорию
mkdir -p /opt/hospital-bot && cd /opt/hospital-bot

# Клонируем репозиторий
git clone https://github.com/your-username/hospital-feedback-bot.git .

# Настраиваем VPS
chmod +x vps-setup.sh && sudo ./vps-setup.sh
```

### 3. Настройка конфигурации
```bash
# Копируем и редактируем .env
cp env.example .env
nano .env
```

**Важные настройки в .env:**
```env
DB_HOST=mysql
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=hospital_feedback
TELEGRAM_BOT_TOKEN=your_token
ADMIN_USER_ID=your_id
EMAIL_FROM=your_email@gmail.com
EMAIL_PASSWORD=your_app_password
EMAIL_TO=admin@hospital.com
```

### 4. Запуск приложения
```bash
# Запускаем
chmod +x deploy.sh && ./deploy.sh

# Проверяем
docker-compose ps
curl http://localhost:8080/health
```

### 5. Настройка веб-сервера
```bash
# Настраиваем Nginx
sudo ./setup-nginx.sh

# Настраиваем SSL (замените your-domain.com)
sudo ./setup-ssl.sh your-domain.com
```

## ✅ Проверка работы

```bash
# Статус контейнеров
docker-compose ps

# Логи приложения
docker-compose logs app

# Health check
curl http://localhost:8080/health

# Проверка Telegram бота
# Отправьте /start вашему боту
```

## 🔧 Полезные команды

```bash
# Перезапуск
docker-compose restart

# Остановка
docker-compose down

# Обновление
git pull && ./deploy.sh

# Логи в реальном времени
docker-compose logs -f
```

## 🎯 Готово!

Ваше приложение доступно по адресу:
- **HTTP**: `http://your-domain.com`
- **HTTPS**: `https://your-domain.com`

Telegram бот готов к работе! 🚀 
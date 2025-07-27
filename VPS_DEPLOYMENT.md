# 🚀 Деплой на VPS - Hospital Feedback Bot

## 📋 Требования

- VPS с Ubuntu 20.04+ или Debian 11+
- Домен (для SSL)
- SSH доступ к серверу

## 🔧 Быстрая настройка

### 1. Подключение к серверу
```bash
ssh root@your-server-ip
```

### 2. Загрузка проекта
```bash
# Создаем директорию
mkdir -p /opt/hospital-bot
cd /opt/hospital-bot

# Клонируем репозиторий
git clone https://github.com/your-username/hospital-feedback-bot.git .
```

### 3. Первоначальная настройка VPS
```bash
# Делаем скрипт исполняемым
chmod +x vps-setup.sh

# Запускаем настройку
sudo ./vps-setup.sh
```

### 4. Настройка переменных окружения
```bash
# Копируем пример конфигурации
cp env.example .env

# Редактируем конфигурацию
nano .env
```

**Важные настройки в .env:**
```env
# Database Configuration (для Docker)
DB_HOST=mysql
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=hospital_feedback

# Telegram Bot Configuration
TELEGRAM_BOT_TOKEN=your_telegram_bot_token
ADMIN_USER_ID=your_admin_user_id

# Email Configuration
EMAIL_FROM=your_email@gmail.com
EMAIL_PASSWORD=your_app_password
EMAIL_TO=admin@hospital.com
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587

# Server Configuration
PORT=8080
```

### 5. Запуск приложения
```bash
# Делаем скрипт исполняемым
chmod +x deploy.sh

# Запускаем приложение
./deploy.sh
```

### 6. Настройка Nginx
```bash
# Редактируем домен в конфигурации
nano setup-nginx.sh

# Запускаем настройку Nginx
sudo ./setup-nginx.sh
```

### 7. Настройка SSL (опционально)
```bash
# Запускаем настройку SSL
sudo ./setup-ssl.sh your-domain.com
```

## 📊 Проверка работы

### Проверка контейнеров
```bash
docker-compose ps
```

### Проверка логов
```bash
# Все логи
docker-compose logs

# Логи приложения
docker-compose logs app

# Логи базы данных
docker-compose logs mysql
```

### Проверка health endpoint
```bash
curl http://localhost:8080/health
```

## 🔧 Управление приложением

### Перезапуск
```bash
docker-compose restart
```

### Остановка
```bash
docker-compose down
```

### Обновление
```bash
# Останавливаем
docker-compose down

# Обновляем код
git pull

# Перезапускаем
./deploy.sh
```

### Просмотр логов в реальном времени
```bash
docker-compose logs -f
```

## 🛠️ Устранение неполадок

### Проблема: Приложение не запускается
```bash
# Проверяем логи
docker-compose logs app

# Проверяем переменные окружения
docker-compose exec app env | grep DB_
```

### Проблема: База данных не подключается
```bash
# Проверяем статус MySQL
docker-compose logs mysql

# Проверяем подключение
docker-compose exec app ping mysql
```

### Проблема: Nginx не работает
```bash
# Проверяем статус
sudo systemctl status nginx

# Проверяем конфигурацию
sudo nginx -t

# Просмотр логов
sudo tail -f /var/log/nginx/error.log
```

## 🔒 Безопасность

### Файрвол
```bash
# Проверяем статус
sudo ufw status

# Открываем порты
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 22/tcp
```

### SSL сертификат
```bash
# Проверяем сертификат
sudo certbot certificates

# Обновляем сертификат
sudo certbot renew
```

## 📈 Мониторинг

### Системные ресурсы
```bash
# CPU и память
htop

# Дисковое пространство
df -h

# Использование Docker
docker system df
```

### Логи приложения
```bash
# Создаем ротацию логов
sudo nano /etc/logrotate.d/hospital-bot
```

## 🔄 Автоматическое обновление

### Создание cron задачи
```bash
# Открываем crontab
crontab -e

# Добавляем задачу (обновление каждый день в 3:00)
0 3 * * * cd /opt/hospital-bot && git pull && ./deploy.sh
```

## 📞 Поддержка

При возникновении проблем:

1. Проверьте логи: `docker-compose logs`
2. Проверьте статус сервисов: `docker-compose ps`
3. Проверьте конфигурацию: `docker-compose config`
4. Перезапустите приложение: `./deploy.sh`

## 🎯 Готово!

После выполнения всех шагов ваше приложение будет доступно по адресу:
- HTTP: `http://your-domain.com`
- HTTPS: `https://your-domain.com` (если настроен SSL)

Telegram бот будет работать и принимать сообщения от пользователей. 
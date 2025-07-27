# 🚀 Развертывание на ps.kz

## Подготовка к развертыванию

### 1. Подготовка файлов

Убедитесь, что у вас есть все необходимые файлы:

```bash
# Структура проекта
├── main.go
├── app.go
├── database.go
├── telegram.go
├── email.go
├── utils.go
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
├── env.example
├── .env (создайте на основе env.example)
├── deploy.sh
└── README.md
```

### 2. Настройка переменных окружения

Создайте файл `.env` на основе `env.example`:

```bash
cp env.example .env
```

Отредактируйте `.env` файл:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=your_db_user
DB_PASSWORD=your_db_password
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

## Развертывание на ps.kz

### 1. Подключение к серверу

```bash
ssh username@your-server-ip
```

### 2. Установка Docker и Docker Compose

```bash
# Обновляем пакеты
sudo apt update

# Устанавливаем Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Добавляем пользователя в группу docker
sudo usermod -aG docker $USER

# Устанавливаем Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Перезагружаем систему
sudo reboot
```

### 3. Загрузка проекта

```bash
# Создаем директорию для проекта
mkdir -p /var/www/hospital-feedback
cd /var/www/hospital-feedback

# Клонируем репозиторий (если используете Git)
git clone <your-repository-url> .

# Или загружаем файлы через SFTP/SCP
```

### 4. Настройка базы данных

```bash
# Создаем базу данных MySQL
mysql -u root -p
```

```sql
CREATE DATABASE hospital_feedback CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'hospital_user'@'localhost' IDENTIFIED BY 'your_secure_password';
GRANT ALL PRIVILEGES ON hospital_feedback.* TO 'hospital_user'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

### 5. Настройка переменных окружения

Отредактируйте `.env` файл с правильными настройками для production:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=hospital_user
DB_PASSWORD=your_secure_password
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

### 6. Запуск приложения

```bash
# Делаем скрипт исполняемым
chmod +x deploy.sh

# Запускаем развертывание
./deploy.sh
```

### 7. Настройка Nginx (опционально)

Создайте конфигурацию Nginx:

```bash
sudo nano /etc/nginx/sites-available/hospital-feedback
```

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Активируйте сайт:

```bash
sudo ln -s /etc/nginx/sites-available/hospital-feedback /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 8. Настройка SSL (рекомендуется)

```bash
# Устанавливаем Certbot
sudo apt install certbot python3-certbot-nginx

# Получаем SSL сертификат
sudo certbot --nginx -d your-domain.com
```

## Мониторинг и обслуживание

### Просмотр логов

```bash
# Логи приложения
docker-compose logs -f app

# Логи базы данных
docker-compose logs -f mysql
```

### Обновление приложения

```bash
# Останавливаем приложение
docker-compose down

# Обновляем код
git pull origin main

# Пересобираем и запускаем
./deploy.sh
```

### Резервное копирование

```bash
# Создаем резервную копию базы данных
docker-compose exec mysql mysqldump -u root -p hospital_feedback > backup_$(date +%Y%m%d_%H%M%S).sql
```

## Устранение неполадок

### Проверка статуса сервисов

```bash
docker-compose ps
```

### Проверка подключения к базе данных

```bash
docker-compose exec app go run main.go
```

### Проверка Telegram бота

1. Убедитесь, что токен бота правильный
2. Проверьте, что бот не заблокирован
3. Протестируйте команду `/start`

### Проверка email

1. Убедитесь, что SMTP настройки правильные
2. Проверьте, что Gmail разрешает доступ для приложений
3. Протестируйте отправку email

## Безопасность

### Firewall

```bash
# Открываем только необходимые порты
sudo ufw allow 22
sudo ufw allow 80
sudo ufw allow 443
sudo ufw enable
```

### Регулярные обновления

```bash
# Обновляем систему
sudo apt update && sudo apt upgrade -y

# Обновляем Docker образы
docker-compose pull
```

## Контакты для поддержки

При возникновении проблем:

1. Проверьте логи приложения
2. Убедитесь в правильности конфигурации
3. Проверьте подключение к базе данных
4. Обратитесь к документации проекта 
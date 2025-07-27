# 🚀 Быстрый старт

## Минимальные требования

- Go 1.21+
- MySQL 8.0+
- Telegram Bot Token
- Gmail аккаунт для отправки email

## Быстрая настройка

### 1. Создание Telegram бота

1. Найдите @BotFather в Telegram
2. Отправьте `/newbot`
3. Следуйте инструкциям
4. Сохраните токен бота

### 2. Настройка Gmail

1. Включите двухфакторную аутентификацию
2. Создайте пароль приложения
3. Сохраните пароль приложения

### 3. Настройка базы данных

```sql
CREATE DATABASE hospital_feedback CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 4. Конфигурация

```bash
# Копируем конфигурацию
cp env.example .env

# Редактируем .env файл
nano .env
```

Заполните `.env` файл:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=hospital_feedback

TELEGRAM_BOT_TOKEN=your_bot_token
ADMIN_USER_ID=your_telegram_user_id

EMAIL_FROM=your_email@gmail.com
EMAIL_PASSWORD=your_app_password
EMAIL_TO=admin@hospital.com
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587

PORT=8080
```

### 5. Запуск

```bash
# Скачиваем зависимости
go mod download

# Запускаем приложение
go run .
```

### 6. Тестирование

1. Найдите вашего бота в Telegram
2. Отправьте `/start`
3. Следуйте инструкциям бота
4. Проверьте получение email

## Docker запуск

```bash
# Запуск с Docker Compose
docker-compose up -d

# Просмотр логов
docker-compose logs -f app
```

## Генерация QR кода

1. Откройте `qr_generator.html` в браузере
2. Введите имя пользователя бота
3. Сгенерируйте QR код
4. Распечатайте и разместите в больнице

## Проверка работы

### Health Check
```bash
curl http://localhost:8080/health
```

### Статистика в боте
```
/stats
```

## Устранение проблем

### Бот не отвечает
- Проверьте токен бота
- Убедитесь, что бот не заблокирован

### Email не отправляется
- Проверьте SMTP настройки
- Убедитесь, что Gmail разрешает доступ

### База данных не подключается
- Проверьте настройки MySQL
- Убедитесь, что база данных создана

## Следующие шаги

1. Настройте домен и SSL
2. Настройте мониторинг
3. Настройте резервное копирование
4. Разместите QR коды в больнице 
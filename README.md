# 🏥 Hospital Feedback Bot

Система обратной связи больницы с интеграцией Telegram бота для сбора жалоб и отзывов.

## 🚀 Быстрый старт

### Локальная разработка
```bash
# Клонируем репозиторий
git clone https://github.com/your-username/hospital-feedback-bot.git
cd hospital-feedback-bot

# Настраиваем переменные окружения
cp env.example .env
# Отредактируйте .env файл

# Запускаем с Docker
docker-compose up --build
```

### Деплой на VPS
```bash
# Подключаемся к серверу
ssh root@your-server-ip

# Быстрая настройка
mkdir -p /opt/hospital-bot && cd /opt/hospital-bot
git clone https://github.com/your-username/hospital-feedback-bot.git .
sudo ./vps-setup.sh

# Настраиваем и запускаем
cp env.example .env
# Отредактируйте .env файл
./deploy.sh
sudo ./setup-nginx.sh
sudo ./setup-ssl.sh your-domain.com
```

## 📋 Требования

- **Локально**: Docker, Docker Compose
- **VPS**: Ubuntu 20.04+ или Debian 11+
- **Telegram**: Bot Token от @BotFather
- **Email**: Gmail App Password

## 🔧 Настройка

### 1. Telegram Bot
1. Создайте бота через @BotFather
2. Получите токен
3. Добавьте токен в `.env` файл

### 2. Email настройки
1. Включите 2FA в Gmail
2. Создайте App Password
3. Добавьте в `.env` файл

### 3. Admin User ID
1. Найдите ваш User ID через @userinfobot
2. Добавьте в `.env` файл

## 📁 Структура проекта

```
hospital-feedback-bot/
├── main.go              # Точка входа
├── app.go               # Основная логика приложения
├── database.go          # Работа с MySQL
├── telegram.go          # Telegram бот
├── email.go             # Email уведомления
├── utils.go             # Утилиты
├── docker-compose.yml   # Docker конфигурация
├── Dockerfile           # Docker образ
├── deploy.sh            # Скрипт деплоя
├── setup-nginx.sh       # Настройка Nginx
├── setup-ssl.sh         # Настройка SSL
├── vps-setup.sh         # Настройка VPS
├── .env                 # Переменные окружения
└── README.md            # Документация
```

## 🌐 API Endpoints

- `GET /health` - Проверка состояния сервиса
- `GET /feedback` - Endpoint для обратной связи

## 📱 Telegram Bot

### Команды
- `/start` - Начать работу с ботом
- `/help` - Справка
- `/stats` - Статистика (только для админа)

### Процесс работы
1. Пользователь сканирует QR код
2. Переходит в Telegram бот
3. Выбирает тип обращения (жалоба/отзыв)
4. Отправляет сообщение
5. Данные сохраняются в БД и отправляются на email

## 🗄️ База данных

### Структура таблицы `feedback`
```sql
CREATE TABLE feedback (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    username VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    message TEXT NOT NULL,
    type ENUM('complaint', 'review') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status ENUM('new', 'processed', 'sent') DEFAULT 'new'
);
```

## 🐳 Docker

### Запуск
```bash
docker-compose up --build
```

### Остановка
```bash
docker-compose down
```

### Логи
```bash
docker-compose logs -f
```

## 🚀 Деплой на VPS

### Автоматическая настройка
```bash
# Настройка VPS
sudo ./vps-setup.sh

# Деплой приложения
./deploy.sh

# Настройка веб-сервера
sudo ./setup-nginx.sh
sudo ./setup-ssl.sh your-domain.com
```

### Ручная настройка
См. [VPS_DEPLOYMENT.md](VPS_DEPLOYMENT.md) для подробных инструкций.

## 📊 Мониторинг

### Проверка состояния
```bash
# Статус контейнеров
docker-compose ps

# Health check
curl http://localhost:8080/health

# Логи
docker-compose logs app
```

### Полезные команды
```bash
# Перезапуск
docker-compose restart

# Обновление
git pull && ./deploy.sh

# Бэкап БД
docker-compose exec mysql mysqldump -u root -p hospital_feedback > backup.sql
```

## 🔒 Безопасность

- SSL сертификат (Let's Encrypt)
- Файрвол (UFW)
- Изолированные Docker контейнеры
- Переменные окружения для секретов

## 📈 Масштабирование

### Горизонтальное масштабирование
```bash
# Увеличить количество экземпляров
docker-compose up --scale app=3
```

### Вертикальное масштабирование
- Увеличить ресурсы VPS
- Настроить кэширование
- Оптимизировать запросы к БД

## 🛠️ Разработка

### Добавление новых функций
1. Создайте feature branch
2. Реализуйте функциональность
3. Добавьте тесты
4. Создайте Pull Request

### Локальная разработка
```bash
# Запуск в режиме разработки
go run main.go

# Тестирование
go test ./...
```

## 📞 Поддержка

При возникновении проблем:

1. Проверьте логи: `docker-compose logs`
2. Проверьте статус: `docker-compose ps`
3. Проверьте конфигурацию: `docker-compose config`
4. Создайте issue в репозитории

## 📄 Лицензия

MIT License

## 🤝 Вклад в проект

1. Fork репозитория
2. Создайте feature branch
3. Внесите изменения
4. Создайте Pull Request

---

**Готово к использованию!** 🚀

Ваше приложение для сбора обратной связи больницы готово к работе с Telegram ботом, базой данных MySQL и email уведомлениями.
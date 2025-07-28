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

# Запускаем локально
./deploy-local.sh
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
├── main.go                    # Точка входа
├── app.go                     # Основная логика приложения
├── database.go                # Работа с MySQL
├── telegram.go                # Telegram бот
├── email.go                   # Email уведомления
├── utils.go                   # Утилиты
├── docker-compose.yml         # Docker конфигурация (продакшен)
├── docker-compose.local.yml   # Docker конфигурация (локально)
├── Dockerfile                 # Docker образ
├── deploy.sh                  # Скрипт деплоя (продакшен)
├── deploy-local.sh            # Скрипт деплоя (локально)
├── mysql-backup.sh            # Резервное копирование MySQL
├── mysql-restore.sh           # Восстановление MySQL
├── mysql-status.sh            # Проверка состояния MySQL
├── setup-nginx.sh             # Настройка Nginx
├── setup-ssl.sh               # Настройка SSL
├── vps-setup.sh               # Настройка VPS
├── mysql/
│   ├── data/                  # Данные MySQL (локально)
│   ├── init/                  # Скрипты инициализации
│   └── backups/               # Резервные копии
├── .env                       # Переменные окружения
└── README.md                  # Документация
```

## 🌐 API Endpoints

- `GET /health` - Проверка состояния сервиса
- `GET /feedback` - Endpoint для обратной связи

## 📱 Telegram Bot

### Команды
- `/start` - Начать работу с ботом
- `/help` - Справка
- `/menu` - Главное меню
- `/stats` - Статистика (только для админа)

### Интерактивные кнопки
Бот использует современный интерфейс с кнопками:
- **📝 Жалоба** - для отправки жалоб
- **⭐ Отзыв** - для положительных отзывов
- **❓ Помощь** - справка по использованию
- **🏠 Главное меню** - навигация

### Процесс работы
1. Пользователь сканирует QR код
2. Переходит в Telegram бот
3. Выбирает тип обращения кнопкой
4. Отправляет сообщение
5. Данные сохраняются в БД и отправляются на email
6. Получает подтверждение с кнопками навигации

## 🗄️ База данных

### Структура таблицы `feedback`
```sql
CREATE TABLE feedback (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    username VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    message TEXT NOT NULL,
    type ENUM('complaint', 'review') NOT NULL,
    status ENUM('new', 'sent', 'processed') DEFAULT 'new',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_type (type),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### Сохранение данных
- **Локально**: Данные сохраняются в `./mysql/data/`
- **Продакшен**: Данные сохраняются в `/opt/mysql_data/`
- **Резервные копии**: Автоматически создаются в `./mysql/backups/`

## 🐳 Docker

### Локальная разработка
```bash
# Используйте локальную конфигурацию
./deploy-local.sh

# Или вручную
docker-compose -f docker-compose.local.yml up --build
```

### Продакшен
```bash
# Используйте продакшен конфигурацию
./deploy.sh

# Или вручную
docker-compose up --build
```

### Управление данными MySQL
```bash
# Резервное копирование (с автоматической очисткой)
./mysql-backup.sh

# Автоматическое резервное копирование
./setup-backup-cron.sh

# Восстановление
./mysql-restore.sh backup_file.sql

# Проверка состояния
./mysql-status.sh

# Мониторинг места на диске
./disk-monitor.sh
```

### Остановка
```bash
# Локально
docker-compose -f docker-compose.local.yml down

# Продакшен
docker-compose down
```

### Логи
```bash
# Локально
docker-compose -f docker-compose.local.yml logs -f

# Продакшен
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

## 📊 Мониторинг и управление данными

### Когда использовать команды MySQL

#### 🔄 Обычное обновление
```bash
git pull && ./deploy.sh
```
**Данные сохраняются автоматически** благодаря настроенным томам Docker.

#### 💾 Резервное копирование
```bash
./mysql-backup.sh
```
**Автоматическая очистка**: удаляет файлы старше 7 дней и оставляет максимум 10 файлов.

#### ⏰ Автоматическое резервное копирование
```bash
./setup-backup-cron.sh
```
**Настраивает cron** для ежедневного бэкапа в 2:00 утра с автоматической очисткой.

#### 🔙 Восстановление
```bash
./mysql-restore.sh backup_file.sql
```
**Только при проблемах** с данными или при переезде на новый сервер.

#### 📊 Проверка состояния
```bash
./mysql-status.sh
```
**Для мониторинга** работы БД и отладки проблем.

#### 💾 Мониторинг места на диске
```bash
./disk-monitor.sh
```
**Проверяет** свободное место и размер резервных копий.

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

# Обновление приложения (данные сохраняются автоматически)
git pull && ./deploy.sh

# Автоматическое резервное копирование (настройка)
./setup-backup-cron.sh

# Ручное резервное копирование
./mysql-backup.sh

# Восстановление (только при проблемах)
./mysql-restore.sh backup_file.sql

# Проверка состояния данных
./mysql-status.sh

# Мониторинг места на диске
./disk-monitor.sh
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

YZEINULLAEV License

## 🤝 Вклад в проект

1. Fork репозитория
2. Создайте feature branch
3. Внесите изменения
4. Создайте Pull Request

---

**Готово к использованию!** 🚀

Ваше приложение для сбора обратной связи больницы готово к работе с Telegram ботом, базой данных MySQL и email уведомлениями.
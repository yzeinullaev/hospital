# 🏥 Hospital Feedback Bot

Telegram бот для сбора жалоб и отзывов о больнице с отправкой на email.

## 🚀 Быстрый старт

### Локальная разработка
```bash
# Клонируйте репозиторий
git clone <repository-url>
cd hospital-bot

# Создайте .env файл
cp env.example .env
# Отредактируйте .env с вашими данными

# Запустите приложение
./deploy-local.sh
```

### Продакшен
```bash
# На сервере
git clone <repository-url>
cd hospital-bot

# Создайте .env файл
cp env.example .env
# Отредактируйте .env с вашими данными

# Запустите приложение
./deploy.sh
```

## 📁 Структура проекта

```
hospital-bot/
├── main.go                 # Точка входа приложения
├── app.go                  # Основная логика приложения
├── database.go             # Работа с базой данных
├── telegram.go             # Telegram бот
├── email.go                # Отправка email
├── utils.go                # Утилиты
├── .env                    # Переменные окружения
├── env.example             # Пример переменных окружения
├── docker-compose.yml      # Docker Compose для продакшена
├── docker-compose.local.yml # Docker Compose для локальной разработки
├── deploy.sh               # Скрипт деплоя для продакшена
├── deploy-local.sh         # Скрипт деплоя для локальной разработки
├── mysql-backup.sh         # Резервное копирование MySQL
├── mysql-restore.sh        # Восстановление MySQL
├── mysql-status.sh         # Проверка состояния MySQL
├── setup-backup-cron.sh    # Настройка автоматических бэкапов
├── disk-monitor.sh         # Мониторинг места на диске
├── qr_generator.html       # Генератор QR кода для бота
├── mysql/
│   ├── init/
│   │   └── 01-init.sql     # Инициализация базы данных
│   ├── data/                # Данные MySQL (локальная разработка)
│   └── backups/             # Резервные копии
└── README.md               # Документация
```

## 🐳 Docker

### Локальная разработка
```bash
# Запуск
./deploy-local.sh

# Остановка
docker-compose -f docker-compose.local.yml down

# Логи
docker-compose -f docker-compose.local.yml logs -f
```

### Продакшен
```bash
# Запуск
./deploy.sh

# Остановка
docker-compose down

# Логи
docker-compose logs -f
```

### Управление данными MySQL

```bash
# Создание резервной копии
./mysql-backup.sh

# Восстановление из резервной копии
./mysql-restore.sh backup_file.sql

# Проверка состояния данных
./mysql-status.sh

# Настройка автоматических бэкапов
./setup-backup-cron.sh

# Мониторинг места на диске
./disk-monitor.sh
```

## 🗄️ База данных

### Структура таблиц

#### Таблица `feedback`
```sql
CREATE TABLE feedback (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    username VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    message TEXT NOT NULL,
    type ENUM('complaint', 'review') NOT NULL DEFAULT 'complaint',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status ENUM('new', 'processed', 'sent') DEFAULT 'new',
    INDEX idx_user_id (user_id),
    INDEX idx_type (type),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### Сохранение данных

- **Локальная разработка**: `./mysql/data/` (bind mount)
- **Продакшен**: `/opt/mysql_data/` (host path)
- **Резервные копии**: `./mysql/backups/`

## 🤖 Telegram Bot

### Команды
- `/start` - Начать работу с ботом
- `/menu` - Показать главное меню
- `/stats` - Показать статистику обращений (только для администратора)

### Интерактивные кнопки
- **📝 Отправить жалобу** - Отправить жалобу
- **⭐ Оставить отзыв** - Оставить отзыв
- **❓ Помощь** - Показать справку по использованию бота
- **📊 Статистика** - Показать статистику обращений (только для администратора)
- **🏥 Новое обращение** - Отправить еще одно обращение (после подтверждения)
- **🏠 Главное меню** - Вернуться в главное меню (из раздела помощи или статистики)

### Процесс работы
1. Пользователь видит главное меню с заголовком "Главное меню системы обратной связи больницы"
2. Пользователь нажимает кнопку "Отправить жалобу", "Оставить отзыв", "Помощь" или "Статистика"
3. Администратор также видит кнопку "Статистика"
4. Для жалобы/отзыва: бот просит описать обращение подробно
5. Пользователь отправляет текстовое сообщение
6. Бот сохраняет данные в базу и отправляет email
7. Пользователь получает подтверждение с кнопками "Новое обращение" и "Помощь"
8. Для статистики (только админ): бот показывает количество жалоб и отзывов
9. Из разделов "Помощь" и "Статистика" можно вернуться в главное меню

### Поддерживаемые типы сообщений
- ✅ Текстовые сообщения
- ❌ Медиафайлы (отключены)

## 📧 Email уведомления

При получении нового обращения система отправляет email с информацией:
- Имя и username отправителя
- Тип обращения (жалоба/отзыв)
- Дата и время
- Текст сообщения

## 🔧 Конфигурация

### Переменные окружения (.env)

```env
# База данных
DB_HOST=mysql
DB_PORT=3306
DB_USER=hospital_user
DB_PASSWORD=hospital_password
DB_NAME=hospital_feedback

# Telegram Bot
TELEGRAM_BOT_TOKEN=your_bot_token
ADMIN_USER_ID=your_user_id  # ID администратора для доступа к статистике

# Email
EMAIL_FROM=your_email@gmail.com
EMAIL_PASSWORD=your_app_password
EMAIL_TO=admin@hospital.com
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587

# Сервер
PORT=8080
```

## 📊 Мониторинг и управление данными

### Когда использовать разные команды:

1. **`git pull && ./deploy.sh`** - Обновление кода приложения
   - Данные сохраняются автоматически
   - Приложение перезапускается

2. **`./mysql-backup.sh`** - Создание резервной копии
   - Рекомендуется перед обновлением
   - Очищает старые бэкапы (7 дней)

3. **`./mysql-restore.sh <file.sql>`** - Восстановление данных
   - Используется при проблемах с данными
   - Требует подтверждения

4. **`./mysql-status.sh`** - Проверка состояния
   - Показывает статистику данных
   - Проверяет подключение к MySQL

5. **`./setup-backup-cron.sh`** - Настройка автоматических бэкапов
   - Запускается один раз при настройке сервера
   - Создает cron задачу для ежедневных бэкапов

6. **`./disk-monitor.sh`** - Мониторинг места на диске
   - Проверяет свободное место
   - Показывает размер бэкапов и данных

## 🚀 Деплой на VPS

### Быстрая настройка
```bash
# Установка Docker и зависимостей
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Клонирование и запуск
git clone <repository-url>
cd hospital-bot
cp env.example .env
# Отредактируйте .env
./deploy.sh
```

### Подробная настройка
См. файлы:
- `VPS_DEPLOYMENT.md` - Подробные инструкции
- `QUICK_DEPLOY.md` - Быстрые команды
- `DEPLOYMENT_CHECKLIST.md` - Чек-лист деплоя

## 🔍 Отладка

### Проверка логов
```bash
# Логи приложения
docker-compose logs app

# Логи MySQL
docker-compose logs mysql

# Все логи
docker-compose logs -f
```

### Проверка состояния
```bash
# Статус контейнеров
docker-compose ps

# Проверка данных
./mysql-status.sh

# Проверка места на диске
./disk-monitor.sh
```

## 📝 Полезные команды

### Мониторинг
```bash
# Просмотр логов
docker-compose logs -f

# Остановка
docker-compose down

# Перезапуск
docker-compose restart

# Обновление
./deploy.sh

# Резервная копия
./mysql-backup.sh

# Проверка данных
./mysql-status.sh

# Восстановление
./mysql-restore.sh <file.sql>
```

### Управление данными
```bash
# Автоматические бэкапы
./setup-backup-cron.sh

# Мониторинг места
./disk-monitor.sh
```

## 📄 Лицензия

MIT License
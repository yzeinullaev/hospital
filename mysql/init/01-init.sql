-- Инициализация базы данных Hospital Feedback
-- Создаем базу данных если её нет
CREATE DATABASE IF NOT EXISTS hospital_feedback CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Используем базу данных
USE hospital_feedback;

-- Создаем таблицу feedback если её нет
CREATE TABLE IF NOT EXISTS feedback (
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

-- Создаем таблицу media_files если её нет
CREATE TABLE IF NOT EXISTS media_files (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    feedback_id BIGINT NOT NULL,
    file_id VARCHAR(255) NOT NULL,
    file_type ENUM('photo', 'video', 'document', 'audio') NOT NULL,
    file_name VARCHAR(255),
    file_size BIGINT,
    mime_type VARCHAR(100),
    url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_feedback_id (feedback_id),
    INDEX idx_file_type (file_type),
    FOREIGN KEY (feedback_id) REFERENCES feedback(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Создаем пользователя для приложения если его нет
CREATE USER IF NOT EXISTS 'hospital_user'@'%' IDENTIFIED BY 'hospital_password';
GRANT ALL PRIVILEGES ON hospital_feedback.* TO 'hospital_user'@'%';
FLUSH PRIVILEGES;

-- Показываем информацию о созданных объектах
SHOW DATABASES;
SHOW TABLES FROM hospital_feedback;
SELECT 'MySQL initialization completed successfully!' as status; 
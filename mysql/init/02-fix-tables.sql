-- Исправление структуры таблиц для совместимости
-- Используем базу данных
USE hospital_feedback;

-- Удаляем существующие таблицы в правильном порядке (сначала зависимые)
DROP TABLE IF EXISTS media_files;
DROP TABLE IF EXISTS feedback;

-- Создаем таблицу feedback с правильным типом id
CREATE TABLE feedback (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    username VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    message TEXT NOT NULL,
    type ENUM('complaint', 'review') NOT NULL DEFAULT 'complaint',
    status ENUM('new', 'processed', 'sent') DEFAULT 'new',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_type (type),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Создаем таблицу media_files с правильным foreign key
CREATE TABLE media_files (
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

SELECT 'Tables fixed successfully!' as status; 
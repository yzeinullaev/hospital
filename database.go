package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Feedback struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Message   string    `json:"message"`
	Type      string    `json:"type"` // "complaint" или "review"
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"` // "new", "processed", "sent"
	// Новые поля для медиафайлов
	MediaFiles []MediaFile `json:"media_files"`
}

type MediaFile struct {
	ID         int64  `json:"id"`
	FeedbackID int64  `json:"feedback_id"`
	FileID     string `json:"file_id"`
	FileType   string `json:"file_type"` // "photo", "video", "document", "audio"
	FileName   string `json:"file_name"`
	FileSize   int64  `json:"file_size"`
	MimeType   string `json:"mime_type"`
	URL        string `json:"url"`
}

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		getEnv("DB_USER", "root"),
		getEnv("DB_PASSWORD", ""),
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "3306"),
		getEnv("DB_NAME", "hospital_feedback"),
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Создаем таблицы если их нет
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &Database{db: db}, nil
}

func createTables(db *sql.DB) error {
	// Создаем таблицу feedback
	feedbackQuery := `
	CREATE TABLE IF NOT EXISTS feedback (
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
	`

	if _, err := db.Exec(feedbackQuery); err != nil {
		return fmt.Errorf("failed to create feedback table: %w", err)
	}

	// Создаем таблицу media_files
	mediaQuery := `
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
	`

	if _, err := db.Exec(mediaQuery); err != nil {
		return fmt.Errorf("failed to create media_files table: %w", err)
	}

	return nil
}

func (d *Database) SaveFeedback(feedback *Feedback) error {
	// Начинаем транзакцию
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Сохраняем основное обращение
	query := `
	INSERT INTO feedback (user_id, username, first_name, last_name, message, type, status)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := tx.Exec(query,
		feedback.UserID,
		feedback.Username,
		feedback.FirstName,
		feedback.LastName,
		feedback.Message,
		feedback.Type,
		feedback.Status,
	)
	if err != nil {
		return fmt.Errorf("failed to save feedback: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	feedback.ID = id

	// Сохраняем медиафайлы
	if len(feedback.MediaFiles) > 0 {
		for i := range feedback.MediaFiles {
			feedback.MediaFiles[i].FeedbackID = id
			if err := d.saveMediaFile(tx, &feedback.MediaFiles[i]); err != nil {
				return fmt.Errorf("failed to save media file: %w", err)
			}
		}
	}

	// Подтверждаем транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (d *Database) saveMediaFile(tx *sql.Tx, media *MediaFile) error {
	query := `
	INSERT INTO media_files (feedback_id, file_id, file_type, file_name, file_size, mime_type, url)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := tx.Exec(query,
		media.FeedbackID,
		media.FileID,
		media.FileType,
		media.FileName,
		media.FileSize,
		media.MimeType,
		media.URL,
	)
	if err != nil {
		return fmt.Errorf("failed to save media file: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get media file id: %w", err)
	}

	media.ID = id
	return nil
}

func (d *Database) GetNewFeedbacks() ([]*Feedback, error) {
	query := `
	SELECT id, user_id, username, first_name, last_name, message, type, created_at, status
	FROM feedback
	WHERE status = 'new'
	ORDER BY created_at ASC
	`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query new feedbacks: %w", err)
	}
	defer rows.Close()

	var feedbacks []*Feedback
	for rows.Next() {
		feedback := &Feedback{}
		err := rows.Scan(
			&feedback.ID,
			&feedback.UserID,
			&feedback.Username,
			&feedback.FirstName,
			&feedback.LastName,
			&feedback.Message,
			&feedback.Type,
			&feedback.CreatedAt,
			&feedback.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan feedback: %w", err)
		}

		// Загружаем медиафайлы для этого обращения
		mediaFiles, err := d.getMediaFiles(feedback.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get media files: %w", err)
		}
		feedback.MediaFiles = mediaFiles

		feedbacks = append(feedbacks, feedback)
	}

	return feedbacks, nil
}

func (d *Database) getMediaFiles(feedbackID int64) ([]MediaFile, error) {
	query := `
	SELECT id, feedback_id, file_id, file_type, file_name, file_size, mime_type, url
	FROM media_files
	WHERE feedback_id = ?
	ORDER BY id
	`

	rows, err := d.db.Query(query, feedbackID)
	if err != nil {
		return nil, fmt.Errorf("failed to query media files: %w", err)
	}
	defer rows.Close()

	var mediaFiles []MediaFile
	for rows.Next() {
		var media MediaFile
		err := rows.Scan(
			&media.ID,
			&media.FeedbackID,
			&media.FileID,
			&media.FileType,
			&media.FileName,
			&media.FileSize,
			&media.MimeType,
			&media.URL,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan media file: %w", err)
		}
		mediaFiles = append(mediaFiles, media)
	}

	return mediaFiles, nil
}

func (d *Database) UpdateFeedbackStatus(id int64, status string) error {
	query := `UPDATE feedback SET status = ? WHERE id = ?`
	_, err := d.db.Exec(query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update feedback status: %w", err)
	}
	return nil
}

func (d *Database) GetFeedbackStats() (map[string]int, error) {
	query := `
	SELECT 
		type,
		COUNT(*) as count
	FROM feedback
	GROUP BY type
	`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get feedback stats: %w", err)
	}
	defer rows.Close()

	stats := make(map[string]int)
	for rows.Next() {
		var feedbackType string
		var count int
		err := rows.Scan(&feedbackType, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stats: %w", err)
		}
		stats[feedbackType] = count
	}

	return stats, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

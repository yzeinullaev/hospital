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
}

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&timeout=10s&readTimeout=30s&writeTimeout=30s",
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

	// Настраиваем пул соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

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

	return nil
}

func (d *Database) SaveFeedback(feedback *Feedback) error {
	query := `
	INSERT INTO feedback (user_id, username, first_name, last_name, message, type, status)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := d.db.Exec(query,
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

		feedbacks = append(feedbacks, feedback)
	}

	return feedbacks, nil
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

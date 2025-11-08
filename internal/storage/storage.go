package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Storage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) SignPrivacyPolicy(ctx context.Context, userId int64) error {
	query := `
		UPDATE users SET privacy_policy_signed = $1 WHERE user_id = $2;
	`
	if _, err := s.db.ExecContext(ctx, query, time.Now().UTC(), userId); err != nil {
		return fmt.Errorf("s.db.ExecContext: %w", err)
	}
	return nil
}

func (s *Storage) PrivacyPolicySigned(ctx context.Context, userId int64) (bool, error) {
	query := `
		SELECT exists(SELECT 1 FROM users WHERE user_id = $1 AND privacy_policy_signed IS NOT NULL);
	`

	var signed bool
	if err := s.db.GetContext(ctx, &signed, query, userId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("s.db.GetContext: %w", err)
	}

	return signed, nil
}

func (s *Storage) CreateUser(ctx context.Context, userId int64) error {
	query := `
		INSERT INTO users (user_id) VALUES ($1) ON CONFLICT DO NOTHING;
	`
	if _, err := s.db.ExecContext(ctx, query, userId); err != nil {
		return fmt.Errorf("s.db.ExecContext: %w", err)
	}
	return nil
}

func (s *Storage) LogVoiceQuestion(ctx context.Context, userId int64, question string, fileNames []string) error {
	return s.logQuestion(ctx, userId, question, fileNames, 1)
}

func (s *Storage) LogTextQuestion(ctx context.Context, userId int64, question string, fileNames []string) error {
	return s.logQuestion(ctx, userId, question, fileNames, 2)
}

func (s *Storage) logQuestion(ctx context.Context, userId int64, question string, fileNames []string, qType int64) error {
	const query = `
		INSERT INTO statistics_log 
		(
			user_id, 
		    question_text, 
		    question_type,
		 	file_used
		) VALUES 
		(
		 	$1,
		 	$2,
		 	$3,
		 	(SELECT ARRAY (SELECT DISTINCT id from files WHERE name = ANY ($4)))
		)
	`

	if _, err := s.db.ExecContext(ctx, query, userId, question, qType, pq.Array(fileNames)); err != nil {
		return fmt.Errorf("s.db.ExecContext: %w", err)
	}
	return nil
}

func (s *Storage) GetStatistics(ctx context.Context) (Statistics, error) {
	const query = `
		SELECT 
			(SELECT COUNT(DISTINCT user_id) FROM users) AS total_users,
			(SELECT COUNT(*) FROM statistics_log WHERE question_type = 1) AS total_voice_questions,
			(SELECT COUNT(*) FROM statistics_log WHERE question_type = 2) AS total_text_questions,
			(SELECT ARRAY( SELECT question_text FROM statistics_log ORDER BY created_at DESC LIMIT 10)) AS last_10_questions
	`

	var stats Statistics
	if err := s.db.GetContext(ctx, &stats, query); err != nil {
		return Statistics{}, fmt.Errorf("s.db.GetContext: %w", err)
	}

	return stats, nil
}

func (s *Storage) GetFilesByNames(ctx context.Context, names []string) ([]File, error) {
	const query = `
		SELECT short_name, url FROM files WHERE name = ANY ($1) ORDER BY short_name
	`

	var files []File
	if err := s.db.SelectContext(ctx, &files, query, pq.Array(names)); err != nil {
		return nil, fmt.Errorf("s.db.SelectContext: %w", err)
	}
	return files, nil
}

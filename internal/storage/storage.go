package storage

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
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
		UPDATE users SET confidentiality_policy_signed = $1 WHERE user_id = $2;
	`
	_, err := s.db.ExecContext(ctx, query, time.Now().UTC(), userId)
	if err != nil {
		return fmt.Errorf("s.db.ExecContext: %w", err)
	}
	return nil
}

func (s *Storage) CreateUser(ctx context.Context, userId int64) error {
	query := `
		INSERT INTO users (user_id) VALUES ($1) ON CONFLICT DO NOTHING;
	`
	_, err := s.db.ExecContext(ctx, query, userId)
	if err != nil {
		return fmt.Errorf("s.db.ExecContext: %w", err)
	}
	return nil
}

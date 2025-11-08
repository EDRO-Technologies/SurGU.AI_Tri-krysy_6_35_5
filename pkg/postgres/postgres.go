package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Option func(*options)

type options struct {
	maxOpenConns    int
	maxIdleConns    int
	connMaxLifetime time.Duration
	connIdleTime    time.Duration
}

func defaultOptions() *options {
	return &options{
		maxOpenConns:    5,
		maxIdleConns:    5,
		connMaxLifetime: 5 * time.Minute,
		connIdleTime:    5 * time.Minute,
	}
}

func WithMaxOpenConns(n int) Option {
	return func(o *options) {
		o.maxOpenConns = n
	}
}

func WithMaxIdleConns(n int) Option {
	return func(o *options) {
		o.maxIdleConns = n
	}
}

func WithConnMaxLifetime(d time.Duration) Option {
	return func(o *options) {
		o.connMaxLifetime = d
	}
}

func WithConnIdleTime(d time.Duration) Option {
	return func(o *options) {
		o.connIdleTime = d
	}
}

func MustNew(dsn string, opts ...Option) *sqlx.DB {
	db, err := New(dsn, opts...)
	if err != nil {
		panic(err)
	}
	return db
}

func New(dsn string, opts ...Option) (*sqlx.DB, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	pgxConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("invalid DSN: %w", err)
	}

	stdlibConnector := stdlib.OpenDB(*pgxConfig)
	db := sqlx.NewDb(stdlibConnector, "pgx")

	db.SetMaxOpenConns(o.maxOpenConns)
	db.SetMaxIdleConns(o.maxIdleConns)
	db.SetConnMaxLifetime(o.connMaxLifetime)

	if setIdleTime := interface{}(db.DB).(interface {
		SetConnMaxIdleTime(d time.Duration)
	}); setIdleTime != nil {
		setIdleTime.SetConnMaxIdleTime(o.connIdleTime)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("could not ping database: %w", err)
	}

	return db, nil
}

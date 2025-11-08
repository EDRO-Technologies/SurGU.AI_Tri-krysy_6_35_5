package migrate

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	defaultAttempts = 10
	defaultTimeout  = time.Second
)

func Migrate(dsn string) error {
	var (
		attempts = defaultAttempts
		m        *migrate.Migrate
		err      error
	)

	for attempts > 0 {
		m, err = migrate.New("file://migrations", dsn)
		if err == nil {
			break
		}

		time.Sleep(defaultTimeout)
		attempts--
	}

	if err != nil {
		return fmt.Errorf("migrate.New: %w", err)
	}

	if m == nil {
		return errors.New("migrate: migrate is nil")
	}

	err = m.Up()
	defer func(m *migrate.Migrate) { _, _ = m.Close() }(m)

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate.Up: %w", err)
	}

	return nil
}

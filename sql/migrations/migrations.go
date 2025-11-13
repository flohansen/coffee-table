package migrations

import (
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

var (
	//go:embed *.sql
	migrationsFS embed.FS
)

func Up(pool *pgxpool.Pool) error {
	sourceDriver, err := iofs.New(migrationsFS, ".")
	if err != nil {
		return fmt.Errorf("iofs new: %w", err)
	}

	db := stdlib.OpenDBFromPool(pool)
	databaseDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("postgres with instance: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "", databaseDriver)
	if err != nil {
		return fmt.Errorf("migrate new with instance: %w", err)
	}

	if err := m.Up(); err != migrate.ErrNoChange {
		return fmt.Errorf("migrate up: %w", err)
	}

	return nil
}

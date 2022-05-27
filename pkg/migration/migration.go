package migration

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migration struct {
	db             *sql.DB
	migrationsPath string
}

func New(db *sql.DB, migrationsFolder string) *Migration {
	return &Migration{db: db, migrationsPath: migrationsFolder}
}

func (m *Migration) Up() error {
	driver, err := postgres.WithInstance(m.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("error on Migration.Up() when resolving postgres driver: %w", err)
	}

	// "file://db/migrations",
	mig, err := migrate.NewWithDatabaseInstance(m.migrationsPath, "postgres", driver)
	if err != nil {
		return fmt.Errorf("error on Migration.Up() when creating Migrate instance: %w", err)
	}

	err = mig.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error on Migration.Up() when running migrations: %w", err)
	}

	return nil
}

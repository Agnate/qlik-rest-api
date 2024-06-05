package migrator

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migrator struct {
	directoryName string
	databaseName  string
}

// Create a Migrator to allow us to migrate the schema(s) and updates.
func New(dirName string, dbName string) *Migrator {
	return &Migrator{
		directoryName: dirName,
		databaseName:  dbName,
	}
}

// Run all schemas up to the latest version.
func (m *Migrator) Run(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("unable to create db instance: %v", err)
	}

	// TODO: Properly sanitize the directory name to add/omit slashes as needed.
	migrator, err := migrate.NewWithDatabaseInstance("file://./migrations", m.databaseName, driver)
	if err != nil {
		return fmt.Errorf("unable to create migration: %v", err)
	}

	// Run migrations.
	if err = migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("unable to apply migrations %v", err)
	}

	return nil
}

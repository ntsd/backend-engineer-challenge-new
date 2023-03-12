package storage

import (
	"fmt"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/ntsd/backend-engineer-challenge-new/internal/model"
	"gorm.io/gorm"
)

func migrateDatabase(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "0",
			Migrate: func(tx *gorm.DB) error {
				return db.Exec(`CREATE TYPE scan_status AS ENUM (
					'Queued',
					'In Progress',
					'Success',
					'Failure'
				);`).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return db.Exec(`DROP TYPE scan_status;`).Error
			},
		},
		{
			ID: "1",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(model.Repository{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("repositories")
			},
		},
		{
			ID: "2",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(model.Scan{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("scans")
			},
		},
	})

	if err := m.Migrate(); err != nil {
		return fmt.Errorf("error to migrate: %w", err)
	}

	return nil
}

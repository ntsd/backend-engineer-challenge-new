package storage

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// newDatabase returns gorm database connection
func newDatabase(postgresURL string) (*gorm.DB, error) {
	var finalErr error

	// create retry to connect to database with delay 3 seconds
	for tires := 0; tires < 3; tires++ {
		db, err := gorm.Open(postgres.Open(postgresURL), &gorm.Config{})
		if err != nil {
			finalErr = err
			time.Sleep(3 * time.Second)
			continue
		}
		return db, nil
	}

	return nil, finalErr
}

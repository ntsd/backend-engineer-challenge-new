package storage

import (
	"fmt"

	"github.com/ntsd/backend-engineer-challenge-new/internal/config"
	"github.com/ntsd/backend-engineer-challenge-new/internal/model"
	"gorm.io/gorm"
)

// Storage is the storage interface
type Storage interface {
	CreateRepository(repo *model.Repository) error
	ListRepositories(repos *[]model.Repository) error
	GetRepositoryByID(repo *model.Repository, id string) error
	DeleteRepositoryByID(id string) error
	UpdateRepositoryByID(repo *model.Repository, id string) error

	CreateScan(scan *model.Scan) error
	UpdateScanByID(scan *model.Scan, id string) error
	ListQueueScans(scans *[]model.Scan) error
}

// storage is the storage struct
type storage struct {
	db *gorm.DB
}

// NewStorage create a new storage
func NewStorage(cfg config.Config) (Storage, error) {
	db, err := newDatabase(cfg.PostgresURL)
	if err != nil {
		return nil, fmt.Errorf("error to create database: %w", err)
	}

	// open database debug log
	db = db.Debug()

	if err := migrateDatabase(db); err != nil {
		return nil, fmt.Errorf("error to migrate database: %w", err)
	}

	return &storage{
		db: db,
	}, nil
}

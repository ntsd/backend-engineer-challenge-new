package storage

import (
	"github.com/ntsd/backend-engineer-challenge-new/internal/model"
	"gorm.io/gorm"
)

// CreateRepository use to save repository to the database
func (s *storage) CreateRepository(repo *model.Repository) error {
	return s.db.Create(repo).Error
}

// ListRepositories use to list all repositories from the database
func (s *storage) ListRepositories(repos *[]model.Repository) error {
	return s.db.Find(repos).Error
}

// GetRepositoryByID use to get one repository by id
func (s *storage) GetRepositoryByID(repo *model.Repository, id string) error {
	return s.db.Where(&model.Repository{ID: id}).
		Preload("Scans", func(db *gorm.DB) *gorm.DB {
			return db.Order("scans.created_at DESC")
		}).
		First(repo).
		Error
}

// DeleteRepositoryByID use to delete one repository by id
func (s *storage) DeleteRepositoryByID(id string) error {
	return s.db.Where(&model.Repository{ID: id}).Delete(&model.Repository{}).Error
}

// UpdateRepositoryByID use to update one repository by id
func (s *storage) UpdateRepositoryByID(repo *model.Repository, id string) error {
	return s.db.Where(&model.Repository{ID: id}).Updates(repo).First(repo).Error
}

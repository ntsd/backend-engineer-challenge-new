package model

import (
	"time"

	"gorm.io/gorm"
)

// Repository the repository struct
type Repository struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name"`
	URL       string         `json:"url"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt,omitempty"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Scans []Scan `json:"scans,omitempty" gorm:"foreignKey:RepositoryID"`
}

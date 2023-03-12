package model

import (
	"time"

	"gorm.io/gorm"
)

// Scan the scan struct
type Scan struct {
	ID            string         `json:"id" gorm:"primaryKey"`
	RepositoryID  string         `json:"repositoryID" gorm:"index:,type:hash"`
	RepositoryURL string         `json:"repositoryURL"`
	Status        ScanStatus     `json:"status" gorm:"type:scan_status;default:'Queued';index:,type:hash"`
	Findings      Findings       `json:"findings,omitempty" gorm:"type:jsonb"`
	StartedAt     time.Time      `json:"startedAt,omitempty"`
	FinishedAt    time.Time      `json:"finishedAt,omitempty"`
	CreatedAt     time.Time      `json:"createdAt" gorm:"index:,sort:desc"`
	UpdatedAt     time.Time      `json:"updatedAt,omitempty"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

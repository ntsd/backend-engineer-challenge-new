package storage

import "github.com/ntsd/backend-engineer-challenge-new/internal/model"

// CreateScan use to create scan to the database
func (s *storage) CreateScan(scan *model.Scan) error {
	return s.db.Create(scan).Error
}

// UpdateScan use to update one scan by id
func (s *storage) UpdateScanByID(scan *model.Scan, id string) error {
	return s.db.Where(&model.Scan{ID: id}).Updates(scan).Error
}

// ListQueueScans list all scans in queue
func (s *storage) ListQueueScans(scans *[]model.Scan) error {
	return s.db.Where(&model.Scan{Status: model.ScanStatusQueue}).Find(scans).Error
}

package scanner

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/ntsd/backend-engineer-challenge-new/internal/config"
	"github.com/ntsd/backend-engineer-challenge-new/internal/model"
	"github.com/ntsd/backend-engineer-challenge-new/internal/storage"
	"github.com/sirupsen/logrus"
)

type Scanner interface {
	CreateScan(repo model.Repository) (*model.Scan, error)
}

type scanner struct {
	store  storage.Storage
	logger *logrus.Logger

	// scanChan FIFO queue for scan
	scanChan chan model.Scan
}

func NewScanner(cfg config.Config, store storage.Storage) (Scanner, error) {
	scanner := &scanner{
		store:  store,
		logger: logrus.New(),
		// buffer channel size need to optimize by the machine ram, 1000 should be enough for now
		scanChan: make(chan model.Scan, 1000),
	}

	scanner.startWorkers(cfg.ScanWorkers)

	return scanner, nil
}

func (s *scanner) startWorkers(maxWorkers int) error {
	// create worker pool
	for workerNum := 1; workerNum <= maxWorkers; workerNum++ {
		go s.newWorker(workerNum)
	}

	// load all queue from the database
	scans := []model.Scan{}
	if err := s.store.ListQueueScans(&scans); err != nil {
		return fmt.Errorf("error to list queue scan: %w", err)
	}
	for _, scan := range scans {
		s.scanChan <- scan
	}

	return nil
}

// CreateScan create scan and put to queue
func (s *scanner) CreateScan(repo model.Repository) (*model.Scan, error) {
	scan := model.Scan{
		ID:            uuid.New().String(),
		RepositoryID:  repo.ID,
		RepositoryURL: repo.URL,
		Status:        model.ScanStatusQueue,
	}
	if err := s.store.CreateScan(&scan); err != nil {
		return nil, fmt.Errorf("error creating scan: %v", err)
	}

	if len(s.scanChan) >= cap(s.scanChan) {
		s.logger.Warn("scan buffer channel is full")
		return &scan, nil
	}

	s.scanChan <- scan
	return &scan, nil
}

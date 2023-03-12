package scanner

import (
	"time"

	"github.com/ntsd/backend-engineer-challenge-new/internal/model"
	"github.com/sirupsen/logrus"
)

// newWorker will start scanner worker and listens to `s.scanChan` channel
func (s *scanner) newWorker(workerNum int) {
	logger := s.logger.WithField("workerNum", workerNum)
	for scan := range s.scanChan {
		startTime := time.Now()
		logger = logger.WithFields(logrus.Fields{
			"startTime":     startTime,
			"repositoryURL": scan.RepositoryURL,
		})

		// update scan status to inprogress
		if err := s.store.UpdateScanByID(&model.Scan{
			StartedAt: startTime,
			Status:    model.ScanStatusInProgress,
		}, scan.ID); err != nil {
			logger.Error("error updating scan")
			continue
		}

		logger.Info("start scanning")
		finding, err := s.scan(logger, scan)

		finishTime := time.Now()
		logger = logger.WithField("finishTime", finishTime)

		if err != nil {
			logger.WithField("error", err).Info("scan failed")

			// update scan status to failure
			if err := s.store.UpdateScanByID(&model.Scan{
				FinishedAt: finishTime,
				Status:     model.ScanStatusFailure,
			}, scan.ID); err != nil {
				logger.Error("error updating scan")
				continue
			}
			continue
		}

		logger.Info("scan success")

		// update scan status to success
		if err := s.store.UpdateScanByID(&model.Scan{
			FinishedAt: finishTime,
			Status:     model.ScanStatusSuccess,
			Findings:   finding,
		}, scan.ID); err != nil {
			logger.Error("error updating scan")
			continue
		}
	}
}

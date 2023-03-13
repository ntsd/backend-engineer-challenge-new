package scanner

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ntsd/backend-engineer-challenge-new/internal/model"
	"github.com/ntsd/backend-engineer-challenge-new/internal/storage"
	"github.com/ntsd/backend-engineer-challenge-new/internal/storage/mock_storage"
	"github.com/sirupsen/logrus"
)

func TestNewWorker(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	t.Log("test new worker")

	s := &scanner{
		store: func(ctrl *gomock.Controller) storage.Storage {
			mockStorage := mock_storage.NewMockStorage(ctrl)
			mockStorage.EXPECT().UpdateScanByID(gomock.Any(), gomock.Any()).Return(nil)
			return mockStorage
		}(mockCtrl),
		scanChan: make(chan model.Scan), // create non buffer channel
		logger:   logrus.New(),
	}

	go s.newWorker(10)

	s.scanChan <- model.Scan{}
}

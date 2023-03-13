package scanner

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/ntsd/backend-engineer-challenge-new/internal/model"
	"github.com/ntsd/backend-engineer-challenge-new/internal/storage"
	"github.com/ntsd/backend-engineer-challenge-new/internal/storage/mock_storage"
	"github.com/sirupsen/logrus"
)

func TestCreateScan(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	tests := []struct {
		name        string
		mockStorage func(ctrl *gomock.Controller) storage.Storage
		repo        model.Repository
		want        *model.Scan
		wantError   error
	}{
		{
			name: "success",
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				mockStorage.EXPECT().CreateScan(gomock.Any()).Return(nil)
				return mockStorage
			},
			repo: model.Repository{
				ID:  "test-repo",
				URL: "repo url",
			},
			want: &model.Scan{
				RepositoryID:  "test-repo",
				RepositoryURL: "repo url",
				Status:        model.ScanStatusQueue,
			},
		},
		{
			name: "database error",
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				mockStorage.EXPECT().CreateScan(gomock.Any()).Return(errors.New("error"))
				return mockStorage
			},
			repo: model.Repository{
				ID:  "test-repo",
				URL: "repo url",
			},
			wantError: errors.New("error creating scan: error"),
		},
	}

	for _, test := range tests {
		t.Log(test.name)

		s := &scanner{
			store:    test.mockStorage(mockCtrl),
			scanChan: make(chan model.Scan, 100),
			logger:   logrus.New(),
		}

		got, err := s.CreateScan(test.repo)

		if fmt.Sprintf("%s", err) != fmt.Sprintf("%s", test.wantError) {
			t.Fatalf("error mismatch want: `%s`, got: `%s`", test.wantError, err)
		}
		if diff := cmp.Diff(
			test.want,
			got,
			cmpopts.IgnoreFields(model.Scan{}, "ID"),
		); diff != "" {
			t.Fatalf("want mismatch(-want +got):\n%s", diff)
		}
	}
}

func TestStartWorker(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	tests := []struct {
		name        string
		mockStorage func(ctrl *gomock.Controller) storage.Storage
		wantError   error
	}{
		{
			name: "success",
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				mockStorage.EXPECT().ListQueueScans(gomock.Any()).Return(nil)
				return mockStorage
			},
		},
		{
			name: "database error",
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				mockStorage.EXPECT().ListQueueScans(gomock.Any()).Return(errors.New("error"))
				return mockStorage
			},
			wantError: errors.New("error to list queue scan: error"),
		},
	}

	for _, test := range tests {
		t.Log(test.name)

		s := &scanner{
			store:    test.mockStorage(mockCtrl),
			scanChan: make(chan model.Scan, 100),
			logger:   logrus.New(),
		}

		err := s.startWorkers(10)

		if fmt.Sprintf("%s", err) != fmt.Sprintf("%s", test.wantError) {
			t.Fatalf("error mismatch want: `%s`, got: `%s`", test.wantError, err)
		}
	}
}

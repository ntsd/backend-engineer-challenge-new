package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/ntsd/backend-engineer-challenge-new/internal/model"
	"github.com/ntsd/backend-engineer-challenge-new/internal/scanner"
	"github.com/ntsd/backend-engineer-challenge-new/internal/scanner/mock_scanner"
	"github.com/ntsd/backend-engineer-challenge-new/internal/storage"
	"github.com/ntsd/backend-engineer-challenge-new/internal/storage/mock_storage"
	"gorm.io/gorm"
)

func TestCreateScan(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockCtrl := gomock.NewController(t)

	tests := []testData{
		{
			name:   "success",
			repoID: "test-id",
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				mockStorage.EXPECT().GetRepositoryByID(gomock.Any(), "test-id").Return(nil)
				return mockStorage
			},
			mockScanner: func(ctrl *gomock.Controller) scanner.Scanner {
				mockScanner := mock_scanner.NewMockScanner(ctrl)
				mockScanner.EXPECT().CreateScan(gomock.Any()).Return(&model.Scan{}, nil)
				return mockScanner
			},
			wantData: &model.Scan{},
		},
		{
			name:   "database not found error",
			repoID: "test-id",
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				mockStorage.EXPECT().GetRepositoryByID(gomock.Any(), "test-id").Return(gorm.ErrRecordNotFound)
				return mockStorage
			},
			mockScanner: func(ctrl *gomock.Controller) scanner.Scanner {
				mockScanner := mock_scanner.NewMockScanner(ctrl)
				return mockScanner
			},
			wantError: &errorResponse{
				Code:    http.StatusNotFound,
				Message: errMessageNotFound,
			},
		},
		{
			name:   "database unknown error",
			repoID: "test-id",
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				mockStorage.EXPECT().GetRepositoryByID(gomock.Any(), "test-id").Return(gorm.ErrDuplicatedKey)
				return mockStorage
			},
			mockScanner: func(ctrl *gomock.Controller) scanner.Scanner {
				mockScanner := mock_scanner.NewMockScanner(ctrl)
				return mockScanner
			},
			wantError: &errorResponse{
				Code:    http.StatusInternalServerError,
				Message: errMessageInternal,
			},
		},
		{
			name:   "scan error",
			repoID: "test-id",
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				mockStorage.EXPECT().GetRepositoryByID(gomock.Any(), "test-id").Return(nil)
				return mockStorage
			},
			mockScanner: func(ctrl *gomock.Controller) scanner.Scanner {
				mockScanner := mock_scanner.NewMockScanner(ctrl)
				mockScanner.EXPECT().CreateScan(gomock.Any()).Return(nil, errors.New("error"))
				return mockScanner
			},
			wantError: &errorResponse{
				Code:    http.StatusInternalServerError,
				Message: errMessageInternal,
			},
		},
	}

	for _, test := range tests {
		t.Log(test.name)

		// create mock handler and gin context
		h, ctx := mockHandlerAndContext(test, mockCtrl)

		// run test
		got, err := h.createScan(ctx)

		// validate response
		validateResponse(t, test, got, err)
	}
}

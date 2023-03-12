package handler

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/ntsd/backend-engineer-challenge-new/internal/model"
	"github.com/ntsd/backend-engineer-challenge-new/internal/storage"
	"github.com/ntsd/backend-engineer-challenge-new/internal/storage/mock_storage"
	"gorm.io/gorm"
)

func TestCreateRepository(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockCtrl := gomock.NewController(t)

	tests := []testData{
		{
			name: "success",
			requestBody: `{
				"name": "scan-test",
				"url": "https://github.com/ntsd/scan-test"
			}`,
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				mockStorage.EXPECT().CreateRepository(gomock.Any()).Return(nil)
				return mockStorage
			},
			wantData: &model.Repository{
				Name: "scan-test",
				URL:  "https://github.com/ntsd/scan-test",
			},
		},
		{
			name: "bining error",
			requestBody: `{
				"name": "scan-test"
			}`,
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				return mockStorage
			},
			wantError: &errorResponse{
				Code:    http.StatusBadRequest,
				Message: errMessageBadRequest,
			},
		},
		{
			name: "database error",
			requestBody: `{
				"name": "scan-test",
				"url": "https://github.com/ntsd/scan-test"
			}`,
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				mockStorage.EXPECT().CreateRepository(gomock.Any()).Return(gorm.ErrDuplicatedKey)
				return mockStorage
			},
			wantError: &errorResponse{
				Code:    http.StatusUnprocessableEntity,
				Message: errMessageDuplicate,
			},
		},
	}

	for _, test := range tests {
		t.Log(test.name)

		// create mock handler and gin context
		h, ctx := mockHandlerAndContext(test, mockCtrl)

		// run test
		got, err := h.createRepository(ctx)

		// validate response
		validateResponse(t, test, got, err)
	}
}

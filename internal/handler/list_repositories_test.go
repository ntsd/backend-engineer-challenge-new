package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/ntsd/backend-engineer-challenge-new/internal/model"
	"github.com/ntsd/backend-engineer-challenge-new/internal/storage"
	"github.com/ntsd/backend-engineer-challenge-new/internal/storage/mock_storage"
)

func TestListRepositories(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockCtrl := gomock.NewController(t)

	tests := []testData{
		{
			name: "success",
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				mockStorage.EXPECT().ListRepositories(gomock.Any()).Return(nil)
				return mockStorage
			},
			wantData: []model.Repository{},
		},
		{
			name: "database unknown error",
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				mockStorage.EXPECT().ListRepositories(gomock.Any()).Return(errors.New(""))
				return mockStorage
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
		got, err := h.listRepositories(ctx)

		// validate response
		validateResponse(t, test, got, err)
	}
}

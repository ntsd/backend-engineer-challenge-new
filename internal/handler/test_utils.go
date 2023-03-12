package handler

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/ntsd/backend-engineer-challenge-new/internal/model"
	"github.com/ntsd/backend-engineer-challenge-new/internal/scanner"
	"github.com/ntsd/backend-engineer-challenge-new/internal/storage"
)

// testData struct for handler test
type testData struct {
	name        string
	requestBody string
	repoID      string
	mockStorage func(ctrl *gomock.Controller) storage.Storage
	mockScanner func(ctrl *gomock.Controller) scanner.Scanner
	wantData    any
	wantError   *errorResponse
}

// mockHandlerAndContext create mock handler and gin context
func mockHandlerAndContext(test testData, mockCtrl *gomock.Controller) (*handler, *gin.Context) {
	// initial handler function and mock
	h := &handler{}
	if test.mockScanner != nil {
		h.scanner = test.mockScanner(mockCtrl)
	}
	if test.mockStorage != nil {
		h.store = test.mockStorage(mockCtrl)
	}

	// init gin context
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	// mock request body
	if test.requestBody != "" {
		ctx.Request = httptest.NewRequest("GET", "/", io.NopCloser(strings.NewReader(test.requestBody)))
	}

	// mock repo id
	if test.repoID != "" {
		ctx.Params = append(ctx.Params, gin.Param{
			Key:   "repoID",
			Value: test.repoID,
		})
	}

	return h, ctx
}

// validateResponse check diff data response and error response, fatal when unmatched
func validateResponse(
	t *testing.T,
	test testData,
	got any,
	err *errorResponse,
) {
	// validate error response
	if diff := cmp.Diff(
		err,
		test.wantError,
		cmpopts.IgnoreFields(errorResponse{}, "Err"),
	); diff != "" {
		t.Fatalf("want error mismatch(-want +got):\n%s", diff)
	}
	// validate data response
	if diff := cmp.Diff(
		test.wantData,
		got,
		cmpopts.IgnoreFields(model.Repository{}, "ID"),
	); diff != "" {
		t.Fatalf("want data mismatch(-want +got):\n%s", diff)
	}
}

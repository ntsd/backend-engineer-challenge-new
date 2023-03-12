package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func TestRespondHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		handler     func(*gin.Context) (any, *errorResponse)
		httpStatus  int
		reponseBody string
	}{
		{
			name: "success",
			handler: func(*gin.Context) (any, *errorResponse) {
				return nil, nil
			},
			httpStatus:  http.StatusOK,
			reponseBody: "{}",
		},
		{
			name: "success with data",
			handler: func(*gin.Context) (any, *errorResponse) {
				return struct {
					Key string `json:"key"`
				}{
					Key: "value",
				}, nil
			},
			httpStatus:  http.StatusOK,
			reponseBody: `{"data":{"key":"value"}}`,
		},
		{
			name: "response error",
			handler: func(*gin.Context) (any, *errorResponse) {
				return nil, &errorResponse{
					Code:    http.StatusBadRequest,
					Message: "error",
				}
			},
			httpStatus:  http.StatusBadRequest,
			reponseBody: `{"message":"error"}`,
		},
		{
			name: "response error no message",
			handler: func(*gin.Context) (any, *errorResponse) {
				return nil, &errorResponse{
					Code: http.StatusInternalServerError,
				}
			},
			httpStatus:  http.StatusInternalServerError,
			reponseBody: `{"message":""}`,
		},
		{
			name: "response error no code",
			handler: func(*gin.Context) (any, *errorResponse) {
				return nil, &errorResponse{
					Message: "error",
				}
			},
			httpStatus:  http.StatusInternalServerError,
			reponseBody: `{"message":"error"}`,
		},
	}

	for _, test := range tests {
		t.Log(test.name)

		// init handler
		h := &handler{
			logger: logrus.New(),
		}

		// init gin
		r := gin.Default()
		r.GET("/test", h.responseHandlerFunc(test.handler))

		req, err := http.NewRequest(http.MethodGet, "/test", nil)
		if err != nil {
			t.Fatalf("error to  create request: %v", err)
		}

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// compare http status code
		if w.Code != test.httpStatus {
			t.Fatalf("status mismatch want: `%d`, got: `%d`", test.httpStatus, w.Code)
		}

		// compare response body
		body, err := io.ReadAll(w.Body)
		if err != nil {
			t.Fatalf("error to read all body: %v", err)
		}
		if string(body) != test.reponseBody {
			t.Fatalf("response body mismatch want: `%s`, got: `%s`", test.reponseBody, body)
		}
	}
}

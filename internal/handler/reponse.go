package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	errMessageBadRequest string = "bad request"
	errMessageNotFound   string = "not found"
	errMessageDuplicate  string = "duplicate"
	errMessageInternal   string = "something went wrong"
)

// successResponse use for response the error
type successResponse struct {
	// Data response data can be array or struct
	Data any `json:"data,omitempty"`
}

// errorResponse use for response the error
type errorResponse struct {
	// Code HTTP status code
	Code int `json:"-"`
	// Message error message
	Message string `json:"message"`
	Err     error  `json:"-"`
}

// Error is error implement of the standard error interface
func (e *errorResponse) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Err.Error())
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// responseHandlerFunc will handle JSON response for Gin
func (h *handler) responseHandlerFunc(handler func(*gin.Context) (any, *errorResponse)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, errResponse := handler(ctx)
		if errResponse != nil {
			h.logger.WithField("error", errResponse.Error()).Error("response error")

			if errResponse.Code == 0 {
				ctx.JSON(http.StatusInternalServerError, errResponse)
				return
			}
			ctx.JSON(errResponse.Code, errResponse)
			return
		}
		if data == nil {
			ctx.JSON(http.StatusOK, successResponse{})
			return
		}
		ctx.JSON(http.StatusOK, successResponse{
			Data: data,
		})
	}
}

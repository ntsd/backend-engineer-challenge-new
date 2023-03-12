package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ntsd/backend-engineer-challenge-new/internal/model"
)

// createRepositoryRequest json request for creating repository
type createRepositoryRequest struct {
	Name string `json:"name" binding:"required"`
	URL  string `json:"url" binding:"required"`
}

// createRepository HTTP handler for creating repository
func (h *handler) createRepository(c *gin.Context) (any, *errorResponse) {
	req := &createRepositoryRequest{}
	if err := c.BindJSON(req); err != nil {
		return nil, &errorResponse{
			Code:    http.StatusBadRequest,
			Message: errMessageBadRequest,
			Err:     err,
		}
	}

	repo := &model.Repository{
		ID:   uuid.New().String(),
		Name: req.Name,
		URL:  req.URL,
	}
	if err := h.store.CreateRepository(repo); err != nil {
		err = fmt.Errorf("error to create repository: %w", err)
		switch {
		case errors.Is(err, gorm.ErrDuplicatedKey):
			return nil, &errorResponse{
				Code:    http.StatusUnprocessableEntity,
				Message: errMessageDuplicate,
				Err:     err,
			}
		}
		return nil, &errorResponse{
			Code:    http.StatusInternalServerError,
			Message: errMessageInternal,
			Err:     err,
		}
	}

	return repo, nil
}

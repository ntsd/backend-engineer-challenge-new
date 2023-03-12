package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntsd/backend-engineer-challenge-new/internal/model"
	"gorm.io/gorm"
)

// createScan HTTP handler to start scan repository
func (h *handler) createScan(c *gin.Context) (any, *errorResponse) {
	repoID := c.Param("repoID")

	repo := model.Repository{}
	if err := h.store.GetRepositoryByID(&repo, repoID); err != nil {
		err = fmt.Errorf("error to get repository by id: %w", err)
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, &errorResponse{
				Code:    http.StatusNotFound,
				Message: errMessageNotFound,
				Err:     err,
			}
		}
		return nil, &errorResponse{
			Code:    http.StatusInternalServerError,
			Message: errMessageInternal,
			Err:     err,
		}
	}

	scan, err := h.scanner.CreateScan(repo)
	if err != nil {
		return nil, &errorResponse{
			Code:    http.StatusInternalServerError,
			Message: errMessageInternal,
			Err:     fmt.Errorf("error to create scan by id: %w", err),
		}
	}

	return scan, nil
}

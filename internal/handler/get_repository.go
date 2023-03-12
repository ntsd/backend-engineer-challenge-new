package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntsd/backend-engineer-challenge-new/internal/model"
	"gorm.io/gorm"
)

func (h *handler) getRepository(c *gin.Context) (any, *errorResponse) {
	repoID := c.Param("repoID")

	repo := &model.Repository{}
	if err := h.store.GetRepositoryByID(repo, repoID); err != nil {
		err = fmt.Errorf("error to get repository by id `%s`: %w", repoID, err)
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

	return repo, nil
}

package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntsd/backend-engineer-challenge-new/internal/model"
	"gorm.io/gorm"
)

// updateRepositoryRequest json request for update repository
type updateRepositoryRequest struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// updateRepository HTTP handler for update repository
func (h *handler) updateRepository(c *gin.Context) (any, *errorResponse) {
	repoID := c.Param("repoID")

	req := &updateRepositoryRequest{}
	if err := c.BindJSON(req); err != nil {
		return nil, &errorResponse{
			Code:    http.StatusBadRequest,
			Message: errMessageBadRequest,
			Err:     err,
		}
	}

	repo := &model.Repository{
		Name: req.Name,
		URL:  req.URL,
	}
	if err := h.store.UpdateRepositoryByID(repo, repoID); err != nil {
		err = fmt.Errorf("error to get repository by id `%s`: %w", repoID, err)
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, &errorResponse{
				Code:    http.StatusNotFound,
				Message: errMessageNotFound,
				Err:     err,
			}
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

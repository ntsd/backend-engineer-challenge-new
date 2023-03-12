package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntsd/backend-engineer-challenge-new/internal/model"
)

func (h *handler) listRepositories(c *gin.Context) (any, *errorResponse) {
	repos := []model.Repository{}
	if err := h.store.ListRepositories(&repos); err != nil {
		return nil, &errorResponse{
			Code:    http.StatusInternalServerError,
			Message: errMessageInternal,
			Err:     fmt.Errorf("error to list repositories: %w", err),
		}
	}
	return repos, nil
}

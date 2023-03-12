package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// deleteRepository HTTP handler for delete repository
func (h *handler) deleteRepository(c *gin.Context) (any, *errorResponse) {
	repoID := c.Param("repoID")

	if err := h.store.DeleteRepositoryByID(repoID); err != nil {
		return nil, &errorResponse{
			Code:    http.StatusInternalServerError,
			Message: errMessageInternal,
			Err:     fmt.Errorf("error to delete repository by id: %w", err),
		}
	}

	return nil, nil
}

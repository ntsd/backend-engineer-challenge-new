package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ntsd/backend-engineer-challenge-new/internal/config"
	"github.com/ntsd/backend-engineer-challenge-new/internal/scanner"
	"github.com/ntsd/backend-engineer-challenge-new/internal/storage"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
)

// Handler is the http server handlers interface
type Handler interface {
	Serve()
}

// handler is the http server handlers struct
type handler struct {
	// addr HTTP listen address
	address string

	store   storage.Storage
	scanner scanner.Scanner

	logger *logrus.Logger
}

// NewHandler create a new handler
func NewHandler(cfg config.Config, store storage.Storage, scan scanner.Scanner) (Handler, error) {
	return &handler{
		address: fmt.Sprintf(":%s", cfg.AppPort),
		store:   store,
		scanner: scan,
		logger:  logrus.New(),
	}, nil
}

// Serve start the http server
func (h *handler) Serve() {

	r := gin.New()
	r.Use(ginlogrus.Logger(h.logger), gin.Recovery())

	r.POST("/repositories", h.responseHandlerFunc(h.createRepository))
	r.GET("/repositories", h.responseHandlerFunc(h.listRepositories))
	r.GET("/repositories/:repoID", h.responseHandlerFunc(h.getRepository))
	r.PATCH("/repositories/:repoID", h.responseHandlerFunc(h.updateRepository))
	r.DELETE("/repositories/:repoID", h.responseHandlerFunc(h.deleteRepository))
	r.POST("/repositories/:repoID/scan", h.responseHandlerFunc(h.createScan))

	if err := r.Run(h.address); err != nil {
		h.logger.Fatal(err)
	}
}

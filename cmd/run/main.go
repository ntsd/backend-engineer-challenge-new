package main

import (
	"log"

	"github.com/ntsd/backend-engineer-challenge-new/internal/config"
	"github.com/ntsd/backend-engineer-challenge-new/internal/handler"
	"github.com/ntsd/backend-engineer-challenge-new/internal/scanner"
	"github.com/ntsd/backend-engineer-challenge-new/internal/storage"
)

func main() {
	cfg := config.NewConfig()

	storage, err := storage.NewStorage(cfg)
	if err != nil {
		log.Fatalf("error creating storage: %v", err)
	}

	scanner, err := scanner.NewScanner(cfg, storage)
	if err != nil {
		log.Fatalf("error new scan: %v", err)
	}

	handler, err := handler.NewHandler(cfg, storage, scanner)
	if err != nil {
		log.Fatalf("error creating handler: %v", err)
	}
	handler.Serve()
}

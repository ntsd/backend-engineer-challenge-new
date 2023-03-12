run:
	go run cmd/run/main.go

lint:
	go fmt ./...

test:
	go test ./...

mockgen:
	mockgen -source ./internal/storage/storage.go -destination ./internal/storage/mock_storage/mock_storage.go
	mockgen -source ./internal/scanner/scanner.go -destination ./internal/scanner/mock_scanner/mock_scanner.go

docker-db:
	docker compose up challenge-postgres

docker-up:
	docker compose up

docker-down:
	docker compose down

docker-remove:
	docker compose down -v

docker-build:
	docker compose build --no-cache challenge-app

APP_NAME: jille

.PHONY: backend
build: 
	go build -o bin/${APP_NAME} ./cmd/api

.PHONY: run
run:
	go run ./cmd/api/main.go

.PHONY: frontend
build:
	cd client && bun run build

dev:
	cd client && bun dev
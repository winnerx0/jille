APP_NAME: jille

.PHONY: build
build: 
	go build -o bin/${APP_NAME} ./cmd/api

.PHONY: run
run:
	go run ./cmd/api/main.go
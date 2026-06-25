run:
	go run cmd/api/main.go
test:
	go test -v ./...
build:
	go build -o bin/main cmd/api/main.go
docker:
	docker build -t payment-service:latest .
tidy:
	go mod tidy

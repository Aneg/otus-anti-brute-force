run:
	docker-compose up -d --build
stop:
	docker-compose down
build:
	docker-compose run app go build -mod=mod cmd/anty-brute-force/main.go
build-console:
	docker-compose run app go build -mod=mod cmd/management_console/main.go
test:
	docker-compose run app go test -race -count 100 ./...
test-integration:
	docker-compose run app go test -race -tags integration ./internal/web/grpc/...
test-dev:
	docker-compose run app go test -race -count 1 ./...
golangci:
	docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v1.26.0 golangci-lint run -v
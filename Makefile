run:
	docker-compose up -d --build
stop:
	docker-compose down
build:
	docker-compose run app go build -mod=mod cmd/anty-brute-force/main.go
test:
	docker-compose run app go test -race -count 100 ./...
test2:
	docker-compose run app go test ./...
golangci:
	docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v1.26.0 golangci-lint run -v
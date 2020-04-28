run:
	docker-compose up -d --build
stop:
	docker-compose down
build:
	docker-compose run app go build -mod=mod cmd/anty-brute-force/main.go
test:
	docker-compose run app go test -race -count 100 ./...

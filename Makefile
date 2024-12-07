.PHONY: create
create:
	migrate create -ext sql -dir migrations -seq $(name)

.PHONY: build.docker
build.docker:
	docker build -t vpbuyanov/gw-backend-go:latest .

.PHONY: run.docker
run.docker:
	docker compose up -d

.PHONY: clean.docker
clean.docker:
	docker stop gw-backend-go gw-postgres gw-redis
	docker rm gw-backend-go gw-postgres gw-redis

.PHONY: build
build:
	go build -o ./bin/app cmd/api/main.go

.PHONY: run
run:
	go run ./bin/app

.PHONY: clean
clean:
	rm -rf ./bin

.PHONY: restart
restart: clean build run

.PHONY: lint
lint:
	golangci-lint run --config .golangci.yml ./...

.PHONY: lint-fix
lint-fix:
	golangci-lint run --config .golangci.yml --fix ./...
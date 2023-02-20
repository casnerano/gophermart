.PHONY: build init start test migrate

build:
	cp .env.dist .env
	docker compose build
init: build migrate
start:
	docker compose up -d
	docker compose logs -f
test:
	go test -v -count=1 -cover -race ./...
migrate:
	migrate -database postgres://user:password@localhost:5432/gophermart?sslmode=disable -source file://migrations/postgres up

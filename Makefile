.PHONY: \
	build \
	up \
	down \
	build_server \
	build_client \
	build_go \
	get_quote \
	tests

build:
	docker-compose --project-directory . build

up:
	docker-compose --project-directory . up -d --no-build
	docker-compose --project-directory . ps

down:
	docker-compose --project-directory . down --remove-orphans

build_server:
	go build -o server/server server/main.go

build_client:
	go build -o client/client client/main.go

build_go: build_server build_client

get_quote:
	client/client --server-addr=:8080

tests:
	go test ./...

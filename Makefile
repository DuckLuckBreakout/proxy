.PHONY: update_db
update_db:
	docker rm -f -v proxy_db
	docker compose up

.PHONY: build
build:
	make server
	make api_server

.PHONY: server
server:
	go build -o bin/proxy -v ./cmd/proxy

.PHONY: api_server
api_server:
	go build -o bin/api -v ./cmd/api

.DEFAULT_GOAL := update_db

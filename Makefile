ifneq (,$(wildcard ./.env))
	include .env
	export
endif

APP_NAME=plant-factory
APP_NAME_SM=${APP_NAME}-small

.PHONY: help run build build/small test fmt lint docs dev dev/install \
		migrate/create migrate/up migrate/down migrate/status

## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## run: run the application
run:
	go run ./cmd/api

## build: build the application binary
build:
	go build -o tmp/$(APP_NAME) ./cmd/api

## build/small: build a small, statically linked binary for Linux
build/small:
	CGO_ENABLED=0
	GOARCH=amd64
	GOOS=linux
	go build -o tmp/$(APP_NAME_SM) -a -ldflags '-s -w' -installsuffix cgo ./cmd/api
	upx -9 tmp/$(APP_NAME_SM)

## test: run all tests
test:
	go test ./...

## fmt: format Go code
fmt:
	go fmt ./...

## lint: run linter
lint:
	swag fmt -d .
	golangci-lint --verbose run ./...

## docs: generate API documentation
docs:
	swag init -g ./cmd/api/main.go -o ./swagger

## dev: run app with live reload (air)
dev:
	air -c .air.toml

## dev/install: install development tools
dev/install:
	go install github.com/air-verse/air@latest
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/v2/cmd/swag@latest
	go install github.com/jackc/tern/v2@latest
	go install golang.org/x/tools/cmd/deadcode@latest

## check/deps: check for unused functions and variables
check/deadcode:
	deadcode ./...

## migrate/create: create a new migration (example = 'make migrate/create name=create_users')
migrate/create:
	@test -n "$(name)" || (echo "usage: make migrate/create name=create_users" && exit 1)
	tern new --migrations migrations $(name)

## migrate/up: apply all up migrations
migrate/up:
	tern migrate --config migrations/tern.conf --migrations migrations

## migrate/down: rollback last migration
migrate/down:
	tern migrate --config migrations/tern.conf --migrations migrations --destination -1

## migrate/status: show migration status
migrate/status:
	tern status --config migrations/tern.conf --migrations migrations

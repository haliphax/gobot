default: help

help:
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done

all: build test # Build and test

init: # Initialize git hooks, etc.
	@$$SHELL -c "source ~/.nvm/nvm.sh && nvm install"
	@npm ci

build: # Compile project binary
	go build

run: # Run project
	@go run .

lint: # Format and lint
	go fmt
	go fix
	npx prettier -luw .

test: # Run unit tests with coverage
	@[ -d coverage ] || mkdir coverage
	@go test -cover ./... -args -test.gocoverdir "$$PWD/coverage"
	@go tool covdata textfmt -i coverage -o coverage/profile
	@go tool cover -func coverage/profile | grep -E '^total:'

cover: # View coverage report
	@[ -f coverage/profile ] || make test
	@go tool cover -html coverage/profile

clean: # Remove cache and project binary
	go clean
	[ ! -f gobot ] || rm gobot

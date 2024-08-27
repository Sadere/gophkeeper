BUILD_DATE="$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')"
VERSION="$(shell git describe --tags --abbrev=0 | tr -d '\n')"
DOCKER_COMPOSE_FILES ?= $(shell find docker -maxdepth 1 -type f -name "*.yaml" -exec printf -- '-f %s ' {} +; echo)
BINARY_NAME=gophkeeper

## Buf build commands ##

# Update buf deps
.PHONY: buf-dep
buf-dep:
	buf dep update

# Generate protobuf .go files
.PHONY: buf-gen
buf-gen:
	buf generate

## Docker commands ##

# Build and run docker containers
.PHONY: up
up:
	docker compose ${DOCKER_COMPOSE_FILES} up --build --detach

# Stop docker containers
.PHONY: down
down:
	docker compose ${DOCKER_COMPOSE_FILES} down

# Show Docker containers info
.PHONY: ps
ps:	
	docker ps --size --all --filter "name=gophkeeper"

## Go commands ##

.PHONY: imports
imports:
	goimports -w $(find . -type f -name '*.go' -not -path "./pkg/proto/*")

.PHONY: build
build:
	GOARCH=amd64 GOOS=darwin go build -o bin/${BINARY_NAME}-darwin -ldflags="-X 'github.com/Sadere/gophkeeper/internal/client/version.buildDate=${BUILD_DATE}' -X 'github.com/Sadere/gophkeeper/internal/client/version.version=${VERSION}'" cmd/client/main.go
	GOARCH=amd64 GOOS=linux go build -o bin/${BINARY_NAME}-linux -ldflags="-X 'github.com/Sadere/gophkeeper/internal/client/version.buildDate=${BUILD_DATE}' -X 'github.com/Sadere/gophkeeper/internal/client/version.version=${VERSION}'" cmd/client/main.go
	GOARCH=amd64 GOOS=windows go build -o bin/${BINARY_NAME}-windows -ldflags="-X 'github.com/Sadere/gophkeeper/internal/client/version.buildDate=${BUILD_DATE}' -X 'github.com/Sadere/gophkeeper/internal/client/version.version=${VERSION}'" cmd/client/main.go
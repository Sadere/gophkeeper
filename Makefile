GOCMD=go
DOCKER_COMPOSE_FILES ?= $(shell find docker -maxdepth 1 -type f -name "*.yaml" -exec printf -- '-f %s ' {} +; echo)

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
	goimports -w .
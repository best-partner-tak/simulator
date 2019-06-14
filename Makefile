NAME:=simulator
DOCKER_REPOSITORY:=controlplane
LAUNCH_DOCKER_IMAGE_NAME:=$(DOCKER_REPOSITORY)/$(NAME)-launch
VERSION:=0.1-dev

SHELL := /usr/bin/env bash

.DEFAULT_GOAL := help

.PHONY: all
all: test

.PHONY: check
check: ## Check required system packages are installed
	@command -v "dgoss" > /dev/null 2>&1 || echo >2 "Couldn't find dgoss - please install goss >= v0.37.0"

.PHONY: run
run: build
	docker run --rm -it $(LAUNCH_DOCKER_IMAGE_NAME):$(VERSION) bash

.PHONY: run
exec: build
	docker run -v $(SIMULATOR_AWS_CREDS_PATH):/app/credentials --rm -it $(LAUNCH_DOCKER_IMAGE_NAME):$(VERSION) $(CMD)

.PHONY: build
build: lint ## Builds the launch container
	@docker build -t $(LAUNCH_DOCKER_IMAGE_NAME):$(VERSION) .

.PHONY: test
test: check ## Run the tests
	@cd test && ./test.sh

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


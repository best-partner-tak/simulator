NAME := simulator
GITHUB_ORG := controlplaneio
DOCKER_HUB_ORG := controlplane
VERSION := 0.1-dev

include prelude.mk

SIMULATOR_CREDS_PATH := $(HOME)/.aws/credentials
SIMULATOR_KEY_PATH := $(HOME)/.ssh/

.DEFAULT_GOAL := help

.PHONY: all
all: test

.PHONY: run
run: build ## Runs the simulator
	docker run                                             \
		-v $(SIMULATOR_AWS_CREDS_PATH):/app/credentials      \
		-v $(SIMULATOR_KEY_PATH):/home/launch-user/.ssh      \
		--rm -it $(CONTAINER_NAME_LATEST)                    \
		bash

.PHONY: build
build: ## Builds the launch container
	@docker build -t $(CONTAINER_NAME_LATEST) .

.PHONY: test
test: build ## Run the tests
	@docker run -v                                   \
		$(SIMULATOR_AWS_CREDS_PATH):/app/credentials   \
		--rm -t $(CONTAINER_NAME_LATEST)               \
		goss validate

.PHONY: infra-init
infra-init:
	@pushd terraform/deployments/AwsSimulatorStandalone; terraform init; popd

.PHONY: infra-checkvars
infra-checkvars:
	@test -f terraform/deployments/AwsSimulatorStandalone/settings/bastion.tfvars || \
		(echo Please create terraform/settings/bastion.tfvars && exit 1)

.PHONY: infra-plan
infra-plan: infra-init infra-checkvars
	@cd terraform/deployments/AwsSimulatorStandalone; terraform plan -var-file=settings/bastion.tfvars;

.PHONY: infra-apply
infra-apply: infra-init infra-checkvars
	@cd terraform/deployments/AwsSimulatorStandalone; terraform apply -var-file=settings/bastion.tfvars -auto-approve;

.PHONY: infra-destroy
infra-destroy: infra-init infra-checkvars
	@cd terraform/deployments/AwsSimulatorStandalone; terraform destroy -var-file=settings/bastion.tfvars;

.PHONY: help
help: ## This message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


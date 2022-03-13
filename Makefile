#
# oerc, alias oer-collector
# Copyright (C) 2021 emschu[aet]mailbox.org
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as
# published by the Free Software Foundation, either version 3 of the
# License, or (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public
# License along with this program.
# If not, see <https://www.gnu.org/licenses/>.
SHELL := /bin/bash

APP_VERSION_DOT = "0.9.17"
APP_VERSION_STR = "0-9-17"

GO := GO111MODULE=on go
GO_PATH = $(shell $(GO) env GOPATH)
GO_REVIVE = $(GO_PATH)/bin/revive
GO_RICE = $(GO_PATH)/bin/rice

SCHEMATHESIS_BIN = ~/.local/bin/schemathesis

OPENAPI_TOOLS_VERSION = 5.1.1

# TODO: add spec conversion script
# TODO add clean target

.PHONY: all
all: help

.DEFAULT_GOAL:=help

.PHONY: help
help: ## show this help
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

.PHONY: clean
clean: ## clean up project
	-rm -rf bin
	-mkdir -p bin

.PHONY: setup
setup: ## install required project and (dev) dependencies
	$(GO) mod download
	$(GO) get -u github.com/GeertJohan/go.rice
	$(GO) get -u github.com/GeertJohan/go.rice/rice
	$(GO) get -u github.com/mgechev/revive
	if [ ! -f openapi-generator-cli.jar ]; then curl -L -o openapi-generator-cli.jar -L https://repo1.maven.org/maven2/org/openapitools/openapi-generator-cli/$(OPENAPI_TOOLS_VERSION)/openapi-generator-cli-$(OPENAPI_TOOLS_VERSION).jar; fi
	pip install --user schemathesis

.PHONY: build-frontend
frontend: ## build the frontend and the static rice box file
	cd client; npm run build-prod
	$(GO_RICE) embed-go

.PHONY: build
build: ## build dev version of application
	$(GO) build -race -o bin/oerc
	export CGO_ENABLED=0 ; export GOOS=linux ; $(GO) build -o bin/oerc-docker

.PHONY: lint
lint: ## linting the code
	@$(GO) fmt ./...
	@$(GO_REVIVE) .

.PHONY: lint-fix
lint-fix: ## lint-fix the code
	$(GO) fix ./...

.PHONY: test
test: ## run unit, integration and api tests
	@$(GO) test -v -race ./...
	@$(GO) test -v -trace=/dev/null .

.PHONY: integration-test-prepare
integration-test-prepare: ## start (local) oerc server to run integration tests against
	@$(GO) build -race -o bin/oerc
	@if [[ -a server.PID ]]; then kill -9 "$$(cat server.PID)" || rm server.PID || true; fi
	@bin/oerc -c config/.oerc.dist.yaml server & echo $$! > server.PID
	-sleep 5

.PHONY: integration-test
integration-test: ## run OpenAPI schema conformity HTTP tests
	@$(SCHEMATHESIS_BIN) run -x --show-errors-tracebacks --hypothesis-deadline 7500 --validate-schema true -c all http://127.0.0.1:8080/spec/openapi3.json
	-if [[ -a server.PID ]]; then kill -9 "$$(cat server.PID)" || rm server.PID || true; fi

.PHONY: cover
cover: ## run unit tests with coverage output
	$(GO) test -race -coverprofile=cover.out -coverpkg=./ ./
	$(GO) tool cover -html=cover.out -o cover.html

.PHONY: spec
spec: ## run openapi spec converter from yaml -> json
	bash convert_spec.sh

# TODO
# CGO_ENABLED=0 go build -ldflags "-w" -a -o oerc .

# TODO add client compilation
.PHONY: release
release: clean ## build release packages for multiple platforms
	mkdir -p bin/windows; mkdir -p bin/linux-arm; mkdir -p bin/linux-arm64; mkdir -p bin/linux-armv7; mkdir -p bin/linux-386; mkdir -p bin/linux-amd64
	GOOS=windows GOARCH=amd64 $(GO) build -o bin/windows/oerc.exe -ldflags "-s -w"
	GOOS=linux GOARCH=arm $(GO) build -o bin/linux-arm/oerc -ldflags "-s -w"
	GOOS=linux GOARCH=arm64 $(GO) build -o bin/linux-arm64/oerc -ldflags "-s -w"
	GOOS=linux GOARCH=arm GOARM=7 $(GO) build -o bin/linux-armv7/oerc -ldflags "-s -w"
	GOOS=linux GOARCH=386 $(GO) build -o bin/linux-386/oerc -ldflags "-s -w"
	GOOS=linux GOARCH=amd64 $(GO) build -o bin/linux-amd64/oerc -ldflags "-s -w"

.PHONY: sonarscan
sonarscan: ## run sonar scanner against local sonarqube
	read -p "Enter SONAR_HOST_URL: " SONAR_HOST_URL && read -p "Enter SONAR_LOGIN: " SONAR_LOGIN && \
						docker run --rm --user="$$(id -u):$$(id -g)" \
						-e SONAR_HOST_URL=$$SONAR_HOST_URL -e SONAR_LOGIN=$$SONAR_LOGIN \
						-v "$$(pwd):/usr/src" sonarsource/sonar-scanner-cli


.PHONY: version
version: ## populate the current version defined in this make file
	sed -r -i 's/version\s*=\s*"([0-9]+.[0-9]+.[0-9]+)"/version       = "'$(APP_VERSION_DOT)'"/g' main.go
	sed -r -i 's/([0-9]+.[0-9]+.[0-9]+, License:)/'$(APP_VERSION_DOT)', License:/g' README.md
	sed -r -i 's/"version": "([0-9]+.[0-9]+.[0-9]+)"/"version": "'$(APP_VERSION_DOT)'"/g' client/package.json
	cd client; npm i;
	make frontend
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

# TODO: add spec conversion script
# TODO add clean target

.PHONY: all
all: help

.DEFAULT_GOAL:=help

.PHONY: help
help: ## show this help
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

.PHONY: install
install: ## install required project and (dev) dependencies
	go mod download
	go install github.com/GeertJohan/go.rice
	go install golang.org/x/lint/golint
	pip install --user schemathesis
	curl -o openapi-generator-cli.jar -L https://repo1.maven.org/maven2/org/openapitools/openapi-generator-cli/5.1.1/openapi-generator-cli-5.1.1.jar

.PHONY: build
build: ## build dev version of application
	cd client; npm run build-prod
	rice embed-go
	go build -o bin/oerc

.PHONY: lint
lint: ## linting the code
	go fmt ./...
	~/go/bin/golint .

.PHONY: test
test: ## run unit, integration and api tests
	go test -v -race ./...
	go test -v -trace=/dev/null .
	rice embed-go
	go build -o bin/oerc
	if [[ -a server.PID ]]; then kill -9 "$$(cat server.PID)" || rm server.PID || true; fi
	./bin/oerc -c ./config/.oerc.dist.yaml server & echo $$! > server.PID
	schemathesis run -x --show-errors-tracebacks --hypothesis-deadline 7500 --validate-schema true -c all http://127.0.0.1:8080/spec/openapi3.json
	if [[ -a server.PID ]]; then kill -9 "$$(cat server.PID)" || rm server.PID || true; fi

.PHONY: cover
cover: ## run unit tests with coverage output
	go test -race -coverprofile=cover.out -coverpkg=./ ./
	go tool cover -html=cover.out -o cover.html

.PHONY: spec
spec: ## run openapi spec converter from yaml -> json
	bash ./convert_spec.sh

# TODO
# CGO_ENABLED=0 go build -ldflags "-w" -a -o oerc .

# TODO add client compilation
.PHONY: release
release: ## build release packages for multiple platforms
	GOOS=windows GOARCH=amd64 go build -o bin/oerc-windows-amd64.exe -ldflags "-s -w"
	GOOS=linux GOARCH=arm go build -o bin/oerc-linux-arm -ldflags "-s -w"
	GOOS=linux GOARCH=arm64 go build -o bin/oerc-linux-arm64 -ldflags "-s -w"
	GOOS=linux GOARCH=arm GOARM=7 go build -o bin/oerc-linux-armv7 -ldflags "-s -w"
	GOOS=linux GOARCH=386 go build -o bin/oerc-linux-386 -ldflags "-s -w"
	GOOS=linux GOARCH=amd64 go build -o bin/oerc-linux-amd64 -ldflags "-s -w"

.PHONY: sonarscan
sonarscan: ## run sonar scanner against local sonarqube
	read -p "Enter SONAR_HOST_URL: " SONAR_HOST_URL && read -p "Enter SONAR_LOGIN: " SONAR_LOGIN && \
						docker run --rm --user="$$(id -u):$$(id -g)" \
						-e SONAR_HOST_URL=$$SONAR_HOST_URL -e SONAR_LOGIN=$$SONAR_LOGIN \
						-v "$$(pwd):/usr/src" sonarsource/sonar-scanner-cli
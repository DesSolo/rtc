PROJECT_NAME=rtc

VERSION=$(shell cat VERSION)

tidy:
	go mod tidy

MAIN_FILE_PATH=cmd/server/main.go
CONFIG_FILE_PATH=examples/config.yaml

run:
	CONFIG_FILE_PATH=${CONFIG_FILE_PATH} go run ${MAIN_FILE_PATH}

LOCAL_BIN := $(CURDIR)/bin
LINT_VERSION := 2.3.1
GOOSE_VERSION := v3.24.2

.prep_bin:
	mkdir -p ${LOCAL_BIN}

.install-lint:
	curl -Ls https://github.com/golangci/golangci-lint/releases/download/v${LINT_VERSION}/golangci-lint-${LINT_VERSION}-linux-amd64.tar.gz | tar xvz --strip-components=1 -C ${LOCAL_BIN} golangci-lint-${LINT_VERSION}-linux-amd64/golangci-lint

.install-goose:
	curl -Ls https://github.com/pressly/goose/releases/download/${GOOSE_VERSION}/goose_linux_x86_64 --output ${LOCAL_BIN}/goose
	chmod +x ${LOCAL_BIN}/goose

install-deps: \
	.prep_bin \
	.install-lint \
	.install-goose

lint: $(LINT_BIN)
	$(LOCAL_BIN)/golangci-lint run

MIGRATIONS_PATH=migrations

migration-up:
	$(LOCAL_BIN)/goose -dir ${MIGRATIONS_PATH} postgres "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=rtc sslmode=disable" up

migration-create:
	$(LOCAL_BIN)/goose -dir ${MIGRATIONS_PATH} create auto sql

RELEASE_DIRECTORY = ${LOCAL_BIN}/release

clear-release:
	rm -rf ${RELEASE_DIRECTORY}

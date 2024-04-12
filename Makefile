# Change these variables as necessary.
MAIN_PACKAGE_PATH := ./cmd/web
BINARY_NAME := snippetbox
.DEFAULT_GOAL := help
# ==================================================================================== #
# HELPERS
# ==================================================================================== #


.PHONY: help
help: ## help: print this help message
	@echo  "\nUsage:\n  make \033[36m<target>\033[0m"
	@awk 'BEGIN {FS = ":.*##"; printf } /^[a-zA-Z_0-9-]+(\/[a-zA-Z_0-9-]+)?:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST) | sort

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

.PHONY: tidy
tidy: ## format code and tidy modfile
	go fmt ./...
	go mod tidy -v


.PHONY: audit
audit: ## run quality control checks
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -race -buildvcs -vet=off ./...


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #


.PHONY: test
test: ## run all tests
	go test -tags test_all -v ${MAIN_PACKAGE_PATH}


.PHONY: test/e2e
test/e2e: ## run end to end tests
	go test -run="E2E" -race -v ${MAIN_PACKAGE_PATH}

.PHONY: test/unit
test/unit: ## run unit tests
	go test -short -race -v ${MAIN_PACKAGE_PATH}

.PHONY: test/cover
test/cover: ## run all tests and display coverage
	go test -race -tags test_all -v -coverprofile=/tmp/coverage.out ${MAIN_PACKAGE_PATH}
	go tool cover -html=/tmp/coverage.out

.PHONY: build
build: ## build the application
	go build -o=/tmp/bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

.PHONY: run/app
run/app:  ## run the  application
	go run ./cmd/web

.PHONY: run/db
run/db:  ## run the  database
	docker-compose up

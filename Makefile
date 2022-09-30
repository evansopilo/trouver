# ============================================================================ #
# HELPERS
# ============================================================================ #

## help: print this help message
.PHONY: help
help:
	@echo "Usage:"
	@sed -n "s/^##//p" ${MAKEFILE_LIST} | column -t -s ":" | sed -e "s/^/ /"

.PHONY: confirm
confirm:
	@echo -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]

# ============================================================================ #
# DEVELOPMENT
# ============================================================================ #

## run/api: run the cmd/api application
.PHONY: cmd
cmd:
	go run ./cmd

# ============================================================================ #
# QUALITY CONTROL
# ============================================================================ #

## audit: tidy and vendor dependencies and format, vet and test all code
.PHONY: audit
audit: vendor
	@echo "Formatting code..."
	go fmt ./...
	@echo "Vetting code..."
	go vet ./...
	staticcheck ./...
	@echo "Running tests..."
	go test -race -vet=off ./...

## coverage: go test coverage
.PHONY: coverage
coverage:
	@echo "Running test coverage ..."
	go test -cover ./...
	go test -covermode=count -coverprofile=/tmp/profile.out ./...
	go tool cover -html=/tmp/profile.out

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo "Tidying and verifying module dependencies..."
	go mod tidy
	go mod verify
	@echo "Vendoring dependencies..."
	go mod vendor

# ============================================================================ #
# BUILD
# ============================================================================ #

## build/api: build the cmd/api application
.PHONY: build/cmd
build/cmd:
	@echo "Building cmd/..."
	go build -o=./bin/ ./cmd/
	GOOS=linux GOARCH=amd64 go build -o=./bin/linux_amd64/ ./cmd/


# Simple Makefile for a Go project

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

no-dirty:
	@test -z "$(shell git status --porcelain)"


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #
audit: test
	go mod tidy -diff
	go mod verify
	test -z "$(shell gofmt -l .)" 
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

# Test the application
test:
	@echo "Testing..."
	@go test -v -race -buildvcs ./...

# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

tidy:
	go mod tidy -v
	go fmt ./...

# Build the application
all: build audit


build:
	@echo "Building..."
	@templ generate
	
	@go build -o tmp/main cmd/api/main.go

# Run the application
run: watch

docker-run:
	@docker compose -f compose.yml up -d --no-build

docker-build:
	@docker compose -f compose.yml up --build -d

docker-down:
	@docker compose -f compose.yml down --remove-orphans


# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

gen-docs:
	@swag init -g ./api/main.go -d cmd,internal && swag fmt

# ==================================================================================== #
# OPERATIONS
# ==================================================================================== #

push: confirm audit no-dirty
	git push

## production/deploy: deploy the application to production
production/deploy: confirm audit no-dirty
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=/tmp/bin/linux_amd64/${binary_name} ${main_package_path}
	upx -5 /tmp/bin/linux_amd64/${binary_name}
	# Include additional deployment steps here...


.PHONY: all build run test clean watch tidy docker-run docker-build docker-down itest templ-install audit confirm no-dirty help test/cover push production/deploy gen-docs


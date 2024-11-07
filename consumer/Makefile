
BINARY_FILE := "main"
COVERAGE_FILE := "cover.out"

.PHONY: all
all: build test run

.PHONY: build
build:
	@echo "Building..."
	@go build -o $(BINARY_FILE) cmd/main.go

.PHONY: test
test:
	@echo "Testing..."
	@go test ./... -v -count=1

.PHONY: run
run:
	@echo "Running..."
	@go run cmd/main.go

.PHONY: cover
cover:
	@echo "Testing..."
	@go test ./... -v -coverprofile $(COVERAGE_FILE)
	@go tool cover -html=$(COVERAGE_FILE)

.PHONY: clean
clean:
	@echo "Removing $(BINARY_FILE)..."
	@rm -f $(BINARY_FILE)
	@echo "Removing $(COVERAGE_FILE)..."
	@rm -f $(COVERAGE_FILE)

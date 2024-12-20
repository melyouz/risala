
BINARY_FILE := "main"
COVERAGE_FILE := "cover.out"

.PHONY: all
all: build test run

.PHONY: build
build:
	@echo "Building..."
	@go build -o $(BINARY_FILE) cmd/api/main.go

.PHONY: test
test:
	@echo "Testing..."
	@go test ./... -v -count=1

.PHONY: race
race:
	@echo "Testing (-race)..."
	@go test ./... -v -race -count=1

.PHONY: run
run:
	@echo "Running..."
ifeq ($(WITH_SAMPLE_DATA),1)
	@go run cmd/api/main.go --with-sample-data
else
	@go run cmd/api/main.go
endif

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

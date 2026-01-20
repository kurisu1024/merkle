.PHONY: test bench fmt tidy clean all help

# Default target
all: fmt tidy test

# Run all tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run all benchmarks
bench:
	go test -bench=. -benchmem ./...

# Run benchmarks at different scales
bench-small:
	go test -bench='_(10|100)$$' -benchmem ./...

bench-medium:
	go test -bench='_1000$$' -benchmem ./...

bench-large:
	go test -bench='_(10000|100000)$$' -benchmem ./...

# Format code
fmt:
	go fmt ./...

# Tidy dependencies
tidy:
	go mod tidy

# Clean build artifacts and coverage files
clean:
	rm -f coverage.out coverage.html
	go clean

# Help
help:
	@echo "Available targets:"
	@echo "  make test              - Run all tests"
	@echo "  make test-coverage     - Run tests with coverage report"
	@echo "  make bench             - Run all benchmarks"
	@echo "  make bench-small       - Run benchmarks with small datasets (10, 100)"
	@echo "  make bench-medium      - Run benchmarks with medium datasets (1K)"
	@echo "  make bench-large       - Run benchmarks with large datasets (10K, 100K)"
	@echo "  make fmt               - Format code with go fmt"
	@echo "  make tidy              - Tidy dependencies with go mod tidy"
	@echo "  make clean             - Clean build artifacts"
	@echo "  make all               - Run fmt, tidy, and test (default)"
	@echo "  make help              - Show this help message"

.PHONY: all build frontend backend clean install-deps run dev test

all: build

# Install all dependencies
install-deps:
	@echo "Installing Go dependencies..."
	go mod tidy
	@echo "Installing Node.js dependencies..."
	cd frontend && npm install

# Build frontend
frontend:
	@echo "Building frontend..."
	cd frontend && npm run build

# Build backend (with embedded frontend)
backend: frontend
	@echo "Building Go binary..."
	go build -o bin/parse-dmarc ./cmd/parse-dmarc

# Full build
build: frontend backend
	@echo "Build complete! Binary available at ./bin/parse-dmarc"

# Run in development mode
dev:
	@echo "Starting development server..."
	go run ./cmd/parse-dmarc -config=config.json

# Generate sample config
config:
	go run ./cmd/parse-dmarc -gen-config

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf internal/api/dist/
	rm -rf frontend/node_modules/
	rm -f dmarc.db
	rm -f config.json

# Run tests
test:
	go test -v ./...

# Install binary to system
install: build
	@echo "Installing to /usr/local/bin..."
	sudo cp bin/parse-dmarc /usr/local/bin/

# Development frontend server
frontend-dev:
	cd frontend && npm run dev

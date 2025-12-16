.PHONY: help dev build docker-build docker-run clean install

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

install: ## Install dependencies for both backend and frontend
	@echo "Installing backend dependencies..."
	cd backend && go mod download
	@echo "Installing frontend dependencies..."
	cd frontend && npm install
	@echo "Dependencies installed!"

dev-backend: ## Run backend in development mode
	cd backend && go run cmd/server/main.go

dev-frontend: ## Run frontend in development mode
	cd frontend && npm run dev

build-backend: ## Build backend binary
	cd backend && go build -o ../bin/server cmd/server/main.go

build-frontend: ## Build frontend for production
	cd frontend && npm run build

build: build-backend build-frontend ## Build both backend and frontend

docker-build: ## Build Docker image
	docker build -t sticktoss:latest .

docker-run: ## Run Docker container
	docker run -p 8080:8080 \
		-e DB_DRIVER=sqlite \
		-e DATABASE_URL=sticktoss.db \
		sticktoss:latest

docker-compose-up: ## Start all services with docker-compose
	docker-compose up --build

docker-compose-down: ## Stop all services
	docker-compose down

clean: ## Clean build artifacts
	rm -rf backend/bin
	rm -rf frontend/dist
	rm -rf frontend/node_modules
	rm -f *.db
	docker-compose down -v

test-backend: ## Run backend tests
	cd backend && go test ./...

fmt: ## Format code
	cd backend && go fmt ./...
	cd frontend && npm run format || true

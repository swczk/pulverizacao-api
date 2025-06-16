# Pulverização API Makefile

# Variables
DOCKER_USERNAME ?= your-dockerhub-username
IMAGE_NAME = pulverizacao-api
TAG ?= latest
COMPOSE_FILE = docker-compose.yml

# Help
.PHONY: help
help: ## Show this help
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Development
.PHONY: dev
dev: ## Run application in development mode
	go run main.go

.PHONY: test
test: ## Run tests
	go test -v ./...

.PHONY: build
build: ## Build the application
	go build -o main .

.PHONY: clean
clean: ## Clean build artifacts
	rm -f main
	go clean

# Docker
.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t $(IMAGE_NAME):$(TAG) .

.PHONY: docker-run
docker-run: ## Run Docker container
	docker run -p 8080:8080 --env-file .env $(IMAGE_NAME):$(TAG)

.PHONY: docker-push
docker-push: ## Push Docker image to registry
	docker tag $(IMAGE_NAME):$(TAG) $(DOCKER_USERNAME)/$(IMAGE_NAME):$(TAG)
	docker push $(DOCKER_USERNAME)/$(IMAGE_NAME):$(TAG)

# Docker Compose
.PHONY: up
up: ## Start services with Docker Compose
	docker compose up -d

.PHONY: down
down: ## Stop services with Docker Compose
	docker compose down

.PHONY: logs
logs: ## Show logs from services
	docker compose logs -f

.PHONY: restart
restart: ## Restart services
	docker compose restart

.PHONY: rebuild
rebuild: ## Rebuild and restart services
	docker compose down
	docker compose build --no-cache
	docker compose up -d

# Database
.PHONY: mongo-shell
mongo-shell: ## Access MongoDB shell
	docker compose exec mongodb mongosh -u admin -p password

.PHONY: mongo-logs
mongo-logs: ## Show MongoDB logs
	docker compose logs -f mongodb

# Cleanup
.PHONY: cleanup
cleanup: ## Remove containers and volumes
	docker compose down -v
	docker system prune -f

.PHONY: reset
reset: ## Reset everything (careful!)
	docker compose down -v --remove-orphans
	docker volume prune -f
	docker image prune -a -f

# Status
.PHONY: status
status: ## Show service status
	docker compose ps

.PHONY: health
health: ## Check service health
	@echo "Checking API health..."
	@curl -f http://localhost:8080/graphql || echo "API not responding"
	@echo "\nChecking MongoDB health..."
	@docker compose exec mongodb mongosh --eval "db.adminCommand('ping')" --quiet || echo "MongoDB not responding"

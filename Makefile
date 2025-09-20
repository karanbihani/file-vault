# This Makefile provides a set of useful commands to manage the application stack and tasks.
.PHONY: up down build logs tidy sqlc migrate-up migrate-down migrate-up-one migrate-down-one
include .env

# ==============================================================================
# Configuration Variables
# ==============================================================================
# The network that docker-compose creates. Default is <project-folder-name>_default.
NETWORK_NAME   := file-vault_default
# The database URL for the migration tool. CRITICAL: uses 'postgres' as the host, not 'localhost'.
DATABASE_URL   := "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable"
MIGRATIONS_DIR := ./sql/migrations
MIGRATE_IMAGE  := migrate/migrate

# ==============================================================================
# Docker Compose Commands
# ==============================================================================

# Starts all application services (postgres, backend) in the background.
up:
	@echo "Starting up services..."
	docker-compose up -d

# Stops and removes all services and the database volume.
down:
	@echo "Stopping and tearing down services..."
	docker-compose down -v

# Forces a rebuild of the application images before starting.
build:
	@echo "Building and starting services..."
	docker-compose up --build -d

# Follows the logs of all running services.
logs:
	@echo "Tailing logs..."
	docker-compose logs -f

# ==============================================================================
# Go & Database Commands
# ==============================================================================

# Tidies up Go module dependencies.
tidy:
	@echo "Tidying go.mod and go.sum..."
	go mod tidy

# Generates Go code from SQL queries.
sqlc:
	@echo "Generating Go code from SQL queries..."
	sqlc generate

# --- Migration Commands ---
# These commands run the migration tool in a temporary container on the correct Docker network.

# Applies ALL 'up' migrations.
migrate-up:
	@echo "Applying all UP migrations..."
	docker run --rm -v $(PWD)/$(MIGRATIONS_DIR):/migrations --network $(NETWORK_NAME) $(MIGRATE_IMAGE) -path=/migrations -database $(DATABASE_URL) up

# Rolls back ALL migrations.
migrate-down:
	@echo "Rolling back all migrations..."
	@docker run --rm -v $(PWD)/$(MIGRATIONS_DIR):/migrations --network $(NETWORK_NAME) $(MIGRATE_IMAGE) -path=/migrations -database $(DATABASE_URL) down -all

# Applies the NEXT single 'up' migration.
migrate-up-one:
	@echo "Applying next UP migration..."
	@docker run --rm -v $(PWD)/$(MIGRATIONS_DIR):/migrations --network $(NETWORK_NAME) $(MIGRATE_IMAGE) -path=/migrations -database $(DATABASE_URL) up 1

# Rolls back the LAST applied migration.
migrate-down-one:
	@echo "Rolling back last migration..."
	@docker run --rm -v $(PWD)/$(MIGRATIONS_DIR):/migrations --network $(NETWORK_NAME) $(MIGRATE_IMAGE) -path=/migrations -database $(DATABASE_URL) down 1

# Seeds the database with initial roles, permissions, and their mappings.
.PHONY: seed
seed:
	@echo "Seeding the database..."
	@docker exec -i file_vault_db psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} < sql/seeds/0001_roles_and_permissions.sql
	@echo "Seeding complete."
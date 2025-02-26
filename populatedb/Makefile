# Configuration
DB_CONTAINER_NAME=go_postgres
DB_NAME=shows
DB_USER=postgres
DB_PASSWORD=mysecretpassword
DB_PORT=5432
DB_DUMP_FILE=data.sql

# Start PostgreSQL in Docker
start-db:
	docker run --name $(DB_CONTAINER_NAME) -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -e POSTGRES_DB=$(DB_NAME) -p $(DB_PORT):5432 -d postgres:latest
	@echo "Waiting for Postgres to start..."
	@sleep 5 # Ensure the database is ready before running the Go app

# Run the Go application (assuming `main.go` populates the DB)
run-app:
	go run main.go

# Dump the database using pg_dump in the container
dump-db:
	docker exec $(DB_CONTAINER_NAME) pg_dump -U $(DB_USER) -d $(DB_NAME) -f /tmp/$(DB_DUMP_FILE)

# Copy dump file from container to host
copy-dump:
	docker cp $(DB_CONTAINER_NAME):/tmp/$(DB_DUMP_FILE) ../sql/$(DB_DUMP_FILE)

# Stop and remove the Postgres container
stop-db:
	docker stop $(DB_CONTAINER_NAME)
	docker rm $(DB_CONTAINER_NAME)

# Full pipeline: Start DB -> Run App -> Dump DB -> Copy Dump -> Stop DB
backup-db: start-db run-app dump-db copy-dump stop-db
	@echo "Database backup completed: $(DB_DUMP_FILE)"

# Clean up everything
clean:
	rm -f $(DB_DUMP_FILE)

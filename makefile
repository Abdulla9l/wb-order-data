.PHONY: run producer build test clean kafka-setup help

run:
	@echo "Starting main application..."
	@echo "Make sure PostgreSQL and Kafka are running!"
	go run ./cmd/app

producer:
	@echo "Starting Kafka producer..."
	go run ./cmd/producer

build:
	@echo "Building binaries..."
	go build -o bin/app ./cmd/app
	go build -o bin/producer ./cmd/producer

test:
	@echo "Running tests..."
	go test ./...

clean:
	@echo "Cleaning..."
	rm -rf bin/
	go clean

kafka-setup:
	@echo "Kafka setup instructions:"
	@echo "1. Download Kafka from https://kafka.apache.org/downloads"
	@echo "2. Extract and run:"
	@echo "   bin/zookeeper-server-start.sh config/zookeeper.properties"
	@echo "   bin/kafka-server-start.sh config/server.properties"
	@echo "3. Create topic:"
	@echo "   bin/kafka-topics.sh --create --topic orders --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1"

help:
	@echo "WB Order Data Service"
	@echo ""
	@echo "Commands:"
	@echo "  make run       - Start main application (requires PostgreSQL & Kafka)"
	@echo "  make producer  - Start Kafka producer"
	@echo "  make build     - Build binaries"
	@echo "  make test      - Run tests"
	@echo "  make clean     - Clean build artifacts"
	@echo "  make kafka-setup - Show Kafka setup instructions"
	@echo ""
	@echo "Prerequisites:"
	@echo "  - PostgreSQL running on localhost:5432"
	@echo "  - Kafka running on localhost:9092"
	@echo "  - .env file with database credentials"
migrate:
	@echo "Running migrations with goose..."
	goose -dir migrations postgres "user=app_user password=12345 dbname=orders_db sslmode=disable" up

migrate-status:
	@echo "Migration status:"
	goose -dir migrations postgres "user=app_user password=12345 dbname=orders_db sslmode=disable" status

migrate-down:
	@echo "Rolling back one migration..."
	goose -dir migrations postgres "user=app_user password=12345 dbname=orders_db sslmode=disable" down

migrate-reset:
	@echo "Resetting all migrations..."
	goose -dir migrations postgres "user=app_user password=12345 dbname=orders_db sslmode=disable" reset

migrate-create:
	@echo "Creating new migration file..."
	goose -dir migrations create $(name) sql
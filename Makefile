.PHONY: run-account run-transaction run-ledger \
        account transaction transaction-consumer \
        ledger ledger-consumer build-base build-up up down rebuild clean

account:
	go run ./cmd/account-service/main.go

run-account:
	$(MAKE) -f account.Makefile account

transaction:
	go run ./cmd/transaction-service/main.go

transaction-consumer:
	go run ./cmd/transaction-consumer/main.go

run-transaction:
	@echo "Starting transaction service..."
	$(MAKE) -f transaction.Makefile transaction &
	@echo "Starting transaction consumer..."
	$(MAKE) -f transaction.Makefile transaction-consumer &
	@echo "Waiting for processes..."
	wait
	@echo "All processes completed"

ledger:
	go run ./cmd/ledger-service/main.go

ledger-consumer:
	go run ./cmd/ledger-consumer/main.go

run-ledger:
	@echo "Starting ledger service..."
	$(MAKE) -f ledger.Makefile ledger &
	@echo "Starting ledger consumer..."
	$(MAKE) -f ledger.Makefile ledger-consumer &
	@echo "Waiting for processes..."
	wait
	@echo "All processes completed"

# Build the base image used in all service Dockerfiles
build-base:
	docker build -f ./deployments/base.Dockerfile -t txsystem-base .

# Start services
build-up:
	$(MAKE) build-base
	docker compose build
	docker compose up -d

up:
	docker compose up -d

# Stop services
down:
	docker compose down

# Force rebuild everything (cleans Docker cache)
rebuild:
	docker build --no-cache -f base.Dockerfile -t txsystem-base .
	docker compose build --no-cache

clean: down
	docker volume rm txsystem_postgres_data || true
	docker volume rm txsystem_mongo_data || true
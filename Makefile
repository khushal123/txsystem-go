.PHONY: run-account run-transaction run-ledger \
        account transaction transaction-consumer \
        ledger ledger-consumer

# --- Account ---
account:
	go run ./cmd/account-service/main.go

run-account:
	$(MAKE) account

# --- Transaction ---
transaction:
	go run ./cmd/transaction-service/main.go

transaction-consumer:
	go run ./cmd/transaction-consumer/main.go

run-transaction:
	$(MAKE) transaction &
	$(MAKE) transaction-consumer &
	wait

# --- Ledger ---
ledger:
	go run ./cmd/ledger-service/main.go

ledger-consumer:
	go run ./cmd/ledger-consumer/main.go

run-ledger:
	$(MAKE) ledger &
	$(MAKE) ledger-consumer &
	wait

# Build the base image used in all service Dockerfiles
build-base:
	docker build -f base.Dockerfile -t txsystem-base .

# Start services
up:
	docker build -f base.Dockerfile -t txsystem-base .
	docker compose build
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



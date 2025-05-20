.PHONY: all account ledger transaction

create-topics:
	@echo "Creating kafka topics..."
	docker exec --workdir /scripts/ kafka bash ./create-topic.sh 
	

# Run account service with reflex watch
account:
	reflex -r '\.go$$' -s -- sh -c "go run ./cmd/account-service/main.go"

# Run ledger service with reflex watch
ledger:
	reflex -r '\.go$$' -s -- sh -c "go run ./cmd/ledger-service/main.go"


transaction-api-doc:
	@echo "Generating transaction API documentation..."
	go generate ./internal/transaction

# Run transaction service with reflex watch
transaction:
	@echo "Generating transaction API documentation..."
	$(MAKE) transaction-api-doc &
	reflex -r '\.go$$' -s -- sh -c "go run ./cmd/transaction-service/main.go"

consumer:
	@echo "Running kafka consumer..."
	sh -c "go run ./cmd/transaction-consumer/main.go"

# Run all services concurrently
all:
	# Use GNU parallel or run in background (POSIX sh)
	$(MAKE) account &
	$(MAKE) ledger &
	$(MAKE) transaction &
	wait
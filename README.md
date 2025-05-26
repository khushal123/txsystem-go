# Project Overview

This project is a microservices-based transaction system designed to handle financial transactions efficiently and reliably. It leverages a combination of modern tools and technologies to ensure scalability, resilience, and maintainability.

# Tools and Technologies

*   **Go:** The primary programming language for building the microservices.
*   **Docker:** Used for containerizing the applications and their dependencies.
*   **Docker Compose:** Facilitates the management of multi-container Docker applications.
*   **Kafka:** A distributed streaming platform used for asynchronous communication between services.
*   **PostgreSQL:** A powerful open-source relational database.
*   **MongoDB:** A NoSQL document database used for flexible data storage.
*   **Echo web framework:** A high-performance, extensible, and minimalist Go web framework.
*   **GORM:** A developer-friendly ORM library for Go.
*   **Swagger:** Used for designing, building, and documenting RESTful APIs.

# Prerequisites

Before you begin, ensure you have the following installed:

*   **Docker:** [Installation Guide](https://docs.docker.com/get-docker/)
*   **Docker Compose:** [Installation Guide](https://docs.docker.com/compose/install/)
*   **Go (for local development):** [Installation Guide](https://golang.org/doc/install)

Additionally, you need to create a `.env` file by copying the example file:

```bash
cp .env.example .env
```

Make sure to update the `.env` file with your specific configuration values.

# Getting Started / How to Run

1.  **Start Services:**
    Use Docker Compose to build and start all services in detached mode:

    ```bash
    docker-compose up -d
    ```

2.  **Uncomment Services (Optional):**
    By default, some services might be commented out in the `docker-compose.yml` file. If you need to run specific services, uncomment them before running `docker-compose up -d`.

3.  **Kafka Topic Creation:**
    Kafka topics can be created in two ways:
    *   **Automatic:** The Kafka container is configured to automatically create topics based on the services that connect to it.
    *   **Manual:** You can use the Makefile target to create topics:
        ```bash
        make create-topics
        ```

# Development

The `Makefile` provides several targets to help with development:

*   `make account`: Builds and runs the account service.
*   `make ledger`: Builds and runs the ledger service.
*   `make transaction`: Builds and runs the transaction service.
*   `make consumer`: Builds and runs the Kafka consumer service.
*   `make all`: Builds and runs all services.


Once the services are running, you can typically access the Swagger UI through one of the services (e.g., the transaction service or an API gateway if implemented) at a path like `/swagger/index.html`. Refer to the specific service's documentation or configuration for the exact URL.

# Project Structure

The project follows a standard Go project layout:

*   `cmd/`: Contains the main applications (entry points) for each microservice.
    *   `cmd/account_service/`
    *   `cmd/transaction_service/`
    *   `cmd/ledger_service/`
    *   `cmd/consumer_service/`
*   `internal/`: Contains the internal logic specific to each service. This code is not meant to be imported by other projects.
    *   `internal/account/`
    *   `internal/transaction/`
    *   `internal/ledger/`
    *   `internal/consumer/`
*   `pkg/`: Contains shared libraries and utilities that can be used across multiple services or projects (e.g., Kafka utilities, database helpers).
*   `scripts/`: Contains utility scripts for tasks like database migrations, Kafka topic creation, etc.
*   `gateway/`: May contain configuration files for an API gateway (e.g., Nginx, Traefik).
*   `docs/`: Contains API documentation files, typically Swagger JSON or YAML.
*   `api/`: Contains protobuf definitions for gRPC services.
```

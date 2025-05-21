# filepath: cmd/transaction-service/Dockerfile 
FROM txsystem-base as builder

# Build transaction-service
WORKDIR /app/cmd/transaction-service
RUN go build -o /transaction-service .

# Build transaction-consumer
WORKDIR /app/cmd/transaction-consumer
RUN go build -o /transaction-consumer .

# --- Final Stage ---
FROM alpine:3.19
WORKDIR /root/

# COPY --from=builder /app/.env .env 
COPY --from=builder /transaction-service .
COPY --from=builder /transaction-consumer .
COPY ./scripts/start.sh . # Assuming start.sh is in scripts/ and configured for these binaries

RUN chmod +x /start.sh

CMD ["/start.sh"]
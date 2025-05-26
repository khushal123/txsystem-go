FROM txsystem-base AS builder

WORKDIR /app/cmd/ledger-service
RUN go build -o /ledger-service .

WORKDIR /app/cmd/ledger-consumer
RUN go build -o /ledger-consumer .

FROM alpine:3.21
RUN apk add --no-cache bash 

WORKDIR /app
COPY --from=builder /ledger-service .
COPY --from=builder /ledger-consumer .
COPY --from=builder /app/.env .env

ENTRYPOINT ["./ledger-service"]
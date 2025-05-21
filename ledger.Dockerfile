FROM txsystem-base AS builder

WORKDIR /app/cmd/ledger-service
RUN go build -o /ledger-service .

WORKDIR /app/cmd/ledger-consumer
RUN go build -o /ledger-consumer .

FROM alpine:3.21
RUN apk add --no-cache make

WORKDIR /app
COPY --from=builder /ledger-service .
COPY --from=builder /ledger-consumer .
COPY --from=builder /app/.env .env

RUN echo -e 'run-ledger:\n\t./ledger-service &\n\t./ledger-consumer &\n\twait' > Makefile

CMD ["make", "run-ledger"]

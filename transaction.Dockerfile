FROM txsystem-base AS builder

WORKDIR /app/cmd/transaction-service
RUN go build -o /transaction-service .

WORKDIR /app/cmd/transaction-consumer
RUN go build -o /transaction-consumer .

FROM alpine:3.21
RUN apk add --no-cache make

WORKDIR /app
COPY --from=builder /transaction-service .
COPY --from=builder /transaction-consumer .
COPY --from=builder /app/.env .env

RUN echo -e 'run-transaction:\n\t./transaction-service &\n\t./transaction-consumer &\n\twait' > Makefile

CMD ["make", "run-transaction"]

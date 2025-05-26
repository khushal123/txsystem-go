FROM txsystem-base AS builder

WORKDIR /app/cmd/account-service
RUN go build -o /account-service .

FROM alpine:3.21
RUN apk add --no-cache bash

WORKDIR /app
COPY --from=builder /account-service .
COPY --from=builder /app/.env .env

ENTRYPOINT ["./account-service"]
FROM txsystem-base as builder

WORKDIR /app/cmd/account-service
RUN go build -o /account-service .

FROM alpine:3.19
WORKDIR /root/

COPY --from=builder /app/.env .env
COPY --from=builder /account-service .

CMD ["./account-service"]

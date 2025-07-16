FROM golang:1.23.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o payment-system ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/payment-system .

COPY --from=builder /app/migrations ./migrations

CMD ["./payment-system"]

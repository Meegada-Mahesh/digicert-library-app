FROM golang:1.24.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

FROM debian:stable-slim

WORKDIR /app
COPY --from=builder /app/main .
COPY .env .env

EXPOSE 8080

CMD ["./main"]
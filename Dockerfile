# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .

RUN go build -o main main.go

# Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY .env .
COPY db/migration ./db/migration

EXPOSE 8080

# Chạy ứng dụng
CMD ["/app/main"]
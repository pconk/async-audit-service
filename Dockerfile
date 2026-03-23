# Stage 1: Build
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Copy dependency files first (untuk caching layer docker)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN go build -o audit-service ./cmd/api/main.go

# Stage 2: Run
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/audit-service .

EXPOSE 50051

CMD ["./audit-service"]
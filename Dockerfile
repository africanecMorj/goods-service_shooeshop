# ---------- Build stage ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/server

# ---------- Runtime stage ----------
FROM alpine:3.19

# Install PostgreSQL in FINAL image (important)
RUN apk add --no-cache \
    postgresql \
    postgresql-client \
    su-exec \
    bash

# Setup Postgres directory
RUN mkdir -p /var/lib/postgresql/data && \
    chown -R postgres:postgres /var/lib/postgresql

# Copy Go binary
WORKDIR /app
COPY --from=builder /app/app .

# Copy startup script
COPY start.sh .
RUN chmod +x start.sh

EXPOSE 5432 8080

CMD ["./start.sh"]
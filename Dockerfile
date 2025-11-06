# Stage 1: Build stage
FROM golang:1.25.4 AS builder

WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Install Air for live reloading
RUN go install github.com/cosmtrek/air@latest

# Build the application (optional but good practice)
RUN go build -o /app/main cmd/server/main.go

# Stage 2: Final stage
FROM alpine:latest

WORKDIR /app

# Copy Air binary, application binary, and configs from the builder stage
COPY --from=builder /go/bin/air /usr/local/bin/
COPY --from=builder /app/main .
COPY .air.toml .
COPY .env .
COPY internal/migrations ./internal/migrations

# Expose port
EXPOSE 8080

# Command to run the application using Air for live reload
CMD ["air"]
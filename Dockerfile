# Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies if needed (e.g. for CGO)
# RUN apk add --no-cache build-base

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application
# CGO_ENABLED=0 for static binary, needed for distroless or scratch
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/api/main.go

# Run Stage
FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/.env.example .env

# Expose the application port
EXPOSE 8000

# Command to run the executable
CMD ["./main"]

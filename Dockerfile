# Build stage
FROM golang:1.24.1-alpine AS builder

# Install git and ca-certificates (needed for go mod download)
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Generate Swagger docs (docs/ is in .dockerignore, so generate inside image)
RUN go install github.com/swaggo/swag/cmd/swag@latest && swag init -g cmd/api/main.go

# Build the application with optimization flags
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o main cmd/api/main.go

# Final stage - distroless (secure and minimal)
FROM gcr.io/distroless/static-debian11:nonroot

# Copy binary
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"] 
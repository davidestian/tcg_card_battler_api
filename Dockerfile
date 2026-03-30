# Stage 1: Build
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

# Copy dependency files first
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build from the subdirectory where main.go lives
# This creates a binary named 'server'
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server/main.go

# Stage 2: Final Image
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /server .

# ADD THIS: Copy the config folder into the image
# This creates /root/config/config.yaml
COPY config/ ./config/

# Expose your API port (change 8080 if your code uses a different one)
EXPOSE 8080

# Run the binary
CMD ["./server"]

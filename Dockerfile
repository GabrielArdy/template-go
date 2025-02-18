# Build stage
FROM golang:1.23.6-alpine AS builder

# Add necessary build tools
RUN apk add --no-cache make gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the project
COPY . .

# Build the application
RUN make build

# Final stage
FROM alpine:latest

# Add CA certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Set timezone
ENV TZ=Asia/Jakarta

# Create non-root user
RUN adduser -D appuser

# Create directories and set permissions
RUN mkdir -p /app/build/config && \
    mkdir -p /app/build/resources && \
    chown -R appuser:appuser /app

# Copy binary and configs to build directory
COPY --from=builder /app/build/main /app/build/main
COPY --from=builder /app/config/config-*.yml /app/build/

# Set working directory to where binary is
WORKDIR /app/build

# Switch to non-root user
USER appuser

# Expose the application port
EXPOSE 8080

# Set environment variable for profile
ENV ACTIVE_PROFILE=local

# Run the application from build directory
CMD ["./main", "--active.profile=local"]
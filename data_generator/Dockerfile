# syntax=docker/dockerfile:1

# --------------------
# Build Stage
# --------------------
  FROM golang:1.23-alpine AS builder

  # Set the working directory inside the container
  WORKDIR /app
  
  # Copy go.mod and go.sum files to leverage Docker cache for dependencies
  COPY go.mod go.sum ./
  
  # Download dependencies
  RUN go mod download
  
  # Copy the rest of the application source code
  COPY . .
  
  # Build the Go binary
  RUN CGO_ENABLED=0 GOOS=linux go build -o data_generator .
  
  # --------------------
  # Production Stage
  # --------------------
  FROM alpine:3.14
  
  # Set the working directory inside the container
  WORKDIR /home/datagenuser
  
  # Install curl for potential healthchecks or debugging
  RUN apk add --no-cache curl
  
  # Create a non-root user for security
  RUN adduser -D datagenuser
  
  # Switch to the non-root user
  USER datagenuser
  
  # Copy the compiled binary from the builder stage
  COPY --from=builder /app/data_generator .
  
  # Define the command to run the data generator (can be overridden in docker-compose)
  CMD ["./data_generator", "-n", "1000", "-url", "http://app:80"]
  
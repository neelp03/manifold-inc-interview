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
  RUN CGO_ENABLED=0 GOOS=linux go build -o main .
  
  # --------------------
  # Production Stage
  # --------------------
  FROM alpine:3.14
  
  # Set the working directory inside the container
  WORKDIR /home/appuser
  
  # Install curl for healthchecks
  RUN apk add --no-cache curl
  
  # Create a non-root user for security
  RUN adduser -D appuser
  
  # Switch to the non-root user
  USER appuser
  
  # Copy the compiled binary from the builder stage
  COPY --from=builder /app/main .
  
  # Expose the application port
  EXPOSE 80
  
  # Define the command to run the application
  CMD ["./main"]
  
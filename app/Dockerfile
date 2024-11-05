# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Runtime stage
FROM alpine:latest
WORKDIR /root/
# Install curl
RUN apk add --no-cache curl
COPY --from=builder /app/main .
EXPOSE 80
CMD ["./main"]

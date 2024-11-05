# syntax=docker/dockerfile:1

FROM golang:1.23-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o data_generator data_generator.go
CMD ["./data_generator", "-n", "1000", "-url", "http://app:80"]

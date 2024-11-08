# Go HTTP Server with InfluxDB Integration

## Overview

This project consists of a simple Go HTTP server that accepts JSON POST requests and inserts the data into an InfluxDB instance. It includes a data generator that simulates log entries for testing purposes. The project is containerized using Docker and supports deployment using Docker Compose.

## Features

- Accepts JSON POST requests at the root endpoint (`/`).
- Inserts received data into InfluxDB with appropriate tags and fields.
- Implements a health check endpoint at `/health`.
- Includes retry logic for writing to InfluxDB.
- Graceful shutdown handling.
- Containerized with Docker, using multi-stage builds for efficient images.
- Configured for zero-downtime deployments and easy rollbacks.
- CI/CD pipeline setup using GitHub Actions.

## Project Structure

- `app/`: Contains the Go HTTP server code and Dockerfile.
- `data_generator/`: Contains the code for generating test data.
- `docker-compose.yml`: Defines the services and configurations for Docker Compose.
- `README.md`: Project documentation.

## Getting Started

### Prerequisites

- Docker and Docker Compose installed on your machine.
- An InfluxDB instance (can be run as a Docker container).

### Environment Variables

The application requires the following environment variables:

- `INFLUXDB_URL`: URL of the InfluxDB instance (e.g., `http://influxdb:8086`).
- `INFLUXDB_TOKEN`: Authentication token for InfluxDB.
- `INFLUXDB_ORG`: InfluxDB organization name.
- `INFLUXDB_BUCKET`: InfluxDB bucket name.

### Running with Docker Compose

1. **Clone the Repository**

   ```bash
   git clone https://github.com/neelp03/manifold-inc-interview.git
   cd manifold-inc-interview
   ```

2. **Configure Environment Variables**

   Update the `docker-compose.yml` file with the appropriate environment variables for setup.

3. **Start Services**

   ```bash
   docker-compose up -d
   ```

4. **Verify Services**

   ```bash
   docker-compose ps
   ```

   Ensure all services are running and healthy.

5. **Test the Application**

   Send a test POST request:

   ```bash
   curl -X POST http://localhost:80 \
     -H 'Content-Type: application/json' \
     -d '{
       "service": "user-service",
       "endpoint": "/api/users",
       "error": "Database connection timeout",
       "traceback": "File \"app.py\", line 42, in get_user\n    raise ConnectionError(\"Database timeout\")"
   }'
   ```

6. **Check InfluxDB**

   Access the InfluxDB UI or use the CLI to verify that the data has been inserted.

### Building and Running Manually

1. **Build the Docker Image**

   ```bash
   cd app
   docker build -t neelp03/app:latest .
   ```

2. **Run the Docker Container**

   ```bash
   docker run -p 80:80 --env-file .env neelp03/app:latest
   ```

   Create an `.env` file with the required environment variables.

### Health Check Endpoint

- **URL:** `http://localhost:80/health`
- **Method:** `GET`
- **Response:** `OK`

### Logging

The application logs important events and errors to stdout, which can be viewed using Docker logs:

```bash
docker-compose logs -f app
```

### Graceful Shutdown

The application handles interrupt signals to gracefully shut down the server, ensuring all connections are properly closed.

## CI/CD Pipeline

The project includes a CI/CD pipeline using GitHub Actions, which automates the build, test, and deployment processes. On every push to the `main` branch:

- Docker images are built and pushed to Docker Hub.
- The application is deployed to the remote server via SSH.
- Zero-downtime deployments are facilitated through careful orchestration.

### Secrets Management

Sensitive information such as Docker Hub credentials and SSH access details are managed using GitHub Secrets.

## Security Considerations

- The application runs inside a Docker container as a non-root user.
- Environment variables are used to pass sensitive configuration data.
- Input data is validated and limited in size to prevent abuse.
- Logs should be monitored and rotated to prevent disk space exhaustion.
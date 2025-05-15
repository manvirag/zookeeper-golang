# Service Discovery with Consul and Nginx

This project demonstrates a simple service discovery setup using Consul and Nginx. It includes:
- Multiple instances of a Go service
- Consul for service discovery and health checks
- Nginx for load balancing

## Architecture

```
┌─────────┐     ┌─────────┐     ┌─────────┐
│ Service │     │ Service │     │ Service │
│  :8080  │     │  :8081  │     │  :8082  │
└────┬────┘     └────┬────┘     └────┬────┘
     │               │               │
     └───────────────┼───────────────┘
                     │
              ┌──────┴──────┐
              │    Consul   │
              │    :8500    │
              └──────┬──────┘
                     │
              ┌──────┴──────┐
              │    Nginx    │
              │    :80      │
              └──────┬──────┘
                     │
                     ▼
                Client Requests
```

## Features

- **Service Discovery**: Services automatically register with Consul
- **Health Checks**: Consul monitors service health
- **Load Balancing**: Nginx distributes traffic to healthy services
- **DNS-based Discovery**: Nginx uses Consul's DNS interface

## Prerequisites

- Docker
- Docker Compose
- Go 1.16 or later

## Running the Project

1. Start the infrastructure:
   ```bash
   docker-compose up -d
   ```

2. Access the services:
   - Main service: http://localhost
   - Consul UI: http://localhost:8500

## How It Works

1. **Service Registration**:
   - Each service instance registers with Consul
   - Services include health check endpoints
   - Consul monitors service health

2. **Load Balancing**:
   - Nginx uses Consul's DNS interface
   - Automatically routes to healthy services
   - No manual configuration needed

3. **Health Checking**:
   - Consul performs health checks
   - Unhealthy services are automatically removed
   - Services are added back when healthy

## Configuration

- **Consul**: Running on port 8500 (UI) and 8600 (DNS)
- **Nginx**: Running on port 80
- **Services**: Running on ports 8080, 8081, and 8082

## Development

1. Build the service:
   ```bash
   go build -o service ./cmd/service
   ```

2. Run a service instance:
   ```bash
   ./service --port 8080 --service-name my-service
   ```

## Cleanup

To stop and remove all containers:
```bash
docker-compose down
```

## License

MIT

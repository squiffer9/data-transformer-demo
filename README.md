# Overview
The service accepts JSON payloads containing a country code and question/answer pairs, then uses in-memory data (loaded from MySQL) to generate additional related question/answer pairs based on predefined mapping rules.

## Project Structure
```
data-transformer-demo/
├── cmd/server/             # Application entrypoint
├── internal/
│   ├── db/                 # Database operations
│   ├── cache/              # In-memory cache
│   └── service/            # Business logic
├── docker/                 # Docker configurations
├── test/                   # Test files
└── docker-compose.yaml     # Local development setup
```
## Technology Stack
- Go 1.23
- MySQL 8.0
- Docker & Docker Compose
- HAProxy
- k6 (for load testing)

### Setup & Development
Prerequisites
- Docker and Docker Compose
- Go 1.23 or later
- MySQL 8.0
Running Locall
Clone the repository
```
git clone https://github.com/squiffer9/data-transformer-demo
cd data-transformer-demo
```
Start all service
```
docker-compose up -d
```
The application will be available at http://localhost:80 (via HAProxy) or http://localhost:8080 (direct access).

For development purposes, you can also run the application directly:
```
# Stop the app container but keep MySQL running
docker-compose stop app

# Run the application locally

go run cmd/server/main.go
```
Current Status
This is a work in progress, focusing on:

Implementing core business logic \
Exploring performance optimization techniques \
Testing different caching strategies \
Experimenting with load balancing configurations \
Performance testing and optimization are ongoing, with current implementation under review for reliability and production readiness.

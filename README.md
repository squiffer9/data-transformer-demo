# Data Transformer Demo

This is a portfolio project implementing a microservice that processes question/answer pairs based on predefined mappings.\
The project aims to explore high-performance data transformation patterns in Go and is sample to apply for the project at the URL below.\
https://www.upwork.com/jobs/Microservice-for-Data-Transformation_~021888370257595849336/?referrer_url_path=%2Fnx%2Fsearch%2Fjobs%2Fsaved%2F

## Overview

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

## Setup & Development

### Prerequisites
- Docker and Docker Compose
- Go 1.23 or later
- MySQL 8.0

### Running Locally

1. Clone the repository
```bash
git clone https://github.com/squiffer9/data-transformer-demo
cd data-transformer-demo
```

2. Start all services
```bash
docker-compose up -d
```

The application will be available at http://localhost:80 (via HAProxy) or http://localhost:8080 (direct access).

For development purposes, you can also run the application directly:
```bash
# Stop the app container but keep MySQL running
docker-compose stop app

# Run the application locally
go run cmd/server/main.go
```

## Current Status

This is a work in progress, focusing on:
- Implementing core business logic
- Exploring performance optimization techniques
- Testing different caching strategies
- Experimenting with load balancing configurations

Performance testing and optimization are ongoing, with current implementation under review for reliability and production readiness.

## Note

This is a portfolio project created for learning and demonstration purposes.\
While it implements the core functionality, some aspects (particularly around performance testing and production deployment) are still experimental and under development.

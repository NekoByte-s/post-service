# Post Service

A RESTful API service for managing posts, built with Go and following standard project layout.

## Features

- ✅ CRUD operations for posts
- ✅ Swagger/OpenAPI documentation
- ✅ Docker containerization
- ✅ Kubernetes deployment (dev/prod)
- ✅ In-memory storage
- ✅ Health checks and probes

## Project Structure

```
.
├── cmd/modular/           # Main application entry point
├── internal/              # Private application code
│   ├── models/           # Data models
│   ├── repository/       # Data access layer
│   ├── service/          # Business logic
│   └── handlers/         # HTTP handlers
├── pkg/                   # Public library code
├── api/                   # OpenAPI specifications
├── configs/               # Configuration files
├── scripts/               # Deployment scripts
├── deployments/k8s/       # Kubernetes manifests
│   ├── dev/              # Development environment
│   └── prod/             # Production environment
└── docs/                  # Generated Swagger docs
```

## API Endpoints

- `GET /api/v1/posts` - Get all posts
- `POST /api/v1/posts` - Create a new post
- `GET /api/v1/posts/{id}` - Get a post by ID
- `PUT /api/v1/posts/{id}` - Update a post
- `DELETE /api/v1/posts/{id}` - Delete a post
- `GET /swagger/*` - Swagger documentation

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker
- kubectl (for Kubernetes deployment)

### Local Development

```bash
# Install dependencies
make deps

# Generate Swagger docs
make swagger

# Run the service
make run
```

The service will be available at:
- API: http://localhost:8080/api/v1/posts
- Swagger UI: http://localhost:8080/swagger/index.html

### Docker

```bash
# Build Docker image
make docker-build

# Run with Docker
make docker-run
```

### Kubernetes Deployment

#### Development Environment

```bash
# Deploy to dev
./scripts/deploy-dev.sh

# Cleanup dev
./scripts/cleanup-dev.sh
```

#### Production Environment

```bash
# Deploy to prod
./scripts/deploy-prod.sh

# Cleanup prod
./scripts/cleanup-prod.sh
```

## Configuration

The service uses environment variables for configuration:

- `PORT` - Server port (default: 8080)
- `GIN_MODE` - Gin mode (debug/release)
- `LOG_LEVEL` - Logging level

## Testing

```bash
make test
```

## Example Usage

### Create a Post

```bash
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My First Post",
    "content": "This is the content of my first post",
    "author": "John Doe"
  }'
```

### Get All Posts

```bash
curl http://localhost:8080/api/v1/posts
```

### Get a Specific Post

```bash
curl http://localhost:8080/api/v1/posts/{post-id}
```

## 🚀 Quick Start

### Deploy to Minikube
```bash
./deploy-dev.sh
```

### Test All Endpoints
```bash
./test-swagger.sh
```

### Access URLs
- **API**: http://localhost:8080/api/v1/posts
- **Swagger**: http://localhost:8080/swagger/index.html

## 🧹 Clean Project Structure

This project follows a clean, minimal structure with only essential files:
- **Core Application**: `cmd/`, `internal/`, `pkg/`
- **Kubernetes Deployment**: `deployments/k8s/dev/` and `deployments/k8s/prod/`
- **Documentation**: `docs/` (auto-generated), `api/swagger.yaml`
- **Scripts**: `deploy-dev.sh`, `test-swagger.sh`

## License

[Add your license here]
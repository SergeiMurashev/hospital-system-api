# Document Service

This is the document service for the hospital system API. It handles document management and search functionality using PostgreSQL and Elasticsearch.

## Features

- Document CRUD operations
- Document search using Elasticsearch
- Patient document history
- Authentication and authorization
- Integration with account service

## Prerequisites

- Go 1.21 or higher
- PostgreSQL
- Elasticsearch
- Account Service running

## Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
# Server
PORT=8004

# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=hospital_system

# Account Service
ACCOUNT_SERVICE_URL=http://account-service:8001

# Elasticsearch
ELASTICSEARCH_URL=http://elasticsearch:9200
```

## API Endpoints

### Documents

- `GET /api/v1/documents/patient/:patientID` - Get patient's documents
- `GET /api/v1/documents/:id` - Get document by ID
- `POST /api/v1/documents` - Create new document
- `PUT /api/v1/documents/:id` - Update document
- `DELETE /api/v1/documents/:id` - Delete document

### Search

- `GET /api/v1/search` - Search documents

## Running the Service

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Run the service:
   ```bash
   go run cmd/main.go
   ```

## Docker

Build and run using Docker:

```bash
docker build -t document-service .
docker run -p 8004:8004 document-service
```

## Development

The service follows a clean architecture pattern with the following structure:

- `cmd/` - Application entry point
- `internal/` - Private application code
  - `domain/` - Business entities and interfaces
  - `repository/` - Data access layer
  - `service/` - Business logic layer
  - `delivery/` - API handlers
- `pkg/` - Public libraries
  - `auth/` - Authentication client
  - `elasticsearch/` - Elasticsearch client 
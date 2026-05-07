# NAS API Documentation

## Overview

This documentation covers the Network Attached Storage (NAS) API endpoints, data models, and testing examples.

## Quick Links

- [API Reference](README.md) - Complete API documentation with all endpoints
- [OpenAPI Specification](openapi.yaml) - Machine-readable API specification (YAML)
- [Examples](examples.md) - Testing examples with curl, PowerShell, and Python
- [Data Models](models.md) - Data structures and configuration schemas

## Base URL

```
http://localhost:8088
```

## Main Features

### File Management (`/simple_upload`)

- **Upload files** - Single and multi-file upload with overwrite control
- **Download files** - Stream or force download as attachment
- **File metadata** - HEAD requests for file information
- **Delete files** - Remove files from storage
- **List directories** - Browse folder contents

### Database Examples (`/example`)

- **SQL examples** - CRUD operations for relational database
- **MongoDB examples** - CRUD operations for document database

## API Structure

All responses follow a standard format:

```json
{
  "TraceID": "unique-request-id",
  "Code": "success|error_code",
  "Message": "Human readable message",
  "Data": {}
}
```

## Getting Started

### 1. Start the Server

```bash
cd backend
go run main.go -config ./rootfs/etc/config.json
```

### 2. Test File Upload

```bash
# Create a test file
echo "Hello, NAS!" > test.txt

# Upload the file
curl -X POST "http://localhost:8088/simple_upload/object/test.txt" \
  -F "File=@test.txt"

# List files
curl "http://localhost:8088/simple_upload/objects/"
```

## Configuration

The server reads configuration from `config.json` or `config.yaml`:

```json
{
  "Http": {
    "ListenPort": 8088,
    "StaticRoot": "./static"
  },
  "Database": {
    "Rdb": {
      "DriverType": "sqlite",
      "DSN": "nas.db?mode=rwc",
      "AutoMigrateLevel": "auto"
    }
  },
  "Nas": {
    "SimpleUploadRoot": "/aaa",
    "RealFilenamePolicy": "uuid"
  }
}
```

## Security Features

- **Path traversal protection** - Blocks `../` and `./` patterns
- **Filename validation** - Normalizes paths and prevents directory escapes
- **CORS enabled** - Allows cross-origin requests from any domain

## Error Handling

Common error codes:

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `success` | 200 | Request successful |
| `invalid_arguments` | 400 | Bad request parameters |
| `not_found` | 404 | Resource not found |
| `file_exists` | 409 | File already exists |

## Additional Resources

- [Backend README](../backend/README.md) - Backend setup and architecture
- [Frontend README](../frontend/README.md) - Frontend development guide

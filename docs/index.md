# NAS API Documentation

## Overview

This documentation covers the Network Attached Storage (NAS) API endpoints, data models, configuration, and testing examples.

## Quick Links

- [API Reference](README.md) — Complete API documentation with all endpoints
- [OpenAPI Specification](openapi.yaml) — Machine-readable OpenAPI 3.0 spec
- [Examples](examples.md) — curl, Python, and WebDAV examples
- [Data Models](models.md) — Data structures and configuration schemas

## Base URL

```
http://localhost:8088
```

LAN devices can also use mDNS: `http://nas.local:8088`

## Main Features

### File Management (`/simple_upload`)
- Upload (single, multi, binary stream)
- Download with optional attachment mode
- File metadata via HEAD
- Delete files and directories (cascade)
- Update file tags (PATCH)
- List directories with tag info
- Create directories
- Search by keyword and tag
- Document preview (Office/PDF → HTML)

### Authentication & Users (`/api/users`)
- Session-based login/logout
- User CRUD (admin only)
- Permission management (CanRead / CanWrite / IsAdmin)
- Change password

### Admin Tools
- Filesystem scan to sync DB (`POST /api/scanfs`)

### WebDAV (`/webdav/`)
- Full WebDAV mount point
- Supports Windows Explorer, macOS Finder, davfs2
- Basic Auth with permission checks

### mDNS
- Automatic local network service discovery
- Broadcasts as `nas.local` via `_http._tcp`

### Frontend (SPA)
- File browser with grid and table views
- Image/video galleries with preview/player
- Tag management
- Responsive mobile support

## API Structure

All responses follow a standard format:

```json
{
  "TraceID": "unique-request-id",
  "Code": "success",
  "Message": "",
  "Data": {}
}
```

## Getting Started

### 1. Start the Server

```bash
cd backend
go run main.go -config ./rootfs/etc/config.json
```

### 2. Login & Test

```bash
# Login
curl -c cookies.txt -X POST http://localhost:8088/api/users/login \
  -H "Content-Type: application/json" \
  -d '{"Username":"admin","Password":"admin"}'

# List root directory
curl -b cookies.txt http://localhost:8088/simple_upload/objects/

# Upload a file
curl -b cookies.txt -X POST http://localhost:8088/simple_upload/object/test.txt \
  -F "File=@test.txt"
```

## Configuration

```json
{
  "Http": {
    "ListenPort": 8088,
    "StaticRoot": "./static",
    "Auth": { "Enabled": true }
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
  },
  "MDNS": {
    "Enabled": true,
    "ServiceName": "_http._tcp",
    "Hostname": "nas.local",
    "Info": "NAS Web Service"
  }
}
```

## Security Features

- Session-based authentication with bcrypt passwords
- Permission levels: CanRead, CanWrite, IsAdmin
- Path traversal protection (blocks `../` and `./`)
- Filename validation and normalization
- WebDAV filesystem-level permission enforcement
- Admin-only API endpoints protected by middleware

## Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `success` | 200 | Request successful |
| `invalid_arguments` | 400 | Bad parameters |
| `not_found` | 404 | Resource not found |
| `file_exists` | 409 | File already exists |
| `403` | 403 | Permission denied |
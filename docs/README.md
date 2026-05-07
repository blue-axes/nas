# NAS API Documentation

## Base URL

```
http://localhost:8088
```

## Response Format

All API responses follow a consistent structure:

```json
{
  "TraceID": "string",
  "Code": "success",
  "Message": "",
  "Data": {}
}
```

### Common Error Codes

| Code | Description |
|------|-------------|
| `success` | Request successful |
| `unknown` | Unknown error |
| `invalid_arguments` | Invalid request parameters |
| `not_found` | Resource not found |
| `file_exists` | File already exists |
| `file_checksum_invalid` | File checksum verification failed |

## File Upload API (`/simple_upload`)

### 1. Get File Metadata

**Endpoint:** `HEAD /simple_upload/object/{filename}`

**Description:** Retrieve file metadata without downloading the file content.

**Path Parameters:**
- `filename` - Path to the file (URL encoded)

**Response Headers:**
- `Content-Length` - File size in bytes
- `Content-Disposition` - Attachment filename
- `Content-Type` - `application/octet-stream`

**Example Request:**
```bash
curl -I "http://localhost:8088/simple_upload/object/myfile.txt"
```

**Success Response:** `200 OK` with headers only

---

### 2. Download File

**Endpoint:** `GET /simple_upload/object/{filename}`

**Description:** Download a file or display it inline.

**Path Parameters:**
- `filename` - Path to the file (URL encoded)

**Query Parameters:**
- `Download` (optional) - Set to `true` to force download as attachment

**Response Headers:**
- `Content-Length` - File size in bytes
- `Content-Type` - Based on file extension or `application/octet-stream`
- `Content-Disposition` - `attachment; filename={name}` (only when Download=true)

**Example Requests:**
```bash
# Download file
curl "http://localhost:8088/simple_upload/object/myfile.txt" -o myfile.txt

# Force download with attachment
curl "http://localhost:8088/simple_upload/object/myfile.txt?Download=true" -o myfile.txt
```

**Success Response:** File content with appropriate headers

---

### 3. Upload File

**Endpoint:** `POST /simple_upload/object/{filename}`

**Description:** Upload a single file to the server.

**Path Parameters:**
- `filename` - Target path and name for the file (URL encoded)

**Request Body:**

Option 1: Multipart form upload
```
Content-Type: multipart/form-data

-- Form field "File": the file to upload
-- Form field "Overwrite" (optional): "true" to overwrite existing files
```

Option 2: Binary stream upload
```
Content-Type: application/octet-stream

-- Raw file content in request body
-- Overwrite is automatically set to true
```

**Example Requests:**
```bash
# Multipart upload
curl -X POST "http://localhost:8088/simple_upload/object/uploads/document.pdf" \
  -F "File=@/path/to/document.pdf" \
  -F "Overwrite=true"

# Binary stream upload
curl -X POST "http://localhost:8088/simple_upload/object/uploads/document.pdf" \
  -H "Content-Type: application/octet-stream" \
  --data-binary @/path/to/document.pdf
```

**Success Response:**
```json
{
  "TraceID": "abc123",
  "Code": "success",
  "Message": "",
  "Data": null
}
```

**Error Responses:**
- `400 Bad Request` - Invalid filename or missing file
- `409 Conflict` - File exists and Overwrite=false (`file_exists`)

---

### 4. Delete File

**Endpoint:** `DELETE /simple_upload/object/{filename}`

**Description:** Delete a file from the server.

**Path Parameters:**
- `filename` - Path to the file to delete (URL encoded)

**Example Request:**
```bash
curl -X DELETE "http://localhost:8088/simple_upload/object/uploads/document.pdf"
```

**Success Response:**
```json
{
  "TraceID": "abc123",
  "Code": "success",
  "Message": "",
  "Data": null
}
```

**Error Responses:**
- `404 Not Found` - File does not exist (`not_found`)

---

### 5. List Directory Contents

**Endpoint:** `GET /simple_upload/objects/{path}`

**Description:** List files and directories at the specified path.

**Path Parameters:**
- `path` - Directory path to list (URL encoded). Use empty string for root.

**Response Structure:**
```json
{
  "TraceID": "abc123",
  "Code": "success",
  "Message": "",
  "Data": {
    "List": [
      {
        "Name": "file1.txt",
        "Size": 1024,
        "FileType": "file"
      },
      {
        "Name": "subdir",
        "Size": 0,
        "FileType": "dir"
      }
    ]
  }
}
```

**FileInfo Object:**
| Field | Type | Description |
|-------|------|-------------|
| `Name` | string | File or directory name |
| `Size` | number | File size in bytes (0 for directories) |
| `FileType` | string | `"file"` or `"dir"` |

**Example Requests:**
```bash
# List root directory
curl "http://localhost:8088/simple_upload/objects/"

# List specific directory
curl "http://localhost:8088/simple_upload/objects/documents/"
```

**Notes:**
- Directories are returned with `FileType: "dir"` and `Size: 0`
- The path parameter is automatically cleaned (trailing slashes handled)

---

### 6. Multi-file Upload

**Endpoint:** `POST /simple_upload/objects/`

**Description:** Upload multiple files to a specified directory in a single request.

**Request Body:** Multipart form with:
- `Dir` - Target directory path (required)
- `Overwrite` - Set to `"true"` to overwrite existing files (optional, default: `false`)
- Multiple `File` fields - Files to upload

**Example Request:**
```bash
curl -X POST "http://localhost:8088/simple_upload/objects/" \
  -F "Dir=uploads/photos" \
  -F "Overwrite=true" \
  -F "File=@/path/to/photo1.jpg" \
  -F "File=@/path/to/photo2.jpg" \
  -F "File=@/path/to/photo3.png"
```

**Success Response:**
```json
{
  "TraceID": "abc123",
  "Code": "success",
  "Message": "",
  "Data": null
}
```

**Notes:**
- All files are uploaded to the directory specified by `Dir`
- File names are preserved from the original upload

---

## Example API (`/example`)

These are example endpoints for testing the database connections.

### 1. List Examples (SQL)

**Endpoint:** `POST /example/list`

**Description:** Retrieve all example records from the SQL database.

**Request Body:** Empty or optional filter parameters

**Example Request:**
```bash
curl -X POST "http://localhost:8088/example/list" \
  -H "Content-Type: application/json"
```

**Success Response:**
```json
{
  "TraceID": "abc123",
  "Code": "success",
  "Message": "",
  "Data": [
    {
      "ID": "1",
      "Name": "Example 1"
    },
    {
      "ID": "2",
      "Name": "Example 2"
    }
  ]
}
```

---

### 2. Create Example (SQL)

**Endpoint:** `POST /example/create`

**Description:** Create a new example record in the SQL database.

**Request Body:**
```json
{
  "Name": "New Example"
}
```

**Example Request:**
```bash
curl -X POST "http://localhost:8088/example/create" \
  -H "Content-Type: application/json" \
  -d '{"Name": "New Example"}'
```

**Success Response:**
```json
{
  "TraceID": "abc123",
  "Code": "success",
  "Message": "",
  "Data": {
    "ID": "new-uuid-here"
  }
}
```

**Validation:**
- `Name` field is required

---

### 3. List Examples (MongoDB)

**Endpoint:** `POST /example/mongo_list`

**Description:** Retrieve all example records from MongoDB.

**Example Request:**
```bash
curl -X POST "http://localhost:8088/example/mongo_list" \
  -H "Content-Type: application/json"
```

**Success Response:** Same as SQL list example

---

### 4. Create Example (MongoDB)

**Endpoint:** `POST /example/mongo_create`

**Description:** Create a new example record in MongoDB.

**Request Body:** Same as SQL create example

**Example Request:**
```bash
curl -X POST "http://localhost:8088/example/mongo_create" \
  -H "Content-Type: application/json" \
  -d '{"Name": "Mongo Example"}'
```

**Success Response:** Same as SQL create example

---

## File Validation Rules

### Filename Validation

Filenames are validated to prevent directory traversal attacks:

**Invalid patterns:**
- Paths containing `./` or `../`
- Paths ending with `.` or `..`
- Windows-style backslashes are converted to forward slashes
- Leading/trailing slashes are trimmed

**Examples:**
| Input | Result | Valid |
|-------|--------|-------|
| `file.txt` | `file.txt` | Yes |
| `dir/file.txt` | `dir/file.txt` | Yes |
| `./file.txt` | Error | No |
| `../file.txt` | Error | No |
| `dir/../file.txt` | Error | No |
| `dir\file.txt` | `dir/file.txt` | Yes (converted) |

---

## Configuration Notes

### File Storage Policies

Files are stored according to `RealFilenamePolicy` in config:

| Policy | Behavior |
|--------|----------|
| `origin` | Keep original filename |
| `_` | Replace path separators with underscores |
| `uuid` | Store with UUID filename (default) |

### Upload Root

Files are stored under the `SimpleUploadRoot` path (default: `/aaa`) within the VFS.

---

## Static Files

The backend serves frontend static files at:
```
GET /static/*
```

Files are served from the `Http.StaticRoot` directory (default: `./static`).

---

## CORS

The API enables CORS with:
- Allowed origins: `*` (all)
- Allowed methods: `*` (all)

---

## Error Handling

Errors return HTTP status codes with JSON body:

```json
{
  "TraceID": "abc123",
  "Code": "invalid_arguments",
  "Message": "filename is invalid",
  "Data": null
}
```

| HTTP Status | Error Code | Description |
|-------------|------------|-------------|
| 400 | `invalid_arguments` | Bad request parameters |
| 404 | `not_found` | Resource not found |
| 409 | `file_exists` | File already exists |
| 503 | (various) | Server error |

---

## Testing with curl

### Quick Test Script

```bash
# Set base URL
BASE_URL="http://localhost:8088"

# Upload a file
curl -X POST "$BASE_URL/simple_upload/object/test.txt" \
  -F "File=@./test.txt"

# List files
curl "$BASE_URL/simple_upload/objects/"

# Download file
curl "$BASE_URL/simple_upload/object/test.txt" -o downloaded.txt

# Get metadata
curl -I "$BASE_URL/simple_upload/object/test.txt"

# Delete file
curl -X DELETE "$BASE_URL/simple_upload/object/test.txt"
```

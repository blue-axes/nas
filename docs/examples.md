# API Testing Examples

## Prerequisites

1. Start the backend server:
   ```bash
   cd backend
   go run main.go -config ./rootfs/etc/config.json
   ```

2. The server will start on `http://localhost:8088`

## Testing File Operations

### 1. Create a test file

```bash
# Create a test file
echo "Hello, NAS!" > test.txt
```

### 2. Upload the file

```bash
# Method 1: Multipart form upload
curl -X POST "http://localhost:8088/simple_upload/object/test.txt" \
  -F "File=@test.txt"

# Method 2: Binary stream upload
curl -X POST "http://localhost:8088/simple_upload/object/test.txt" \
  -H "Content-Type: application/octet-stream" \
  --data-binary @test.txt

# Method 3: Upload to subdirectory
curl -X POST "http://localhost:8088/simple_upload/object/documents/myfile.txt" \
  -F "File=@test.txt"
```

### 3. Get file metadata

```bash
# Get file info via HEAD request
curl -I "http://localhost:8088/simple_upload/object/test.txt"

# Response headers:
# Content-Length: 13
# Content-Type: application/octet-stream
# Content-Disposition: attachment; filename=test.txt
```

### 4. Download the file

```bash
# Download and save
curl "http://localhost:8088/simple_upload/object/test.txt" -o downloaded.txt

# Force download as attachment (triggers browser download)
curl "http://localhost:8088/simple_upload/object/test.txt?Download=true" -o downloaded.txt

# Display in terminal
curl "http://localhost:8088/simple_upload/object/test.txt"
```

### 5. List directory contents

```bash
# List root directory
curl "http://localhost:8088/simple_upload/objects/"

# List subdirectory
curl "http://localhost:8088/simple_upload/objects/documents/"

# Pretty print JSON response
curl -s "http://localhost:8088/simple_upload/objects/" | jq .
```

### 6. Upload multiple files

```bash
# Create multiple test files
echo "File 1" > file1.txt
echo "File 2" > file2.txt
echo "File 3" > file3.txt

# Upload multiple files to a directory
curl -X POST "http://localhost:8088/simple_upload/objects/" \
  -F "Dir=uploads" \
  -F "Overwrite=true" \
  -F "File=@file1.txt" \
  -F "File=@file2.txt" \
  -F "File=@file3.txt"
```

### 7. Update/overwrite a file

```bash
# Upload with overwrite flag
curl -X POST "http://localhost:8088/simple_upload/object/test.txt" \
  -F "File=@updated_test.txt" \
  -F "Overwrite=true"

# Try without overwrite (will fail if file exists)
curl -X POST "http://localhost:8088/simple_upload/object/test.txt" \
  -F "File=@updated_test.txt"
# Returns: 409 Conflict with code "file_exists"
```

### 8. Delete a file

```bash
# Delete a single file
curl -X DELETE "http://localhost:8088/simple_upload/object/test.txt"

# Delete a file in subdirectory
curl -X DELETE "http://localhost:8088/simple_upload/object/documents/myfile.txt"
```

## Testing Example Endpoints

### SQL Database Examples

```bash
# Create an example
curl -X POST "http://localhost:8088/example/create" \
  -H "Content-Type: application/json" \
  -d '{"Name": "My First Example"}'

# Response:
# {
#   "TraceID": "abc123",
#   "Code": "success",
#   "Message": "",
#   "Data": {
#     "ID": "550e8400-e29b-41d4-a716-446655440000"
#   }
# }

# List all examples
curl -X POST "http://localhost:8088/example/list" \
  -H "Content-Type: application/json"
```

### MongoDB Examples

```bash
# Create MongoDB example
curl -X POST "http://localhost:8088/example/mongo_create" \
  -H "Content-Type: application/json" \
  -d '{"Name": "MongoDB Example"}'

# List MongoDB examples
curl -X POST "http://localhost:8088/example/mongo_list" \
  -H "Content-Type: application/json"
```

## Testing Error Cases

### Invalid filename

```bash
# Directory traversal attempt (blocked)
curl -X POST "http://localhost:8088/simple_upload/object/../../../etc/passwd" \
  -F "File=@test.txt"
# Returns: 400 Bad Request with code "invalid_arguments"

# Path with ./ (blocked)
curl -X POST "http://localhost:8088/simple_upload/object/./test.txt" \
  -F "File=@test.txt"
# Returns: 400 Bad Request
```

### File not found

```bash
# Try to download non-existent file
curl "http://localhost:8088/simple_upload/object/nonexistent.txt"
# Returns: 404 Not Found with code "not_found"

# Try to delete non-existent file
curl -X DELETE "http://localhost:8088/simple_upload/object/nonexistent.txt"
# Returns: 404 Not Found
```

### Missing required fields

```bash
# Create example without Name field
curl -X POST "http://localhost:8088/example/create" \
  -H "Content-Type: application/json" \
  -d '{}'
# Returns: 400 Bad Request
```

## File Path Handling

The API automatically handles paths:

```bash
# These all work correctly:
curl -X POST "http://localhost:8088/simple_upload/object/file.txt" \
  -F "File=@test.txt"

curl -X POST "http://localhost:8088/simple_upload/object//file.txt" \
  -F "File=@test.txt"  # Double slashes are normalized

curl -X POST "http://localhost:8088/simple_upload/object/dir\file.txt" \
  -F "File=@test.txt"  # Backslashes converted to forward slashes

# But these are blocked (security):
# curl -X POST "http://localhost:8088/simple_upload/object/../file.txt"
# curl -X POST "http://localhost:8088/simple_upload/object/./file.txt"
```

## Response Format

All responses follow this format:

```json
{
  "TraceID": "unique-request-id",
  "Code": "success|error_code",
  "Message": "Human readable message",
  "Data": {}  // Response payload or null
}
```

### Success Response

```json
{
  "TraceID": "abc123def456",
  "Code": "success",
  "Message": "",
  "Data": null
}
```

### Error Response

```json
{
  "TraceID": "abc123def456",
  "Code": "file_exists",
  "Message": "file has already exists",
  "Data": null
}
```

## PowerShell Examples

For Windows PowerShell users:

```powershell
# Upload file
Invoke-RestMethod -Uri "http://localhost:8088/simple_upload/object/test.txt" `
  -Method Post `
  -Form @{File=Get-Item "test.txt"}

# Download file
Invoke-WebRequest -Uri "http://localhost:8088/simple_upload/object/test.txt" `
  -OutFile "downloaded.txt"

# List directory
Invoke-RestMethod -Uri "http://localhost:8088/simple_upload/objects/" | ConvertTo-Json -Depth 5

# Create example
Invoke-RestMethod -Uri "http://localhost:8088/example/create" `
  -Method Post `
  -ContentType "application/json" `
  -Body '{"Name": "PowerShell Example"}'
```

## Python Examples

```python
import requests

BASE_URL = "http://localhost:8088"

# Upload file
with open("test.txt", "rb") as f:
    response = requests.post(
        f"{BASE_URL}/simple_upload/object/test.txt",
        files={"File": f}
    )
    print(response.json())

# Download file
response = requests.get(f"{BASE_URL}/simple_upload/object/test.txt")
with open("downloaded.txt", "wb") as f:
    f.write(response.content)

# List directory
response = requests.get(f"{BASE_URL}/simple_upload/objects/")
print(response.json())

# Create example
response = requests.post(
    f"{BASE_URL}/example/create",
    json={"Name": "Python Example"}
)
print(response.json())
```

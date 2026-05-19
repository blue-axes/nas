# API Testing Examples

## Authentication

```bash
BASE="http://localhost:8088"

# Login (get session cookie)
curl -c cookies.txt -X POST "$BASE/api/users/login" \
  -H "Content-Type: application/json" \
  -d '{"Username":"admin","Password":"admin"}'

# Get current user
curl -b cookies.txt "$BASE/api/users/me"

# Logout
curl -b cookies.txt -X POST "$BASE/api/users/logout"
```

## File Operations

### Upload

```bash
# Single file (multipart)
curl -b cookies.txt -X POST "$BASE/simple_upload/object/img/photo.jpg" \
  -F "File=@photo.jpg"

# Single file (binary stream)
curl -b cookies.txt -X POST "$BASE/simple_upload/object/img/photo.jpg" \
  -H "Content-Type: application/octet-stream" \
  --data-binary @photo.jpg

# Multi-file upload
curl -b cookies.txt -X POST "$BASE/simple_upload/objects/" \
  -F "Dir=uploads" -F "Overwrite=true" \
  -F "File=@photo1.jpg" -F "File=@photo2.jpg"
```

### Download

```bash
curl -b cookies.txt "$BASE/simple_upload/object/img/photo.jpg" -o photo.jpg
curl -b cookies.txt "$BASE/simple_upload/object/img/photo.jpg?Download=true" -o photo.jpg
```

### Metadata

```bash
curl -b cookies.txt -I "$BASE/simple_upload/object/img/photo.jpg"
```

### Directory Listing

```bash
# Root
curl -b cookies.txt "$BASE/simple_upload/objects/"

# Subdirectory
curl -b cookies.txt "$BASE/simple_upload/objects/img/"
curl -b cookies.txt "$BASE/simple_upload/objects/img/vacation/"
```

### Create Directory

```bash
curl -b cookies.txt -X POST "$BASE/simple_upload/mkdir/img/vacation"
```

### Update Tags

```bash
curl -b cookies.txt -X PATCH "$BASE/simple_upload/object/img/photo.jpg" \
  -H "Content-Type: application/json" \
  -d '{"Tags":["vacation","family","2024"]}'
```

### Search

```bash
# By keyword
curl -b cookies.txt "$BASE/simple_upload/search?Keyword=photo"

# By tag
curl -b cookies.txt "$BASE/simple_upload/search?Tag=vacation"

# Combined
curl -b cookies.txt "$BASE/simple_upload/search?Keyword=photo&Tag=vacation"
```

### Delete

```bash
# Delete file
curl -b cookies.txt -X DELETE "$BASE/simple_upload/object/img/photo.jpg"

# Delete directory (cascades)
curl -b cookies.txt -X DELETE "$BASE/simple_upload/object/img/vacation/"
```

## Admin Operations

### User Management

```bash
# List users
curl -b cookies.txt "$BASE/api/users"

# Create user
curl -b cookies.txt -X POST "$BASE/api/users" \
  -H "Content-Type: application/json" \
  -d '{"Username":"guest","Password":"123456","CanRead":true,"CanWrite":false}'

# Update permissions
curl -b cookies.txt -X PUT "$BASE/api/users/guest" \
  -H "Content-Type: application/json" \
  -d '{"CanRead":true,"CanWrite":true}'

# Delete user
curl -b cookies.txt -X DELETE "$BASE/api/users/guest"

# Change password
curl -b cookies.txt -X PUT "$BASE/api/users/admin/password" \
  -H "Content-Type: application/json" \
  -d '{"CurrentPassword":"admin","NewPassword":"newpass"}'
```

### File Scan

```bash
curl -b cookies.txt -X POST "$BASE/api/scanfs"
# Returns: { "Data": { "Total": 150, "New": 12, "Skipped": 138 } }
```

## WebDAV Access

```bash
# List directory via PROPFIND
curl -u admin:admin -X PROPFIND "$BASE/webdav/"

# Upload file
curl -u admin:admin -T photo.jpg "$BASE/webdav/img/photo.jpg"

# Download file  
curl -u admin:admin "$BASE/webdav/img/photo.jpg" -o photo.jpg

# Create directory via MKCOL
curl -u admin:admin -X MKCOL "$BASE/webdav/img/albums/"

# Delete file
curl -u admin:admin -X DELETE "$BASE/webdav/img/photo.jpg"
```

## Error Cases

```bash
# Path traversal (blocked)
curl -b cookies.txt "$BASE/simple_upload/object/../../../etc/passwd"
# → 400: "invalid_arguments"

# File not found
curl -b cookies.txt "$BASE/simple_upload/object/nonexistent.txt"
# → 404: "not_found"

# Unauthorized
curl "$BASE/api/users"
# → 401: "unauthorized"

# Permission denied (user without CanRead)
curl -b guest_cookies.txt "$BASE/simple_upload/objects/"
# → 403: "read permission required"
```

## Python Examples

```python
import requests

BASE = "http://localhost:8088"
session = requests.Session()

# Login
session.post(f"{BASE}/api/users/login", json={"Username": "admin", "Password": "admin"})

# Upload
with open("photo.jpg", "rb") as f:
    session.post(f"{BASE}/simple_upload/object/img/photo.jpg", files={"File": f})

# List directory
resp = session.get(f"{BASE}/simple_upload/objects/img/")
for item in resp.json()["Data"]["List"]:
    print(f"  {item['FileType']:4s} {item['Name']}")

# Search
resp = session.get(f"{BASE}/simple_upload/search", params={"Keyword": "photo", "Tag": "vacation"})
print(resp.json())

# Scan filesystem
resp = session.post(f"{BASE}/api/scanfs")
print(f"Scanned: {resp.json()['Data']}")
```
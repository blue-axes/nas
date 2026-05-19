# NAS API Documentation

## Base URL

```
http://localhost:8088
```

## Authentication

Authentication uses session cookies. The default admin account is `admin` / `admin`.

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/users/login` | POST | Login, sets session cookie |
| `/api/users/logout` | POST | Logout, clears session |
| `/api/users/me` | GET | Get current user info |

### Login

**Endpoint:** `POST /api/users/login`

```json
{ "Username": "admin", "Password": "admin" }
```

**Response:**
```json
{
  "TraceID": "...",
  "Code": "success",
  "Data": { "Username": "admin", "CanRead": true, "CanWrite": true, "IsAdmin": true }
}
```

### Get Current User

**Endpoint:** `GET /api/users/me`

Returns the current authenticated user info. Requires valid session cookie.

---

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

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `success` | 200 | Request successful |
| `unknown` | 503 | Unknown error |
| `invalid_arguments` | 400 | Invalid request parameters |
| `not_found` | 404 | Resource not found |
| `file_exists` | 409 | File already exists |
| `403` | 403 | Permission denied |

---

## File API (`/simple_upload`)

### 1. Get File Metadata

**Endpoint:** `HEAD /simple_upload/object/{filename}`

**Path Parameters:**
- `filename` - Path to the file (URL encoded)

**Response Headers:**
- `Content-Length` - File size in bytes
- `Content-Disposition` - Attachment filename
- `Content-Type` - `application/octet-stream`

**Example:**
```bash
curl -I "http://localhost:8088/simple_upload/object/myfile.txt"
```

---

### 2. Download File

**Endpoint:** `GET /simple_upload/object/{filename}`

**Path Parameters:**
- `filename` - Path to the file (URL encoded)

**Query Parameters:**
- `Download` (optional, boolean) - Force download as attachment

**Example:**
```bash
curl "http://localhost:8088/simple_upload/object/myfile.txt" -o myfile.txt
```

---

### 3. Upload File

**Endpoint:** `POST /simple_upload/object/{filename}`

**Path Parameters:**
- `filename` - Target path for the file

**Request Body:**

Option 1 — Multipart form:
```
Content-Type: multipart/form-data
File: binary file
Overwrite: (optional) "true" to overwrite
```

Option 2 — Binary stream:
```
Content-Type: application/octet-stream
Body: raw file content (overwrite is automatic)
```

**Example:**
```bash
curl -X POST "http://localhost:8088/simple_upload/object/photo.jpg" -F "File=@photo.jpg"
```

**Errors:**
- `400` — Invalid filename
- `409` — File exists and Overwrite=false

---

### 4. Delete File / Directory

**Endpoint:** `DELETE /simple_upload/object/{filename}`

**Description:** Delete a file or directory. Deleting a directory cascades to all files under it.

**Example:**
```bash
curl -X DELETE "http://localhost:8088/simple_upload/object/photo.jpg"
curl -X DELETE "http://localhost:8088/simple_upload/object/subdir/"
```

**Errors:**
- `404` — File not found

---

### 5. Update File Tags

**Endpoint:** `PATCH /simple_upload/object/{filename}`

**Description:** Update tags on a file.

**Request Body:**
```json
{ "Tags": ["tag1", "tag2", "tag3"] }
```

**Example:**
```bash
curl -X PATCH "http://localhost:8088/simple_upload/object/photo.jpg" \
  -H "Content-Type: application/json" \
  -d '{"Tags":["vacation","family"]}'
```

---

### 6. List Directory Contents

**Endpoint:** `GET /simple_upload/objects/{path}`

**Description:** List files and directories at the specified path. Use `/objects` (no path) for root.

**Path Parameters:**
- `path` - Directory path (URL encoded). Empty for root.

**Response:**
```json
{
  "Data": {
    "List": [
      {
        "Name": "photo.jpg",
        "Size": 102400,
        "FileType": "file",
        "Tags": ["vacation"]
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
| `Size` | uint64 | File size (0 for directories) |
| `FileType` | string | `"file"` or `"dir"` |
| `Tags` | []string | File tags (only on files) |

**Example:**
```bash
curl "http://localhost:8088/simple_upload/objects/img/"
```

---

### 7. Multi-file Upload

**Endpoint:** `POST /simple_upload/objects/`

**Request Body:** Multipart form
- `Dir` — Target directory path (required)
- `Overwrite` — `"true"` to overwrite (optional)
- `File` — One or more files

**Example:**
```bash
curl -X POST "http://localhost:8088/simple_upload/objects/" \
  -F "Dir=uploads" -F "Overwrite=true" \
  -F "File=@photo1.jpg" -F "File=@photo2.jpg"
```

---

### 8. Create Directory

**Endpoint:** `POST /simple_upload/mkdir/{path}`

**Description:** Create a new directory.

**Example:**
```bash
curl -X POST "http://localhost:8088/simple_upload/mkdir/img/vacation"
```

---

### 9. Search Files

**Endpoint:** `GET /simple_upload/search`

**Query Parameters:**
- `Keyword` — Search by filename keyword
- `Tag` — Filter by tag

**Example:**
```bash
curl "http://localhost:8088/simple_upload/search?Keyword=photo&Tag=vacation"
```

---

### 10. Preview File

**Endpoint:** `GET /simple_upload/preview/{path}`

**Description:** Convert Office/PDF documents to HTML preview. Requires `libreoffice` on the server.

**Example:**
```bash
curl "http://localhost:8088/simple_upload/preview/docs/report.docx"
```

---

## User Management API (`/api/users`)

**Requires admin permission.**

### List Users

**Endpoint:** `GET /api/users`

**Response:**
```json
{
  "Data": [
    { "Username": "admin", "CanRead": true, "CanWrite": true, "IsAdmin": true },
    { "Username": "guest", "CanRead": true, "CanWrite": false, "IsAdmin": false }
  ]
}
```

### Create User

**Endpoint:** `POST /api/users`

```json
{ "Username": "guest", "Password": "123456", "CanRead": true, "CanWrite": false }
```

### Update User Permissions

**Endpoint:** `PUT /api/users/{username}`

```json
{ "CanRead": true, "CanWrite": true }
```

### Delete User

**Endpoint:** `DELETE /api/users/{username}`

### Change Password

**Endpoint:** `PUT /api/users/{username}/password`

```json
{ "CurrentPassword": "oldpass", "NewPassword": "newpass" }
```

---

## File Scanning API

### Scan Filesystem

**Endpoint:** `POST /api/scanfs`

**Requires admin permission.**

Scans the filesystem under `SimpleUploadRoot` and syncs missing file records into the database.

**Response:**
```json
{
  "Data": { "Total": 150, "New": 12, "Skipped": 138 }
}
```

---

## WebDAV

The server exposes a WebDAV mount at:

```
http://localhost:8088/webdav/
```

Windows mount: `\\nas.local@8088\webdav\`  
macOS: Finder → Go → Connect to Server → `http://nas.local:8088/webdav/`  
Linux: `mount -t davfs http://nas.local:8088/webdav/ /mnt/nas`

Supports Basic Auth. Permission checks (`CanRead` / `CanWrite`) are enforced at the filesystem level.

---

## mDNS Discovery

The server broadcasts its service on the local network via mDNS with the hostname `nas.local`. Configuration:

```json
"MDNS": {
    "Enabled": true,
    "ServiceName": "_http._tcp",
    "Hostname": "nas.local",
    "Info": "NAS Web Service"
}
```

LAN devices can access the NAS at `http://nas.local:8088` without knowing the IP.

---

## Configuration

### File Storage Policies

| Policy | Behavior |
|--------|----------|
| `origin` | Keep original filename |
| `_` | Replace path separators with underscores |
| `uuid` | Store with UUID filename (default) |

### Default Directories

On startup, three default directories are created under `SimpleUploadRoot`:
- `img` — Image files
- `video` — Video files
- `other` — Other files

### File Validation

| Input | Result | Valid |
|-------|--------|-------|
| `file.txt` | `file.txt` | Yes |
| `dir/file.txt` | `dir/file.txt` | Yes |
| `./file.txt` | Error | No |
| `../file.txt` | Error | No |
| `dir\file.txt` | `dir/file.txt` | Yes (converted) |

---

## Web Frontend

The frontend SPA is served at `/`. It includes:
- **所有文件** — Root directory browser, table view with file management
- **图片文件** — Image gallery (grid + preview modal)
- **视频文件** — Video browser (grid + player modal)
- **普通文件** — File table (sorted view of `/other/` directory)
- **用户管理** — Admin user CRUD (admin only)
- **扫描文件** — Filesystem sync button (admin only)

Static files are served from `Http.StaticRoot` (default: `./static`).

---

## CORS

The API enables CORS with all origins and methods allowed.
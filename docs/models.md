# Data Models

## Request/Response Structures

### Standard Response Wrapper

All API responses are wrapped in this structure:

```go
type respStruct struct {
    TraceID string      `json:"TraceID"`  // Unique request identifier
    Code    string      `json:"Code"`     // Status code (e.g., "success")
    Message string      `json:"Message"`  // Human-readable message
    Data    interface{} `json:"Data"`     // Response payload
}
```

**Example:**
```json
{
  "TraceID": "550e8400-e29b-41d4-a716-446655440000",
  "Code": "success",
  "Message": "",
  "Data": {}
}
```

---

## File Models

### File

Internal representation of a file in the database:

```go
type File struct {
    ID   uint   `json:"ID"`   // Auto-increment ID
    Name string `json:"Name"` // Original filename with path
    Ext  string `json:"Ext"`  // File extension
    Path string `json:"Path"` // Physical storage path
    Size uint64 `json:"Size"` // File size in bytes
    Md5  string `json:"Md5"`  // MD5 checksum
}
```

**Example:**
```json
{
  "ID": 1,
  "Name": "documents/report.pdf",
  "Ext": ".pdf",
  "Path": "/aaa/documents/a1b2c3d4.pdf",
  "Size": 102400,
  "Md5": "d41d8cd98f00b204e9800998ecf8427e"
}
```

---

### FileInfo

Directory listing item (used in ReadDir response):

```go
type FileInfo struct {
    Name     string `json:"Name"`     // File or directory name
    Size     uint64 `json:"Size"`     // File size (0 for directories)
    FileType string `json:"FileType"` // "file" or "dir"
}
```

**Examples:**

File entry:
```json
{
  "Name": "document.txt",
  "Size": 1024,
  "FileType": "file"
}
```

Directory entry:
```json
{
  "Name": "subfolder",
  "Size": 0,
  "FileType": "dir"
}
```

---

### Filename

Path parameter structure for file operations:

```go
type Filename struct {
    Name string `param:"*"` // Captures the entire path after /object/
}
```

---

## Example Models

### Example

Simple example record for database testing:

```go
type Example struct {
    ID   string `json:"ID"`   // UUID identifier
    Name string `json:"Name"` // Example name
}
```

**Example:**
```json
{
  "ID": "550e8400-e29b-41d4-a716-446655440000",
  "Name": "My Example"
}
```

---

## Error Codes

### ErrorCode Type

```go
type ErrorCode = string
```

### Defined Error Codes

| Code | Constant | HTTP Status | Description |
|------|----------|-------------|-------------|
| `success` | `ErrCodeSuccess` | 200 | Request successful |
| `unknown` | `ErrCodeUnknown` | 503 | Unknown error |
| `invalid_arguments` | `ErrCodeInvalidArgs` | 400 | Invalid request parameters |
| `not_found` | `ErrCodeNotFound` | 404 | Resource not found |
| `file_exists` | `ErrCodeFileExists` | 409 | File already exists |
| `file_checksum_invalid` | `ErrFileCheckSumInvalid` | 503 | File checksum verification failed |

---

## Configuration Models

### Config (JSON/YAML)

Root configuration structure:

```go
type Config struct {
    Http     HttpConfig     `json:"Http" yaml:"Http"`
    Log      LogConfig      `json:"Log" yaml:"Log"`
    Database DatabaseConfig `json:"Database" yaml:"Database"`
    Nas      NasConfig      `json:"Nas" yaml:"Nas"`
}
```

### HttpConfig

```go
type HttpConfig struct {
    ListenAddress string `json:"ListenAddress" yaml:"ListenAddress"` // e.g., "0.0.0.0"
    ListenPort    uint16 `json:"ListenPort" yaml:"ListenPort"`       // e.g., 8088
    StaticRoot    string `json:"StaticRoot" yaml:"StaticRoot"`       // e.g., "./static"
}
```

### LogConfig

```go
type LogConfig struct {
    Level string `json:"Level" yaml:"Level"` // "debug", "info", "warn", "error"
}
```

### DatabaseConfig

```go
type DatabaseConfig struct {
    Rdb   *RdbConfig   `json:"Rdb" yaml:"Rdb"`
    Mongo *MongoConfig `json:"Mongo" yaml:"Mongo"`
}
```

### RdbConfig

```go
type RdbConfig struct {
    DriverType            RdbDriverType `json:"DriverType" yaml:"DriverType"`     // "sqlite" or "postgres"
    Debug                 bool          `json:"Debug" yaml:"Debug"`               // Enable SQL logging
    DSN                   string        `json:"DSN" yaml:"DSN"`                   // Connection string
    MaxIdleConnCount      int           `json:"MaxIdleConnCount" yaml:"MaxIdleConnCount"`
    MaxConnCount          int           `json:"MaxConnCount" yaml:"MaxConnCount"`
    ConnMaxIdleTimeSecond int           `json:"ConnMaxIdleTimeSecond" yaml:"ConnMaxIdleTimeSecond"`
    AutoMigrateLevel      string        `json:"AutoMigrateLevel" yaml:"AutoMigrateLevel"` // "auto" or "must"
}
```

### MongoConfig

```go
type MongoConfig struct {
    Debug                 bool   `json:"Debug" yaml:"Debug"`
    Address               string `json:"Address" yaml:"Address"`
    Port                  uint16 `json:"Port" yaml:"Port"`
    Username              string `json:"Username" yaml:"Username"`
    Password              string `json:"Password" yaml:"Password"`
    Database              string `json:"Database" yaml:"Database"`
    MaxIdleConnCount      int    `json:"MaxIdleConnCount" yaml:"MaxIdleConnCount"`
    MaxConnCount          int    `json:"MaxConnCount" yaml:"MaxConnCount"`
    ConnMaxIdleTimeSecond int    `json:"ConnMaxIdleTimeSecond" yaml:"ConnMaxIdleTimeSecond"`
    AutoMigrateLevel      string `json:"AutoMigrateLevel" yaml:"AutoMigrateLevel"`
}
```

### NasConfig

```go
type NasConfig struct {
    SimpleUploadRoot   string             `json:"SimpleUploadRoot" yaml:"SimpleUploadRoot"`     // e.g., "/aaa"
    RealFilenamePolicy RealFilenamePolicy `json:"RealFilenamePolicy" yaml:"RealFilenamePolicy"` // "origin", "_", or "uuid"
}
```

### RdbDriverType

```go
type RdbDriverType string

const (
    DriverTypeSqlite   RdbDriverType = "sqlite"
    DriverTypePostgres RdbDriverType = "postgres"
)
```

### RealFilenamePolicy

```go
type RealFilenamePolicy string

const (
    RFNP_Origin    RealFilenamePolicy = "origin"    // Keep original filename
    RFNP_Underline RealFilenamePolicy = "_"         // Replace / with _
    RFNP_UUID      RealFilenamePolicy = "uuid"      // Use UUID as filename
)
```

---

## API Schema Definitions

### Request Schemas

**Create Example Request:**
```json
{
  "Name": "required string"
}
```

**Multi Upload Form:**
```
Dir: string (required)
Overwrite: "true" or "false" (optional)
File: binary file(s) (required)
```

**Single Upload Form:**
```
File: binary file (required)
Overwrite: boolean (optional)
```

---

### Response Schemas

**Success Response (no data):**
```json
{
  "TraceID": "string",
  "Code": "success",
  "Message": "",
  "Data": null
}
```

**Create Example Response:**
```json
{
  "TraceID": "string",
  "Code": "success",
  "Message": "",
  "Data": {
    "ID": "uuid"
  }
}
```

**List Examples Response:**
```json
{
  "TraceID": "string",
  "Code": "success",
  "Message": "",
  "Data": [
    {
      "ID": "uuid",
      "Name": "string"
    }
  ]
}
```

**Directory Listing Response:**
```json
{
  "TraceID": "string",
  "Code": "success",
  "Message": "",
  "Data": {
    "List": [
      {
        "Name": "string",
        "Size": 0,
        "FileType": "file|dir"
      }
    ]
  }
}
```

---

## Database Tables

### files

Table for storing file metadata (SQL):

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY, AUTOINCREMENT | Auto-increment ID |
| name | TEXT | NOT NULL, UNIQUE | Original file path |
| ext | TEXT | | File extension |
| path | TEXT | NOT NULL | Physical storage path |
| size | INTEGER | | File size in bytes |
| md5 | TEXT | | MD5 checksum |

### examples

Table for example data (SQL):

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | TEXT | PRIMARY KEY | UUID |
| name | TEXT | | Example name |

---

## VFS Models

### MountFs

Virtual filesystem with mount points:

```go
type MountFs interface {
    Stat(name string) (fs.FileInfo, error)
    Remove(name string) error
    RemoveAll(path string) error
    OpenFile(name string, flag int, perm fs.FileMode) (File, error)
    Mkdir(dir string, perm fs.FileMode) error
    MkdirAll(path string, perm fs.FileMode) error
    ReadDir(dir string) ([]fs.DirEntry, error)
    TempDir() string
    Mount(dir string, fs VFS) error
    Umount(dir string) error
}
```

### OsFs

OS filesystem implementation:

```go
type OsFsConf struct {
    RootDir string  // Root directory for the filesystem
}
```

---

## Request Flow

1. Request arrives at Echo server
2. Middleware adds context with TraceID
3. Handler validates and binds parameters
4. Service layer processes business logic
5. Store layer interacts with database
6. VFS layer handles file operations
7. Response is wrapped and returned

# Data Models

## File Model

Internal representation of a file/directory in the database:

```go
type File struct {
    ID    uint     `json:"ID"`
    Name  string   `json:"Name"`   // Virtual path (e.g., "img/photo.jpg")
    Ext   string   `json:"Ext"`    // File extension
    Path  string   `json:"Path"`   // Real filesystem path (e.g., "/aaa/img/uuid.jpg")
    Size  uint64   `json:"Size"`   // File size in bytes
    Md5   string   `json:"Md5"`    // MD5 checksum
    IsDir bool     `json:"IsDir"`  // Whether this is a directory
    Tags  []string `json:"Tags"`   // User-defined tags
}
```

### FileInfo (API Response)

Directory listing item:

```go
type FileInfo struct {
    Name     string   `json:"Name"`
    Size     uint64   `json:"Size"`
    FileType string   `json:"FileType"` // "file" or "dir"
    Tags     []string `json:"Tags,omitempty"`
}
```

**Example:**
```json
{
  "Name": "photo.jpg",
  "Size": 102400,
  "FileType": "file",
  "Tags": ["vacation", "family"]
}
```

---

## User Model

```go
type UserInfo struct {
    Username string `json:"Username"`
    CanRead  bool   `json:"CanRead"`
    CanWrite bool   `json:"CanWrite"`
    IsAdmin  bool   `json:"IsAdmin"`
}
```

---

## Scan Result

```go
type ScanResult struct {
    Total   int `json:"Total"`
    New     int `json:"New"`
    Skipped int `json:"Skipped"`
}
```

---

## Error Codes

| Code | Constant | HTTP Status | Description |
|------|----------|-------------|-------------|
| `success` | `ErrCodeSuccess` | 200 | Request successful |
| `unknown` | `ErrCodeUnknown` | 503 | Unknown error |
| `invalid_arguments` | `ErrCodeInvalidArgs` | 400 | Invalid request parameters |
| `not_found` | `ErrCodeNotFound` | 404 | Resource not found |
| `file_exists` | `ErrCodeFileExists` | 409 | File already exists |
| `file_checksum_invalid` | `ErrFileCheckSumInvalid` | 503 | Checksum verification failed |

---

## Configuration Models

### Config

```go
type Config struct {
    Http     HttpConfig     `json:"Http" yaml:"Http"`
    Log      LogConfig      `json:"Log" yaml:"Log"`
    Database DatabaseConfig `json:"Database" yaml:"Database"`
    Nas      NasConfig      `json:"Nas" yaml:"Nas"`
    MDNS     MDNSConfig     `json:"MDNS" yaml:"MDNS"`
}
```

### HttpConfig

```go
type HttpConfig struct {
    ListenAddress string     `json:"ListenAddress" yaml:"ListenAddress"`
    ListenPort    uint16     `json:"ListenPort" yaml:"ListenPort"`
    StaticRoot    string     `json:"StaticRoot" yaml:"StaticRoot"`
    CertFile      string     `json:"CertFile" yaml:"CertFile"`
    KeyFile       string     `json:"KeyFile" yaml:"KeyFile"`
    Auth          AuthConfig `json:"Auth" yaml:"Auth"`
}
```

### AuthConfig

```go
type AuthConfig struct {
    Enabled            bool   `json:"Enabled" yaml:"Enabled"`
    CookieName         string `json:"CookieName" yaml:"CookieName"`           // default: "nas_session"
    SessionExpireHours int    `json:"SessionExpireHours" yaml:"SessionExpireHours"` // default: 24
}
```

### MDNSConfig

```go
type MDNSConfig struct {
    Enabled     bool   `json:"Enabled" yaml:"Enabled"`
    ServiceName string `json:"ServiceName" yaml:"ServiceName"` // default: "_http._tcp"
    Hostname    string `json:"Hostname" yaml:"Hostname"`       // default: "nas.local"
    Info        string `json:"Info" yaml:"Info"`               // default: "NAS Web Service"
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

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `DriverType` | string | `"sqlite"` | `"sqlite"` or `"postgres"` |
| `Debug` | bool | false | Enable SQL logging |
| `DSN` | string | `"nas.db?mode=rwc"` | Connection string |
| `MaxConnCount` | int | 5 | Max open connections |
| `MaxIdleConnCount` | int | 2 | Max idle connections |
| `ConnMaxIdleTimeSecond` | int | 300 | Idle timeout |
| `AutoMigrateLevel` | string | - | `"auto"` or `"must"` |

### NasConfig

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `SimpleUploadRoot` | string | - | Root directory for uploaded files |  
| `RealFilenamePolicy` | string | `"uuid"` | `"origin"`, `"_"`, or `"uuid"` |

---

## Database Tables

### file_object

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY | Auto-increment |
| name | TEXT | UNIQUE, NOT NULL | Virtual file path |
| ext | VARCHAR(50) | | File extension |
| path | VARCHAR(1024) | | Real filesystem path |
| size | INTEGER | DEFAULT 0 | File size |
| md5_sum | TEXT | | MD5 checksum |
| is_dir | BOOLEAN | DEFAULT false | Directory flag |
| tags | TEXT | | JSON array of tags |
| created_at | DATETIME | | Auto timestamp |
| updated_at | DATETIME | | Auto timestamp |

### users

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY | Auto-increment |
| username | VARCHAR(128) | UNIQUE, NOT NULL | Login name |
| password_hash | VARCHAR(256) | NOT NULL | bcrypt hash |
| can_read | BOOLEAN | DEFAULT true | Read permission |
| can_write | BOOLEAN | DEFAULT false | Write permission |
| is_admin | BOOLEAN | DEFAULT false | Admin flag |

---

## Default Values

| Config Field | Default |
|-------------|---------|
| `Http.ListenAddress` | `"0.0.0.0"` |
| `Http.ListenPort` | 80 |
| `Http.StaticRoot` | `"./"` |
| `Http.Auth.CookieName` | `"nas_session"` |
| `Http.Auth.SessionExpireHours` | 24 |
| `Log.Level` | `"info"` |
| `Database.Rdb.DriverType` | `"sqlite"` |
| `Nas.RealFilenamePolicy` | `"uuid"` |
| `MDNS.ServiceName` | `"_http._tcp"` |
| `MDNS.Hostname` | `"nas.local"` |
| `MDNS.Info` | `"NAS Web Service"` |
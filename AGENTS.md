# AGENTS.md

## API Documentation
- Complete API docs in `docs/` directory
- Main docs: `docs/README.md`
- OpenAPI spec: `docs/openapi.yaml`
- Examples: `docs/examples.md`

## Repository structure
- `frontend/` – React + Vite SPA (Ant Design, React Router hash router)
- `backend/` – Go Echo HTTP server (GORM, VFS, SQLite/Postgres/Mongo)
- `scripts/` – SMB (Samba) and NFS setup scripts

## Key commands
**Frontend**
- `npm run dev` – Vite dev server (proxies `/simple_upload` to backend :8088)
- `npm run build` – builds to `dist/` (served by backend as static)
- `npm run lint` – ESLint

**Backend**
- `go run main.go -config ./config.json` (default config path)
- `go test ./...` – only one test file (`vfs/mountfs_test.go`)
- Build with `go build`; no Makefile

## Architecture notes
- Backend serves static files from `Http.StaticRoot` (default `./static`). Place frontend build output there.
- Frontend base path is `/static/` (see `vite.config.js`). Backend must serve static root at that path.
- Database migrations run automatically based on `Database.Rdb.AutoMigrateLevel` (`auto` or `must`).
- The `file` model includes `Tags` (JSON text, stored as `["tag1","tag2"]`) and `IsDir` flag.
- VFS root is OS `/` (see `backend/vfs.go`). Be careful with file operations.
- Frontend uses hash router (`/#/image` etc.).

## Frontend components
- `SearchBar` — keyword + tag filter input, auto-searches on change
- `TagEditor` — Popover with tag add/remove UI, saves via PATCH API
- `DocPreviewer` — Modal with iframe for PDF/Office document preview
- `PathTravel` — Breadcrumb navigation component
- `useScreenWidth` hook — reactive screen width for responsive layout

## Features
- **Dark tech theme** — CSS variables in `tech-theme.css` (custom scrollbar, glow effects, glassmorphism)
- **Mobile responsive** — 768px/480px breakpoints, hamburger menu overlay, adaptive grids
- **Tags + Search** — backend `GET /simple_upload/search?Keyword=&Tag=` + store `SearchFiles`
- **Folder support** — `POST /simple_upload/mkdir/*` creates dir entries, delete cascades
- **Document preview** — `GET /simple_upload/preview/*` uses libreoffice to convert Office docs to PDF, cached in `/tmp/nas_preview_cache/`

## API endpoints
| Method | Path | Description |
|--------|------|-------------|
| HEAD | `/simple_upload/object/*` | File metadata |
| GET | `/simple_upload/object/*` | Download/display file |
| POST | `/simple_upload/object/*` | Upload file |
| DELETE | `/simple_upload/object/*` | Delete file/directory |
| PATCH | `/simple_upload/object/*` | Update file tags (`{"Tags":["a","b"]}`) |
| GET | `/simple_upload/objects/*` | List directory contents (includes Tags) |
| POST | `/simple_upload/objects/` | Multi-file upload |
| GET | `/simple_upload/search` | Search (query: Keyword, Tag) |
| POST | `/simple_upload/mkdir/*` | Create directory |
| GET | `/simple_upload/preview/*` | Preview file (PDF/Office conversion) |

## Config files
- Backend config: `config.json` or `config.yaml` (same schema). Example in `backend/rootfs/etc/`.
- Frontend config: `vite.config.js`, `eslint.config.js`.

## Testing
- No frontend test framework. Backend only has `vfs/mountfs_test.go`.
- Run `go test ./...` from `backend/`.

## Gotchas
- Backend logs to stdout; level set in config.
- MongoDB config is commented out by default; enable in config.
- Frontend proxy target is `http://127.0.0.1:8088` – ensure backend runs on that port.
- `AutoMigrateLevel` must be `auto` or `must`; otherwise migrations are skipped.
- RealFilenamePolicy controls file storage naming (`origin`, `_`, `uuid`). Default is `uuid`.
- Database driver defaults to SQLite (`DriverType: sqlite`). Config example shows postgres but sqlite is default.
- SimpleUploadRoot defines where uploaded files are stored relative to VFS root (default `/aaa`).
- Document preview requires `libreoffice` installed on the server. Conversion can take 2-5 seconds for large files.
- Tags are stored as a JSON-serialized string array in the DB (not a separate table).
- Directory entries are stored as DB records with `IsDir=true`; deleting a directory cascades to all files under it.
- The `file` model auto-migrates — adding new fields requires updating `ToEntity()`/`FromEntity()` and the GORM struct tags.

## Work flow
- thinking
- coding
- commit
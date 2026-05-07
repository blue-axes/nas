# AGENTS.md

## API Documentation
- Complete API docs in `docs/` directory
- Main docs: `docs/README.md`
- OpenAPI spec: `docs/openapi.yaml`
- Examples: `docs/examples.md`

## Repository structure
- `frontend/` – React + Vite SPA (Ant Design, React Router hash router)
- `backend/` – Go Echo HTTP server (GORM, VFS, SQLite/Postgres/Mongo)

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
- VFS root is OS `/` (see `backend/vfs.go`). Be careful with file operations.
- Frontend uses hash router (`/#/image` etc.).

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

## work flow
- thinking
- coding
- commit
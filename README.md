# NAS — 轻量个人云存储

基于 **Go + React** 的 NAS 系统，支持文件管理、在线预览、标签检索、SMB/NFS 网络挂载。

## 功能特性

### 文件管理
- **图片管理** (`/#/image`) — 缩略图网格、懒加载、无限滚动
- **视频管理** (`/#/video`) — 缩略图网格、内置播放器（支持倍速、快捷键）
- **通用文件** (`/#/file`) — 表格视图、支持任意文件类型
- **文件夹组织** — 创建文件夹、浏览子目录、面包屑导航

### 文件检索
- **标签系统** — 为文件添加/删除标签，标签显示在卡片叠加层
- **搜索** — 按文件名关键字 + 标签组合筛选文件

### 在线预览
- **PDF** — 浏览器原生渲染
- **Office 文档** — Word (.doc/.docx)、Excel (.xls/.xlsx)、PowerPoint (.ppt/.pptx) 通过 LibreOffice 转为 PDF 后在线预览
- **文本** — CSV、RTF、TXT 等也支持预览
- 转换结果缓存至 `/tmp/nas_preview_cache/`，避免重复转换

### 网络挂载
- **SMB (Samba)** — `sudo bash scripts/setup_smb.sh` 一键配置共享
- **NFS** — `sudo bash scripts/setup_nfs.sh` 一键配置导出
- 支持 Linux / macOS / Windows 挂载访问

### 界面设计
- 深色科技风主题（霓虹青 `#00d4ff` 主色调、毛玻璃卡片、发光边框）
- 响应式布局，适配安卓手机 (768px / 480px 断点)
- 汉堡菜单、自适应网格列数

## 快速开始

```bash
# 前端
cd frontend
npm install
npm run build      # 输出到 dist/

# 后端
cd backend
go build
cp -r ../frontend/dist ./static    # 将前端放入静态目录
./tmpl -config ./config.json
```

## 技术栈

| 层 | 技术 |
|---|---|
| 前端 | React 18, Vite 6, Ant Design 5, React Router 7 |
| 后端 | Go, Echo, GORM (SQLite/Postgres) |
| 存储 | VFS 抽象层（OS 文件系统） |
| 文档转换 | LibreOffice headless |
| 网络共享 | Samba / NFS |

## 项目结构

```
nas/
├── frontend/             # React SPA
│   └── src/
│       ├── apis/         # API 请求封装
│       ├── components/   # 可复用组件 (SearchBar, TagEditor, DocPreviewer, PathTravel)
│       ├── hooks/        # 自定义 Hook (useScreenWidth)
│       ├── pages/        # 页面 (ImageManage, VideoManage, FileManage)
│       └── reducers/     # Immer 状态管理
├── backend/
│   ├── main.go
│   ├── http/api/simple_upload/   # 文件上传/下载/搜索/预览 API
│   ├── service/                  # 业务逻辑层
│   ├── store/rdb/                # GORM 数据持久化
│   └── types/                    # 领域类型定义
├── scripts/              # SMB/NFS 配置脚本
└── docs/                 # API 文档
```

## API 概览

| 方法 | 端点 | 说明 |
|---|---|---|
| `GET` | `/simple_upload/objects/*` | 列出目录 |
| `GET` | `/simple_upload/object/*` | 下载/显示文件 |
| `POST` | `/simple_upload/object/*` | 上传文件 |
| `DELETE` | `/simple_upload/object/*` | 删除文件/目录 |
| `PATCH` | `/simple_upload/object/*` | 更新文件标签 |
| `POST` | `/simple_upload/mkdir/*` | 创建目录 |
| `GET` | `/simple_upload/search` | 搜索文件 (Keyword + Tag) |
| `GET` | `/simple_upload/preview/*` | 文档预览（PDF/Office） |
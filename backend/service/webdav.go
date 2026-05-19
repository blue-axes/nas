package service

import (
	"context"
	"os"
	"path"
	"strings"

	pkgCtx "github.com/blue-axes/tmpl/pkg/context"
	"github.com/blue-axes/tmpl/types"
	"github.com/blue-axes/tmpl/vfs"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/webdav"
)

func NewWebDavHandler(prefix string, fs webdav.FileSystem) *webdav.Handler {
	h := &webdav.Handler{
		Prefix:     prefix,
		FileSystem: fs,
		LockSystem: webdav.NewMemLS(),
	}
	return h
}

type (
	webDavFileSystem struct {
		fs  vfs.VFS
		svc *Service
	}
	trackedFile struct {
		vfs.File
		fs      *webDavFileSystem
		name    string
		isWrite bool
		written int64
	}
)

func NewWebDevFileSystem(fs vfs.VFS, svc *Service) webdav.FileSystem {
	res := &webDavFileSystem{
		fs:  fs,
		svc: svc,
	}
	return res
}

func (fs *webDavFileSystem) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	err := fs.fs.Mkdir(name, perm)
	if err != nil {
		return err
	}
	dbCtx := pkgCtx.New(pkgCtx.WithCtx(ctx))
	normalized := normalizeName(name)
	realDirPath := path.Join(fs.svc.cfg.Nas.SimpleUploadRoot, normalized)
	if err := fs.svc.store.RDB().CreateDir(dbCtx, normalized, realDirPath); err != nil {
		log.WithError(err).Warnf("webdav mkdir: failed to create db record for %s", normalized)
	}
	return nil
}

func (fs *webDavFileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	f, err := fs.fs.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	isWrite := flag&(os.O_CREATE|os.O_WRONLY|os.O_RDWR) != 0
	if !isWrite {
		return f, nil
	}
	return &trackedFile{
		File:    f,
		fs:      fs,
		name:    name,
		isWrite: true,
	}, nil
}

func (f *trackedFile) Write(p []byte) (n int, err error) {
	n, err = f.File.Write(p)
	f.written += int64(n)
	return
}

func (f *trackedFile) Close() error {
	err := f.File.Close()
	if err != nil {
		log.WithError(err).Warnf("webdav close file error: %s", err)
	}
	if f.isWrite && f.written > 0 {
		f.syncToDB()
	}
	return err
}

func (f *trackedFile) syncToDB() {
	info, err := f.fs.fs.Stat(f.name)
	if err != nil {
		log.WithError(err).Warnf("webdav stat after write failed: %s", f.name)
		return
	}
	normalized := normalizeName(f.name)
	realPath := path.Join(f.fs.svc.cfg.Nas.SimpleUploadRoot, normalized)
	fileRecord := &types.File{
		Name:  normalized,
		Ext:   path.Ext(normalized),
		Path:  realPath,
		Size:  uint64(info.Size()),
		IsDir: info.IsDir(),
	}
	dbCtx := pkgCtx.New()
	if err := f.fs.svc.store.RDB().UpsertFileByName(dbCtx, normalized, fileRecord); err != nil {
		log.WithError(err).Warnf("webdav upsert db record failed: %s", normalized)
	}
}

func (fs *webDavFileSystem) RemoveAll(ctx context.Context, name string) error {
	info, statErr := fs.fs.Stat(name)
	err := fs.fs.RemoveAll(name)
	if err != nil {
		return err
	}
	dbCtx := pkgCtx.New(pkgCtx.WithCtx(ctx))
	normalized := normalizeName(name)
	if statErr == nil && info.IsDir() {
		if err := fs.svc.store.RDB().DeleteDir(dbCtx, normalized); err != nil {
			log.WithError(err).Warnf("webdav delete dir db record failed: %s", normalized)
		}
	} else {
		if err := fs.svc.store.RDB().DeleteByName(dbCtx, normalized); err != nil {
			log.WithError(err).Warnf("webdav delete file db record failed: %s", normalized)
		}
	}
	return nil
}

func (fs *webDavFileSystem) Rename(ctx context.Context, oldName, newName string) error {
	info, statErr := fs.fs.Stat(oldName)
	err := fs.fs.Rename(oldName, newName)
	if err != nil {
		return err
	}
	dbCtx := pkgCtx.New(pkgCtx.WithCtx(ctx))
	oldNorm := normalizeName(oldName)
	newNorm := normalizeName(newName)
	if statErr == nil && info.IsDir() {
		if err := fs.svc.store.RDB().DeleteDir(dbCtx, oldNorm); err != nil {
			log.WithError(err).Warnf("webdav rename: delete old dir db record failed: %s", oldNorm)
		}
		fs.syncDirToDB(dbCtx, newNorm)
	} else {
		if err := fs.svc.store.RDB().RenameFileByName(dbCtx, oldNorm, newNorm); err != nil {
			log.WithError(err).Warnf("webdav rename file db record failed: %s -> %s", oldNorm, newNorm)
		}
	}
	return nil
}

func (fs *webDavFileSystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	return fs.fs.Stat(name)
}

func normalizeName(name string) string {
	return strings.TrimLeft(name, "/")
}

func (fs *webDavFileSystem) syncDirToDB(ctx *pkgCtx.Context, dirPath string) {
	dirPath = normalizeName(strings.TrimRight(dirPath, "/"))
	if dirPath == "" {
		dirPath = "."
	}
	basePath := path.Join(fs.svc.cfg.Nas.SimpleUploadRoot, dirPath)
	if err := fs.svc.store.RDB().CreateDir(ctx, dirPath, basePath); err != nil {
		log.WithError(err).Warnf("webdav sync dir: create dir db record failed: %s", dirPath)
	}
	entries, err := fs.svc.vfs.ReadDir(basePath)
	if err != nil {
		log.WithError(err).Warnf("webdav sync dir: read dir failed: %s", basePath)
		return
	}
	for _, entry := range entries {
		childName := dirPath + "/" + entry.Name()
		if entry.IsDir() {
			fs.syncDirToDB(ctx, childName)
		} else {
			info, _ := entry.Info()
			realPath := path.Join(fs.svc.cfg.Nas.SimpleUploadRoot, childName)
			fileRecord := &types.File{
				Name:  childName,
				Ext:   path.Ext(entry.Name()),
				Path:  realPath,
				Size:  uint64(info.Size()),
				IsDir: false,
			}
			if err := fs.svc.store.RDB().UpsertFileByName(ctx, childName, fileRecord); err != nil {
				log.WithError(err).Warnf("webdav sync dir: upsert file db record failed: %s", childName)
			}
		}
	}
}

func (svc *Service) GetWebDavHandler(prefix string) (*webdav.Handler, error) {
	fs, err := svc.vfs.Chroot(svc.cfg.Nas.SimpleUploadRoot)
	if err != nil {
		return nil, err
	}
	return NewWebDavHandler(prefix, NewWebDevFileSystem(fs, svc)), nil
}

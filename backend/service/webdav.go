package service

import (
	"context"
	"os"

	"github.com/blue-axes/tmpl/vfs"
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
		fs vfs.VFS
	}
	webDavFile struct {
		file vfs.File
	}
)

func NewWebDevFileSystem(fs vfs.VFS) webdav.FileSystem {
	res := &webDavFileSystem{
		fs: fs,
	}
	return res
}

func (fs *webDavFileSystem) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	return fs.fs.Mkdir(name, perm)
}

func (fs *webDavFileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	return fs.fs.OpenFile(name, flag, perm)
}

func (fs *webDavFileSystem) RemoveAll(ctx context.Context, name string) error {
	return fs.fs.RemoveAll(name)
}

func (fs *webDavFileSystem) Rename(ctx context.Context, oldName, newName string) error {
	return fs.fs.Rename(oldName, newName)
}

func (fs *webDavFileSystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	return fs.fs.Stat(name)
}

func (svc *Service) GetWebDavHandler(prefix string) (*webdav.Handler, error) {
	fs, err := svc.vfs.Chroot(svc.cfg.Nas.SimpleUploadRoot)
	if err != nil {
		return nil, err
	}
	return NewWebDavHandler(prefix, NewWebDevFileSystem(fs)), nil
}

package vfs

import (
	iofs "io/fs"
	"os"
	"path"
)

type (
	osFs struct {
		rootDir string
	}
)

func NewOsFS(rootDir string) *osFs {
	rootDir = path.Clean(rootDir)
	return &osFs{
		rootDir: rootDir,
	}
}

func (o *osFs) getPath(name string) string {
	return path.Join(o.rootDir, name)
}

func (o *osFs) Stat(name string) (iofs.FileInfo, error) {
	return os.Stat(o.getPath(name))
}

func (o *osFs) Remove(name string) error {
	return os.Remove(o.getPath(name))
}

func (o *osFs) RemoveAll(path string) error {
	return os.RemoveAll(o.getPath(path))
}

func (o *osFs) OpenFile(name string, flag int, perm iofs.FileMode) (File, error) {
	return os.OpenFile(o.getPath(name), flag, perm)
}

func (o *osFs) Mkdir(name string, perm iofs.FileMode) error {
	return os.Mkdir(o.getPath(name), perm)
}

func (o *osFs) MkdirAll(path string, perm iofs.FileMode) error {
	return os.MkdirAll(o.getPath(path), perm)
}

func (o *osFs) ReadDir(name string) ([]iofs.DirEntry, error) {
	return os.ReadDir(o.getPath(name))
}

func (o *osFs) TempDir() string {
	if o.rootDir == "/" || o.rootDir == "" {
		return os.TempDir()
	}
	return path.Join(o.rootDir, "/temp")
}

package vfs

import (
	"io"
	"io/fs"
)

type (
	VFS interface {
		Stat(name string) (fs.FileInfo, error)
		Remove(name string) error
		RemoveAll(path string) error

		OpenFile(name string, flag int, perm fs.FileMode) (File, error)
		Mkdir(name string, perm fs.FileMode) error
		MkdirAll(path string, perm fs.FileMode) error
		ReadDir(name string) ([]fs.DirEntry, error)
		TempDir() string
	}
	File interface {
		io.Reader
		io.Writer
		io.Closer
		io.Seeker
		Truncate(size int64) error
		Name() string
	}
	MountFs interface {
		Mount(dir string, fs VFS) error
		Umount(dir string) error
		VFS
	}
)

type (
	VFSType string
	VFSConf interface {
		TypeName() VFSType
	}
)

const (
	VFS_OS VFSType = "os"
)

package vfs

import (
	"errors"
	"io"
	"io/fs"
	iofs "io/fs"
	"os"
	"path"
	"strings"
	"sync"
)

type (
	mountFs struct {
		sync.RWMutex
		rootfs VFS
		mounts map[string]VFS
	}
)

var (
	ErrAbsolutePathOnly   = errors.New("only support absolutely path")
	ErrDirIsMounted       = errors.New("this dir is mounted")
	ErrDirNotMounted      = errors.New("this dir not mount any fs")
	ErrRootfsCannotUmount = errors.New("rootfs cannot umount")
	ErrNilVFS             = errors.New("nil vfs")
)

func NewMountFs(rootfs VFS) MountFs {
	if rootfs == nil {
		panic(ErrNilVFS.Error())
	}
	return &mountFs{
		rootfs: rootfs,
		mounts: make(map[string]VFS),
	}
}

func (m *mountFs) Mount(dir string, fs VFS) error {
	if fs == nil {
		return ErrNilVFS
	}
	// dir只能是绝对路径
	dir = path.Clean(dir)
	if !path.IsAbs(dir) {
		return ErrAbsolutePathOnly
	}
	if dir == "/" {
		// replace rootfs
		m.rootfs = fs
		return nil
	}

	m.Lock()
	defer m.Unlock()
	_, ok := m.mounts[dir]
	if ok {
		return ErrDirIsMounted
	}
	m.mounts[dir] = fs
	return nil
}

func (m *mountFs) Umount(dir string) error {
	// dir只能是绝对路径
	dir = path.Clean(dir)
	if !path.IsAbs(dir) {
		return ErrAbsolutePathOnly
	}
	if dir == "/" {
		// unmount rootfs
		return ErrRootfsCannotUmount
	}
	m.Lock()
	defer m.Unlock()
	_, ok := m.mounts[dir]
	if !ok {
		return ErrDirNotMounted
	}
	delete(m.mounts, dir)
	return nil
}

func (m *mountFs) findFs(dir string) (VFS, string) {
	dir = path.Clean(dir)
	if !path.IsAbs(dir) {
		return m.rootfs, path.Join("/", dir)
	}
	if dir == "/" {
		return m.rootfs, "/"
	}
	// 一级一级匹配，找到匹配最长的那个
	m.RLock()
	defer m.RUnlock()
	var currentDir = dir
	for {
		currentDir = path.Dir(currentDir)
		fs, ok := m.mounts[currentDir]
		if ok {
			// 去掉前缀
			return fs, strings.TrimPrefix(dir, currentDir)
		}
		if currentDir == "/" {
			break
		}
	}
	return m.rootfs, dir
}

func (m *mountFs) Stat(name string) (iofs.FileInfo, error) {
	fs, dir := m.findFs(name)
	return fs.Stat(dir)
}

func (m *mountFs) Remove(name string) error {
	fs, dir := m.findFs(name)
	return fs.Remove(dir)
}

func (m *mountFs) RemoveAll(path string) error {
	fs, dir := m.findFs(path)
	return fs.RemoveAll(dir)
}

func (m *mountFs) Rename(oldName, newName string) error {
	oldFs, oldDir := m.findFs(oldName)
	newFs, newDir := m.findFs(newName)
	if oldFs != newFs { // 跨文件系统，则先复制再删除
		oldStat, _ := oldFs.Stat(oldDir)
		perm := fs.FileMode(0600)
		if oldStat != nil {
			perm = oldStat.Mode()
		}
		d, err := newFs.OpenFile(newDir, os.O_CREATE|os.O_RDWR|os.O_TRUNC, perm)
		if err != nil {
			return err
		}
		defer d.Close()
		s, err := oldFs.OpenFile(oldDir, os.O_RDONLY, 0600)
		if err != nil {
			return err
		}
		defer s.Close()
		_, err = io.Copy(d, s)
		return err
	}
	// 不跨文件系统
	return oldFs.Rename(oldDir, newDir)
}

func (m *mountFs) OpenFile(name string, flag int, perm iofs.FileMode) (File, error) {
	fs, filename := m.findFs(name)
	return fs.OpenFile(filename, flag, perm)
}

func (m *mountFs) Mkdir(dir string, perm iofs.FileMode) error {
	fs, dirPath := m.findFs(dir)
	return fs.Mkdir(dirPath, perm)
}

func (m *mountFs) MkdirAll(path string, perm iofs.FileMode) error {
	fs, dirPath := m.findFs(path)
	return fs.MkdirAll(dirPath, perm)
}

func (m *mountFs) ReadDir(dir string) ([]iofs.DirEntry, error) {
	fs, dirPath := m.findFs(dir)
	return fs.ReadDir(dirPath)
}

func (m *mountFs) TempDir() string {
	return m.rootfs.TempDir()
}

func (m *mountFs) Chroot(dir string) (MountFs, error) {
	dir = path.Clean(dir)
	if !path.IsAbs(dir) {
		return nil, ErrAbsolutePathOnly
	}
	info, err := m.Stat(dir)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, errors.New("not a directory")
	}
	return &chrootFs{parent: m, rootDir: dir}, nil
}

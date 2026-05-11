package vfs

import (
	iofs "io/fs"
	"path"
)

// chrootFs 将所有路径操作限制在 rootDir 下

type chrootFs struct {
	parent  MountFs
	rootDir string
}

func (c *chrootFs) resolve(name string) string {
	cleaned := path.Clean("/" + name)
	return path.Join(c.rootDir, cleaned)
}

func (c *chrootFs) Stat(name string) (iofs.FileInfo, error) {
	return c.parent.Stat(c.resolve(name))
}

func (c *chrootFs) Remove(name string) error {
	return c.parent.Remove(c.resolve(name))
}

func (c *chrootFs) RemoveAll(name string) error {
	return c.parent.RemoveAll(c.resolve(name))
}

func (c *chrootFs) Rename(oldName, newName string) error {
	return c.parent.Rename(c.resolve(oldName), c.resolve(newName))
}

func (c *chrootFs) OpenFile(name string, flag int, perm iofs.FileMode) (File, error) {
	return c.parent.OpenFile(c.resolve(name), flag, perm)
}

func (c *chrootFs) Mkdir(name string, perm iofs.FileMode) error {
	return c.parent.Mkdir(c.resolve(name), perm)
}

func (c *chrootFs) MkdirAll(name string, perm iofs.FileMode) error {
	return c.parent.MkdirAll(c.resolve(name), perm)
}

func (c *chrootFs) ReadDir(name string) ([]iofs.DirEntry, error) {
	return c.parent.ReadDir(c.resolve(name))
}

func (c *chrootFs) TempDir() string {
	return path.Join(c.rootDir, "/temp")
}

func (c *chrootFs) Mount(dir string, fs VFS) error {
	return c.parent.Mount(c.resolve(dir), fs)
}

func (c *chrootFs) Umount(dir string) error {
	return c.parent.Umount(c.resolve(dir))
}

func (c *chrootFs) Chroot(dir string) (MountFs, error) {
	return c.parent.Chroot(c.resolve(dir))
}

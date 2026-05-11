package vfs

import (
	"testing"
)

func TestMountFs(t *testing.T) {
	fs := NewMountFs(NewOsFS(OsFsConf{RootDir: "/"}))
	_ = fs.Mount("/mnt", NewOsFS(OsFsConf{RootDir: "/"}))

	err := fs.MkdirAll("/a/b", 0666)
	if err != nil {
		t.Error(err)
		return
	}
	err = fs.MkdirAll("/mnt/b/c", 0666)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestChroot(t *testing.T) {
	root := NewMountFs(NewOsFS(OsFsConf{RootDir: "/"}))

	// 先创建一个测试目录
	err := root.MkdirAll("/tmp/chroot_test/subdir", 0755)
	if err != nil {
		t.Fatal("MkdirAll:", err)
	}

	// chroot 到 /tmp/chroot_test
	chrooted, err := root.Chroot("/tmp/chroot_test")
	if err != nil {
		t.Fatal("Chroot:", err)
	}

	// chroot 内 Stat "/" 应该对应父 fs 的 /tmp/chroot_test
	info, err := chrooted.Stat("/")
	if err != nil {
		t.Fatal("Stat /:", err)
	}
	if !info.IsDir() {
		t.Error("expected root to be a directory")
	}

	// chroot 内能访问子目录
	_, err = chrooted.Stat("/subdir")
	if err != nil {
		t.Fatal("Stat /subdir:", err)
	}

	// chroot 内能 ReadDir
	entries, err := chrooted.ReadDir("/")
	if err != nil {
		t.Fatal("ReadDir:", err)
	}
	if len(entries) != 1 || entries[0].Name() != "subdir" {
		t.Errorf("expected [subdir], got %v", entries)
	}

	// 递归 chroot
	chrooted2, err := chrooted.Chroot("/subdir")
	if err != nil {
		t.Fatal("Chroot nested:", err)
	}
	info2, err := chrooted2.Stat("/")
	if err != nil {
		t.Fatal("Stat nested /:", err)
	}
	if !info2.IsDir() {
		t.Error("expected nested root to be a directory")
	}

	// Chroot 到不存在的目录应该报错
	_, err = root.Chroot("/nonexistent_dir_xyz")
	if err == nil {
		t.Error("expected error for nonexistent dir")
	}

	// Chroot 到非绝对路径应该报错
	_, err = root.Chroot("relative/path")
	if err == nil {
		t.Error("expected error for relative path")
	}
}

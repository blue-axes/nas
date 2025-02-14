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

package main

import "github.com/blue-axes/tmpl/vfs"

func NewMountFs() vfs.MountFs {
	// 暂时内置
	rootFs := vfs.NewMountFs(vfs.NewOsFS(vfs.OsFsConf{RootDir: "/"}))
	return rootFs
}

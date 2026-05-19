package service

import (
	"path"

	"github.com/blue-axes/tmpl/pkg/context"
	"github.com/blue-axes/tmpl/store"
	log "github.com/sirupsen/logrus"
)

var defaultDirs = []string{"img", "video", "other"}

func (svc *Service) InitDefaultDirs() error {
	ctx := context.New()
	root := svc.cfg.Nas.SimpleUploadRoot

	for _, dir := range defaultDirs {
		if _, err := svc.store.RDB().GetFileByName(ctx, dir); err == nil {
			continue
		}

		realPath := path.Join(root, dir)
		if err := svc.vfs.MkdirAll(realPath, 0700); err != nil {
			log.WithError(err).Warnf("init default dir: mkdir %s failed", realPath)
		}

		if err := svc.store.RDB().CreateDir(ctx, dir, realPath); err != nil {
			log.WithError(err).Warnf("init default dir: db insert %s failed", dir)
		} else {
			log.Infof("init default dir: created %s", dir)
		}
	}
	return nil
}

func (svc *Service) GetStore() *store.Store {
	return svc.store
}
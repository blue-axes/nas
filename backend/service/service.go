package service

import (
	"github.com/hashicorp/mdns"
	"github.com/blue-axes/tmpl/store"
	"github.com/blue-axes/tmpl/types"
	"github.com/blue-axes/tmpl/vfs"
)

type (
	Service struct {
		cfg        *types.Config
		store      *store.Store
		vfs        vfs.MountFs
		session    *SessionStore
		mdnsServer *mdns.Server
	}
	Option func(svc *Service) error
)

func WithConfig(cfg *types.Config) Option {
	return func(svc *Service) error {
		svc.cfg = cfg
		return nil
	}
}

func WithMountFs(mountFs vfs.MountFs) Option {
	return func(svc *Service) error {
		svc.vfs = mountFs
		return nil
	}
}

func New(store *store.Store, options ...Option) (*Service, error) {
	svc := &Service{
		store: store,
	}
	for _, opt := range options {
		err := opt(svc)
		if err != nil {
			return nil, err
		}
	}
	return svc, nil
}

func (svc *Service) Config() *types.Config {
	if svc.cfg == nil {
		panic("initial service WithConfig Option nil")
	}
	return svc.cfg
}

func (svc *Service) Session() *SessionStore {
	return svc.session
}

func (svc *Service) InitSession() {
	cfg := svc.Config()
	if cfg.Http.Auth.Enabled {
		svc.session = NewSessionStore(cfg.Http.Auth.SessionExpireHours)
	}
}

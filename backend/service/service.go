package service

import (
	"github.com/blue-axes/tmpl/store"
	"github.com/blue-axes/tmpl/vfs"
)

type (
	Service struct {
		store *store.Store
		vfs   vfs.MountFs
	}
	Option func(svc *Service) error
)

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

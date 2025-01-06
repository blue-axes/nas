package store

import (
	"github.com/blue-axes/tmpl/store/mongo"
	"github.com/blue-axes/tmpl/store/postgres"
	"github.com/blue-axes/tmpl/types"
)

type (
	Config = types.DatabaseConfig
	Store  struct {
		postgresStore *postgres.Store
		mongoStore    *mongo.Store
	}
)

func New(cfg Config) (*Store, error) {
	var (
		pgStore *postgres.Store
		mgStore *mongo.Store
		err     error
	)
	if cfg.Postgres != nil {
		pgStore, err = postgres.New(*cfg.Postgres)
		if err != nil {
			return nil, err
		}
	}
	if cfg.Mongo != nil {
		mgStore, err = mongo.New(*cfg.Mongo)
		if err != nil {
			return nil, err
		}
	}

	s := &Store{
		postgresStore: pgStore,
		mongoStore:    mgStore,
	}

	return s, nil
}

func (s *Store) Postgres() *postgres.Store {
	return s.postgresStore
}

func (s *Store) Mongo() *mongo.Store {
	return s.mongoStore
}

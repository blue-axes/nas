package store

import (
	"github.com/blue-axes/tmpl/store/mongo"
	"github.com/blue-axes/tmpl/store/rdb"
	"github.com/blue-axes/tmpl/types"
)

type (
	Config = types.DatabaseConfig
	Store  struct {
		rdb        *rdb.Store
		mongoStore *mongo.Store
	}
)

func New(cfg Config) (*Store, error) {
	var (
		rdbStore *rdb.Store
		mgStore  *mongo.Store
		err      error
	)
	if cfg.Rdb != nil {
		rdbStore, err = rdb.New(*cfg.Rdb)
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
		rdb:        rdbStore,
		mongoStore: mgStore,
	}

	return s, nil
}

func (s *Store) RDB() *rdb.Store {
	return s.rdb
}

func (s *Store) Mongo() *mongo.Store {
	return s.mongoStore
}

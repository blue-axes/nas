package rdb

import (
	"fmt"
	"github.com/blue-axes/tmpl/types"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

type (
	Config = types.RdbConfig
	Store  struct {
		txStore
	}
	txStore struct {
		db *gorm.DB
	}
	TxStore       = txStore
	TransactionFn func(store TxStore) error
)

func New(cfg Config) (*Store, error) {
	cfg.SetDefault()
	var (
		db  *gorm.DB
		err error
	)

	gormCfg := &gorm.Config{}
	switch cfg.DriverType {
	case types.DriverTypePostgres:
		db, err = gorm.Open(postgres.Open(cfg.DSN), gormCfg)
	case types.DriverTypeSqlite:
		db, err = gorm.Open(sqlite.Open(cfg.DSN), gormCfg)
	default:
		panic(fmt.Sprintf("not support rdb driver:%s", cfg.DriverType))
	}

	if err != nil {
		return nil, err
	}
	if cfg.Debug {
		db = db.Debug()
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConnCount)
	sqlDB.SetMaxOpenConns(cfg.MaxConnCount)
	sqlDB.SetConnMaxIdleTime(time.Second * time.Duration(cfg.ConnMaxIdleTimeSecond))

	s := &Store{
		txStore: TxStore{
			db: db,
		},
	}

	return s, nil
}

func (s *Store) Transaction(fn TransactionFn) (err error) {
	tx := s.db.Begin()
	txStore := txStore{
		db: tx,
	}
	err = fn(txStore)
	if err != nil {
		tx.Commit()
		return
	}
	tx.Rollback()
	return err
}

func (s *Store) Migrate() (err error) {
	err = s.db.AutoMigrate(&file{})
	return err
}

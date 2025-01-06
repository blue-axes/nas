package postgres

import (
	"fmt"
	"github.com/blue-axes/tmpl/types"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/url"
	"time"
)

type (
	Config = types.PostgresConfig
	Store  struct {
		txStore
	}
	txStore struct {
		db *gorm.DB
	}
	TxStore       = txStore
	TransactionFn func(store TxStore) error
)

func getDsn(cfg Config) string {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.Username, cfg.Password, cfg.Address, cfg.Port, cfg.Database)
	dsnUrl, _ := url.Parse(dsn)
	dsnUrl.Query().Set("sslmode", "disable")
	dsnUrl.Query().Set("TimeZone", "Asia/Shanghai")
	return dsnUrl.String()
}

func New(cfg Config) (*Store, error) {
	cfg.SetDefault()

	db, err := gorm.Open(postgres.Open(getDsn(cfg)))
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
	err = s.db.AutoMigrate(&example{})

	return err
}

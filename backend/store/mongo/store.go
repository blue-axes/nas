package mongo

import (
	"context"
	"fmt"
	"github.com/blue-axes/tmpl/pkg/log"
	"github.com/blue-axes/tmpl/types"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/url"
	"time"
)

type (
	Config = types.MongoConfig
	Store  struct {
		txStore
	}
	txStore struct {
		db *mongo.Client
	}
	TxStore       = txStore
	TransactionFn func(store TxStore) error

	Model interface {
		TableName() string
		DatabaseName() string
	}
)

func getDsn(cfg Config) string {
	dsn := fmt.Sprintf("mongodb://%s:%d/%s", cfg.Address, cfg.Port, cfg.Database)
	if cfg.Username != "" {
		dsn = fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", cfg.Username, cfg.Password, cfg.Address, cfg.Port, cfg.Database)
	}
	dsnUrl, _ := url.Parse(dsn)
	dsnUrl.Query().Set("ReadEncoding", "UTF8Encoding")
	dsnUrl.Query().Set("WriteEncoding", "UTF8Encoding")
	//dsnUrl.Query().Set("TimeZone", "Asia/Shanghai")
	return dsnUrl.String()
}

func New(cfg Config) (*Store, error) {
	cfg.SetDefault()

	connOption := options.Client().ApplyURI(getDsn(cfg))
	connOption.SetMaxConnecting(uint64(cfg.MaxConnCount))
	connOption.SetMaxPoolSize(uint64(cfg.MaxConnCount + cfg.MaxIdleConnCount))
	connOption.SetMaxConnIdleTime(time.Second * time.Duration(cfg.ConnMaxIdleTimeSecond))
	connOption.SetConnectTimeout(time.Second * 60)
	if cfg.Debug {
		connOption.Monitor = &event.CommandMonitor{
			Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
				log.Infof("mongo sql:%s", startedEvent.Command.String())
			},
		}
	}

	cli, err := mongo.Connect(context.Background(), connOption)
	if err != nil {
		return nil, err
	}

	s := &Store{
		txStore: TxStore{
			db: cli,
		},
	}

	return s, nil
}

func (s *Store) Transaction(fn TransactionFn) (err error) {
	sess, err := s.db.StartSession(options.Session())
	if err != nil {
		return err
	}
	txStore := txStore{
		db: sess.Client(),
	}
	ctx := context.Background()
	defer sess.EndSession(ctx)

	err = fn(txStore)
	if err != nil {
		return sess.CommitTransaction(ctx)
	}
	return sess.AbortTransaction(ctx)
}

func (s *Store) Migrate() (err error) {
	s.migrate(&example{})
	return err
}

func (s *Store) migrate(mdl Model) error {
	return s.db.Database(mdl.DatabaseName()).CreateCollection(context.Background(), mdl.TableName())
}

func (s *txStore) collection(mdl Model) *mongo.Collection {
	return s.db.Database(mdl.DatabaseName()).Collection(mdl.TableName())
}

package main

import (
	"flag"
	"github.com/blue-axes/tmpl/http"
	"github.com/blue-axes/tmpl/pkg/config"
	"github.com/blue-axes/tmpl/pkg/log"
	"github.com/blue-axes/tmpl/service"
	"github.com/blue-axes/tmpl/store"
	"github.com/blue-axes/tmpl/types"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	var (
		cfgPath = ""
		cfg     types.Config
	)
	flag.StringVar(&cfgPath, "config", "./config.json", "the config file path")
	flag.Parse()

	// 加载配置文件
	err := config.Load(cfgPath, &cfg)
	if err != nil {
		log.Errorf("%s", err.Error())
		os.Exit(1)
		return
	}
	cfg.SetDefault()

	//设置日志等级
	logLvl, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		log.Error(err)
		os.Exit(1)
		return
	}
	log.SetLevel(logLvl)

	// 初始化store
	sor, err := store.New(cfg.Database)
	if err != nil {
		log.Error(err)
		os.Exit(1)
		return
	}
	// 初始化表
	err = databaseMigrate(sor, cfg.Database.Postgres.AutoMigrateLevel)
	if err != nil {
		log.Errorf("database migrate err:%s", err.Error())
		os.Exit(1)
		return
	}

	// 初始化service
	svc, err := service.New(sor)
	if err != nil {
		log.Error(err)
		os.Exit(1)
		return
	}
	// 初始化server
	srv, err := http.New(cfg.Http, svc)
	if err != nil {
		log.Error(err)
		os.Exit(1)
		return
	}
	go func() {
		err = srv.Start()
		if err != nil {
			log.Error(err)
			return
		}
	}()

	// 系统信号
	var (
		sig = make(chan os.Signal)
	)
	signal.Notify(sig, syscall.SIGINT)
	<-sig
	srv.Shutdown()
	os.Exit(0)
}

func databaseMigrate(sor *store.Store, level string) error {
	switch strings.ToLower(level) {
	case "auto":
		_ = sor.Postgres().Migrate()
	case "must":
		err := sor.Postgres().Migrate()
		if err != nil {
			return err
		}
	default:
	}
	return nil
}

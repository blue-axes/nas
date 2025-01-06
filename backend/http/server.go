package http

import (
	"context"
	"fmt"
	"github.com/blue-axes/tmpl/http/api"
	"github.com/blue-axes/tmpl/pkg/log"
	"github.com/blue-axes/tmpl/service"
	"github.com/blue-axes/tmpl/types"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	Config = types.HttpConfig
	Server struct {
		e   *echo.Echo
		cfg Config
		svc *service.Service
	}
	Option func(s *Server) error
)

func New(cfg Config, svc *service.Service, options ...Option) (*Server, error) {
	s := &Server{
		e:   echo.New(),
		cfg: cfg,
		svc: svc,
	}
	for _, opt := range options {
		err := opt(s)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *Server) Start() error {
	e := s.e
	//e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper:          middleware.DefaultLoggerConfig.Skipper,
		Format:           middleware.DefaultLoggerConfig.Format,
		CustomTimeFormat: middleware.DefaultLoggerConfig.CustomTimeFormat,
		CustomTagFunc:    middleware.DefaultLoggerConfig.CustomTagFunc,
		Output:           log.GetOutput(),
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultCORSConfig.Skipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
	}))
	e.Pre(Pre)
	e.HTTPErrorHandler = api.ErrorHandler
	e.Binder = &Binder{}
	// 初始化路由
	initRouter(s.svc, s.e)

	return e.Start(fmt.Sprintf("%s:%d", s.cfg.ListenAddress, s.cfg.ListenPort))
}

func (s *Server) Shutdown() {
	_ = s.e.Shutdown(context.Background())
}

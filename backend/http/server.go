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
		Format:           middleware.DefaultLoggerConfig.Format,
		CustomTimeFormat: middleware.DefaultLoggerConfig.CustomTimeFormat,
		CustomTagFunc:    middleware.DefaultLoggerConfig.CustomTagFunc,
		Output:           log.GetOutput(),
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
	}))
	if s.cfg.Auth.Enabled {
		e.Use(BasicAuth(func(username, password string) (*types.UserInfo, bool) {
			info, ok := s.svc.ValidateUser(username, password)
			if !ok {
				return nil, false
			}
			return &types.UserInfo{
				Username: info.Username,
				CanRead:  info.CanRead,
				CanWrite: info.CanWrite,
				IsAdmin:  info.IsAdmin,
			}, true
		}))
	}
	e.Pre(Pre)
	e.HTTPErrorHandler = api.ErrorHandler
	e.Binder = &Binder{}
	// 初始化路由
	initRouter(s.svc, s.e)

	addr := fmt.Sprintf("%s:%d", s.cfg.ListenAddress, s.cfg.ListenPort)
	if s.cfg.CertFile != "" && s.cfg.KeyFile != "" {
		return e.StartTLS(addr, s.cfg.CertFile, s.cfg.KeyFile)
	}
	return e.Start(addr)
}

func (s *Server) Shutdown() {
	_ = s.e.Shutdown(context.Background())
}

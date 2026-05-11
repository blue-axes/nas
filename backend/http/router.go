package http

import (
	"github.com/blue-axes/tmpl/http/api/example"
	"github.com/blue-axes/tmpl/http/api/simple_upload"
	"github.com/blue-axes/tmpl/http/api/webdav"
	"github.com/blue-axes/tmpl/service"
	"github.com/labstack/echo/v4"
)

func initRouter(svc *service.Service, e *echo.Echo) {
	example.InitRouter(svc, e.Group("/example"))
	e.Static("/static", svc.Config().Http.StaticRoot)
	simple_upload.InitRouter(svc, e.Group("/simple_upload"))
	webdav.InitRouter(svc, e.Group("/webdav"))
}

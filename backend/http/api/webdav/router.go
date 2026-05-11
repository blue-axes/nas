package webdav

import (
	"github.com/blue-axes/tmpl/service"
	"github.com/labstack/echo/v4"
)

func InitRouter(svc *service.Service, e *echo.Group) {
	handler := New(svc)
	e.Any("", handler.Webdav)
	e.Any("/*", handler.Webdav)
}

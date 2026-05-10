package http

import (
	"github.com/blue-axes/tmpl/http/api/example"
	"github.com/blue-axes/tmpl/http/api/simple_upload"
	"github.com/blue-axes/tmpl/service"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/webdav"
)

func initRouter(svc *service.Service, e *echo.Echo) {
	example.InitRouter(svc, e.Group("/example"))
	e.Static("/static", svc.Config().Http.StaticRoot)
	simple_upload.InitRouter(svc, e.Group("/simple_upload"))

	webdavHandler := &webdav.Handler{
		FileSystem: webdav.Dir(svc.Config().Nas.SimpleUploadRoot),
		LockSystem: webdav.NewMemLS(),
	}
	e.Any("/webdav", func(c echo.Context) error {
		w := c.Response()
		r := c.Request()
		webdavHandler.ServeHTTP(w, r)
		return nil
	})
	e.Any("/webdav/*", func(c echo.Context) error {
		w := c.Response()
		r := c.Request()
		webdavHandler.ServeHTTP(w, r)
		return nil
	})
}

package http

import (
	"github.com/blue-axes/tmpl/http/api"
	"github.com/blue-axes/tmpl/http/api/example"
	"github.com/blue-axes/tmpl/http/api/simple_upload"
	"github.com/blue-axes/tmpl/http/api/user"
	"github.com/blue-axes/tmpl/http/api/webdav"
	"github.com/blue-axes/tmpl/service"
	"github.com/labstack/echo/v4"
)

func initRouter(svc *service.Service, e *echo.Echo) {
	example.InitRouter(svc, e.Group("/example"))
	su := e.Group("/simple_upload", FilePermissionCheck)
	simple_upload.InitRouter(svc, su)

	wd := e.Group("/webdav", FilePermissionCheck)
	webdav.InitRouter(svc, wd)

	userApiPrefix := "/api/users"
	ug := e.Group(userApiPrefix, RequireAdmin([]string{
		userApiPrefix + "/me",
		userApiPrefix + "/login",
		userApiPrefix + "/logout",
	}))
	user.InitRouter(svc, ug)

	staticRoot := svc.Config().Http.StaticRoot
	e.Static("/", staticRoot)

	scanHandler := api.NewFileScannerHandler(svc)
	e.POST("/api/scanfs", scanHandler.Scan, RequireAdmin([]string{}))
}

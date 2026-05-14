package http

import (
	"github.com/blue-axes/tmpl/http/api/example"
	"github.com/blue-axes/tmpl/http/api/simple_upload"
	"github.com/blue-axes/tmpl/http/api/user"
	"github.com/blue-axes/tmpl/http/api/webdav"
	"github.com/blue-axes/tmpl/service"
	"github.com/labstack/echo/v4"
)

func initRouter(svc *service.Service, e *echo.Echo) {
	example.InitRouter(svc, e.Group("/example"))
	e.Static("/static", svc.Config().Http.StaticRoot)

	su := e.Group("/simple_upload")
	su.Use(RequireRead)
	simple_upload.InitRouter(svc, su, RequireWrite)

	webdav.InitRouter(svc, e.Group("/webdav"))

	userHandler := user.New(svc)

	ug := e.Group("/api/users")
	ug.Use(RequireAdmin)
	user.InitRouter(svc, ug, userHandler)

	e.GET("/api/users/me", userHandler.Me)

	pwGroup := e.Group("/api/users")
	pwGroup.PUT("/:username/password", userHandler.ChangePassword)
}
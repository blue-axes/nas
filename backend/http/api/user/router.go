package user

import (
	"github.com/blue-axes/tmpl/service"
	"github.com/labstack/echo/v4"
)

func InitRouter(svc *service.Service, e *echo.Group) {
	handler := New(svc)
	e.POST("/login", handler.Login)
	e.POST("/logout", handler.Logout)
	e.GET("/me", handler.Me)

	e.GET("", handler.List)
	e.POST("", handler.Create)
	e.PUT("/:username", handler.Update)
	e.DELETE("/:username", handler.Delete)
	e.PUT("/:username/password", handler.ChangePassword)
}

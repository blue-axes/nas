package user

import (
	"github.com/blue-axes/tmpl/service"
	"github.com/labstack/echo/v4"
)

func InitRouter(svc *service.Service, e *echo.Group, handler *UserHandler) {
	if handler == nil {
		handler = New(svc)
	}

	e.GET("", handler.List)
	e.POST("", handler.Create)
	e.PUT("/:username", handler.Update)
	e.DELETE("/:username", handler.Delete)
}
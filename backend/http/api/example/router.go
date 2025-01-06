package example

import (
	"github.com/blue-axes/tmpl/service"
	"github.com/labstack/echo/v4"
)

func InitRouter(svc *service.Service, e *echo.Group) {
	handler := New(svc)
	e.POST("/list", handler.ListExample)
	e.POST("/create", handler.CreateExample)
	e.POST("/mongo_list", handler.MgListExample)
	e.POST("/mongo_create", handler.MgCreateExample)
}

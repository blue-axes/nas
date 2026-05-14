package simple_upload

import (
	"github.com/blue-axes/tmpl/service"
	"github.com/labstack/echo/v4"
)

func InitRouter(svc *service.Service, e *echo.Group, writeMw ...echo.MiddlewareFunc) {
	handler := New(svc)

	e.HEAD("/object/*", handler.Schema)
	e.GET("/object/*", handler.Download)
	e.POST("/object/*", handler.Upload, writeMw...)
	e.DELETE("/object/*", handler.Delete, writeMw...)
	e.PATCH("/object/*", handler.UpdateTags, writeMw...)

	e.GET("/objects/*", handler.ReadDir)
	e.POST("/objects/", handler.MultiUpload, writeMw...)

	e.GET("/search", handler.Search)
	e.POST("/mkdir/*", handler.Mkdir, writeMw...)
}
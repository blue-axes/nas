package simple_upload

import (
	"github.com/blue-axes/tmpl/service"
	"github.com/labstack/echo/v4"
)

func InitRouter(svc *service.Service, e *echo.Group) {
	handler := New(svc)

	// 获取元数据
	e.HEAD("/object/*", handler.Schema)
	// 文件下载
	e.GET("/object/*", handler.Download)
	// 文件上传
	e.POST("/object/*", handler.Upload)
	// 文件删除
	e.DELETE("/object/*", handler.Delete)

	// 文件列表
	e.GET("/objects/*", handler.ReadDir)
	// 文件批量上传
	e.POST("/objects/", handler.MultiUpload)
}

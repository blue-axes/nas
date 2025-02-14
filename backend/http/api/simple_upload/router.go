package simple_upload

import (
	"github.com/blue-axes/tmpl/service"
	"github.com/labstack/echo/v4"
)

func InitRouter(svc *service.Service, e *echo.Group) {
	handler := New(svc)
	routePrefix := "/object/:name"
	// 获取对象元数据
	e.GET(routePrefix+"/info", handler.Schema)
	e.HEAD(routePrefix, handler.Schema)
	// 读取目录信息
	e.GET(routePrefix+"/read_dir", handler.ReadDir)
	// 获取一个对象
	e.GET(routePrefix, handler.Download)
	// 上传新文件
	e.POST(routePrefix, handler.Upload)
	// 上传多文件
	e.POST("/objects", handler.MultiUpload)
	// 删除文件
	e.DELETE(routePrefix, handler.Delete)
}

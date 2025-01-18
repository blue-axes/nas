package file_object

import (
	"github.com/blue-axes/tmpl/service"
	"github.com/labstack/echo/v4"
)

func InitRouter(svc *service.Service, e *echo.Group) {
	handler := New(svc)
	routePrefix := "/object/:name"
	// 获取对象元数据
	e.HEAD(routePrefix, handler.Upload)
	// 获取一个对象
	e.GET(routePrefix, handler.Upload)
	// 上传新文件
	e.POST(routePrefix, handler.Upload)
	// 更新文件
	e.PUT(routePrefix, handler.Upload)
	// 更新文件片段
	e.PATCH(routePrefix, handler.Upload)
	// 删除文件
	e.DELETE(routePrefix, handler.Upload)
}

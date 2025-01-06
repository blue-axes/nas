package http

import (
	"github.com/blue-axes/tmpl/http/api/example"
	"github.com/blue-axes/tmpl/service"
	"github.com/labstack/echo/v4"
)

func initRouter(svc *service.Service, e *echo.Echo) {
	// 集中式，方便查看
	//exampleHdl := example.New(svc)
	//exampleGrp := e.Group("/example")
	//exampleGrp.POST("/list", exampleHdl.ListExample)

	//分散式，便于管理
	example.InitRouter(svc, e.Group("/example"))

}

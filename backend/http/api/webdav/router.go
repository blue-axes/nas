package webdav

import (
	"net/http"

	"github.com/blue-axes/tmpl/service"
	"github.com/labstack/echo/v4"
)

func InitRouter(svc *service.Service, e *echo.Group, writeMw ...echo.MiddlewareFunc) {
	handler := New(svc)
	methods := []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodOptions,
		echo.PROPFIND,
		echo.REPORT,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodTrace,
		"MKCOL",
		"COPY",
		"MOVE",
		"LOCK",
		"UNLOCK",
		"PROPPATCH",
	}
	e.Match(methods, "", handler.Webdav)
	e.Match(methods, "/*", handler.Webdav)
}

package api

import (
	"github.com/blue-axes/tmpl/service"
	"github.com/labstack/echo/v4"
)

type FileScannerHandler struct {
	Handler
}

func NewFileScannerHandler(svc *service.Service) *FileScannerHandler {
	return &FileScannerHandler{
		Handler: *New(svc),
	}
}

func (h FileScannerHandler) Scan(c echo.Context) error {
	ctx := h.Ctx(c)
	result, err := h.svc.ScanFiles(ctx)
	return h.RespJson(c, result, err)
}
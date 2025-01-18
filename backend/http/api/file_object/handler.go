package file_object

import (
	"github.com/blue-axes/tmpl/http/api"
	"github.com/blue-axes/tmpl/service"
	"github.com/labstack/echo/v4"
)

type (
	FileObjectHandler struct {
		*api.Handler
		svc *service.Service
	}
)

func New(svc *service.Service) *FileObjectHandler {
	h := &FileObjectHandler{
		Handler: api.New(svc),
		svc:     svc,
	}
	return h
}

func (h FileObjectHandler) Upload(c echo.Context) error {

	return nil
}

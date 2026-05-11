package webdav

import (
	api "github.com/blue-axes/tmpl/http/api"
	"github.com/blue-axes/tmpl/service"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/webdav"
)

type (
	WebDavHandler struct {
		*api.Handler
		svc       *service.Service
		webHander *webdav.Handler
	}
)

func New(svc *service.Service) *WebDavHandler {
	webdavHdl, err := svc.GetWebDavHandler("/webdav")
	if err != nil {
		panic(err)
	}
	h := &WebDavHandler{
		Handler:   api.New(svc),
		svc:       svc,
		webHander: webdavHdl,
	}
	return h
}

func (h *WebDavHandler) Webdav(c echo.Context) error {
	r := c.Request()
	w := c.Response()
	h.webHander.ServeHTTP(w, r)
	return nil
}

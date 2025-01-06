package api

import (
	"github.com/blue-axes/tmpl/pkg/constants"
	"github.com/blue-axes/tmpl/pkg/context"
	"github.com/blue-axes/tmpl/service"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type (
	Handler struct {
		svc *service.Service
	}
)

func New(svc *service.Service) *Handler {
	h := &Handler{
		svc: svc,
	}
	return h
}

func (h Handler) Ctx(c echo.Context) *context.Context {
	ctx, _ := c.Get(constants.CtxKeyContext).(*context.Context)
	return ctx
}

func (h Handler) Log(ctx *context.Context) *log.Entry {
	if ctx == nil {
		return log.WithField("trace_id", "")
	}
	return log.WithField("trace_id", ctx.TraceID)
}

func (h Handler) RespJson(c echo.Context, data interface{}, err error) error {
	if err != nil {
		ErrorHandler(err, c)
		return nil
	}
	ctx, _ := c.Get(constants.CtxKeyContext).(*context.Context)

	resp := respStruct{
		Code:    constants.ErrCodeSuccess,
		Message: "",
		TraceID: ctx.TraceID,
		Data:    data,
	}
	return c.JSON(http.StatusOK, resp)
}

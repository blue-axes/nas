package api

import (
	"github.com/blue-axes/tmpl/pkg/constants"
	"github.com/blue-axes/tmpl/pkg/context"
	"github.com/blue-axes/tmpl/pkg/errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	respStruct struct {
		TraceID string      `json:"TraceID"`
		Code    string      `json:"Code"`
		Message string      `json:"Message"`
		Data    interface{} `json:"Data"`
	}
)

func ErrorHandler(err error, c echo.Context) {
	ctx, _ := c.Get(constants.CtxKeyContext).(*context.Context)
	resp := respStruct{
		Code:    constants.ErrCodeUnknown,
		Message: "unknown error",
		TraceID: ctx.TraceID,
	}
	switch verr := err.(type) {
	case *errors.Error:
		resp.Code = verr.Code()
		resp.Message = verr.Message()
	default:
		resp.Message = err.Error()
	}
	c.JSON(http.StatusOK, resp)
}

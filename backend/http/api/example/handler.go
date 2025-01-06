package example

import (
	"github.com/blue-axes/tmpl/http/api"
	"github.com/blue-axes/tmpl/pkg/constants"
	"github.com/blue-axes/tmpl/pkg/context"
	"github.com/blue-axes/tmpl/service"
	"github.com/blue-axes/tmpl/types"
	"github.com/labstack/echo/v4"
)

type (
	Example struct {
		*api.Handler
		svc *service.Service
	}
)

func New(svc *service.Service) *Example {
	h := &Example{
		Handler: api.New(svc),
		svc:     svc,
	}

	return h
}

func (h Example) ListExample(c echo.Context) error {
	ctx, _ := c.Get(constants.CtxKeyContext).(*context.Context)
	res, err := h.svc.ListExample(ctx)
	if err != nil {
		return err
	}
	return h.RespJson(c, res, nil)
}

func (h Example) CreateExample(c echo.Context) error {
	ctx, _ := c.Get(constants.CtxKeyContext).(*context.Context)
	var (
		params = struct {
			Name string `json:"Name"  validate:"required"`
		}{}
		response = struct {
			ID string `json:"ID"`
		}{}
	)
	if err := c.Bind(&params); err != nil {
		return h.RespJson(c, nil, err)
	}
	empl := &types.Example{
		Name: params.Name,
	}

	err := h.svc.CreateExample(ctx, empl)
	if err != nil {
		return err
	}
	response.ID = empl.ID
	return h.RespJson(c, response, nil)
}

func (h Example) MgListExample(c echo.Context) error {
	ctx, _ := c.Get(constants.CtxKeyContext).(*context.Context)
	res, err := h.svc.MgListExample(ctx)
	if err != nil {
		return err
	}
	return h.RespJson(c, res, nil)
}

func (h Example) MgCreateExample(c echo.Context) error {
	ctx, _ := c.Get(constants.CtxKeyContext).(*context.Context)
	var (
		params = struct {
			Name string `json:"Name"  validate:"required"`
		}{}
		response = struct {
			ID string `json:"ID"`
		}{}
	)
	if err := c.Bind(&params); err != nil {
		return h.RespJson(c, nil, err)
	}
	empl := &types.Example{
		Name: params.Name,
	}

	err := h.svc.MgCreateExample(ctx, empl)
	if err != nil {
		return err
	}
	response.ID = empl.ID
	return h.RespJson(c, response, nil)
}

package user

import (
	api "github.com/blue-axes/tmpl/http/api"
	"github.com/blue-axes/tmpl/pkg/constants"
	"github.com/blue-axes/tmpl/service"
	"github.com/blue-axes/tmpl/types"
	"github.com/labstack/echo/v4"
)

type (
	UserHandler struct {
		*api.Handler
		svc *service.Service
	}
)

func New(svc *service.Service) *UserHandler {
	return &UserHandler{
		Handler: api.New(svc),
		svc:     svc,
	}
}

func (h UserHandler) Me(c echo.Context) error {
	u, _ := c.Get(constants.CtxKeyUser).(*types.UserInfo)
	if u == nil {
		return h.RespJson(c, nil, nil)
	}
	return h.RespJson(c, map[string]interface{}{
		"Username": u.Username,
		"CanRead":  u.CanRead,
		"CanWrite": u.CanWrite,
		"IsAdmin":  u.IsAdmin,
	}, nil)
}

func (h UserHandler) List(c echo.Context) error {
	users, err := h.svc.ListUsers()
	return h.RespJson(c, users, err)
}

func (h UserHandler) Create(c echo.Context) error {
	var req struct {
		Username string `json:"Username"`
		Password string `json:"Password"`
		CanRead  bool   `json:"CanRead"`
		CanWrite bool   `json:"CanWrite"`
	}
	if err := c.Bind(&req); err != nil {
		return err
	}
	if req.Username == "" || req.Password == "" {
		return h.RespJson(c, nil, echo.NewHTTPError(400, "Username and Password are required"))
	}
	err := h.svc.CreateUser(req.Username, req.Password, req.CanRead, req.CanWrite)
	return h.RespJson(c, nil, err)
}

func (h UserHandler) Update(c echo.Context) error {
	username := c.Param("username")
	var req struct {
		CanRead  bool `json:"CanRead"`
		CanWrite bool `json:"CanWrite"`
	}
	if err := c.Bind(&req); err != nil {
		return err
	}
	err := h.svc.UpdateUserPermissions(username, req.CanRead, req.CanWrite)
	return h.RespJson(c, nil, err)
}

func (h UserHandler) Delete(c echo.Context) error {
	username := c.Param("username")
	err := h.svc.DeleteUser(username)
	return h.RespJson(c, nil, err)
}

func (h UserHandler) ChangePassword(c echo.Context) error {
	username := c.Param("username")
	var req struct {
		CurrentPassword string `json:"CurrentPassword"`
		NewPassword     string `json:"NewPassword"`
	}
	if err := c.Bind(&req); err != nil {
		return err
	}
	if req.CurrentPassword == "" || req.NewPassword == "" {
		return h.RespJson(c, nil, echo.NewHTTPError(400, "CurrentPassword and NewPassword are required"))
	}
	err := h.svc.ChangePassword(username, req.CurrentPassword, req.NewPassword)
	return h.RespJson(c, nil, err)
}
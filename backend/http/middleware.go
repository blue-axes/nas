package http

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/blue-axes/tmpl/pkg/constants"
	"github.com/blue-axes/tmpl/pkg/context"
	"github.com/blue-axes/tmpl/pkg/errors"
	"github.com/blue-axes/tmpl/types"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
)

type (
	Binder struct {
	}
)

func GetUserInfo(c echo.Context) *types.UserInfo {
	u, _ := c.Get(constants.CtxKeyUser).(*types.UserInfo)
	return u
}

func (Binder) Bind(i interface{}, c echo.Context) error {
	b := echo.DefaultBinder{}
	err := b.Bind(i, c)
	if err != nil {
		return err
	}
	val := reflect.ValueOf(i)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}

	valid := validator.New()
	switch val.Kind() {
	case reflect.Struct:
		err = valid.Struct(val.Interface())
		switch verr := err.(type) {
		case validator.FieldError:
			return errors.WithCode(constants.ErrCodeInvalidArgs, fmt.Sprintf("field:%s is invalid. %s", verr.Field(), verr.Error()))
		case validator.ValidationErrors:
			var (
				fields = make([]string, 0)
				errMsg = make([]string, 0)
			)
			for _, v := range verr {
				fields = append(fields, v.Field())
				errMsg = append(errMsg, v.Error())
			}
			return errors.WithCode(constants.ErrCodeInvalidArgs, fmt.Sprintf("fields:%s is invalid. %s", strings.Join(fields, ","), strings.Join(errMsg, ",")))
		default:
			return err
		}

	default:
		return nil
	}
}

func Pre(next echo.HandlerFunc) echo.HandlerFunc {
	traceID := uuid.New().String()
	ctx := context.New(context.WithTraceID(traceID))
	return func(c echo.Context) error {
		c.Set(constants.CtxKeyContext, ctx)

		return next(c)
	}
}

func BasicAuth(validator func(username, password string) (*types.UserInfo, bool)) echo.MiddlewareFunc {
	return echomw.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		userInfo, ok := validator(username, password)
		if ok {
			c.Set(constants.CtxKeyUser, userInfo)
			return true, nil
		}
		return false, nil
	})
}

func RequireRead(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		u := GetUserInfo(c)
		if u != nil && !u.CanRead {
			return echo.NewHTTPError(http.StatusForbidden, "read permission required")
		}
		return next(c)
	}
}

func RequireWrite(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		u := GetUserInfo(c)
		if u != nil && !u.CanWrite {
			return echo.NewHTTPError(http.StatusForbidden, "write permission required")
		}
		return next(c)
	}
}

func RequireAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		u := GetUserInfo(c)
		if u != nil && !u.IsAdmin {
			return echo.NewHTTPError(http.StatusForbidden, "admin permission required")
		}
		return next(c)
	}
}
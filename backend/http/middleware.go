package http

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/blue-axes/tmpl/pkg/constants"
	"github.com/blue-axes/tmpl/pkg/context"
	"github.com/blue-axes/tmpl/pkg/errors"
	"github.com/blue-axes/tmpl/pkg/utils"
	"github.com/blue-axes/tmpl/service"
	"github.com/blue-axes/tmpl/types"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
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

func Auth(svc *service.Service) echo.MiddlewareFunc {
	cookieName := svc.Config().Http.Auth.CookieName
	skipUrl := []string{
		"^/api/users/login$",
		"^/index[^/]+",
		"^/$",
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Request().URL.Path
			fmt.Printf("-----------------%+v\n", path)

			if utils.StrInArray(path, skipUrl, func(src, dst string) bool {
				r := regexp.MustCompile(dst)
				return r.MatchString(src)
			}) {
				return next(c)
			}

			if cookie, err := c.Cookie(cookieName); err == nil && cookie.Value != "" {
				if userInfo := svc.Session().Validate(cookie.Value); userInfo != nil {
					c.Set(constants.CtxKeyUser, userInfo)
					return next(c)
				}
			}

			if strings.HasPrefix(path, "/webdav") {
				if username, password, ok := c.Request().BasicAuth(); ok {
					if info, ok := svc.ValidateUser(username, password); ok {
						ui := &types.UserInfo{
							Username: info.Username,
							CanRead:  info.CanRead,
							CanWrite: info.CanWrite,
							IsAdmin:  info.IsAdmin,
						}
						c.Set(constants.CtxKeyUser, ui)
						return next(c)
					}
				}
				c.Response().Header().Set("WWW-Authenticate", `Basic realm="NAS"`)
			}

			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		}
	}
}

func FilePermissionCheck(next echo.HandlerFunc) echo.HandlerFunc {
	readMethod := []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodOptions,
		echo.PROPFIND,
		echo.REPORT,
	}
	writeMethod := []string{
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
	return func(c echo.Context) error {
		method := c.Request().Method
		u := GetUserInfo(c)
		if u == nil {
			return echo.NewHTTPError(http.StatusForbidden, "require login")
		}
		if utils.StrInArray(method, readMethod, nil) { // 读操作
			if !u.CanRead && !u.IsAdmin {
				return echo.NewHTTPError(http.StatusForbidden, "read permission required")
			}
		} else if utils.StrInArray(method, writeMethod, nil) { // 写操作
			if !u.CanWrite && !u.IsAdmin {
				return echo.NewHTTPError(http.StatusForbidden, "write permission required")
			}
		} else { // 未知操作
			if !u.IsAdmin { // 只有admin允许
				return echo.NewHTTPError(http.StatusForbidden, "unknown permission:"+method)
			}
		}
		return next(c)
	}
}

func RequireAdmin(excludeUrlPattern []string) echo.MiddlewareFunc {
	excludePatterns := []*regexp.Regexp{}
	for _, v := range excludeUrlPattern {
		excludePatterns = append(excludePatterns, regexp.MustCompile(v))
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			url := c.Request().URL.Path
			for _, v := range excludePatterns {
				if v.MatchString(url) {
					return next(c)
				}
			}
			u := GetUserInfo(c)
			if u != nil && !u.IsAdmin {
				return echo.NewHTTPError(http.StatusForbidden, "admin permission required")
			}
			return next(c)
		}
	}
}

package http

import (
	"fmt"
	"github.com/blue-axes/tmpl/pkg/constants"
	"github.com/blue-axes/tmpl/pkg/context"
	"github.com/blue-axes/tmpl/pkg/errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"reflect"
	"strings"
)

type (
	Binder struct {
	}
)

func (Binder) Bind(i interface{}, c echo.Context) error {
	b := echo.DefaultBinder{}
	err := b.Bind(i, c)
	if err != nil {
		return err
	}
	// validate
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

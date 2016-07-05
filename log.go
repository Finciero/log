package log

// TODO(jaguirre): add documentation

import (
	"os"
	"reflect"

	"github.com/Finciero/errors"
	"github.com/Finciero/utils/strings"
	kitlog "github.com/go-kit/kit/log"
)

// Logger ...
type Logger interface {
	Log(keyvalues ...interface{}) error
	Info(keyvalues ...interface{}) error
	Warn(keyvalues ...interface{}) error
	Error(err error, keyvalues ...interface{}) error

	With(keyvalues ...interface{}) Logger
}

// Context ...
type Context struct {
	kitlog.Context
}

// NewContext ...
func NewContext(keyvalues ...interface{}) *Context {
	logger := kitlog.NewJSONLogger(os.Stderr)
	return &Context{*kitlog.NewContext(logger).With(keyvalues...)}
}

// NewRequestContext ...
func NewRequestContext(requestID string, keyvalues ...interface{}) *Context {
	logger := kitlog.NewJSONLogger(os.Stderr)
	ctx := kitlog.NewContext(logger).With(keyvalues...)
	return &Context{*ctx.With("request_id", requestID)}
}

// With returns a new Context with keyvals appended to those of the receiver.
func (ctx *Context) With(keyvals ...interface{}) Logger {
	return &Context{*ctx.Context.With(keyvals...)}
}

// Info ...
func (ctx *Context) Info(keyvals ...interface{}) error {
	return ctx.Context.With("level", "info").Log(keyvals...)
}

func serializeStruct(val interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	el := reflect.ValueOf(val)

	if el.Kind() == reflect.Ptr {
		el = el.Elem()
	}

	tp := el.Type()

	for i := 0; i < el.NumField(); i++ {
		sf := tp.Field(i)
		fv := el.Field(i)

		if sf.PkgPath != "" {
			continue
		}

		key := strings.ToSnake(sf.Name)

		switch fv.Kind() {
		case reflect.Struct:
			res[key] = serializeStruct(fv.Interface())
		case reflect.Ptr:
			subEl := fv.Elem()
			if subEl.Kind() == reflect.Struct {
				res[key] = serializeStruct(fv.Interface())
			} else {
				res[key] = fv.Interface()
			}
		default:
			res[key] = fv.Interface()
		}
	}

	return res
}

// Error ...
func (ctx *Context) Error(err error, keyvals ...interface{}) error {
	if err == nil {
		return nil
	}

	if val, ok := (err).(*errors.Error); ok {
		for k, v := range val.Meta {
			keyvals = append(keyvals, k, v)
		}

		if val.InternalError != nil {
			errMap := serializeStruct(val.InternalError)
			keyvals = append(keyvals, "error", errMap)
		}

		return ctx.Context.With(
			"level", "error",
			"msg", val.Message,
		).Log(keyvals...)
	}

	return ctx.Context.With("level", "error", "msg", err).Log(keyvals...)
}

// Warn ...
func (ctx *Context) Warn(keyvals ...interface{}) error {
	return ctx.Context.With("level", "warning").Log(keyvals...)
}

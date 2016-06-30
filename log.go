package log

// TODO(jaguirre): add documentation

import (
	"os"

	"github.com/Finciero/errors"
	kitlog "github.com/go-kit/kit/log"
)

// Logger ...
type Logger interface {
	Log(keyvalues ...interface{}) error
	Info(keyvalues ...interface{}) error
	Warn(keyvalues ...interface{}) error
	Error(err error, keyvalues ...interface{}) error
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
func (ctx *Context) With(keyvals ...interface{}) *Context {
	return &Context{*ctx.Context.With(keyvals...)}
}

// Info ...
func (ctx *Context) Info(keyvals ...interface{}) error {
	return ctx.Context.With("level", "info").Log(keyvals...)
}

// Error ...
func (ctx *Context) Error(err error, keyvals ...interface{}) error {
	if err == nil {
		return nil
	}

	var desc string

	if val, ok := (err).(*errors.Error); ok {
		for k, v := range val.Meta {
			keyvals = append(keyvals, k, v)
		}
		desc = val.Description
	} else {
		desc = err.Error()
	}

	if len(desc) > 0 {
		keyvals = append(keyvals, "desc", desc)
	}

	return ctx.Context.With("level", "error").Log(keyvals...)
}

// Warn ...
func (ctx *Context) Warn(keyvals ...interface{}) error {
	return ctx.Context.With("level", "warning").Log(keyvals...)
}

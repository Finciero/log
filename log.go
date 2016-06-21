package log

import (
	"os"

	kitlog "github.com/go-kit/kit/log"
)

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

func (ctx *Context) Info(keyvals ...interface{}) {
	ctx.With("level", "info").Log(keyvals...)
}

func (ctx *Context) Error(keyvals ...interface{}) {
	ctx.With("level", "error").Log(keyvals...)
}

func (ctx *Context) Fatal(keyvals ...interface{}) {
	ctx.With("level", "fatal").Log(keyvals...)
	os.Exit(1)
}

func (ctx *Context) Warn(keyvals ...interface{}) {
	ctx.With("level", "warning").Log(keyvals...)
}

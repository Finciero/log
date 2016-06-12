package log

import (
	"os"

	kitlog "github.com/go-kit/kit/log"
)

type Context struct {
	kitlog.Context
}

// NewContext ...
func NewContext() *Context {
	logger := kitlog.NewJSONLogger(os.Stderr)
	return &Context{*kitlog.NewContext(logger)}
}

// NewRequestContext ...
func NewRequestContext(requestID string) *Context {
	logger := kitlog.NewJSONLogger(os.Stderr)
	return &Context{*kitlog.NewContext(logger).With("request_id", requestID)}
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

package log

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

var (
	reader, writer *os.File
	buffer         bytes.Buffer
)

func TestInfo(t *testing.T) {
	var (
		reqid = "test-id"
	)

	tests := []struct {
		in  []interface{}
		out string
	}{
		{[]interface{}{"foo", "bar"}, `{"foo":"bar","level":"info","request_id":"test-id"}`},
		{[]interface{}{"foo", 1}, `{"foo":1,"level":"info","request_id":"test-id"}`},
		{[]interface{}{"foo", true}, `{"foo":true,"level":"info","request_id":"test-id"}`},
		{[]interface{}{"foo", true, "bar", "bar", "baz", 1}, `{"bar":"bar","baz":1,"foo":true,"level":"info","request_id":"test-id"}`},
	}

	for _, tt := range tests {
		old := os.Stderr
		reader, writer, err := os.Pipe()
		ok(t, err)
		os.Stderr = writer

		ctx := NewRequestContext(reqid)
		ctx.Info(tt.in...)

		writer.Close()
		os.Stderr = old
		r := bufio.NewReader(reader)

		out, _, err := r.ReadLine()
		ok(t, err)
		equals(t, tt.out, string(out))
	}
}

func TestError(t *testing.T) {
	var (
		reqid = "test-id"
	)

	tests := []struct {
		in  []interface{}
		out string
	}{
		{[]interface{}{"foo", "bar"}, `{"foo":"bar","level":"error","request_id":"test-id"}`},
		{[]interface{}{"foo", 1}, `{"foo":1,"level":"error","request_id":"test-id"}`},
		{[]interface{}{"foo", true}, `{"foo":true,"level":"error","request_id":"test-id"}`},
		{[]interface{}{"foo", true, "bar", "bar", "baz", 1}, `{"bar":"bar","baz":1,"foo":true,"level":"error","request_id":"test-id"}`},
	}

	for _, tt := range tests {
		old := os.Stderr
		reader, writer, err := os.Pipe()
		ok(t, err)
		os.Stderr = writer

		ctx := NewRequestContext(reqid)
		ctx.Error(tt.in...)

		writer.Close()
		os.Stderr = old
		r := bufio.NewReader(reader)

		out, _, err := r.ReadLine()
		ok(t, err)
		equals(t, tt.out, string(out))
	}
}

func TestWarn(t *testing.T) {
	var (
		reqid = "test-id"
	)

	tests := []struct {
		in  []interface{}
		out string
	}{
		{[]interface{}{"foo", "bar"}, `{"foo":"bar","level":"warning","request_id":"test-id"}`},
		{[]interface{}{"foo", 1}, `{"foo":1,"level":"warning","request_id":"test-id"}`},
		{[]interface{}{"foo", true}, `{"foo":true,"level":"warning","request_id":"test-id"}`},
		{[]interface{}{"foo", true, "bar", "bar", "baz", 1}, `{"bar":"bar","baz":1,"foo":true,"level":"warning","request_id":"test-id"}`},
	}

	for _, tt := range tests {
		old := os.Stderr
		reader, writer, err := os.Pipe()
		ok(t, err)
		os.Stderr = writer

		ctx := NewRequestContext(reqid)
		ctx.Warn(tt.in...)

		writer.Close()
		os.Stderr = old
		r := bufio.NewReader(reader)

		out, _, err := r.ReadLine()
		ok(t, err)
		equals(t, tt.out, string(out))
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}
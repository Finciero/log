package log

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	finciero_errors "github.com/Finciero/errors"
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

type errParams struct {
	Foo string
	Bar int
}

type errTest struct {
	Foo errParams
}

func (e *errTest) Error() string {
	return "YOU CAN'T SEE ME"
}

func TestError(t *testing.T) {
	var (
		reqid = "test-id"
	)

	tests := []struct {
		err    error
		params []interface{}
		out    string
	}{
		{
			err:    nil,
			params: nil,
			out:    "",
		},
		{
			err: finciero_errors.NewFromError(finciero_errors.StatusNotFound, &errTest{
				Foo: errParams{Foo: "hi", Bar: 2},
			}, "finciero error"),
			params: []interface{}{"foo", "bar"},
			out:    `{"error":{"foo":{"bar":2,"foo":"hi"}},"foo":"bar","level":"error","msg":"finciero error","request_id":"test-id"}`,
		},
		{
			err:    finciero_errors.New(finciero_errors.StatusNotFound, "finciero error"),
			params: []interface{}{"foo", "bar"},
			out:    `{"foo":"bar","level":"error","msg":"finciero error","request_id":"test-id"}`,
		},
		{
			err:    finciero_errors.InternalServer("finciero error", finciero_errors.SetMeta(finciero_errors.Meta{"hi": "ho"})),
			params: []interface{}{"foo", "bar"},
			out:    `{"foo":"bar","hi":"ho","level":"error","msg":"finciero error","request_id":"test-id"}`,
		},
		{
			err: finciero_errors.InternalServer("finciero error", finciero_errors.SetMeta(finciero_errors.Meta{
				"hi": "ho",
				"ho": finciero_errors.Meta{
					"foo": "bar",
				},
			})),
			params: []interface{}{"foo", "bar"},
			out:    `{"foo":"bar","hi":"ho","ho":{"foo":"bar"},"level":"error","msg":"finciero error","request_id":"test-id"}`,
		},
		{errors.New(""), []interface{}{"foo", "bar"}, `{"foo":"bar","level":"error","msg":"","request_id":"test-id"}`},
		{errors.New(""), []interface{}{"foo", 1}, `{"foo":1,"level":"error","msg":"","request_id":"test-id"}`},
		{errors.New("error"), []interface{}{"foo", true}, `{"foo":true,"level":"error","msg":"error","request_id":"test-id"}`},
		{errors.New(""), []interface{}{"foo", true, "bar", "bar", "baz", 1}, `{"bar":"bar","baz":1,"foo":true,"level":"error","msg":"","request_id":"test-id"}`},
	}

	for _, tt := range tests {
		old := os.Stderr
		reader, writer, err := os.Pipe()
		ok(t, err)
		os.Stderr = writer

		ctx := NewRequestContext(reqid)
		ctx.Error(tt.err, tt.params...)

		writer.Close()
		os.Stderr = old
		r := bufio.NewReader(reader)

		out, _, err := r.ReadLine()
		if tt.out == "" && err == io.EOF {
			continue
		}
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

func TestSerializeStruct(t *testing.T) {
	type testStruct struct {
		A int
	}

	tests := []struct {
		val interface{}
		exp map[string]interface{}
	}{
		{
			val: struct{ A *testStruct }{A: &testStruct{3}},
			exp: map[string]interface{}{"a": map[string]interface{}{"a": 3}},
		},
		{
			val: struct{ A *testStruct }{A: nil},
			exp: map[string]interface{}{"a": (*testStruct)(nil)},
		},
		{
			val: struct{ a *testStruct }{a: &testStruct{3}},
			exp: map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		res := serializeStruct(tt.val)
		equals(t, res, tt.exp)
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

package log

import (
	"bufio"
	"errors"
	"io"
	"os"
	"testing"

	finciero_errors "github.com/Finciero/errors"
)

func TestNOOPNOOPInfo(t *testing.T) {
	tests := []struct {
		in  []interface{}
		out string
	}{
		{[]interface{}{"foo", "bar"}, ""},
		{[]interface{}{"foo", 1}, ""},
		{[]interface{}{"foo", true}, ""},
		{[]interface{}{"foo", true, "bar", "bar", "baz", 1}, ""},
	}

	for _, tt := range tests {
		old := os.Stderr
		reader, writer, err := os.Pipe()
		ok(t, err)
		os.Stderr = writer

		ctx := NewNOOPContext()
		ctx.Info(tt.in...)

		writer.Close()
		os.Stderr = old
		r := bufio.NewReader(reader)

		out, _, err := r.ReadLine()
		if string(out) != tt.out {
			t.Errorf("out expected %q, but got %q", tt.out, out)
		}
		if err != io.EOF {
			t.Errorf("an `io.EOF` error was expected but got %q", err.Error())
		}
	}
}

func TestNOOPError(t *testing.T) {
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
			err:    finciero_errors.New(finciero_errors.StatusNotFound, "finciero error"),
			params: []interface{}{"foo", "bar"},
			out:    "",
		},
		{
			err:    finciero_errors.InternalServer("finciero error", finciero_errors.SetMeta(finciero_errors.Meta{"hi": "ho"})),
			params: []interface{}{"foo", "bar"},
			out:    "",
		},
		{
			err: finciero_errors.InternalServer("finciero error", finciero_errors.SetMeta(finciero_errors.Meta{
				"hi": "ho",
				"ho": finciero_errors.Meta{
					"foo": "bar",
				},
			})),
			params: []interface{}{"foo", "bar"},
			out:    "",
		},
		{errors.New(""), []interface{}{"foo", "bar"}, ""},
		{errors.New(""), []interface{}{"foo", 1}, ""},
		{errors.New("error"), []interface{}{"foo", true}, ""},
		{errors.New(""), []interface{}{"foo", true, "bar", "bar", "baz", 1}, ""},
	}

	for _, tt := range tests {
		old := os.Stderr
		reader, writer, err := os.Pipe()
		ok(t, err)
		os.Stderr = writer

		ctx := NewNOOPContext()
		ctx.Error(tt.err, tt.params...)

		writer.Close()
		os.Stderr = old
		r := bufio.NewReader(reader)

		out, _, err := r.ReadLine()
		if string(out) != tt.out {
			t.Errorf("out expected %q, but got %q", tt.out, out)
		}
		if err != io.EOF {
			t.Errorf("an `io.EOF` error was expected but got %q", err.Error())
		}
	}
}

func TestNOOPWarn(t *testing.T) {
	tests := []struct {
		in  []interface{}
		out string
	}{
		{[]interface{}{"foo", "bar"}, ""},
		{[]interface{}{"foo", 1}, ""},
		{[]interface{}{"foo", true}, ""},
		{[]interface{}{"foo", true, "bar", "bar", "baz", 1}, ""},
	}

	for _, tt := range tests {
		old := os.Stderr
		reader, writer, err := os.Pipe()
		ok(t, err)
		os.Stderr = writer

		ctx := NewNOOPContext()
		ctx.Warn(tt.in...)

		writer.Close()
		os.Stderr = old
		r := bufio.NewReader(reader)

		out, _, err := r.ReadLine()
		if string(out) != tt.out {
			t.Errorf("out expected %q, but got %q", tt.out, out)
		}
		if err != io.EOF {
			t.Errorf("an `io.EOF` error was expected but got %q", err.Error())
		}
	}
}

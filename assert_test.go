package flags_test

import (
	"github.com/jessevdk/go-flags"
	"testing"
)

func assertString(t *testing.T, a string, b string) {
	if a != b {
		t.Errorf("Expected %#v, but got %#v", b, a)
	}
}
func assertStringArray(t *testing.T, a []string, b []string) {
	if len(a) != len(b) {
		t.Errorf("Expected %#v, but got %#v", b, a)
		return
	}

	for i, v := range a {
		if b[i] != v {
			t.Errorf("Expected %#v, but got %#v", b, a)
			return
		}
	}
}

func assertBoolArray(t *testing.T, a []bool, b []bool) {
	if len(a) != len(b) {
		t.Errorf("Expected %#v, but got %#v", b, a)
		return
	}

	for i, v := range a {
		if b[i] != v {
			t.Errorf("Expected %#v, but got %#v", b, a)
			return
		}
	}
}

func assertParserSuccess(t *testing.T, data interface{}, args ...string) (*flags.Parser, []string) {
	parser := flags.NewParser(data, flags.Default&^flags.PrintErrors)
	ret, err := parser.ParseArgs(args)

	if err != nil {
		t.Fatalf("Unexpected parse error: %s", err)
		return nil, nil
	}

	return parser, ret
}

func assertParseSuccess(t *testing.T, data interface{}, args ...string) []string {
	_, ret := assertParserSuccess(t, data, args...)
	return ret
}

func assertParseFail(t *testing.T, typ flags.ErrorType, msg string, data interface{}, args ...string) {
	parser := flags.NewParser(data, flags.Default&^flags.PrintErrors)
	_, err := parser.ParseArgs(args)

	if err == nil {
		t.Fatalf("Expected error: %s", msg)
		return
	}

	if e, ok := err.(*flags.Error); !ok {
		t.Fatalf("Expected Error type, but got %#v", err)
		return
	} else {
		if e.Type != typ {
			t.Errorf("Expected error type {%s}, but got {%s}", typ, e.Type)
		}

		if e.Message != msg {
			t.Errorf("Expected error message %#v, but got %#v", msg, e.Message)
		}
	}
}
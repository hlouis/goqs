package goqs

// This test file is follow https://github.com/ljharb/qs/blob/main/test/parse.js

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	Input  string
	Result *QSType
}

func TestParseSimpleString(t *testing.T) {
	cases := []testCase{
		{"0=foo", &QSType{0: "foo"}},
		{"foo=bar", &QSType{"foo": "bar"}},
		{"foo=c++", &QSType{"foo": "c  "}},
		{"a[>=]=23", &QSType{"a": QSType{">=": "23"}}},
		{"a[<=>]==23", &QSType{"a": QSType{"<=>": "=23"}}},
		{"a[==]=23", &QSType{"a": QSType{"==": "23"}}},
	}

	d := NewDecoder()

	for _, c := range cases {
		res, err := d.Parse(c.Input)
		assert.NoError(t, err, "parse %v failed %v", c.Input, err)
		assert.Equal(t, c.Result, res, "parse %v not equal.", c.Input)
	}
}

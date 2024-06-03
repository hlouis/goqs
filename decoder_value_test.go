package goqs

// This test file is follow https://github.com/ljharb/qs/blob/main/test/parse.js

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type valueTestCase struct {
	Input  string
	Result map[string]interface{}
}

func TestParseValue(t *testing.T) {
	cases := []valueTestCase{
		{"0=foo", map[string]interface{}{"0": "foo"}},
		{"foo=bar", map[string]interface{}{"foo": "bar"}},
		{"foo=c++", map[string]interface{}{"foo": "c  "}},
		{"a[>=]=23", map[string]interface{}{"a[>=]": "23"}},
		{"a[<=>]==23", map[string]interface{}{"a[<=>]": "=23"}},
		{"a[==]=23", map[string]interface{}{"a[==]": "23"}},
		{"a[1]=3&a[2]=4", map[string]interface{}{"a[1]": "3", "a[2]": "4"}},
		{"a.b=3&a[]=4", map[string]interface{}{"a.b": "3", "a[]": "4"}},
	}

	d := NewDecoder()

	for _, c := range cases {
		res := d.parseValues(c.Input)
		assert.Equal(t, c.Result, res, "parse %v not equal.", c.Input)
		t.Logf("parse from %v\t to %v\n", c.Input, res)
	}
}

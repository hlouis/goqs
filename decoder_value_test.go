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

func _doTest(d *Decoder, t *testing.T, cases []valueTestCase) {
	for _, c := range cases {
		res := d.parseValues(c.Input)
		assert.Equal(t, c.Result, res, "parse %v not equal.", c.Input)
		t.Logf("parse from %v\t to %v\n", c.Input, res)
	}
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
		{"foo=a,b", map[string]interface{}{"foo": "a,b"}},
		{"foo[]=a,b", map[string]interface{}{"foo[]": "a,b"}},
	}

	d := NewDecoder()
	_doTest(d, t, cases)
}

func TestParseValueWithComma(t *testing.T) {
	cases := []valueTestCase{
		{"foo=a,b", map[string]interface{}{"foo": []interface{}{"a", "b"}}},
		{"foo=a,b&foo=c", map[string]interface{}{"foo": []interface{}{"a", "b", "c"}}},
		{"foo[]=a,b", map[string]interface{}{"foo[]": []interface{}{[]interface{}{"a", "b"}}}},
		{"foo[]=a,b&foo[]=c", map[string]interface{}{"foo[]": []interface{}{[]interface{}{"a", "b"}, "c"}}},
	}

	d := NewDecoder(WithComma(true))
	_doTest(d, t, cases)
}

type keyTestCase struct {
	Key    string
	Val    interface{}
	Result QSType
}

func TestParseKeys(t *testing.T) {
	cases := []keyTestCase{
		{"foo[]", []interface{}{"a", "b"}, QSType{"foo": []interface{}{"a", "b"}}},
		{"foo[1]", []interface{}{"a", "b"}, QSType{"foo": map[interface{}]interface{}{"1": []interface{}{"a", "b"}}}},
		{"a.b", "apple", QSType{"a.b": "apple"}},
		{"a.b[]", "apple", QSType{"a.b": []interface{}{"apple"}}},
	}

	d := NewDecoder()
	for _, c := range cases {
		res := d.parseKeys(c.Key, c.Val)
		assert.Equal(t, c.Result, res, "parse %v not equal!", c.Key)
		t.Logf("parse key from %v:%v\t to %v\n", c.Key, c.Val, res)
	}
}

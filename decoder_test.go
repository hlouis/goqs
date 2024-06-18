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

func _test(d *Decoder, t *testing.T, cases []testCase) {
	for _, c := range cases {
		res, err := d.Parse(c.Input)
		assert.NoError(t, err, "parse %v failed %v", c.Input, err)
		assert.Equal(t, c.Result, res, "parse %v not equal.", c.Input)
		t.Logf("parse from %v\t to %v\n", c.Input, res)
	}
}

func TestUseage(t *testing.T) {
	d := NewDecoder()
	res, err := d.Parse("foo=bar")
	if err != nil {
		t.Errorf("decode failed: %v", err)
	}
	bar := (*res)["foo"]
	assert.Equal(t, "bar", bar, "Get data from res: %v", bar)
}

func TestParseSimpleString(t *testing.T) {
	cases := []testCase{
		{"0=foo", &QSType{"0": "foo"}},
		{"foo=c++", &QSType{"foo": "c  "}},
		{"a[>=]=23", &QSType{"a": QSType{">=": "23"}}},
		{"a[<=>]==23", &QSType{"a": QSType{"<=>": "=23"}}},
		{"a[==]=23", &QSType{"a": QSType{"==": "23"}}},
		{"foo", &QSType{"foo": ""}},
		{"foo=", &QSType{"foo": ""}},
		{"foo=bar", &QSType{"foo": "bar"}},
		{" foo = bar = baz ", &QSType{" foo ": " bar = baz "}},
		{"foo=bar&bar=baz", &QSType{"foo": "bar", "bar": "baz"}},
		{"foo2=bar2&baz2=", &QSType{"foo2": "bar2", "baz2": ""}},
		{"foo=bar&baz", &QSType{"foo": "bar", "baz": ""}},
		{"cht=p3&chd=t:60,40&chs=250x100&chl=Hello|World", &QSType{
			"cht": "p3",
			"chd": "t:60,40",
			"chs": "250x100",
			"chl": "Hello|World",
		}},
	}

	d := NewDecoder()
	_test(d, t, cases)
}

func TestParseNoComma(t *testing.T) {
	cases := []testCase{
		{"a[]=b&a[]=c", &QSType{"a": []interface{}{"b", "c"}}},
		{"a[0]=b&a[1]=c", &QSType{"a": []interface{}{"b", "c"}}},
		{"a=b,c", &QSType{"a": "b,c"}},
		{"a=b&a=c", &QSType{"a": []interface{}{"b", "c"}}},
	}

	d := NewDecoder()
	_test(d, t, cases)
}

func TestParseWithComma(t *testing.T) {
	cases := []testCase{
		{"a[]=b&a[]=c", &QSType{"a": []interface{}{"b", "c"}}},
		{"a[0]=b&a[1]=c", &QSType{"a": []interface{}{"b", "c"}}},
		{"a=b,c", &QSType{"a": []interface{}{"b", "c"}}},
		{"a=b&a=c", &QSType{"a": []interface{}{"b", "c"}}},
		{"a[0]=a,b&a[1]=c,d", &QSType{"a": []interface{}{[]interface{}{"a", "b"}, []interface{}{"c", "d"}}}},
	}

	d := NewDecoder(WithComma(true))
	_test(d, t, cases)
}

func TestDebug(t *testing.T) {
	cases := []testCase{
		{"a[0]=a,b&a[1]=c,d", &QSType{"a": []interface{}{[]interface{}{"a", "b"}, []interface{}{"c", "d"}}}},
	}

	d := NewDecoder(WithComma(true))
	_test(d, t, cases)
}

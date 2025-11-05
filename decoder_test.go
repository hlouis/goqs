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
		{"a.b=bar", &QSType{"a.b": "bar"}},
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

func TestParseWithAllowDots(t *testing.T) {
	cases := []testCase{
		{"a.b=bar", &QSType{"a": QSType{"b": "bar"}}},
	}

	d := NewDecoder(WithAllowDots(true))
	_test(d, t, cases)
}

func TestDebug(t *testing.T) {
	cases := []testCase{
		{"a[0]=a,b&a[1]=c,d", &QSType{"a": []interface{}{[]interface{}{"a", "b"}, []interface{}{"c", "d"}}}},
	}

	d := NewDecoder(WithComma(true))
	_test(d, t, cases)
}

// Test fixes for the 5 known differences between Go and JS implementations

func TestFix1_EmptyKeyHandling(t *testing.T) {
	// Issue: Trailing delimiter creates empty key
	// Expected: {"a": "1"} not {"a": "1", "": ""}
	cases := []testCase{
		{"a=1&", &QSType{"a": "1"}},
		{"a=1&b=2&", &QSType{"a": "1", "b": "2"}},
		{"&a=1", &QSType{"a": "1"}},
		{"a=1&&b=2", &QSType{"a": "1", "b": "2"}},
	}

	d := NewDecoder()
	_test(d, t, cases)
}

func TestFix2_EmptyArrays(t *testing.T) {
	// Issue: Empty array contains empty string instead of being truly empty
	// Expected: {"foo": []} not {"foo": [""]}
	d := NewDecoder(WithAllowEmptyArrays(true))

	res, err := d.Parse("foo[]")
	assert.NoError(t, err)

	fooVal := (*res)["foo"]
	assert.NotNil(t, fooVal, "foo should not be nil")

	fooArr, ok := fooVal.([]interface{})
	assert.True(t, ok, "foo should be an array")
	assert.Equal(t, 0, len(fooArr), "foo array should be empty, got: %v", fooArr)

	// Also test with = sign
	res2, err := d.Parse("foo[]=")
	assert.NoError(t, err)
	fooVal2 := (*res2)["foo"]
	fooArr2, ok := fooVal2.([]interface{})
	assert.True(t, ok, "foo should be an array")
	assert.Equal(t, 0, len(fooArr2), "foo array should be empty")
}

func TestFix3_PreEncodedCommaHandling(t *testing.T) {
	// Issue: Pre-encoded comma (%2C) is decoded then split
	// Expected: {"foo": "a,b"} not {"foo": ["a", "b"]}
	d := NewDecoder(WithComma(true))

	// Test pre-encoded comma - should NOT split
	res1, err := d.Parse("foo=a%2Cb")
	assert.NoError(t, err)
	assert.Equal(t, "a,b", (*res1)["foo"], "Pre-encoded comma should not be split")

	// Test raw comma - should split
	res2, err := d.Parse("foo=a,b")
	assert.NoError(t, err)
	expected := []interface{}{"a", "b"}
	assert.Equal(t, expected, (*res2)["foo"], "Raw comma should be split")

	// Test mixed - one encoded, one raw
	res3, err := d.Parse("foo=a%2Cb,c")
	assert.NoError(t, err)
	expected3 := []interface{}{"a,b", "c"}
	assert.Equal(t, expected3, (*res3)["foo"])
}

func TestFix4_DepthLimitOffByOne(t *testing.T) {
	// Issue: Depth counting is off by one
	// Expected: { a: { b: { '[c][d]': 'e' } } } not {"a": {"b": {"c": {"[d]": "e"}}}}

	// Test depth=2
	d2 := NewDecoder(WithDepth(2))
	res2, err := d2.Parse("a[b][c][d]=e")
	assert.NoError(t, err)

	aVal := (*res2)["a"].(QSType)
	bVal := aVal["b"].(QSType)

	// With depth=2, '[c][d]' should be a literal key, not parsed further
	_, hasC := bVal["c"]
	assert.False(t, hasC, "Should not have key 'c' at depth 2")

	_, hasCDKey := bVal["[c][d]"]
	assert.True(t, hasCDKey, "Should have literal key '[c][d]'")
	assert.Equal(t, "e", bVal["[c][d]"])

	// Test depth=1
	d1 := NewDecoder(WithDepth(1))
	res1, err := d1.Parse("a[b][c][d]=e")
	assert.NoError(t, err)

	aVal1 := (*res1)["a"].(QSType)
	// With depth=1, '[b][c][d]' should be a literal key
	_, hasB := aVal1["b"]
	assert.False(t, hasB, "Should not have key 'b' at depth 1")

	_, hasBCDKey := aVal1["[b][c][d]"]
	assert.True(t, hasBCDKey, "Should have literal key '[b][c][d]'")
	assert.Equal(t, "e", aVal1["[b][c][d]"])

	// Test depth=3
	d3 := NewDecoder(WithDepth(3))
	res3, err := d3.Parse("a[b][c][d]=e")
	assert.NoError(t, err)

	aVal3 := (*res3)["a"].(QSType)
	bVal3 := aVal3["b"].(QSType)
	cVal3 := bVal3["c"].(QSType)

	// With depth=3, only [d] should be unparsed
	_, hasDKey := cVal3["[d]"]
	assert.True(t, hasDKey, "Should have literal key '[d]'")
	assert.Equal(t, "e", cVal3["[d]"])
}

func TestFix5_ArrayLimitKeyType(t *testing.T) {
	// Issue: When exceeding arrayLimit, key is string instead of integer
	// Expected: { a: { 100: 'b' } } (integer key) not {"a": {"100": "b"}} (string key)
	d := NewDecoder(WithArrayLimit(20))

	res, err := d.Parse("a[100]=b")
	assert.NoError(t, err)

	aVal := (*res)["a"].(QSType)

	// Check that the key is an integer, not a string
	_, hasStringKey := aVal["100"]
	assert.False(t, hasStringKey, "Should not have string key '100'")

	val, hasIntKey := aVal[100]
	assert.True(t, hasIntKey, "Should have integer key 100")
	assert.Equal(t, "b", val)

	// Verify it's an object, not an array (because index > arrayLimit)
	_, isArray := (*res)["a"].([]interface{})
	assert.False(t, isArray, "Should be an object, not an array")

	// Test that indices within arrayLimit still work and create arrays
	res2, err := d.Parse("a[0]=x&a[1]=y")
	assert.NoError(t, err)
	aVal2, isArray2 := (*res2)["a"].([]interface{})
	assert.True(t, isArray2, "Should be an array for continuous indices within limit")
	assert.Equal(t, "x", aVal2[0])
	assert.Equal(t, "y", aVal2[1])
}

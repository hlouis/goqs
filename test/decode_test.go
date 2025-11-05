package test

import (
	"testing"

	"github.com/hlouis/goqs"
	"github.com/stretchr/testify/assert"
)

// TestParseSimpleString tests basic string parsing
func TestParseSimpleString(t *testing.T) {
	d := goqs.NewDecoder()

	tests := []struct {
		name     string
		input    string
		opts     []goqs.DecoderOption
		expected *goqs.QSType
	}{
		{
			name:     "numeric key",
			input:    "0=foo",
			expected: &goqs.QSType{"0": "foo"},
		},
		{
			name:     "url encoded plus",
			input:    "foo=c++",
			expected: &goqs.QSType{"foo": "c  "},
		},
		{
			name:     "bracket notation with special chars",
			input:    "a[>=]=23",
			expected: &goqs.QSType{"a": goqs.QSType{">=": "23"}},
		},
		{
			name:     "bracket notation with arrow",
			input:    "a[<=>]==23",
			expected: &goqs.QSType{"a": goqs.QSType{"<=>": "=23"}},
		},
		{
			name:     "bracket notation with equals",
			input:    "a[==]=23",
			expected: &goqs.QSType{"a": goqs.QSType{"==": "23"}},
		},
		{
			name:     "simple key value",
			input:    "foo=bar",
			expected: &goqs.QSType{"foo": "bar"},
		},
		{
			name:     "preserves spaces",
			input:    " foo = bar = baz ",
			expected: &goqs.QSType{" foo ": " bar = baz "},
		},
		{
			name:     "equals in value",
			input:    "foo=bar=baz",
			expected: &goqs.QSType{"foo": "bar=baz"},
		},
		{
			name:     "multiple params",
			input:    "foo=bar&bar=baz",
			expected: &goqs.QSType{"foo": "bar", "bar": "baz"},
		},
		{
			name:     "empty value",
			input:    "foo2=bar2&baz2=",
			expected: &goqs.QSType{"foo2": "bar2", "baz2": ""},
		},
		{
			name:     "complex chart params",
			input:    "cht=p3&chd=t:60,40&chs=250x100&chl=Hello|World",
			expected: &goqs.QSType{"cht": "p3", "chd": "t:60,40", "chs": "250x100", "chl": "Hello|World"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.opts) > 0 {
				d = goqs.NewDecoder(tt.opts...)
			} else {
				d = goqs.NewDecoder()
			}
			result, err := d.Parse(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestParseStrictNullHandling tests null handling
func TestParseStrictNullHandling(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		opts     []goqs.DecoderOption
		expected *goqs.QSType
	}{
		{
			name:     "key without value - strict null",
			input:    "foo",
			opts:     []goqs.DecoderOption{goqs.WithStrictNullHandling(true)},
			expected: &goqs.QSType{"foo": nil},
		},
		{
			name:     "key without value - default",
			input:    "foo",
			expected: &goqs.QSType{"foo": ""},
		},
		{
			name:     "key with equals - strict null",
			input:    "foo=",
			opts:     []goqs.DecoderOption{goqs.WithStrictNullHandling(true)},
			expected: &goqs.QSType{"foo": ""},
		},
		{
			name:     "mixed with strict null",
			input:    "foo=bar&baz",
			opts:     []goqs.DecoderOption{goqs.WithStrictNullHandling(true)},
			expected: &goqs.QSType{"foo": "bar", "baz": nil},
		},
		{
			name:     "mixed without strict null",
			input:    "foo=bar&baz",
			expected: &goqs.QSType{"foo": "bar", "baz": ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := goqs.NewDecoder(tt.opts...)
			result, err := d.Parse(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestParseArraysWithoutComma tests array parsing without comma option
func TestParseArraysWithoutComma(t *testing.T) {
	d := goqs.NewDecoder()

	tests := []struct {
		name     string
		input    string
		expected *goqs.QSType
	}{
		{
			name:     "explicit array syntax",
			input:    "a[]=b&a[]=c",
			expected: &goqs.QSType{"a": []interface{}{"b", "c"}},
		},
		{
			name:     "indexed array",
			input:    "a[0]=b&a[1]=c",
			expected: &goqs.QSType{"a": []interface{}{"b", "c"}},
		},
		{
			name:     "comma not parsed",
			input:    "a=b,c",
			expected: &goqs.QSType{"a": "b,c"},
		},
		{
			name:     "duplicate keys combine",
			input:    "a=b&a=c",
			expected: &goqs.QSType{"a": []interface{}{"b", "c"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := d.Parse(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestParseArraysWithComma tests array parsing with comma option enabled
func TestParseArraysWithComma(t *testing.T) {
	d := goqs.NewDecoder(goqs.WithComma(true))

	tests := []struct {
		name     string
		input    string
		expected *goqs.QSType
	}{
		{
			name:     "explicit array syntax",
			input:    "a[]=b&a[]=c",
			expected: &goqs.QSType{"a": []interface{}{"b", "c"}},
		},
		{
			name:     "indexed array",
			input:    "a[0]=b&a[1]=c",
			expected: &goqs.QSType{"a": []interface{}{"b", "c"}},
		},
		{
			name:     "comma parsed as array",
			input:    "a=b,c",
			expected: &goqs.QSType{"a": []interface{}{"b", "c"}},
		},
		{
			name:     "duplicate keys combine",
			input:    "a=b&a=c",
			expected: &goqs.QSType{"a": []interface{}{"b", "c"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := d.Parse(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestParseDotNotation tests dot notation support
func TestParseDotNotation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		opts     []goqs.DecoderOption
		expected *goqs.QSType
	}{
		{
			name:     "dot notation disabled",
			input:    "a.b=c",
			expected: &goqs.QSType{"a.b": "c"},
		},
		{
			name:     "dot notation enabled",
			input:    "a.b=c",
			opts:     []goqs.DecoderOption{goqs.WithAllowDots(true)},
			expected: &goqs.QSType{"a": goqs.QSType{"b": "c"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := goqs.NewDecoder(tt.opts...)
			result, err := d.Parse(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestParseDecodeDotInKeys tests decoding dots in keys
func TestParseDecodeDotInKeys(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		opts     []goqs.DecoderOption
		expected *goqs.QSType
	}{
		{
			name:     "allowDots false, decodeDotInKeys false",
			input:    "name%252Eobj.first=John&name%252Eobj.last=Doe",
			opts:     []goqs.DecoderOption{goqs.WithAllowDots(false), goqs.WithDecodeDotInKeys(false)},
			expected: &goqs.QSType{"name%2Eobj.first": "John", "name%2Eobj.last": "Doe"},
		},
		{
			name:     "allowDots true, decodeDotInKeys false",
			input:    "name.obj.first=John&name.obj.last=Doe",
			opts:     []goqs.DecoderOption{goqs.WithAllowDots(true), goqs.WithDecodeDotInKeys(false)},
			expected: &goqs.QSType{"name": goqs.QSType{"obj": goqs.QSType{"first": "John", "last": "Doe"}}},
		},
		{
			name:     "allowDots true, decodeDotInKeys false, encoded dot",
			input:    "name%252Eobj.first=John&name%252Eobj.last=Doe",
			opts:     []goqs.DecoderOption{goqs.WithAllowDots(true), goqs.WithDecodeDotInKeys(false)},
			expected: &goqs.QSType{"name%2Eobj": goqs.QSType{"first": "John", "last": "Doe"}},
		},
		{
			name:     "allowDots true, decodeDotInKeys true",
			input:    "name%252Eobj.first=John&name%252Eobj.last=Doe",
			opts:     []goqs.DecoderOption{goqs.WithAllowDots(true), goqs.WithDecodeDotInKeys(true)},
			expected: &goqs.QSType{"name.obj": goqs.QSType{"first": "John", "last": "Doe"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := goqs.NewDecoder(tt.opts...)
			result, err := d.Parse(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestParseNestedObjects tests nested object parsing
func TestParseNestedObjects(t *testing.T) {
	d := goqs.NewDecoder()

	tests := []struct {
		name     string
		input    string
		expected *goqs.QSType
	}{
		{
			name:     "single level nesting",
			input:    "a[b]=c",
			expected: &goqs.QSType{"a": goqs.QSType{"b": "c"}},
		},
		{
			name:     "double level nesting",
			input:    "a[b][c]=d",
			expected: &goqs.QSType{"a": goqs.QSType{"b": goqs.QSType{"c": "d"}}},
		},
		{
			name:     "triple level nesting",
			input:    "a[b][c][d]=e",
			expected: &goqs.QSType{"a": goqs.QSType{"b": goqs.QSType{"c": goqs.QSType{"d": "e"}}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := d.Parse(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestParseDepthLimit tests depth limiting
func TestParseDepthLimit(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		opts     []goqs.DecoderOption
		expected *goqs.QSType
	}{
		{
			name:  "exceeds default depth (5)",
			input: "a[b][c][d][e][f][g][h]=i",
			expected: &goqs.QSType{
				"a": goqs.QSType{
					"b": goqs.QSType{
						"c": goqs.QSType{
							"d": goqs.QSType{
								"e": goqs.QSType{
									"f": goqs.QSType{
										"[g][h]": "i",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "custom depth limit",
			input: "a[b][c][d]=e",
			opts:  []goqs.DecoderOption{goqs.WithDepth(2)},
			expected: &goqs.QSType{
				"a": goqs.QSType{
					"b": goqs.QSType{
						"[c][d]": "e",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := goqs.NewDecoder(tt.opts...)
			result, err := d.Parse(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestParseArrayLimit tests array limit enforcement
func TestParseArrayLimit(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		opts     []goqs.DecoderOption
		expected *goqs.QSType
	}{
		{
			name:  "within array limit",
			input: "a[0]=b&a[1]=c&a[2]=d",
			opts:  []goqs.DecoderOption{goqs.WithArrayLimit(20)},
			expected: &goqs.QSType{
				"a": []interface{}{"b", "c", "d"},
			},
		},
		{
			name:  "exceeds array limit",
			input: "a[100]=b",
			opts:  []goqs.DecoderOption{goqs.WithArrayLimit(20)},
			expected: &goqs.QSType{
				"a": goqs.QSType{
					100: "b",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := goqs.NewDecoder(tt.opts...)
			result, err := d.Parse(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestParseParameterLimit tests parameter limit enforcement
func TestParseParameterLimit(t *testing.T) {
	d := goqs.NewDecoder(goqs.WithParameterLimit(2))

	result, err := d.Parse("a=b&c=d&e=f")
	assert.NoError(t, err)
	// Should only parse first 2 parameters
	assert.Equal(t, 2, len(*result))
	assert.Equal(t, "b", (*result)["a"])
	assert.Equal(t, "d", (*result)["c"])
}

// TestParseCustomDelimiter tests custom delimiter support
func TestParseCustomDelimiter(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		delimiter string
		expected  *goqs.QSType
	}{
		{
			name:      "semicolon delimiter",
			input:     "a=b;c=d",
			delimiter: ";",
			expected:  &goqs.QSType{"a": "b", "c": "d"},
		},
		{
			name:      "comma delimiter",
			input:     "a=b,c=d",
			delimiter: ",",
			expected:  &goqs.QSType{"a": "b", "c": "d"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := goqs.NewDecoder(goqs.WithDelimiter(tt.delimiter))
			result, err := d.Parse(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestParseIgnoreQueryPrefix tests ignoring query prefix
func TestParseIgnoreQueryPrefix(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		opts     []goqs.DecoderOption
		expected *goqs.QSType
	}{
		{
			name:     "with query prefix - ignored",
			input:    "?foo=bar&baz=qux",
			opts:     []goqs.DecoderOption{goqs.WithIgnoreQueryPrefix(true)},
			expected: &goqs.QSType{"foo": "bar", "baz": "qux"},
		},
		{
			name:     "with query prefix - not ignored",
			input:    "?foo=bar",
			expected: &goqs.QSType{"?foo": "bar"},
		},
		{
			name:     "without query prefix - ignored option",
			input:    "foo=bar",
			opts:     []goqs.DecoderOption{goqs.WithIgnoreQueryPrefix(true)},
			expected: &goqs.QSType{"foo": "bar"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := goqs.NewDecoder(tt.opts...)
			result, err := d.Parse(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestParseURLEncoding tests URL encoding/decoding
func TestParseURLEncoding(t *testing.T) {
	d := goqs.NewDecoder()

	tests := []struct {
		name     string
		input    string
		expected *goqs.QSType
	}{
		{
			name:     "encoded spaces in key",
			input:    "a[b%20c]=d",
			expected: &goqs.QSType{"a": goqs.QSType{"b c": "d"}},
		},
		{
			name:     "encoded special chars",
			input:    "a=%2B%3D%26",
			expected: &goqs.QSType{"a": "+=&"},
		},
		{
			name:     "plus to space",
			input:    "a=hello+world",
			expected: &goqs.QSType{"a": "hello world"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := d.Parse(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestParseEmptyKeys tests empty key handling
func TestParseEmptyKeys(t *testing.T) {
	d := goqs.NewDecoder()

	tests := []struct {
		name     string
		input    string
		expected *goqs.QSType
	}{
		{
			name:     "trailing ampersand",
			input:    "_r=1&",
			expected: &goqs.QSType{"_r": "1"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: &goqs.QSType{},
		},
		{
			name:     "only ampersands",
			input:    "&&&",
			expected: &goqs.QSType{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := d.Parse(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestParseAllowEmptyArrays tests empty array support
func TestParseAllowEmptyArrays(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		opts     []goqs.DecoderOption
		expected *goqs.QSType
	}{
		{
			name:     "empty array with option enabled",
			input:    "foo[]",
			opts:     []goqs.DecoderOption{goqs.WithAllowEmptyArrays(true)},
			expected: &goqs.QSType{"foo": []interface{}{}},
		},
		{
			name:     "empty array with option disabled",
			input:    "foo[]",
			expected: &goqs.QSType{"foo": []interface{}{""}},
		},
		{
			name:     "multiple empty arrays",
			input:    "foo[]&bar[]=a",
			opts:     []goqs.DecoderOption{goqs.WithAllowEmptyArrays(true)},
			expected: &goqs.QSType{"foo": []interface{}{}, "bar": []interface{}{"a"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := goqs.NewDecoder(tt.opts...)
			result, err := d.Parse(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestParseMixedArraysAndObjects tests mixed array and object structures
func TestParseMixedArraysAndObjects(t *testing.T) {
	d := goqs.NewDecoder()

	tests := []struct {
		name     string
		input    string
		expected *goqs.QSType
	}{
		{
			name:  "array inside object",
			input: "a[b][]=c&a[b][]=d",
			expected: &goqs.QSType{
				"a": goqs.QSType{
					"b": []interface{}{"c", "d"},
				},
			},
		},
		{
			name:  "object inside array",
			input: "a[0][b]=c&a[1][d]=e",
			expected: &goqs.QSType{
				"a": []interface{}{
					goqs.QSType{"b": "c"},
					goqs.QSType{"d": "e"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := d.Parse(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestParseWithCommaInValue tests handling commas in values with comma option
func TestParseWithCommaInValue(t *testing.T) {
	d := goqs.NewDecoder(goqs.WithComma(true))

	tests := []struct {
		name     string
		input    string
		expected *goqs.QSType
	}{
		{
			name:     "encoded comma not split",
			input:    "foo=a%2Cb",
			expected: &goqs.QSType{"foo": "a,b"},
		},
		{
			name:     "regular comma split",
			input:    "foo=a,b",
			expected: &goqs.QSType{"foo": []interface{}{"a", "b"}},
		},
		{
			name:     "array with comma values",
			input:    "a[]=1,2,3&a[]=4,5,6",
			expected: &goqs.QSType{"a": []interface{}{[]interface{}{"1", "2", "3"}, []interface{}{"4", "5", "6"}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := d.Parse(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

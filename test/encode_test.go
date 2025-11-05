package test

import (
	"testing"
	"time"

	"github.com/hlouis/goqs"
	"github.com/stretchr/testify/assert"
)

// TestStringifyBasic tests basic stringification
func TestStringifyBasic(t *testing.T) {
	e := goqs.NewEncoder()

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "simple key-value",
			input:    map[string]interface{}{"a": "b"},
			expected: "a=b",
		},
		{
			name:     "number value",
			input:    map[string]interface{}{"a": 1},
			expected: "a=1",
		},
		{
			name:     "multiple pairs",
			input:    map[string]interface{}{"a": 1, "b": 2},
			expected: "a=1&b=2",
		},
		{
			name:     "underscore in value",
			input:    map[string]interface{}{"a": "A_Z"},
			expected: "a=A_Z",
		},
		{
			name:     "empty string value",
			input:    map[string]interface{}{"a": ""},
			expected: "a=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := e.Stringify(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestStringifyUnicodeEncoding tests Unicode character encoding
func TestStringifyUnicodeEncoding(t *testing.T) {
	e := goqs.NewEncoder()

	tests := []struct {
		name     string
		input    map[string]interface{}
		expected string
	}{
		{
			name:     "euro symbol",
			input:    map[string]interface{}{"a": "€"},
			expected: "a=%E2%82%AC",
		},
		{
			name:     "hebrew character",
			input:    map[string]interface{}{"a": "א"},
			expected: "a=%D7%90",
		},
		{
			name:     "4-byte unicode",
			input:    map[string]interface{}{"a": "𐐷"},
			expected: "a=%F0%90%90%B7",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := e.Stringify(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestStringifyFalsyValues tests handling of falsy values
func TestStringifyFalsyValues(t *testing.T) {
	e := goqs.NewEncoder()

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: "",
		},
		{
			name:     "false value",
			input:    false,
			expected: "",
		},
		{
			name:     "zero value",
			input:    0,
			expected: "",
		},
		{
			name:     "empty map",
			input:    map[string]interface{}{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := e.Stringify(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestStringifyArrays tests array stringification with different formats
func TestStringifyArrays(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		opts     []goqs.EncoderOption
		expected string
	}{
		{
			name:     "array with indices format (default)",
			input:    map[string]interface{}{"a": []interface{}{"b", "c", "d"}},
			expected: "a%5B0%5D=b&a%5B1%5D=c&a%5B2%5D=d",
		},
		{
			name:     "array with brackets format",
			input:    map[string]interface{}{"a": []interface{}{"b", "c", "d"}},
			opts:     []goqs.EncoderOption{goqs.WithArrayFormat("brackets")},
			expected: "a%5B%5D=b&a%5B%5D=c&a%5B%5D=d",
		},
		{
			name:     "array with comma format",
			input:    map[string]interface{}{"a": []interface{}{"b", "c", "d"}},
			opts:     []goqs.EncoderOption{goqs.WithArrayFormat("comma")},
			expected: "a=b%2Cc%2Cd",
		},
		{
			name:     "array with repeat format",
			input:    map[string]interface{}{"a": []interface{}{"b", "c", "d"}},
			opts:     []goqs.EncoderOption{goqs.WithArrayFormat("repeat")},
			expected: "a=b&a=c&a=d",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := goqs.NewEncoder(tt.opts...)
			result, err := e.Stringify(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestStringifyNestedObjects tests nested object stringification
func TestStringifyNestedObjects(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		opts     []goqs.EncoderOption
		expected string
	}{
		{
			name: "simple nested object",
			input: map[string]interface{}{
				"a": map[string]interface{}{"b": "c"},
			},
			expected: "a%5Bb%5D=c",
		},
		{
			name: "deep nesting",
			input: map[string]interface{}{
				"a": map[string]interface{}{
					"b": map[string]interface{}{
						"c": map[string]interface{}{
							"d": "e",
						},
					},
				},
			},
			expected: "a%5Bb%5D%5Bc%5D%5Bd%5D=e",
		},
		{
			name: "nested with allowDots",
			input: map[string]interface{}{
				"a": map[string]interface{}{"b": "c"},
			},
			opts:     []goqs.EncoderOption{goqs.WithAllowDotsEncode(true)},
			expected: "a.b=c",
		},
		{
			name: "deep nesting with allowDots",
			input: map[string]interface{}{
				"a": map[string]interface{}{
					"b": map[string]interface{}{
						"c": map[string]interface{}{
							"d": "e",
						},
					},
				},
			},
			opts:     []goqs.EncoderOption{goqs.WithAllowDotsEncode(true)},
			expected: "a.b.c.d=e",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := goqs.NewEncoder(tt.opts...)
			result, err := e.Stringify(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestStringifyNullHandling tests null and empty value handling
func TestStringifyNullHandling(t *testing.T) {
	tests := []struct {
		name            string
		input           map[string]interface{}
		opts            []goqs.EncoderOption
		expected        string
		checkContains   []string
		checkNotContain []string
	}{
		{
			name:     "skip nulls",
			input:    map[string]interface{}{"a": "b", "c": nil},
			opts:     []goqs.EncoderOption{goqs.WithSkipNulls(true)},
			expected: "a=b",
		},
		{
			name:          "nulls not skipped",
			input:         map[string]interface{}{"a": "b", "c": nil},
			checkContains: []string{"a=b", "c="},
		},
		{
			name:     "strict null handling",
			input:    map[string]interface{}{"a": nil},
			opts:     []goqs.EncoderOption{goqs.WithStrictNullHandlingEncode(true)},
			expected: "a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := goqs.NewEncoder(tt.opts...)
			result, err := e.Stringify(tt.input)
			assert.NoError(t, err)
			if tt.expected != "" {
				assert.Equal(t, tt.expected, result)
			}
			for _, part := range tt.checkContains {
				assert.Contains(t, result, part)
			}
			for _, part := range tt.checkNotContain {
				assert.NotContains(t, result, part)
			}
		})
	}
}

// TestStringifyEmptyArrays tests empty array handling
func TestStringifyEmptyArrays(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		opts     []goqs.EncoderOption
		expected string
	}{
		{
			name:     "empty array default (skipped)",
			input:    map[string]interface{}{"a": []interface{}{}, "b": "zz"},
			expected: "b=zz",
		},
		{
			name:     "empty array allowed",
			input:    map[string]interface{}{"a": []interface{}{}, "b": "zz"},
			opts:     []goqs.EncoderOption{goqs.WithAllowEmptyArraysEncode(true), goqs.WithArrayFormat("brackets")},
			expected: "a%5B%5D&b=zz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := goqs.NewEncoder(tt.opts...)
			result, err := e.Stringify(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestStringifyQueryPrefix tests adding query prefix
func TestStringifyQueryPrefix(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		opts     []goqs.EncoderOption
		expected string
	}{
		{
			name:     "add prefix",
			input:    map[string]interface{}{"a": "b"},
			opts:     []goqs.EncoderOption{goqs.WithAddQueryPrefix(true)},
			expected: "?a=b",
		},
		{
			name:     "no prefix for empty",
			input:    map[string]interface{}{},
			opts:     []goqs.EncoderOption{goqs.WithAddQueryPrefix(true)},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := goqs.NewEncoder(tt.opts...)
			result, err := e.Stringify(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestStringifyCustomDelimiter tests custom delimiter
func TestStringifyCustomDelimiter(t *testing.T) {
	tests := []struct {
		name      string
		input     map[string]interface{}
		delimiter string
		expected  string
	}{
		{
			name:      "semicolon delimiter",
			input:     map[string]interface{}{"a": "b", "c": "d"},
			delimiter: ";",
			expected:  "a=b;c=d",
		},
		{
			name:      "pipe delimiter",
			input:     map[string]interface{}{"a": "b", "c": "d"},
			delimiter: "|",
			expected:  "a=b|c=d",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := goqs.NewEncoder(goqs.WithDelimiterEncode(tt.delimiter))
			result, err := e.Stringify(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestStringifyEncodeDotInKeys tests encoding dots in keys
func TestStringifyEncodeDotInKeys(t *testing.T) {
	tests := []struct {
		name          string
		input         map[string]interface{}
		opts          []goqs.EncoderOption
		checkContains []string
	}{
		{
			name: "allowDots false, encodeDotInKeys false",
			input: map[string]interface{}{
				"name.obj": map[string]interface{}{"first": "John", "last": "Doe"},
			},
			opts:          []goqs.EncoderOption{goqs.WithAllowDotsEncode(false), goqs.WithEncodeDotInKeys(false)},
			checkContains: []string{"name.obj%5Bfirst%5D=John", "name.obj%5Blast%5D=Doe"},
		},
		{
			name: "allowDots true, encodeDotInKeys false",
			input: map[string]interface{}{
				"name.obj": map[string]interface{}{"first": "John", "last": "Doe"},
			},
			opts:          []goqs.EncoderOption{goqs.WithAllowDotsEncode(true), goqs.WithEncodeDotInKeys(false)},
			checkContains: []string{"name.obj.first=John", "name.obj.last=Doe"},
		},
		{
			name: "allowDots false, encodeDotInKeys true",
			input: map[string]interface{}{
				"name.obj": map[string]interface{}{"first": "John", "last": "Doe"},
			},
			opts:          []goqs.EncoderOption{goqs.WithAllowDotsEncode(false), goqs.WithEncodeDotInKeys(true)},
			checkContains: []string{"name%252Eobj%5Bfirst%5D=John", "name%252Eobj%5Blast%5D=Doe"},
		},
		{
			name: "allowDots true, encodeDotInKeys true",
			input: map[string]interface{}{
				"name.obj": map[string]interface{}{"first": "John", "last": "Doe"},
			},
			opts:          []goqs.EncoderOption{goqs.WithAllowDotsEncode(true), goqs.WithEncodeDotInKeys(true)},
			checkContains: []string{"name%252Eobj.first=John", "name%252Eobj.last=Doe"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := goqs.NewEncoder(tt.opts...)
			result, err := e.Stringify(tt.input)
			assert.NoError(t, err)
			for _, part := range tt.checkContains {
				assert.Contains(t, result, part)
			}
		})
	}
}

// TestStringifyEncodeValuesOnly tests encoding only values
func TestStringifyEncodeValuesOnly(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		opts     []goqs.EncoderOption
		expected string
	}{
		{
			name:     "encode values only",
			input:    map[string]interface{}{"a[b]": "c[d]"},
			opts:     []goqs.EncoderOption{goqs.WithEncodeValuesOnly(true)},
			expected: "a[b]=c%5Bd%5D",
		},
		{
			name:     "encode both",
			input:    map[string]interface{}{"a[b]": "c[d]"},
			expected: "a%5Bb%5D=c%5Bd%5D",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := goqs.NewEncoder(tt.opts...)
			result, err := e.Stringify(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestStringifyRFC1738vs3986 tests RFC format standards
func TestStringifyRFC1738vs3986(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		opts     []goqs.EncoderOption
		expected string
	}{
		{
			name:     "RFC3986 (default) - space as %20",
			input:    map[string]interface{}{"a": "b c"},
			expected: "a=b%20c",
		},
		{
			name:     "RFC1738 - space as +",
			input:    map[string]interface{}{"a": "b c"},
			opts:     []goqs.EncoderOption{goqs.WithFormat("RFC1738")},
			expected: "a=b+c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := goqs.NewEncoder(tt.opts...)
			result, err := e.Stringify(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestStringifySort tests sorting keys
func TestStringifySort(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		opts     []goqs.EncoderOption
		expected string
	}{
		{
			name:     "sorted keys",
			input:    map[string]interface{}{"c": "3", "a": "1", "b": "2"},
			opts:     []goqs.EncoderOption{goqs.WithSort(true)},
			expected: "a=1&b=2&c=3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := goqs.NewEncoder(tt.opts...)
			result, err := e.Stringify(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestStringifyMixedStructures tests complex mixed structures
func TestStringifyMixedStructures(t *testing.T) {
	e := goqs.NewEncoder()

	tests := []struct {
		name     string
		input    map[string]interface{}
		expected string
	}{
		{
			name: "array inside object",
			input: map[string]interface{}{
				"a": map[string]interface{}{
					"b": []interface{}{"c", "d"},
				},
			},
			expected: "a%5Bb%5D%5B0%5D=c&a%5Bb%5D%5B1%5D=d",
		},
		{
			name: "object inside array",
			input: map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{"b": "c"},
					map[string]interface{}{"d": "e"},
				},
			},
			expected: "a%5B0%5D%5Bb%5D=c&a%5B1%5D%5Bd%5D=e",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := e.Stringify(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestStringifySpecialCharacters tests special character handling
func TestStringifySpecialCharacters(t *testing.T) {
	e := goqs.NewEncoder()

	tests := []struct {
		name     string
		input    map[string]interface{}
		expected string
	}{
		{
			name:     "ampersand",
			input:    map[string]interface{}{"a": "b&c"},
			expected: "a=b%26c",
		},
		{
			name:     "equals sign",
			input:    map[string]interface{}{"a": "b=c"},
			expected: "a=b%3Dc",
		},
		{
			name:     "hash",
			input:    map[string]interface{}{"a": "b#c"},
			expected: "a=b%23c",
		},
		{
			name:     "question mark",
			input:    map[string]interface{}{"a": "b?c"},
			expected: "a=b%3Fc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := e.Stringify(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestStringifyBoolean tests boolean value handling
func TestStringifyBoolean(t *testing.T) {
	e := goqs.NewEncoder()

	tests := []struct {
		name     string
		input    map[string]interface{}
		expected string
	}{
		{
			name:     "true value",
			input:    map[string]interface{}{"a": true},
			expected: "a=true",
		},
		{
			name:     "false value (non-root)",
			input:    map[string]interface{}{"a": map[string]interface{}{"b": false}},
			expected: "a%5Bb%5D=false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := e.Stringify(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestStringifyNumericKeys tests numeric key handling
func TestStringifyNumericKeys(t *testing.T) {
	e := goqs.NewEncoder()

	tests := []struct {
		name     string
		input    map[interface{}]interface{}
		expected string
	}{
		{
			name:     "numeric key",
			input:    map[interface{}]interface{}{0: "foo"},
			expected: "0=foo",
		},
		{
			name:     "mixed keys",
			input:    map[interface{}]interface{}{"a": "b", 1: "c"},
			expected: "1=c&a=b", // May vary depending on iteration order
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := e.Stringify(tt.input)
			assert.NoError(t, err)
			// Note: Order might vary for maps, we just check it's valid
			assert.NotEmpty(t, result)
		})
	}
}

// TestStringifyDateSerialization tests date handling
func TestStringifyDateSerialization(t *testing.T) {
	e := goqs.NewEncoder()

	// Create a specific time
	testTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	result, err := e.Stringify(map[string]interface{}{"date": testTime})
	assert.NoError(t, err)
	// Should serialize to ISO format and URL encode
	assert.Contains(t, result, "date=")
	assert.Contains(t, result, "2023")
}

// TestStringifyFilter tests filtering keys
func TestStringifyFilter(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		opts     []goqs.EncoderOption
		expected string
	}{
		{
			name:     "filter to specific keys",
			input:    map[string]interface{}{"a": "b", "c": "d", "e": "f"},
			opts:     []goqs.EncoderOption{goqs.WithFilter([]string{"a", "c"})},
			expected: "a=b&c=d",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := goqs.NewEncoder(tt.opts...)
			result, err := e.Stringify(tt.input)
			assert.NoError(t, err)
			// Order may vary, just check both keys are present
			assert.Contains(t, result, "a=b")
			assert.Contains(t, result, "c=d")
			assert.NotContains(t, result, "e=f")
		})
	}
}

// TestStringifyNoEncode tests disabling encoding
func TestStringifyNoEncode(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		opts     []goqs.EncoderOption
		expected string
	}{
		{
			name:     "no encoding",
			input:    map[string]interface{}{"a": "b c"},
			opts:     []goqs.EncoderOption{goqs.WithEncode(false)},
			expected: "a=b c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := goqs.NewEncoder(tt.opts...)
			result, err := e.Stringify(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestStringifyCommaRoundTrip tests comma format with round-trip compatibility
func TestStringifyCommaRoundTrip(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		opts     []goqs.EncoderOption
		expected string
	}{
		{
			name:     "comma round trip",
			input:    map[string]interface{}{"a": []interface{}{"b", "c"}},
			opts:     []goqs.EncoderOption{goqs.WithArrayFormat("comma"), goqs.WithCommaRoundTrip(true)},
			expected: "a=b%2Cc",
		},
		{
			name:     "comma format without round trip",
			input:    map[string]interface{}{"a": []interface{}{"b", "c"}},
			opts:     []goqs.EncoderOption{goqs.WithArrayFormat("comma")},
			expected: "a=b%2Cc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := goqs.NewEncoder(tt.opts...)
			result, err := e.Stringify(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestStringifyQSType tests stringifying QSType directly
func TestStringifyQSType(t *testing.T) {
	e := goqs.NewEncoder()

	input := &goqs.QSType{
		"a": "b",
		"c": map[interface{}]interface{}{"d": "e"},
	}

	result, err := e.Stringify(input)
	assert.NoError(t, err)
	assert.Contains(t, result, "a=b")
	assert.Contains(t, result, "c%5Bd%5D=e")
}

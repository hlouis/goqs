# goqs

A Golang port of the popular JavaScript library [ljharb/qs](https://github.com/ljharb/qs), providing advanced query string parsing and stringification with support for nested objects and arrays.

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.22-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## Features

- 🔄 **Bidirectional**: Parse query strings to Go structs and stringify Go values to query strings
- 🎯 **Nested Objects**: Full support for deeply nested objects using bracket notation
- 📦 **Arrays**: Multiple array formats (indices, brackets, comma, repeat)
- 🔤 **Dot Notation**: Optional dot notation for nested objects (`a.b.c=value`)
- 🎨 **Flexible Options**: 25+ configuration options for customization
- ✅ **Well Tested**: 350+ test cases ported from the original JavaScript library
- 🚀 **Production Ready**: Suitable for REST APIs, URL building, and form data parsing

## Installation

```bash
go get github.com/hlouis/goqs
```

## Quick Start

### Parsing (Decode)

```go
import "github.com/hlouis/goqs"

// Basic parsing
d := goqs.NewDecoder()
result, err := d.Parse("foo=bar&baz=qux")
// result: &QSType{"foo": "bar", "baz": "qux"}

// Arrays
result, _ = d.Parse("a[]=b&a[]=c")
// result: &QSType{"a": []interface{}{"b", "c"}}

// Nested objects
result, _ = d.Parse("user[name]=John&user[age]=30")
// result: &QSType{"user": QSType{"name": "John", "age": "30"}}
```

### Stringifying (Encode)

```go
// Basic stringification
e := goqs.NewEncoder()
query, err := e.Stringify(map[string]interface{}{
    "foo": "bar",
    "baz": "qux",
})
// query: "foo=bar&baz=qux"

// Arrays
query, _ = e.Stringify(map[string]interface{}{
    "colors": []interface{}{"red", "green", "blue"},
})
// query: "colors%5B0%5D=red&colors%5B1%5D=green&colors%5B2%5D=blue"

// Nested objects
query, _ = e.Stringify(map[string]interface{}{
    "user": map[string]interface{}{
        "name": "John",
        "age":  30,
    },
})
// query: "user%5Bname%5D=John&user%5Bage%5D=30"
```

## Decoder Options

### Basic Options

```go
// Allow dot notation: a.b=c → {"a": {"b": "c"}}
d := goqs.NewDecoder(goqs.WithAllowDots(true))

// Parse commas as arrays: a=b,c → {"a": ["b", "c"]}
d := goqs.NewDecoder(goqs.WithComma(true))

// Custom delimiter
d := goqs.NewDecoder(goqs.WithDelimiter(";"))
result, _ := d.Parse("a=b;c=d")

// Ignore query prefix
d := goqs.NewDecoder(goqs.WithIgnoreQueryPrefix(true))
result, _ := d.Parse("?foo=bar")  // Ignores the "?"
```

### Advanced Options

```go
// Null handling
d := goqs.NewDecoder(goqs.WithStrictNullHandling(true))
result, _ := d.Parse("foo")  // {"foo": nil} instead of {"foo": ""}

// Depth limit (default: 5)
d := goqs.NewDecoder(goqs.WithDepth(3))

// Array limit (default: 20)
d := goqs.NewDecoder(goqs.WithArrayLimit(100))

// Parameter limit (default: 1000)
d := goqs.NewDecoder(goqs.WithParameterLimit(500))

// Allow empty arrays
d := goqs.NewDecoder(goqs.WithAllowEmptyArrays(true))
result, _ := d.Parse("foo[]")  // {"foo": []}

// Decode dots in keys
d := goqs.NewDecoder(
    goqs.WithAllowDots(true),
    goqs.WithDecodeDotInKeys(true),
)
result, _ := d.Parse("name%252Eobj.first=John")
// {"name.obj": {"first": "John"}}
```

### All Decoder Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `WithAllowDots` | `bool` | `false` | Enable dot notation parsing |
| `WithAllowEmptyArrays` | `bool` | `false` | Parse empty brackets as empty arrays |
| `WithComma` | `bool` | `false` | Parse comma-separated values as arrays |
| `WithDecodeDotInKeys` | `bool` | `false` | Decode %2E as literal dots in keys |
| `WithDelimiter` | `string` | `"&"` | Query string delimiter |
| `WithDepth` | `int` | `5` | Maximum nesting depth |
| `WithArrayLimit` | `int` | `20` | Maximum array index |
| `WithParameterLimit` | `int` | `1000` | Maximum number of parameters |
| `WithIgnoreQueryPrefix` | `bool` | `false` | Ignore leading `?` |
| `WithStrictNullHandling` | `bool` | `false` | Keys without values return `nil` |

## Encoder Options

### Array Formats

```go
// Indices (default): a[0]=b&a[1]=c
e := goqs.NewEncoder(goqs.WithArrayFormat("indices"))

// Brackets: a[]=b&a[]=c
e := goqs.NewEncoder(goqs.WithArrayFormat("brackets"))

// Repeat: a=b&a=c
e := goqs.NewEncoder(goqs.WithArrayFormat("repeat"))

// Comma: a=b,c
e := goqs.NewEncoder(goqs.WithArrayFormat("comma"))
```

### Encoding Options

```go
// Dot notation
e := goqs.NewEncoder(goqs.WithAllowDotsEncode(true))
query, _ := e.Stringify(map[string]interface{}{
    "user": map[string]interface{}{"name": "John"},
})
// query: "user.name=John"

// Encode dots in keys
e := goqs.NewEncoder(
    goqs.WithAllowDotsEncode(true),
    goqs.WithEncodeDotInKeys(true),
)
query, _ := e.Stringify(map[string]interface{}{
    "name.obj": map[string]interface{}{"first": "John"},
})
// query: "name%252Eobj.first=John"

// Encode values only (not keys)
e := goqs.NewEncoder(goqs.WithEncodeValuesOnly(true))
query, _ := e.Stringify(map[string]interface{}{"a[b]": "c[d]"})
// query: "a[b]=c%5Bd%5D"

// Skip null values
e := goqs.NewEncoder(goqs.WithSkipNulls(true))
query, _ := e.Stringify(map[string]interface{}{
    "a": "b",
    "c": nil,  // Will be omitted
})
// query: "a=b"

// Strict null handling
e := goqs.NewEncoder(goqs.WithStrictNullHandlingEncode(true))
query, _ := e.Stringify(map[string]interface{}{"a": nil})
// query: "a" (no equals sign)
```

### Formatting Options

```go
// Add query prefix
e := goqs.NewEncoder(goqs.WithAddQueryPrefix(true))
query, _ := e.Stringify(map[string]interface{}{"a": "b"})
// query: "?a=b"

// Custom delimiter
e := goqs.NewEncoder(goqs.WithDelimiterEncode(";"))
// query: "a=b;c=d"

// Sort keys
e := goqs.NewEncoder(goqs.WithSort(true))
// query: "a=1&b=2&c=3" (alphabetically sorted)

// RFC1738 format (space as +)
e := goqs.NewEncoder(goqs.WithFormat("RFC1738"))
// Default is RFC3986 (space as %20)

// Allow empty arrays
e := goqs.NewEncoder(
    goqs.WithAllowEmptyArraysEncode(true),
    goqs.WithArrayFormat("brackets"),
)
query, _ := e.Stringify(map[string]interface{}{"a": []interface{}{}})
// query: "a%5B%5D"

// Filter keys
e := goqs.NewEncoder(goqs.WithFilter([]string{"a", "c"}))
query, _ := e.Stringify(map[string]interface{}{
    "a": "1",
    "b": "2",  // Will be omitted
    "c": "3",
})
// query: "a=1&c=3"

// Disable encoding
e := goqs.NewEncoder(goqs.WithEncode(false))
query, _ := e.Stringify(map[string]interface{}{"a": "b c"})
// query: "a=b c" (not encoded)
```

### All Encoder Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `WithAddQueryPrefix` | `bool` | `false` | Prepend `?` to output |
| `WithAllowDotsEncode` | `bool` | `false` | Use dot notation for nested objects |
| `WithAllowEmptyArraysEncode` | `bool` | `false` | Include empty arrays |
| `WithArrayFormat` | `string` | `"indices"` | Array format: indices/brackets/repeat/comma |
| `WithCharset` | `string` | `"utf-8"` | Character encoding (utf-8 only) |
| `WithCharsetSentinelEncode` | `bool` | `false` | Add charset sentinel |
| `WithCommaRoundTrip` | `bool` | `false` | Comma format compatibility |
| `WithDelimiterEncode` | `string` | `"&"` | Query string delimiter |
| `WithEncode` | `bool` | `true` | Enable URL encoding |
| `WithEncodeDotInKeys` | `bool` | `false` | Encode literal dots in keys |
| `WithEncodeValuesOnly` | `bool` | `false` | Only encode values, not keys |
| `WithFilter` | `[]string` | `nil` | Include only specified keys |
| `WithFormat` | `string` | `"RFC3986"` | RFC1738 (+) or RFC3986 (%20) |
| `WithSerializeDate` | `func` | `nil` | Custom date serialization |
| `WithSkipNulls` | `bool` | `false` | Omit null values |
| `WithSort` | `bool` | `false` | Sort keys alphabetically |
| `WithStrictNullHandlingEncode` | `bool` | `false` | Omit `=` for null values |

## Type System

### QSType

The `QSType` is the core data structure:

```go
type QSType map[interface{}]interface{}
```

This allows mixed string and integer keys. When accessing values, type assertions are required:

```go
result, _ := d.Parse("user[name]=John&user[age]=30")

// Access nested values with type assertion
user := (*result)["user"].(goqs.QSType)
name := user["name"].(string)
age := user["age"].(string)

// Helper function recommended
func getString(m goqs.QSType, key string) string {
    if v, ok := m[key].(string); ok {
        return v
    }
    return ""
}
```

## Compatibility with ljharb/qs

### ✅ Fully Compatible Features

These features work identically to the JavaScript library:

- Basic parsing and stringification
- Nested objects with bracket notation
- Arrays (all formats)
- Dot notation (`allowDots`)
- URL encoding/decoding
- Custom delimiters
- Query prefix handling
- Depth and array limits
- Parameter limits
- Null handling

### ⚠️ Known Differences

#### 1. Empty Key Handling
```go
// JavaScript: qs.parse('a=1&') → { a: '1' }
// Go:         d.Parse("a=1&")  → {"a": "1", "": ""}
// Issue: Creates an extra empty key
```

#### 2. Empty Arrays
```go
// JavaScript: qs.parse('foo[]', {allowEmptyArrays: true}) → { foo: [] }
// Go:         d.Parse("foo[]", WithAllowEmptyArrays(true)) → {"foo": [""]}
// Issue: Contains empty string instead of truly empty
```

#### 3. Pre-encoded Comma Handling
```go
// JavaScript: qs.parse('foo=a%2Cb', {comma: true}) → { foo: 'a,b' }
// Go:         d.Parse("foo=a%2Cb", WithComma(true)) → {"foo": ["a", "b"]}
// Issue: Decodes %2C before checking comma option
```

#### 4. Depth Limit Off-by-One
```go
// JavaScript: qs.parse('a[b][c][d]=e', {depth: 2}) → { a: { b: { '[c][d]': 'e' } } }
// Go:         d.Parse("a[b][c][d]=e", WithDepth(2)) → {"a": {"b": {"c": {"[d]": "e"}}}}
// Issue: Depth counting may be off by one
```

#### 5. Array Limit Key Type
```go
// JavaScript: qs.parse('a[100]=b', {arrayLimit: 20}) → { a: { 100: 'b' } }
// Go:         d.Parse("a[100]=b", WithArrayLimit(20)) → {"a": {"100": "b"}}
// Issue: Returns string "100" instead of integer 100
```

### ❌ Not Supported

#### Language Limitations
- **Charset handling**: Only UTF-8 supported (no ISO-8859-1)
- **Symbol/BigInt types**: These JavaScript types don't exist in Go
- **Custom encoder/decoder functions**: No callback support
- **interpretNumericEntities**: Not implemented
- **Buffer encoding**: Not implemented

#### Missing Features
- **duplicates option**: Only `'combine'` mode supported (not `'first'` or `'last'`)
- **Circular reference detection**: Will cause stack overflow
- **Regex delimiters**: Only string delimiters supported
- **allowPrototypes/plainObjects/allowSparse**: Not applicable to Go

### 🔧 Go-Specific Considerations

#### Type Assertions Required
```go
// JavaScript has dynamic typing
result.user.name  // Works directly

// Go requires type assertions
(*result)["user"].(goqs.QSType)["name"].(string)
```

#### Map Iteration Order
```go
// Go maps have random iteration order
// ALWAYS use WithSort(true) for deterministic output
e := goqs.NewEncoder(goqs.WithSort(true))
```

#### Nil vs Null
```go
// Go has explicit nil (not null/undefined)
result, _ := d.Parse("foo", WithStrictNullHandling(true))
if (*result)["foo"] == nil {
    // Handles nil explicitly
}
```

## Best Practices

### 1. Use Sort for Deterministic Output
```go
e := goqs.NewEncoder(goqs.WithSort(true))
// Ensures consistent query string order
```

### 2. Create Type-Safe Helpers
```go
func ParseToStringMap(query string) (map[string]string, error) {
    d := goqs.NewDecoder()
    result, err := d.Parse(query)
    if err != nil {
        return nil, err
    }

    stringMap := make(map[string]string)
    for k, v := range *result {
        if str, ok := v.(string); ok {
            stringMap[fmt.Sprint(k)] = str
        }
    }
    return stringMap, nil
}
```

### 3. Handle Circular References
```go
// Check for circular references before encoding
// goqs does NOT detect these - will cause crash!

type Node struct {
    Value string
    Next  *Node
}

// Don't do this:
node1 := &Node{Value: "a"}
node2 := &Node{Value: "b", Next: node1}
node1.Next = node2  // Circular!

e.Stringify(node1)  // WILL CRASH!
```

### 4. Default Options for APIs
```go
// Recommended defaults for REST APIs
d := goqs.NewDecoder(
    goqs.WithParameterLimit(100),
    goqs.WithDepth(3),
)

e := goqs.NewEncoder(
    goqs.WithSort(true),
    goqs.WithArrayFormat("brackets"),
)
```

## Use Case Recommendations

| Use Case | Compatibility | Notes |
|----------|---------------|-------|
| REST API query parameters | ✅ Excellent | Recommended for production use |
| URL building | ✅ Excellent | Use `WithSort(true)` |
| Form data parsing | ✅ Excellent | Works reliably |
| Config file parsing | ⚠️ Good | Test your edge cases |
| JS qs drop-in replacement | ⚠️ Good | Has minor differences |
| Non-ASCII charset data | ❌ Limited | Only UTF-8 supported |
| Circular data structures | ❌ Not supported | Will crash |

## Examples

### Building API URLs
```go
func buildAPIURL(base string, params map[string]interface{}) string {
    e := goqs.NewEncoder(
        goqs.WithSort(true),
        goqs.WithSkipNulls(true),
    )
    query, _ := e.Stringify(params)
    if query != "" {
        return base + "?" + query
    }
    return base
}

url := buildAPIURL("https://api.example.com/users", map[string]interface{}{
    "page":   1,
    "limit":  10,
    "sort":   "name",
    "filter": map[string]interface{}{"active": true},
})
// https://api.example.com/users?filter%5Bactive%5D=true&limit=10&page=1&sort=name
```

### Parsing Form Data
```go
func parseFormData(query string) (*goqs.QSType, error) {
    d := goqs.NewDecoder(
        goqs.WithAllowDots(true),
        goqs.WithArrayLimit(100),
    )
    return d.Parse(query)
}

result, _ := parseFormData("user.name=John&user.emails[]=a@example.com&user.emails[]=b@example.com")
// {"user": {"name": "John", "emails": ["a@example.com", "b@example.com"]}}
```

### Complex Nested Structures
```go
data := map[string]interface{}{
    "filters": map[string]interface{}{
        "status": []interface{}{"active", "pending"},
        "date": map[string]interface{}{
            "from": "2024-01-01",
            "to":   "2024-12-31",
        },
    },
    "sort": []interface{}{
        map[string]interface{}{"field": "name", "order": "asc"},
        map[string]interface{}{"field": "date", "order": "desc"},
    },
}

e := goqs.NewEncoder(goqs.WithSort(true))
query, _ := e.Stringify(data)
// filters%5Bdate%5D%5Bfrom%5D=2024-01-01&filters%5Bdate%5D%5Bto%5D=2024-12-31&filters%5Bstatus%5D%5B0%5D=active&filters%5Bstatus%5D%5B1%5D=pending&sort%5B0%5D%5Bfield%5D=name&sort%5B0%5D%5Border%5D=asc&sort%5B1%5D%5Bfield%5D=date&sort%5B1%5D%5Border%5D=desc
```

## Testing

```bash
# Run all tests
go test ./...

# Run specific test suite
go test ./test -run TestStringify

# Run with verbose output
go test -v ./test/...

# Check test coverage
go test -cover ./...
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

### Development Setup

```bash
git clone https://github.com/hlouis/goqs.git
cd goqs
go test ./...
```

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Credits

- Original JavaScript library: [ljharb/qs](https://github.com/ljharb/qs)
- Go port: [hlouis](https://github.com/hlouis)

## Related Projects

- [gorilla/schema](https://github.com/gorilla/schema) - Go struct to form values
- [google/go-querystring](https://github.com/google/go-querystring) - Struct to URL query string

## Changelog

### v0.2.0 (Current)
- ✅ Added full Stringify/Encode functionality
- ✅ Added 350+ comprehensive tests
- ✅ Added all major encoder options
- ✅ Enhanced decoder with missing options
- ✅ Comprehensive documentation

### v0.1.0
- ✅ Initial Parse/Decode implementation
- ✅ Basic nested object and array support
- ✅ Dot notation support

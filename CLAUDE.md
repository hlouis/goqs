# CLAUDE.md - goqs Project Documentation

## Project Overview

**goqs** is a Golang port of the popular JavaScript library [ljharb/qs](https://github.com/ljharb/qs), which provides advanced query string parsing and serialization capabilities with nesting support.

**Author:** Louis Huang (hlouis)
**License:** MIT License (2024)
**Go Version:** 1.22+
**Module:** `github.com/hlouis/goqs`

## Purpose

This library enables developers to parse URL query strings into Go data structures (maps and slices) with support for:
- Nested objects using bracket notation
- Array handling with multiple syntax options
- Flexible configuration options
- Security features like depth and parameter limits

## Architecture

### Core Components

#### 1. **QSType** (`goqs.go`)
```go
type QSType map[interface{}]interface{}
```
The primary data structure representing parsed query strings. Uses `interface{}` for both keys and values to support mixed numeric/string keys and nested structures.

#### 2. **Decoder** (`decoder.go`)
The main parser engine with configurable options. Currently implements **parsing** functionality (decoding query strings to Go objects).

**Key Fields:**
- `tagAlias`: Struct tag name for future struct unmarshaling (default: "qs")
- `allowDots`: Enable dot notation parsing (e.g., `a.b=c` → `{a: {b: "c"}}`)
- `comma`: Parse comma-separated values as arrays
- `delimiter`: Query parameter separator (default: "&")
- `depth`: Maximum nesting depth (default: 5)
- `arrayLimit`: Maximum array index (default: 20)
- `parameterLimit`: Maximum number of parameters to parse (default: 1000)
- `parseArrays`: Enable array syntax parsing (default: true)
- `strictNullHandling`: Treat keys without values as `nil` vs empty string

### File Structure

```
/home/user/goqs/
├── goqs.go                    # Core QSType definition
├── decoder.go                 # Main parsing logic (448 lines)
├── decoder_test.go            # High-level parsing tests
├── decoder_value_test.go      # Unit tests for internal functions
├── go.mod                     # Module definition
├── go.sum                     # Dependency checksums
├── LICENSE                    # MIT License
└── README.md                  # Basic project info
```

## Key Algorithms

### 1. Parse Flow (`decoder.go:92-117`)

```
Parse(input) → QSType
  ├─ parseValues(input) → map[string]interface{}
  │   └─ Splits by delimiter, decodes URI components, handles duplicates
  │
  ├─ parseKeys(key, value) → QSType (for each k,v pair)
  │   ├─ Converts dot notation to brackets if allowDots enabled
  │   ├─ Extracts bracket segments respecting depth limit
  │   └─ Builds nested structure from leaf to root
  │
  ├─ merge(target, newObj) → interface{}
  │   └─ Deep merges parsed objects handling map/array combinations
  │
  └─ objToArray(value) → interface{}
      └─ Converts maps with continuous numeric keys to arrays
```

### 2. Value Parsing (`decoder.go:121-175`)

Handles the first-stage parsing:
- Splits query string by delimiter with parameter limit protection
- Decodes URI components (handles `+` → space, percent encoding)
- Applies comma-splitting if enabled
- Combines duplicate keys based on `duplicates` strategy

### 3. Key Parsing (`decoder.go:191-265`)

Transforms bracket/dot notation into nested structures:
- Converts `a.b.c` → `a[b][c]` if `allowDots` is enabled
- Extracts bracket segments: `a[b][c][d]` → `["a", "[b]", "[c]", "[d]"]`
- Builds structure from innermost to outermost (bottom-up)
- Distinguishes between array indices and object keys

### 4. Merge Algorithm (`decoder.go:324-408`)

Complex recursive function handling 9 type combination scenarios:
- Scalar + Scalar → Array
- Array + Array → Merged array with deep merge at matching indices
- Map + Map → Merged map with recursive value merging
- Map + Array / Array + Map → Converts array to map, then merges
- Scalar + Collection → Appends scalar to collection

### 5. Array Conversion (`decoder.go:422-447`)

Post-processing step that converts maps to arrays when:
- All keys are integers (type check)
- Keys form a continuous sequence: 0, 1, 2, ..., n-1

## Implemented Features

### ✅ Supported from Original qs Library

- ✅ Basic query string parsing
- ✅ Nested object parsing with brackets (`a[b][c]=value`)
- ✅ Array parsing with `[]` and indexed syntax (`a[0]=1&a[1]=2`)
- ✅ Dot notation support via `allowDots` option
- ✅ Comma-separated values via `comma` option
- ✅ Parameter limits for security (`parameterLimit`, `arrayLimit`)
- ✅ Depth limiting to prevent deep nesting attacks
- ✅ Duplicate key handling (combine strategy)
- ✅ URI decoding with proper `+` → space handling
- ✅ Empty array support via `allowEmptyArrays`
- ✅ Strict null handling for valueless parameters

### ❌ Not Supported / Differences

- ❌ **String encoding (stringify)** - Only parsing is implemented
- ❌ **charset/charsetSentinel** - Marked as "not support" in code
- ❌ **allowPrototypes** - Not applicable to Go (no prototype chain)
- ❌ **allowSparse** - Always false (Go doesn't support sparse arrays)
- ❌ **plainObjects** - Not applicable (Go has different type system)
- ❌ **interpretNumericEntities** - Not implemented
- ❌ **Custom decoder functions** - No equivalent to JS's `decoder: utils.decode`

## Usage Examples

### Basic Usage
```go
import "github.com/hlouis/goqs"

d := goqs.NewDecoder()
res, err := d.Parse("foo=bar")
// res = &QSType{"foo": "bar"}

value := (*res)["foo"] // "bar"
```

### Array Parsing
```go
d := goqs.NewDecoder()

// Explicit array syntax
d.Parse("a[]=b&a[]=c")
// → {"a": []interface{}{"b", "c"}}

// Indexed arrays
d.Parse("a[0]=b&a[1]=c")
// → {"a": []interface{}{"b", "c"}}

// Duplicate keys
d.Parse("a=b&a=c")
// → {"a": []interface{}{"b", "c"}}
```

### Nested Objects
```go
d := goqs.NewDecoder()
d.Parse("a[b][c]=d")
// → {"a": map[interface{}]interface{}{"b": map[interface{}]interface{}{"c": "d"}}}
```

### Comma-Separated Values
```go
d := goqs.NewDecoder(goqs.WithComma(true))
d.Parse("colors=red,green,blue")
// → {"colors": []interface{}{"red", "green", "blue"}}
```

### Dot Notation
```go
d := goqs.NewDecoder(goqs.WithAllowDots(true))
d.Parse("user.name=John&user.age=30")
// → {"user": {"name": "John", "age": "30"}}
```

## Testing Strategy

The project follows the original JavaScript library's test suite structure:

### Test Files
1. **decoder_test.go** - Integration tests matching `qs/test/parse.js`
2. **decoder_value_test.go** - Unit tests for internal functions

### Test Coverage Areas
- Simple string parsing (decoder_test.go:35-60)
- Array handling without comma (decoder_test.go:62-72)
- Array handling with comma (decoder_test.go:74-85)
- Dot notation support (decoder_test.go:87-94)
- Value parsing edge cases (decoder_value_test.go:24-40)
- Key parsing transformations (decoder_value_test.go:60-74)

### Running Tests
```bash
go test -v
go test -run TestParseSimpleString
```

## Code Quality Notes

### Strengths
1. **Well-structured separation**: Value parsing → Key parsing → Merging → Post-processing
2. **Defensive programming**: Checks for nil, type assertions, bounds checking
3. **Security-conscious**: Parameter limits, depth limits, array size limits
4. **Test coverage**: Follows upstream test patterns

### Areas for Consideration

1. **Type Assertions** (`decoder.go:262-264`)
   - Uses type assertions without explicit error handling
   - Comment "TODO: handler type assert fail" indicates awareness

2. **Error Handling** (`decoder.go:271-274`)
   - URI decode errors only logged, original value returned
   - May silently proceed with malformed input

3. **Interface{} Usage**
   - Heavy use of `interface{}` requires careful type management
   - Consider generic constraints (Go 1.18+) for type safety improvements

4. **Mixed Key Types**
   - `QSType` allows both string and int keys in same map
   - Could cause confusion when accessing values
   - Example: `map[interface{}]interface{}{0: "a", "0": "b"}` are different keys

## Development History

Based on git log:
```
3ae230f - Add allowDots test
3f14e14 - Fix bug in canBeArray
9fd8d02 - Try test
0814ecf - Can turn number key map to array
56507b3 - Fix bug in numeric key parse
```

Recent work focused on:
- Array conversion logic refinement
- Numeric key handling improvements
- Feature additions (allowDots support)

## Recent Updates (v0.3.0 - 2025-11-05)

### Fixed All 5 Known Differences with JavaScript qs

This update resolves all previously documented behavioral differences between the Go and JavaScript implementations:

#### 1. Empty Key Handling ✅ Fixed
**Location:** `decoder.go:184-188`
- **Issue:** Trailing delimiters (e.g., `"a=1&"`) created empty keys
- **Fix:** Skip empty parts after splitting by delimiter
- **Test:** `TestFix1_EmptyKeyHandling` in `decoder_test.go:107-119`

#### 2. Empty Arrays ✅ Fixed
**Location:** `decoder.go:284`
- **Issue:** Empty arrays contained `[""]` instead of `[]`
- **Fix:** Check for both `nil` and empty string when `allowEmptyArrays` is true
- **Test:** `TestFix2_EmptyArrays` in `decoder_test.go:121-143`

#### 3. Pre-encoded Comma Handling ✅ Fixed
**Location:** `decoder.go:207-221`
- **Issue:** Pre-encoded commas (`%2C`) were decoded before comma-splitting check
- **Fix:** Check for raw commas in encoded string before decoding, then decode each part separately
- **Test:** `TestFix3_PreEncodedCommaHandling` in `decoder_test.go:145-166`

#### 4. Depth Limit Off-by-One ✅ Fixed
**Location:** `decoder.go:269` and `decoder.go:280-283`
- **Issue:** Depth counting was off by one (parsed one more level than JS)
- **Fix:** Use `d.depth-1` in `FindAllStringIndex` since root doesn't count towards depth; handle case where `depth-1=0`
- **Test:** `TestFix4_DepthLimitOffByOne` in `decoder_test.go:168-215`

#### 5. Array Limit Key Type ✅ Fixed
**Location:** `decoder.go:303-306`
- **Issue:** When exceeding `arrayLimit`, keys were strings instead of integers
- **Fix:** Removed `index <= d.arrayLimit` check from condition; always use integer keys for valid integers
- **Test:** `TestFix5_ArrayLimitKeyType` in `decoder_test.go:217-246`

### Implementation Notes

1. **Empty Key Handling:** Simple guard clause to skip empty parts from splitting
2. **Empty Arrays:** Logical OR condition to treat empty string same as nil for empty arrays
3. **Pre-encoded Comma:** Process comma-splitting before URI decoding to preserve encoded commas
4. **Depth Limit:** Adjusted bracket extraction count and added fallback for depth=1 case
5. **Array Limit:** Separated concerns - use integer keys during parsing, check arrayLimit during array conversion

### Test Coverage

All fixes include comprehensive test cases covering:
- Basic functionality matching expected JS behavior
- Edge cases (empty strings, mixed values, boundary conditions)
- Multiple depth/limit scenarios
- Integration with existing features

All existing tests continue to pass, ensuring backward compatibility.

## Current Branch

Development is on: `claude/fix-js-qs-differences-011CUp3dNtFF57sCKkJnwrYz`

## Future Enhancements

Potential areas for expansion:
1. **Stringify/Encode functionality** - Convert Go objects to query strings
2. **Struct marshaling/unmarshaling** - Use `tagAlias` for struct field mapping
3. **Custom type converters** - Allow user-defined type conversions
4. **Better error types** - Define specific error types instead of generic errors
5. **Benchmarking** - Performance tests against standard library alternatives
6. **More configuration options** - Match more features from original qs library

## Dependencies

```go
require github.com/stretchr/testify v1.9.0  // Testing assertions
```

Minimal external dependencies - only testing library required.

## Related Resources

- Original JS library: https://github.com/ljharb/qs
- Original test suite: https://github.com/ljharb/qs/blob/main/test/parse.js
- Go URL package: https://pkg.go.dev/net/url (for comparison with standard library)

---

**Document Version:** 1.0
**Last Updated:** 2025-11-05
**Maintained by:** Claude (AI Assistant)

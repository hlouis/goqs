package goqs

import (
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type Decoder struct {
	tagAlias                 string
	allowDots                bool
	allowEmptyArrays         bool
	allowPrototypes          bool // not support
	allowSparse              bool
	arrayLimit               int
	charset                  string // not support
	charsetSentinel          bool   // not support
	comma                    bool
	decodeDotInKeys          bool
	delimiter                string
	depth                    int
	duplicates               string
	ignoreQueryPrefix        bool
	interpretNumericEntities bool
	parameterLimit           int
	parseArrays              bool
	plainObjects             bool // not support
	strictNullHandling       bool
	// decoder: utils.decode, // not support
}

var defaultDecoder Decoder = Decoder{
	tagAlias:                 "qs",
	allowDots:                false,
	allowEmptyArrays:         false,
	allowPrototypes:          false,
	allowSparse:              false,
	arrayLimit:               20,
	charset:                  "utf-8",
	charsetSentinel:          false,
	comma:                    false,
	decodeDotInKeys:          false,
	delimiter:                "&",
	depth:                    5,
	duplicates:               "combine",
	ignoreQueryPrefix:        false,
	interpretNumericEntities: false,
	parameterLimit:           1000,
	parseArrays:              true,
	plainObjects:             false,
	strictNullHandling:       false,
}

type DecoderOption func(encoder *Decoder)

func WithTagAlias(tagAlias string) DecoderOption {
	return func(d *Decoder) {
		d.tagAlias = tagAlias
	}
}

// WithComma will enable/disable comman in value
// if enabled, comma in value will parse to array
// e.g: v=a,b,c => v:[a,b,c]
// default: false
func WithComma(comma bool) DecoderOption {
	return func(d *Decoder) {
		d.comma = comma
	}
}

func NewDecoder(options ...DecoderOption) *Decoder {
	d := defaultDecoder

	for _, opt := range options {
		opt(&d)
	}

	return &d
}

func (d *Decoder) Parse(input string) (*QSType, error) {
	obj := &QSType{}
	if len(input) == 0 {
		return obj, nil
	}

	// tempObj := d.parseValues(input)

	// Iterate over the keys and setup the new object
	// for k, v := range tempObj {
	// 	newObj := d.parseKeys(k, v)
	// 	obj = merge(obj, newObj)
	// }

	// if d.allowSparse {
	// 	return obj, nil
	// }
	//
	// return compact(obj), nil
	return obj, nil
}

func (d *Decoder) parseValues(str string) map[string]interface{} {
	// clear first query prefix if any
	if d.ignoreQueryPrefix && str[0] == '?' {
		str = str[1:]
	}

	// split and keep limit number in parts
	parts := strings.SplitN(str, d.delimiter, d.parameterLimit)
	if len(parts) == d.parameterLimit {
		last := parts[d.parameterLimit-1]
		last, _, _ = strings.Cut(last, d.delimiter)
		parts[d.parameterLimit-1] = last
	}

	result := make(map[string]interface{})
	for _, part := range parts {
		bracketEqualsPos := strings.Index(part, "]=")
		pos := bracketEqualsPos + 1
		if bracketEqualsPos == -1 {
			pos = strings.Index(part, "=")
		}

		var key string
		var val interface{}
		if pos == -1 {
			key = decodeURI(part)
			if !d.strictNullHandling {
				val = ""
			} else {
				val = nil
			}
		} else {
			key = decodeURI(part[0:pos])
			strv := decodeURI(part[pos+1:])
			if d.comma && strings.Contains(strv, ",") {
				val = split(strv, ",")
			} else {
				val = strv
			}
		}

		if strings.Contains(part, "[]=") && IsArrayLike(val) {
			val = []interface{}{val}
		}

		ev, existing := result[key]
		if existing && d.duplicates == "combine" {
			result[key] = combineValue(ev, val)
		} else {
			result[key] = val
		}
	}

	return result
}

func split(str string, sep string) []interface{} {
	val := strings.Split(str, sep)
	ret := make([]interface{}, len(val))
	for i, v := range val {
		ret[i] = v
	}
	return ret
}

var (
	dotReg     = regexp.MustCompile(`\.([^.[]+)`)
	bracketReg = regexp.MustCompile(`(\[[^[\]]*])`)
)

func (d *Decoder) parseKeys(key string, val interface{}) interface{} {
	if d.allowDots {
		// convert dot string to bracket format (a.b.c => a[b][c])
		key = dotReg.ReplaceAllString(key, "[$1]")
	}

	var keys []string
	// deal with parent (a[b][c][d] => a)
	loc := bracketReg.FindStringIndex(key)
	if d.depth > 0 && loc != nil {
		// push header
		keys = append(keys, key[0:loc[0]])

		// deal with brackets
		locs := bracketReg.FindAllStringIndex(key, d.depth)
		if locs != nil {
			for _, l := range locs {
				keys = append(keys, key[l[0]:l[1]])
			}

			// add any reminder as it is
			lastLoc := locs[len(locs)-1]
			if lastLoc[1] >= len(key)-1 {
				keys = append(keys, fmt.Sprintf("[%v]", key[lastLoc[1]:]))
			}
		}
	} else {
		// if depth is zero or can't find any bracket, add all
		keys = append(keys, key)
	}

	// convert string bracket to map
	leaf := val
	for i := len(keys) - 1; i >= 0; i-- {
		var obj interface{}
		root := keys[i]
		if root == "[]" && d.parseArrays {
			if d.allowEmptyArrays && leaf == nil {
				obj = []interface{}{}
			} else {
				obj = concat([]interface{}{}, leaf)
			}
		} else {
			cleanRoot := root
			if root[0] == '[' && root[len(root)-1] == ']' {
				cleanRoot = root[1 : len(root)-1]
				decodedRoot := cleanRoot
				if d.decodeDotInKeys {
					decodedRoot = strings.ReplaceAll(cleanRoot, "%2E", ".")
				}
				index, err := strconv.ParseInt(decodedRoot, 10, 32)
				if !d.parseArrays && decodedRoot == "" {
					obj = map[int]interface{}{0: leaf}
				} else if err != nil && root != decodedRoot && index > 0 && (d.parseArrays && index <= int64(d.arrayLimit)) {
					obj = map[int64]interface{}{index: leaf}
				} else {
					obj = map[string]interface{}{decodedRoot: leaf}
				}
			}
		}

		leaf = obj
	}

	return leaf
}

func decodeURI(v string) string {
	// in query string replace all + to space
	v = strings.ReplaceAll(v, "+", " ")
	ret, err := url.QueryUnescape(v)
	if err != nil {
		fmt.Printf("url query unescape failed: %v", err)
		return v
	}

	return ret
}

func IsArrayLike(v interface{}) bool {
	k := reflect.TypeOf(v).Kind()
	if k == reflect.Slice || k == reflect.Array {
		return true
	} else {
		return false
	}
}

func combineValue(v1 interface{}, v2 interface{}) []interface{} {
	isArr1 := IsArrayLike(v1)
	isArr2 := IsArrayLike(v2)

	if isArr1 {
		a1 := v1.([]interface{})
		if isArr2 {
			a2 := v2.([]interface{})
			return slices.Concat(a1, a2)
		}
		return append(a1, v2)
	} else {
		if isArr2 {
			a2 := v2.([]interface{})
			return append([]interface{}{v1}, a2...)
		} else {
			return []interface{}{v1, v2}
		}
	}
}

// Concat mult element to target slice
// if source is array or slice, we push elements for source to target
// if source is not array, we push it directly
func concat(target []interface{}, sources ...interface{}) []interface{} {
	for _, s := range sources {
		if IsArrayLike(s) {
			ss := s.([]interface{})
			target = append(target, ss...)
		} else {
			target = append(target, s)
		}
	}
	return target
}

func merge(target interface{}, source interface{}) interface{} {
	if source == nil {
		return target
	}

	tk := reflect.TypeOf(target).Kind()
	sk := reflect.TypeOf(source).Kind()

	// if source is not a map
	if sk != reflect.Map {
		if IsArrayLike(target) {
			tArr := target.([]interface{})
			target = append(tArr, source)
		} else if tk == reflect.Map {
			tMap := target.(map[interface{}]interface{})
			if _, exist := tMap[source]; !exist {
				tMap[source] = true
			}
		} else {
			return []interface{}{target, source}
		}

		return target
	}

	mergeTarget := target
	// target is not exist or not map
	if target == nil || tk != reflect.Map {
		return concat([]interface{}{target}, source)
	}

	// target is a array but source is not
	if tk == reflect.Slice && sk != reflect.Slice {
		tArr := target.([]interface{})
		mergeTarget = arrayToObj(tArr)
	}

	// both array
	if tk == reflect.Slice && sk == reflect.Slice {
		tArr := target.([]interface{})
		tSrc := source.([]interface{})
		for i, item := range tSrc {
			if i < len(tArr) {
				// if both object, deep merge
				// else append
				// targetItem := tArr[i]
				tArr[i] = merge(tArr[i], item)
			} else {
				// TODO: crashed here
				tArr[i] = item
			}
		}
	}

	return mergeTarget
}

func arrayToObj(arr []interface{}) map[int]interface{} {
	ret := make(map[int]interface{}, len(arr))
	for i := 0; i < len(arr); i++ {
		ret[i] = arr[i]
	}
	return ret
}

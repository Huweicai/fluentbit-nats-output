package main

import (
	"errors"
	"strconv"
)

// The presentation style of keys.
type SeparatorStyle int

const (
	_ SeparatorStyle = iota

	// Separate nested key components with dots, e.g. "a.b.1.c.d"
	DotStyle

	// Separate ala Rails, e.g. "a[b][c][1][d]"
	RailsStyle

	// Separate with underscore, e.g. "a_b_1_c_d"
	UnderscoreStyle
)

// Nested input must be a map or slice
var NotValidInputError = errors.New("Not a valid input: map or slice")

// Flatten generates a flat map from a nested one.  The original may include values of type map, slice and scalar,
// but not struct.  Keys in the flat map will be a compound of descending map keys and slice iterations.
// The presentation of keys is set by style.  A prefix is joined to each key.
// copied from: https://github.com/iyacontrol/fluent-bit-clickhouse/blob/master/flatten.go
func Flatten(nested map[interface{}]interface{}, prefix string, style SeparatorStyle) (map[string]interface{}, error) {
	flatmap := make(map[string]interface{})

	err := flatten(true, flatmap, nested, prefix, style)
	if err != nil {
		return nil, err
	}

	return flatmap, nil
}

func flatten(top bool, flatMap map[string]interface{}, nested interface{}, prefix string, style SeparatorStyle) error {
	assign := func(newKey string, v interface{}) error {
		switch v.(type) {
		case map[interface{}]interface{}, []interface{}:
			if err := flatten(false, flatMap, v, newKey, style); err != nil {
				return err
			}
		default:
			flatMap[newKey] = v
		}

		return nil
	}

	switch nested.(type) {
	case map[interface{}]interface{}:
		for k, v := range nested.(map[interface{}]interface{}) {
			newKey := enkey(top, prefix, k.(string), style)
			assign(newKey, v)
		}
	case []interface{}:
		for i, v := range nested.([]interface{}) {
			newKey := enkey(top, prefix, strconv.Itoa(i), style)
			assign(newKey, v)
		}
	default:
		return NotValidInputError
	}

	return nil
}

func enkey(top bool, prefix, subkey string, style SeparatorStyle) string {
	key := prefix

	if top {
		key += subkey
	} else {
		switch style {
		case DotStyle:
			key += "." + subkey
		case RailsStyle:
			key += "[" + subkey + "]"
		case UnderscoreStyle:
			key += "_" + subkey
		}
	}

	return key
}

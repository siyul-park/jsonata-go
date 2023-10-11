package jsonata

import (
	"math"
	"reflect"
)

// isNumeric check if value is a finite number.
func isNumeric(v any) bool {
	switch v.(type) {
	case int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
		return true
	case float32, float64:
		f64 := 0.0
		if f, ok := v.(float64); ok {
			f64 = f
		} else if f, ok := v.(float32); ok {
			f64 = float64(f)
		}
		return !math.IsNaN(f64) && !math.IsInf(f64, 0)
	default:
		return false
	}
}

// isArrayOfStrings returns true if the arg is an array of strings.
func isArrayOfStrings(v any) bool {
	if v == nil {
		return true
	} else if _, ok := v.([]string); ok {
		return true
	} else if arr, ok := v.([]any); ok {
		for _, e := range arr {
			if _, ok := e.(string); !ok {
				return false
			}
		}
		return true
	}
	return false
}

// isArrayOfNumbers true if the arg is an array of numbers
func isArrayOfNumbers(v any) bool {
	if v == nil {
		return true
	}

	switch v.(type) {
	case []any, []int, []uint, []int8, []int16, []int32, []int64, []uint8, []uint16, []uint32, []uint64, []float32, []float64:
		v := reflect.ValueOf(v)
		for i := 0; i < v.Len(); i++ {
			if !isNumeric(v.Index(i).Interface()) {
				return false
			}
		}
		return true
	}

	return false
}

func isNil(v any) bool {
	defer func() { _ = recover() }()
	return v == nil || reflect.ValueOf(v).IsNil()
}

func forEach(v any, f func(any, any) bool) {
	rv := reflect.ValueOf(v)

	for rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			if !f(i, rv.Index(i).Interface()) {
				return
			}
		}
	case reflect.Map:
		iter := rv.MapRange()
		for iter.Next() {
			if !f(iter.Key().Interface(), iter.Value().Interface()) {
				return
			}
		}
	case reflect.Struct:
		for i := 0; i < rv.NumField(); i++ {
			if !rv.Type().Field(i).IsExported() {
				continue
			}
			if !f(rv.Type().Field(i).Name, rv.Field(i).Interface()) {
				return
			}
		}
	}
}

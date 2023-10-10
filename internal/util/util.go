package util

import (
	"math"
	"reflect"
)

// IsNumeric check if value is a finite number.
func IsNumeric(v any) bool {
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

// IsArrayOfStrings returns true if the arg is an array of strings.
func IsArrayOfStrings(v any) bool {
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

// IsArrayOfNumbers true if the arg is an array of numbers
func IsArrayOfNumbers(v any) bool {
	if v == nil {
		return true
	}

	switch v.(type) {
	case []int, []uint, []int8, []int16, []int32, []int64, []uint8, []uint16, []uint32, []uint64, []float32, []float64:
		v := reflect.ValueOf(v)
		for i := 0; i < v.Len(); i++ {
			if !IsNumeric(v.Index(i).Interface()) {
				return false
			}
		}
		return true
	}

	if arr, ok := v.([]any); ok {
		for _, e := range arr {
			if !IsNumeric(e) {
				return false
			}
		}
		return true
	}

	return false
}

package util

import "math"

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

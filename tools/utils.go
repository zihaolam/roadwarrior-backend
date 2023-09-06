package tools

import (
	"strconv"
)

func FloatSliceToStringSlice(floatSlice []float64) []string {
	stringSlice := make([]string, len(floatSlice))

	for i, v := range floatSlice {
		// Convert the float value to a string using strconv.FormatFloat
		stringSlice[i] = strconv.FormatFloat(v, 'f', -1, 64)
	}

	return stringSlice
}

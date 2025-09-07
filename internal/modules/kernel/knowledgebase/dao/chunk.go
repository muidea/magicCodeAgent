package dao

import (
	"fmt"
	"strings"
)

// vec -> DuckDB array literal like [0.1, 0.2, ...]
func floatSliceToArrayLiteral(vec []float64) string {
	if len(vec) == 0 {
		return "ARRAY[]::FLOAT[]"
	}
	var result strings.Builder
	result.WriteString("ARRAY[")
	for i, v := range vec {
		if i > 0 {
			result.WriteString(",")
		}
		result.WriteString(fmt.Sprintf("%g", v))
	}
	result.WriteString("]::FLOAT[]")
	return result.String()
}

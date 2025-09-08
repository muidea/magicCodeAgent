package dao

import (
	"fmt"
	"strings"
)

// vec -> DuckDB array literal like [0.1, 0.2, ...]
func floatSliceToArrayLiteral(vec []float32) string {
	var result strings.Builder
	result.WriteString("[")
	for i, v := range vec {
		if i > 0 {
			result.WriteString(",")
		}
		result.WriteString(fmt.Sprintf("%g", v))
	}
	result.WriteString(fmt.Sprintf("]::FLOAT[%d]", len(vec)))
	return result.String()
}

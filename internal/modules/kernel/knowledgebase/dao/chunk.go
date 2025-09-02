package dao

import (
	"fmt"
	"strings"
)

// vec -> DuckDB array literal like [0.1, 0.2, ...]
func floatSliceToArrayLiteral(vec []float32) string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, v := range vec {
		if i > 0 {
			sb.WriteString(", ")
		}
		// use %g to avoid trailing zeros
		sb.WriteString(fmt.Sprintf("%g", v))
	}
	sb.WriteString("]")
	return sb.String()
}

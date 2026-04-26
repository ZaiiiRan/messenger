package postgresimpl

import (
	"fmt"
	"strings"
)

func appendEqual[T any](sb *strings.Builder, col string, val *T, args *[]any, argPos *int) {
	if val == nil {
		return
	}
	fmt.Fprintf(sb, " AND %s = $%d", col, *argPos)
	*args = append(*args, *val)
	*argPos++
}

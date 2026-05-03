package postgresimpl

import (
	"fmt"
	"strings"
)

func appendILike(sb *strings.Builder, col string, values []string, args *[]any, argPos *int) {
	if len(values) == 0 {
		return
	}
	sb.WriteString(" AND (")
	for i, s := range values {
		if i > 0 {
			sb.WriteString(" OR ")
		}
		fmt.Fprintf(sb, "%s ILIKE $%d", col, *argPos)
		*args = append(*args, fmt.Sprintf("%%%s%%", s))
		*argPos++
	}
	sb.WriteString(")")
}

func appendIPrefix(sb *strings.Builder, col string, values []string, args *[]any, argPos *int) {
	if len(values) == 0 {
		return
	}
	sb.WriteString(" AND (")
	for i, s := range values {
		if i > 0 {
			sb.WriteString(" OR ")
		}
		fmt.Fprintf(sb, "%s ILIKE $%d", col, *argPos)
		*args = append(*args, fmt.Sprintf("%s%%", s))
		*argPos++
	}
	sb.WriteString(")")
}

func appendPartialNames(sb *strings.Builder, firstNameCol, lastNameCol string, values []string, args *[]any, argPos *int) {
	if len(values) == 0 {
		return
	}
	sb.WriteString(" AND (")
	for i, s := range values {
		if i > 0 {
			sb.WriteString(" OR ")
		}
		fmt.Fprintf(sb,
			"(%s || ' ' || %s) ILIKE $%d OR (%s || ' ' || %s) ILIKE $%d",
			lastNameCol, firstNameCol, *argPos,
			firstNameCol, lastNameCol, *argPos,
		)
		*args = append(*args, fmt.Sprintf("%s%%", s))
		*argPos++
	}
	sb.WriteString(")")
}

func appendEqual[T any](sb *strings.Builder, col string, val *T, args *[]any, argPos *int) {
	if val == nil {
		return
	}
	fmt.Fprintf(sb, " AND %s = $%d", col, *argPos)
	*args = append(*args, *val)
	*argPos++
}

func appendAnyEqual[T any](sb *strings.Builder, col string, values []T, args *[]any, argPos *int) {
	if len(values) == 0 {
		return
	}
	fmt.Fprintf(sb, " AND %s = ANY($%d)", col, *argPos)
	*args = append(*args, values)
	*argPos++
}

func appendIn[T any](sb *strings.Builder, col string, values []T, args *[]any, argPos *int) {
	if len(values) == 0 {
		return
	}
	sb.WriteString(" AND (")
	for i, v := range values {
		if i > 0 {
			sb.WriteString(", ")
		}
		fmt.Fprintf(sb, "$%d", *argPos)
		*args = append(*args, v)
		*argPos++
	}
	fmt.Fprintf(sb, ") = ANY(%s)", col)
}

func appendRange[T any](sb *strings.Builder, col string, from, to *T, args *[]any, argPos *int) {
	if from != nil && to != nil {
		fmt.Fprintf(sb, " AND %s BETWEEN $%d AND $%d", col, *argPos, *argPos+1)
		*args = append(*args, *from, *to)
		*argPos += 2
		return
	}
	if from != nil {
		fmt.Fprintf(sb, " AND %s >= $%d", col, *argPos)
		*args = append(*args, *from)
		*argPos++
	}
	if to != nil {
		fmt.Fprintf(sb, " AND %s <= $%d", col, *argPos)
		*args = append(*args, *to)
		*argPos++
	}
}

func appendBool(sb *strings.Builder, col string, val *bool, args *[]any, argPos *int) {
	if val == nil {
		return
	}
	fmt.Fprintf(sb, " AND %s = $%d", col, *argPos)
	*args = append(*args, *val)
	*argPos++
}

func appendIsNotNull(sb *strings.Builder, col string, val *bool) {
	if val == nil {
		return
	}
	if *val {
		fmt.Fprintf(sb, " AND %s IS NOT NULL", col)
	} else {
		fmt.Fprintf(sb, " AND %s IS NULL", col)
	}
}

func appendOrder(sb *strings.Builder, col string, asc bool) {
	dir := "ASC"
	if !asc {
		dir = "DESC"
	}
	fmt.Fprintf(sb, " ORDER BY %s %s", col, dir)
}

func appendLimitOffset(sb *strings.Builder, limit, offset int, args *[]any, argPos *int) {
	fmt.Fprintf(sb, " LIMIT $%d OFFSET $%d", *argPos, *argPos+1)
	*args = append(*args, limit, offset)
	*argPos += 2
}
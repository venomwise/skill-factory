package output

import (
	"fmt"
	"strings"
)

// RenderTable renders rows as a simple aligned text table.
func RenderTable(columns []string, rows [][]interface{}) string {
	if len(columns) == 0 {
		return "(no result)"
	}
	widths := make([]int, len(columns))
	for i, col := range columns {
		widths[i] = len(col)
	}
	stringRows := stringifyRows(rows)
	for _, row := range stringRows {
		for i, value := range row {
			if i < len(widths) && len(value) > widths[i] {
				widths[i] = len(value)
			}
		}
	}
	var b strings.Builder
	writePaddedRow(&b, columns, widths)
	for i, width := range widths {
		if i > 0 {
			b.WriteString("-+-")
		}
		b.WriteString(strings.Repeat("-", width))
	}
	if len(stringRows) == 0 {
		b.WriteString("\n(0 rows)")
		return b.String()
	}
	for _, row := range stringRows {
		b.WriteByte('\n')
		writePaddedRow(&b, row, widths)
	}
	return b.String()
}

func writePaddedRow(b *strings.Builder, row []string, widths []int) {
	for i, width := range widths {
		if i > 0 {
			b.WriteString(" | ")
		}
		value := ""
		if i < len(row) {
			value = row[i]
		}
		b.WriteString(value)
		if padding := width - len(value); padding > 0 {
			b.WriteString(strings.Repeat(" ", padding))
		}
	}
}

func stringifyRows(rows [][]interface{}) [][]string {
	out := make([][]string, len(rows))
	for i, row := range rows {
		out[i] = make([]string, len(row))
		for j, value := range row {
			if value == nil {
				out[i][j] = "NULL"
			} else {
				out[i][j] = fmt.Sprint(value)
			}
		}
	}
	return out
}

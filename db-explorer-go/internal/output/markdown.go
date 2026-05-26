package output

import "strings"

// RenderMarkdown renders rows as a Markdown table.
func RenderMarkdown(columns []string, rows [][]interface{}) string {
	if len(columns) == 0 {
		return "(no result)"
	}
	var b strings.Builder
	b.WriteString("| ")
	b.WriteString(strings.Join(columns, " | "))
	b.WriteString(" |\n| ")
	b.WriteString(strings.Join(repeat("---", len(columns)), " | "))
	b.WriteString(" |")
	for _, row := range stringifyRows(rows) {
		b.WriteString("\n| ")
		cells := make([]string, len(columns))
		for i := range columns {
			if i < len(row) {
				cells[i] = row[i]
			}
		}
		b.WriteString(strings.Join(cells, " | "))
		b.WriteString(" |")
	}
	return b.String()
}

func repeat(value string, count int) []string {
	items := make([]string, count)
	for i := range items {
		items[i] = value
	}
	return items
}

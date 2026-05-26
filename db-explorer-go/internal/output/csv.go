package output

import (
	"bytes"
	"encoding/csv"
)

// RenderCSV renders rows as CSV.
func RenderCSV(columns []string, rows [][]interface{}) (string, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	if len(columns) > 0 {
		if err := writer.Write(columns); err != nil {
			return "", err
		}
	}
	for _, row := range stringifyRows(rows) {
		if err := writer.Write(row); err != nil {
			return "", err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}
	return buf.String(), nil
}

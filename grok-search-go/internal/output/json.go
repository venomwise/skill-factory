package output

import (
	"encoding/json"
	"io"
)

// RenderJSON writes pretty JSON output.
func RenderJSON(w io.Writer, value any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(value)
}

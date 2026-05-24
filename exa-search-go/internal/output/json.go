package output

import (
	"encoding/json"
	"fmt"
	"os"
)

// RenderJSON outputs data as formatted JSON to stdout
func RenderJSON(data OutputData) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}
	return nil
}

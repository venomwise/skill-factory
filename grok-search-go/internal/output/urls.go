package output

import (
	"fmt"
	"io"
)

// RenderURLs writes source URLs, one per line.
func RenderURLs(w io.Writer, sources []Source) error {
	for _, source := range sources {
		if source.URL == "" {
			continue
		}
		if _, err := fmt.Fprintln(w, source.URL); err != nil {
			return err
		}
	}
	return nil
}

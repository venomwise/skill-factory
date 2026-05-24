package output

import "fmt"

// RenderURLs outputs only URLs, one per line, to stdout
func RenderURLs(data OutputData) error {
	for _, result := range data.Results {
		if result.URL != "" {
			fmt.Println(result.URL)
		}
	}
	return nil
}

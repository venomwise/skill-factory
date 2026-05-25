package prompts

import (
	"strings"
	"testing"
)

func TestForMode(t *testing.T) {
	cases := []string{NewsMode, SocialMode, ResearchMode, DocsCompareMode}
	for _, mode := range cases {
		t.Run(mode, func(t *testing.T) {
			prompt := ForMode(mode)
			if prompt == "" {
				t.Fatalf("expected prompt for %s", mode)
			}
			if !strings.Contains(prompt, "Return ONLY a single JSON object") {
				t.Fatalf("prompt for %s does not require JSON output: %q", mode, prompt)
			}
		})
	}
}

func TestDocsComparePromptLabels(t *testing.T) {
	prompt := ForMode(DocsCompareMode)
	for _, label := range []string{"Official docs:", "Community interpretation:", "Agreement/conflict:", "Bottom line:"} {
		if !strings.Contains(prompt, label) {
			t.Fatalf("docs-compare prompt missing label %q", label)
		}
	}
}

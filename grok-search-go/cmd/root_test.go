package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func executeForTest(args ...string) (string, string, error) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&errOut)
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	rootCmd.SetArgs(nil)
	rootCmd.SetOut(nil)
	rootCmd.SetErr(nil)
	return out.String(), errOut.String(), err
}

func TestResearchCommandsRequireQuery(t *testing.T) {
	cases := []string{"news", "social", "research", "docs-compare"}
	for _, name := range cases {
		t.Run(name, func(t *testing.T) {
			newsQuery = ""
			socialQuery = ""
			researchQuery = ""
			docsCompareQuery = ""

			_, _, err := executeForTest(name)
			if err == nil {
				t.Fatalf("expected missing query error")
			}
			if !strings.Contains(err.Error(), "--query is required") {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestVersionRunsWithoutConfigOrAPIKey(t *testing.T) {
	SetVersionInfo("test-version", "test-commit", "test-date", "test-go")

	out, _, err := executeForTest("version")
	if err != nil {
		t.Fatalf("version returned error: %v", err)
	}
	for _, want := range []string{
		"grok-search version test-version",
		"commit: test-commit",
		"built: test-date",
		"go: test-go",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("version output missing %q in %q", want, out)
		}
	}
}

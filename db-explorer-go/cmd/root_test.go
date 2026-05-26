package cmd

import "testing"

func TestRootCommandDefaults(t *testing.T) {
	if rootCmd.Use != "db-explorer" {
		t.Fatalf("unexpected root command use: %s", rootCmd.Use)
	}

	cases := map[string]string{
		"format":  "json",
		"timeout": "30",
	}
	for name, want := range cases {
		flag := rootCmd.PersistentFlags().Lookup(name)
		if flag == nil {
			t.Fatalf("missing flag %q", name)
		}
		if flag.DefValue != want {
			t.Fatalf("flag %q default = %q, want %q", name, flag.DefValue, want)
		}
	}
}

func TestRootCommandHasGlobalFlags(t *testing.T) {
	for _, name := range []string{"config", "profile", "db", "url", "url-env", "format", "timeout", "debug"} {
		if rootCmd.PersistentFlags().Lookup(name) == nil {
			t.Fatalf("missing global flag %q", name)
		}
	}
}

func TestVersionCommandRegistered(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"version"})
	if err != nil {
		t.Fatal(err)
	}
	if cmd == nil || cmd.Use != "version" {
		t.Fatalf("version command not registered")
	}
}

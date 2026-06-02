package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/venomwise/skill-factory/db-explorer/internal/introspect"
	"github.com/venomwise/skill-factory/db-explorer/internal/output"
)

func TestWriteFormattedTabular(t *testing.T) {
	qr := introspect.QueryResult{
		Columns: []string{"id", "email"},
		Rows:    [][]interface{}{{1, "a@x.com"}, {2, "b@x.com"}},
	}
	cases := map[string]struct {
		format string
		want   []string
	}{
		"table":    {format: "table", want: []string{"id", "email", "a@x.com"}},
		"markdown": {format: "markdown", want: []string{"| id | email |", "a@x.com"}},
		"csv":      {format: "csv", want: []string{"id,email", "a@x.com"}},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			format = tc.format
			defer func() { format = "json" }()

			var buf bytes.Buffer
			writer := output.NewWriter(&buf)
			if err := writeFormatted(writer, "query", commandResult{Data: qr}, time.Now()); err != nil {
				t.Fatalf("writeFormatted returned error: %v", err)
			}
			got := buf.String()
			if strings.Contains(got, "schema_version") {
				t.Fatalf("%s output must not be a JSON envelope: %q", tc.format, got)
			}
			for _, want := range tc.want {
				if !strings.Contains(got, want) {
					t.Fatalf("%s output missing %q: %q", tc.format, want, got)
				}
			}
		})
	}
}

func TestWriteFormattedNonTabularErrors(t *testing.T) {
	format = "table"
	defer func() { format = "json" }()

	var buf bytes.Buffer
	writer := output.NewWriter(&buf)
	// schema/tables payloads are not QueryResult; they cannot be tabularized.
	data := map[string]any{"tables": []string{"users"}}
	err := writeFormatted(writer, "tables", commandResult{Data: data}, time.Now())
	if err != errCommandFailed {
		t.Fatalf("expected errCommandFailed, got %v", err)
	}

	var env output.Envelope
	if jerr := json.Unmarshal(buf.Bytes(), &env); jerr != nil {
		t.Fatalf("error output is not a JSON envelope: %v (%q)", jerr, buf.String())
	}
	if env.OK || env.Error == nil {
		t.Fatalf("expected failed envelope with error, got %+v", env)
	}
	if env.Error.Code != "FORMAT_UNSUPPORTED" {
		t.Fatalf("expected FORMAT_UNSUPPORTED, got %q", env.Error.Code)
	}
}

func TestValidateFormat(t *testing.T) {
	cases := map[string]struct {
		format  string
		command string
		wantErr bool
	}{
		"json tables":       {format: "json", command: "tables"},
		"table query":       {format: "table", command: "query"},
		"markdown data":     {format: "markdown", command: "data"},
		"csv data":          {format: "csv", command: "data"},
		"table tables":      {format: "table", command: "tables", wantErr: true},
		"markdown schema":   {format: "markdown", command: "schema", wantErr: true},
		"unknown query fmt": {format: "pretty", command: "query", wantErr: true},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			format = tc.format
			defer func() { format = "json" }()

			err := validateFormat(tc.command)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

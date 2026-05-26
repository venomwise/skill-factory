package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestWriteSuccessEnvelope(t *testing.T) {
	var buf bytes.Buffer
	writer := NewWriter(&buf)
	err := writer.WriteSuccess("tables", "sqlite", "local", map[string]any{"tables": []string{"users"}}, Meta{DurationMS: 12, Truncated: true})
	if err != nil {
		t.Fatal(err)
	}

	var got Envelope
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.SchemaVersion != "1" || !got.OK || got.Command != "tables" || got.DB != "sqlite" || got.Profile != "local" {
		t.Fatalf("unexpected envelope: %+v", got)
	}
	if got.Error != nil {
		t.Fatalf("success envelope has error: %+v", got.Error)
	}
	if got.Meta.DurationMS != 12 || !got.Meta.Truncated {
		t.Fatalf("unexpected meta: %+v", got.Meta)
	}
}

func TestWriteErrorEnvelope(t *testing.T) {
	var buf bytes.Buffer
	writer := NewWriter(&buf)
	err := writer.WriteError("query", ErrorBody{Code: "SQL_NOT_READONLY", Message: "Only read-only SQL is allowed"}, Meta{DurationMS: 1})
	if err != nil {
		t.Fatal(err)
	}

	var got Envelope
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.SchemaVersion != "1" || got.OK || got.Command != "query" || got.Error == nil {
		t.Fatalf("unexpected envelope: %+v", got)
	}
	if got.Error.Code != "SQL_NOT_READONLY" || got.Error.Message == "" {
		t.Fatalf("unexpected error: %+v", got.Error)
	}
}

func TestRenderers(t *testing.T) {
	columns := []string{"id", "name"}
	rows := [][]interface{}{{1, "Alice"}, {2, nil}}
	if out := RenderTable(columns, rows); !strings.Contains(out, "Alice") || !strings.Contains(out, "NULL") {
		t.Fatalf("unexpected table output: %s", out)
	}
	if out := RenderMarkdown(columns, rows); !strings.HasPrefix(out, "| id | name |") || !strings.Contains(out, "Alice") {
		t.Fatalf("unexpected markdown output: %s", out)
	}
	csvOut, err := RenderCSV(columns, rows)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(csvOut, "id,name") || !strings.Contains(csvOut, "Alice") {
		t.Fatalf("unexpected csv output: %s", csvOut)
	}
}

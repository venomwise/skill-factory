package output

import (
	"encoding/json"
	"io"
)

const SchemaVersion = "1"

// Meta contains execution metadata shared by all responses.
type Meta struct {
	DurationMS int64 `json:"duration_ms,omitempty"`
	Truncated  bool  `json:"truncated"`
}

// ErrorBody is the structured error payload for failed commands.
type ErrorBody struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

// Envelope is the stable top-level JSON response object.
type Envelope struct {
	SchemaVersion string     `json:"schema_version"`
	OK            bool       `json:"ok"`
	Command       string     `json:"command"`
	DB            string     `json:"db,omitempty"`
	Profile       string     `json:"profile,omitempty"`
	Data          any        `json:"data,omitempty"`
	Error         *ErrorBody `json:"error,omitempty"`
	Meta          Meta       `json:"meta"`
}

// Writer emits command responses.
type Writer struct {
	out io.Writer
}

// NewWriter creates a response writer.
func NewWriter(out io.Writer) *Writer {
	return &Writer{out: out}
}

// WriteSuccess writes a successful JSON envelope.
func (w *Writer) WriteSuccess(command, dbType, profile string, data any, meta Meta) error {
	return w.write(Envelope{
		SchemaVersion: SchemaVersion,
		OK:            true,
		Command:       command,
		DB:            dbType,
		Profile:       profile,
		Data:          data,
		Meta:          meta,
	})
}

// WriteError writes a failed JSON envelope.
func (w *Writer) WriteError(command string, errBody ErrorBody, meta Meta) error {
	return w.write(Envelope{
		SchemaVersion: SchemaVersion,
		OK:            false,
		Command:       command,
		Error:         &errBody,
		Meta:          meta,
	})
}

func (w *Writer) write(envelope Envelope) error {
	encoder := json.NewEncoder(w.out)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(envelope)
}

// WriteRaw writes a pre-rendered string (e.g. table/markdown/csv) followed by a
// trailing newline. It bypasses the JSON envelope for human-oriented formats.
func (w *Writer) WriteRaw(s string) error {
	if _, err := io.WriteString(w.out, s); err != nil {
		return err
	}
	_, err := io.WriteString(w.out, "\n")
	return err
}

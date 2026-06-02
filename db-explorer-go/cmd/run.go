package cmd

import (
	"context"
	"errors"
	"time"

	"github.com/spf13/cobra"
	"github.com/venomwise/skill-factory/db-explorer/internal/config"
	"github.com/venomwise/skill-factory/db-explorer/internal/db"
	"github.com/venomwise/skill-factory/db-explorer/internal/introspect"
	"github.com/venomwise/skill-factory/db-explorer/internal/output"
	"github.com/venomwise/skill-factory/db-explorer/internal/safety"
)

var errCommandFailed = errors.New("command failed")

type commandResult struct {
	Data      any
	Meta      output.Meta
	Truncated bool
}

func runWithConnection(cmd *cobra.Command, command string, run func(context.Context, *db.Connection, config.ResolvedConnection) (commandResult, error)) error {
	started := time.Now()
	writer := output.NewWriter(cmd.OutOrStdout())
	if err := validateFormat(command); err != nil {
		return writeCommandError(writer, command, err, started)
	}
	ctx, cancel := context.WithTimeout(cmd.Context(), time.Duration(timeout)*time.Second)
	defer cancel()

	resolved, err := config.ResolveConnection(config.Options{
		ConfigPath: cfgFile,
		ProfileID:  profile,
		DB:         dbType,
		URL:        url,
		URLEnv:     urlEnv,
	})
	if err != nil {
		return writeCommandError(writer, command, err, started)
	}

	conn, err := db.Open(ctx, resolved.DB, resolved.URL)
	if err != nil {
		return writeCommandError(writer, command, err, started)
	}
	defer conn.Close()

	result, err := run(ctx, conn, resolved)
	if err != nil {
		return writeCommandError(writer, command, err, started)
	}
	result.Meta.DurationMS = elapsedMS(started)
	if result.Truncated {
		result.Meta.Truncated = true
	}
	if format != "json" {
		if err := writeFormatted(writer, command, result, started); err != nil {
			return err
		}
		return nil
	}
	if err := writer.WriteSuccess(command, string(resolved.DB), resolved.Profile, result.Data, result.Meta); err != nil {
		return err
	}
	return nil
}

func validateFormat(command string) error {
	switch format {
	case "json":
		return nil
	case "table", "markdown", "csv":
		if command == "query" || command == "data" {
			return nil
		}
		return &safety.Error{
			Code:   "FORMAT_UNSUPPORTED",
			Detail: "--format " + format + " is only supported for query and data commands; use the default JSON output for " + command,
		}
	default:
		return &safety.Error{
			Code:   "FORMAT_UNSUPPORTED",
			Detail: "unknown --format value: " + format,
		}
	}
}

// writeFormatted renders tabular command results (query/data) in the requested
// non-JSON format. Commands whose payload is not row/column shaped cannot be
// represented this way and return a structured FORMAT_UNSUPPORTED error.
func writeFormatted(writer *output.Writer, command string, result commandResult, started time.Time) error {
	qr, ok := result.Data.(introspect.QueryResult)
	if !ok {
		return writeCommandError(writer, command, &safety.Error{
			Code:   "FORMAT_UNSUPPORTED",
			Detail: "--format " + format + " is only supported for query and data commands; use the default JSON output for " + command,
		}, started)
	}

	var (
		rendered string
		err      error
	)
	switch format {
	case "table":
		rendered = output.RenderTable(qr.Columns, qr.Rows)
	case "markdown":
		rendered = output.RenderMarkdown(qr.Columns, qr.Rows)
	case "csv":
		rendered, err = output.RenderCSV(qr.Columns, qr.Rows)
	default:
		return writeCommandError(writer, command, &safety.Error{
			Code:   "FORMAT_UNSUPPORTED",
			Detail: "unknown --format value: " + format,
		}, started)
	}
	if err != nil {
		return writeCommandError(writer, command, err, started)
	}
	return writer.WriteRaw(rendered)
}

func writeValidationError(cmd *cobra.Command, command string, err error) error {
	writer := output.NewWriter(cmd.OutOrStdout())
	return writeCommandError(writer, command, err, time.Now())
}

func writeCommandError(writer *output.Writer, command string, err error, started time.Time) error {
	body := errorBody(err)
	if writeErr := writer.WriteError(command, body, output.Meta{DurationMS: elapsedMS(started)}); writeErr != nil {
		return writeErr
	}
	return errCommandFailed
}

func errorBody(err error) output.ErrorBody {
	var cfgErr *config.Error
	if errors.As(err, &cfgErr) {
		return output.ErrorBody{Code: cfgErr.Code, Message: config.MaskSecrets(cfgErr.Detail)}
	}
	var dbErr *db.Error
	if errors.As(err, &dbErr) {
		return output.ErrorBody{Code: dbErr.Code, Message: config.MaskSecrets(dbErr.Detail)}
	}
	var safetyErr *safety.Error
	if errors.As(err, &safetyErr) {
		return output.ErrorBody{Code: safetyErr.Code, Message: safetyErr.Detail}
	}
	return output.ErrorBody{Code: "QUERY_FAILED", Message: config.MaskSecrets(err.Error())}
}

func elapsedMS(started time.Time) int64 {
	return time.Since(started).Milliseconds()
}

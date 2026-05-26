package cmd

import (
	"context"
	"errors"
	"time"

	"github.com/spf13/cobra"
	"github.com/venomwise/skill-factory/db-explorer/internal/config"
	"github.com/venomwise/skill-factory/db-explorer/internal/db"
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
	if err := writer.WriteSuccess(command, string(resolved.DB), resolved.Profile, result.Data, result.Meta); err != nil {
		return err
	}
	return nil
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

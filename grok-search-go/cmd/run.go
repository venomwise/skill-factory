package cmd

import (
	"context"
	"errors"
	"path/filepath"
	"time"

	"github.com/venomwise/skill-factory/grok-search/internal/client"
	"github.com/venomwise/skill-factory/grok-search/internal/config"
	"github.com/venomwise/skill-factory/grok-search/internal/cooldown"
	"github.com/venomwise/skill-factory/grok-search/internal/output"
	"github.com/venomwise/skill-factory/grok-search/internal/prompts"
)

func runResearchMode(mode, query string) error {
	started := time.Now()

	cfg, err := config.Load(config.Options{
		ConfigPath:       cfgFile,
		APIKey:           apiKey,
		BaseURL:          baseURL,
		Model:            model,
		Timeout:          timeout,
		ProfileID:        profileID,
		ExtraBodyJSON:    extraBodyJSON,
		ExtraHeadersJSON: extraHeadersJSON,
	})
	if err != nil {
		return renderRuntimeError(output.FromError(err), err)
	}

	statePath := resolveCooldownPath(cfg.Cooldown.StateFile)
	state := cooldown.State{Profiles: map[string]cooldown.Entry{}}
	stateChanged := false
	if cfg.Cooldown.Enabled {
		var err error
		state, err = cooldown.LoadState(statePath)
		if err != nil {
			return renderRuntimeError(output.NewErrorResponse("request_failed", err.Error()), err)
		}
		now := time.Now()
		stateChanged = cooldown.PruneExpired(&state, now)
	}

	attempts := []output.Attempt{}
	clientImpl := client.New(cfg.Timeout)
	var lastErr error
	var lastReqErr *client.RequestError

	for idx, profile := range cfg.Profiles {
		if entry, active := cooldown.Active(state, profile.ID, time.Now()); cfg.Cooldown.Enabled && active && !ignoreCooldown {
			attempts = append(attempts, output.Attempt{
				ProfileID:        profile.ID,
				OK:               false,
				Cooldown:         true,
				RemainingSeconds: max(1, int(entry.Until-float64(time.Now().Unix()))),
				UntilText:        entry.UntilText,
				Detail:           entry.Reason,
			})
			continue
		}

		modelName := profile.Model
		if modelName == "" {
			modelName = cfg.Model
		}
		req := client.ResearchRequest{
			BaseURL:       profile.BaseURL,
			APIKey:        profile.APIKey,
			Model:         modelName,
			Mode:          mode,
			Query:         query,
			SystemPrompt:  prompts.ForMode(mode),
			Timeout:       cfg.Timeout,
			ExtraBody:     cfg.ExtraBody,
			ExtraHeaders:  cfg.ExtraHeaders,
			ProfileID:     profile.ID,
			ProfileSource: profile.Source,
		}

		resp, err := clientImpl.DoResearch(context.Background(), req)
		if err == nil {
			attempts = append(attempts, output.Attempt{ProfileID: profile.ID, OK: true})
			if cooldown.Clear(&state, profile.ID) {
				stateChanged = true
			}
			if stateChanged && cfg.Cooldown.Enabled {
				if err := cooldown.SaveState(statePath, state); err != nil {
					return renderRuntimeError(output.NewErrorResponse("request_failed", err.Error()), err)
				}
			}
			if resp.Model != "" {
				modelName = resp.Model
			}
			return renderSuccess(mode, query, profile, attempts, cfg, resp, modelName, started)
		}

		lastErr = err
		var reqErr *client.RequestError
		if errors.As(err, &reqErr) {
			lastReqErr = reqErr
			failover := client.ShouldFailover(reqErr.StatusCode, reqErr.Detail)
			cooldownSeconds := 0
			var cooldownEntry cooldown.Entry
			if failover {
				cooldownSeconds = cooldown.SecondsForFailure(reqErr.StatusCode, reqErr.Detail, cfg.Cooldown)
				if cooldownSeconds > 0 {
					cooldownEntry = cooldown.Set(&state, profile.ID, cooldownSeconds, reqErr.Detail, reqErr.StatusCode, time.Now())
					stateChanged = true
					if err := cooldown.SaveState(statePath, state); err != nil {
						return renderRuntimeError(output.NewErrorResponse("request_failed", err.Error()), err)
					}
				}
			}
			attempts = append(attempts, output.Attempt{
				ProfileID:       profile.ID,
				OK:              false,
				Status:          reqErr.StatusCode,
				Failover:        failover,
				CooldownSeconds: cooldownSeconds,
				UntilText:       cooldownEntry.UntilText,
				Detail:          shorten(reqErr.Detail, 240),
			})
			if failover && idx < len(cfg.Profiles)-1 {
				continue
			}
			break
		}

		attempts = append(attempts, output.Attempt{
			ProfileID: profile.ID,
			OK:        false,
			Failover:  false,
			Detail:    shorten(err.Error(), 240),
		})
		break
	}

	if allCooling(attempts) {
		resp := output.ErrorResponse{
			OK:                false,
			Error:             "all_profiles_in_cooldown",
			Detail:            "All configured profiles are cooling down. Wait for cooldown expiry or retry with --ignore-cooldown.",
			Attempts:          attempts,
			CooldownStateFile: statePath,
			ConfigPath:        cfg.ConfigPath,
			ConfigPaths:       cfg.ConfigPaths,
			BaseURL:           cfg.BaseURL,
			Model:             cfg.Model,
			ElapsedMS:         int(time.Since(started).Milliseconds()),
		}
		return renderRuntimeError(resp, output.NewCommandError(resp.Error, resp.Detail, nil))
	}

	if lastReqErr != nil && client.ShouldFailover(lastReqErr.StatusCode, lastReqErr.Detail) && len(attempts) > 1 {
		resp := output.ErrorResponse{
			OK:                false,
			Error:             "all_profiles_failed",
			Detail:            lastReqErr.Detail,
			FailoverExhausted: true,
			Attempts:          attempts,
			CooldownStateFile: statePath,
			ConfigPath:        cfg.ConfigPath,
			ConfigPaths:       cfg.ConfigPaths,
			BaseURL:           cfg.BaseURL,
			Model:             cfg.Model,
			ElapsedMS:         int(time.Since(started).Milliseconds()),
		}
		return renderRuntimeError(resp, output.NewCommandError(resp.Error, resp.Detail, lastErr))
	}
	if lastErr != nil {
		resp := output.FromError(lastErr)
		resp.Attempts = attempts
		resp.ConfigPath = cfg.ConfigPath
		resp.ConfigPaths = cfg.ConfigPaths
		resp.BaseURL = cfg.BaseURL
		resp.Model = cfg.Model
		resp.ElapsedMS = int(time.Since(started).Milliseconds())
		return renderRuntimeError(resp, lastErr)
	}

	resp := output.ErrorResponse{OK: false, Error: "all_profiles_failed", Attempts: attempts}
	return renderRuntimeError(resp, output.NewCommandError(resp.Error, resp.Detail, nil))
}

func renderSuccess(mode, query string, profile config.ResolvedProfile, attempts []output.Attempt, cfg *config.ResolvedConfig, resp *client.ResearchResponse, modelName string, started time.Time) error {
	result := output.Result{
		OK:            true,
		Mode:          mode,
		Query:         query,
		ProfileID:     profile.ID,
		ProfileSource: profile.Source,
		Attempts:      attempts,
		ConfigPath:    cfg.ConfigPath,
		ConfigPaths:   cfg.ConfigPaths,
		BaseURL:       profile.BaseURL,
		Model:         modelName,
		Content:       resp.AssistantContent.Content,
		Sources:       convertSources(resp.AssistantContent.Sources),
		Raw:           resp.AssistantContent.Raw,
		Usage:         resp.Usage,
		ElapsedMS:     int(time.Since(started).Milliseconds()),
	}

	switch getOutputFormat() {
	case "plain":
		return output.RenderPlain(rootCmd.OutOrStdout(), result)
	case "urls":
		return output.RenderURLs(rootCmd.OutOrStdout(), result.Sources)
	default:
		return output.RenderJSON(rootCmd.OutOrStdout(), result)
	}
}

func renderRuntimeError(response output.ErrorResponse, err error) error {
	if getOutputFormat() == "plain" {
		_ = output.RenderPlainError(rootCmd.OutOrStdout(), response)
	} else {
		_ = output.RenderJSON(rootCmd.OutOrStdout(), response)
	}
	if cmdErr, ok := err.(*output.CommandError); ok {
		return cmdErr
	}
	return output.NewCommandError(response.Error, response.Detail, err)
}

func convertSources(sources []client.Source) []output.Source {
	out := make([]output.Source, 0, len(sources))
	for _, source := range sources {
		out = append(out, output.Source{URL: source.URL, Title: source.Title, Snippet: source.Snippet})
	}
	return out
}

func resolveCooldownPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Clean(path)
}

func allCooling(attempts []output.Attempt) bool {
	if len(attempts) == 0 {
		return false
	}
	for _, attempt := range attempts {
		if !attempt.Cooldown {
			return false
		}
	}
	return true
}

func shorten(text string, limit int) string {
	if len(text) <= limit {
		return text
	}
	if limit <= 3 {
		return text[:limit]
	}
	return text[:limit-3] + "..."
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

package cooldown

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Entry is one profile cooldown record.
type Entry struct {
	Until     float64 `json:"until"`
	UntilText string  `json:"untilText"`
	Seconds   int     `json:"seconds"`
	Reason    string  `json:"reason"`
	Status    int     `json:"status,omitempty"`
	SetAt     float64 `json:"setAt"`
	SetAtText string  `json:"setAtText"`
}

// State is the persisted cooldown state.
type State struct {
	Profiles map[string]Entry `json:"profiles"`
}

// LoadState reads cooldown state from path.
func LoadState(path string) (State, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return State{Profiles: map[string]Entry{}}, nil
	}
	if err != nil {
		return State{Profiles: map[string]Entry{}}, err
	}
	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return State{Profiles: map[string]Entry{}}, nil
	}
	if state.Profiles == nil {
		state.Profiles = map[string]Entry{}
	}
	return state, nil
}

// SaveState writes cooldown state to path.
func SaveState(path string, state State) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// PruneExpired removes expired or malformed entries.
func PruneExpired(state *State, now time.Time) bool {
	ensureProfiles(state)
	changed := false
	for id, entry := range state.Profiles {
		if entry.Until <= float64(now.Unix()) {
			delete(state.Profiles, id)
			changed = true
		}
	}
	return changed
}

// Active returns the active cooldown entry for profileID, if any.
func Active(state State, profileID string, now time.Time) (Entry, bool) {
	entry, ok := state.Profiles[profileID]
	if !ok || entry.Until <= float64(now.Unix()) {
		return Entry{}, false
	}
	return entry, true
}

// Clear removes profileID from cooldown state.
func Clear(state *State, profileID string) bool {
	ensureProfiles(state)
	if _, ok := state.Profiles[profileID]; !ok {
		return false
	}
	delete(state.Profiles, profileID)
	return true
}

// Set writes a cooldown entry for profileID.
func Set(state *State, profileID string, seconds int, reason string, status int, now time.Time) Entry {
	ensureProfiles(state)
	if seconds < 0 {
		seconds = 0
	}
	until := now.Add(time.Duration(seconds) * time.Second)
	entry := Entry{
		Until:     float64(until.Unix()),
		UntilText: formatUTC(until),
		Seconds:   seconds,
		Reason:    reason,
		Status:    status,
		SetAt:     float64(now.Unix()),
		SetAtText: formatUTC(now),
	}
	state.Profiles[profileID] = entry
	return entry
}

func ensureProfiles(state *State) {
	if state.Profiles == nil {
		state.Profiles = map[string]Entry{}
	}
}

func formatUTC(t time.Time) string {
	return t.UTC().Format("2006-01-02 15:04:05 UTC")
}

package sensor

import (
	"encoding/json"
	"fmt"
	"os"
)

// Result is the top-level PRTG Script v2 JSON output (Schema v3).
type Result struct {
	Version  int       `json:"version"`
	Status   string    `json:"status"`
	Message  string    `json:"message"`
	Channels []Channel `json:"channels,omitempty"`
}

// Channel represents a single PRTG sensor channel.
type Channel struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	Type        string       `json:"type"`
	Kind        string       `json:"kind"`
	Value       interface{}  `json:"value"`
	DisplayUnit string       `json:"display_unit,omitempty"`
	Limits      *Limits      `json:"limits,omitempty"`
}

// Limits defines warning/error thresholds for a channel.
type Limits struct {
	Warning *Threshold `json:"warning,omitempty"`
	Error   *Threshold `json:"error,omitempty"`
}

// Threshold defines upper and/or lower bounds.
type Threshold struct {
	Upper *float64 `json:"upper,omitempty"`
	Lower *float64 `json:"lower,omitempty"`
}

// NewResult creates a successful result with the given message.
func NewResult(message string) *Result {
	return &Result{
		Version: 3,
		Status:  "ok",
		Message: message,
	}
}

// NewErrorResult creates an error result with the given message.
func NewErrorResult(message string) *Result {
	return &Result{
		Version: 3,
		Status:  "error",
		Message: message,
	}
}

// AddIntChannel adds an integer channel with a standard kind.
func (r *Result) AddIntChannel(id int, name string, kind string, value int64) {
	r.Channels = append(r.Channels, Channel{
		ID:    id,
		Name:  name,
		Type:  "integer",
		Kind:  kind,
		Value: value,
	})
}

// AddFloatChannel adds a float channel with a standard kind.
func (r *Result) AddFloatChannel(id int, name string, kind string, value float64) {
	r.Channels = append(r.Channels, Channel{
		ID:    id,
		Name:  name,
		Type:  "float",
		Kind:  kind,
		Value: roundFloat(value, 2),
	})
}

// AddFloatChannelWithLimits adds a float channel with warning/error limits.
func (r *Result) AddFloatChannelWithLimits(id int, name, kind string, value, warnUpper, errUpper float64) {
	r.Channels = append(r.Channels, Channel{
		ID:    id,
		Name:  name,
		Type:  "float",
		Kind:  kind,
		Value: roundFloat(value, 2),
		Limits: &Limits{
			Warning: &Threshold{Upper: &warnUpper},
			Error:   &Threshold{Upper: &errUpper},
		},
	})
}

// AddCustomFloatChannel adds a float channel with kind "custom" and a display unit.
func (r *Result) AddCustomFloatChannel(id int, name, displayUnit string, value float64) {
	r.Channels = append(r.Channels, Channel{
		ID:          id,
		Name:        name,
		Type:        "float",
		Kind:        "custom",
		Value:       roundFloat(value, 2),
		DisplayUnit: displayUnit,
	})
}

// Print serializes the result as JSON to stdout and exits with code 0.
func (r *Result) Print() {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(r)
}

// PrintAndExit serializes the result and exits. Always exit 0 per PRTG convention.
func (r *Result) PrintAndExit() {
	r.Print()
	os.Exit(0)
}

// Fatal prints an error result and exits with code 0.
func Fatal(format string, args ...interface{}) {
	NewErrorResult(fmt.Sprintf(format, args...)).PrintAndExit()
}

func roundFloat(val float64, precision int) float64 {
	p := 1.0
	for i := 0; i < precision; i++ {
		p *= 10
	}
	return float64(int64(val*p+0.5)) / p
}

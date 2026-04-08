package params

import (
	"fmt"
	"strings"
)

// Params holds the parsed sensor parameters.
type Params struct {
	Endpoint   string
	APIToken   string
	Scope      string
	AuthMode   string
	APIVersion string
	Insecure   bool
	Warning    float64
	Critical   float64
	Volumes     []string
	SnapAgeDays int
}

// Parse parses a parameter string (as received on stdin from PRTG Script v2).
// The string is split shell-style (respecting double quotes) into key-value pairs.
func Parse(input string) (*Params, error) {
	args := shellSplit(strings.TrimSpace(input))

	p := &Params{
		AuthMode:    "session",
		Insecure:    true,
		Warning:     80,
		Critical:    90,
		SnapAgeDays: 30,
	}

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--endpoint":
			i++
			if i >= len(args) {
				return nil, fmt.Errorf("--endpoint requires a value")
			}
			p.Endpoint = args[i]
		case "--apitoken":
			i++
			if i >= len(args) {
				return nil, fmt.Errorf("--apitoken requires a value")
			}
			p.APIToken = args[i]
		case "--scope":
			i++
			if i >= len(args) {
				return nil, fmt.Errorf("--scope requires a value")
			}
			p.Scope = args[i]
		case "--auth-mode":
			i++
			if i >= len(args) {
				return nil, fmt.Errorf("--auth-mode requires a value")
			}
			p.AuthMode = args[i]
		case "--apiversion":
			i++
			if i >= len(args) {
				return nil, fmt.Errorf("--apiversion requires a value")
			}
			p.APIVersion = args[i]
		case "--insecure":
			p.Insecure = true
		case "--secure":
			p.Insecure = false
		case "--warning":
			i++
			if i >= len(args) {
				return nil, fmt.Errorf("--warning requires a value")
			}
			val, err := parseFloat(args[i])
			if err != nil {
				return nil, fmt.Errorf("--warning: %w", err)
			}
			p.Warning = val
		case "--critical":
			i++
			if i >= len(args) {
				return nil, fmt.Errorf("--critical requires a value")
			}
			val, err := parseFloat(args[i])
			if err != nil {
				return nil, fmt.Errorf("--critical: %w", err)
			}
			p.Critical = val
		case "--volumes":
			i++
			if i >= len(args) {
				return nil, fmt.Errorf("--volumes requires a value")
			}
			p.Volumes = strings.Split(args[i], ",")
		case "--snap-age-days":
			i++
			if i >= len(args) {
				return nil, fmt.Errorf("--snap-age-days requires a value")
			}
			val, err := parseInt(args[i])
			if err != nil {
				return nil, fmt.Errorf("--snap-age-days: %w", err)
			}
			p.SnapAgeDays = val
		default:
			// Ignore unknown flags for forward compatibility
		}
	}

	if p.Endpoint == "" {
		return nil, fmt.Errorf("--endpoint is required")
	}
	if p.APIToken == "" {
		return nil, fmt.Errorf("--apitoken is required")
	}
	if p.Scope == "" {
		return nil, fmt.Errorf("--scope is required")
	}

	validScopes := map[string]bool{
		"capacity":    true,
		"performance": true,
		"hardware":    true,
		"volumes":     true,
		"snapshots":   true,
	}
	if !validScopes[p.Scope] {
		return nil, fmt.Errorf("invalid scope %q (valid: capacity, performance, hardware, volumes)", p.Scope)
	}

	return p, nil
}

// shellSplit splits a string respecting double-quoted segments, similar to shlex.split.
func shellSplit(s string) []string {
	var args []string
	var current strings.Builder
	inQuotes := false

	for i := 0; i < len(s); i++ {
		ch := s[i]
		switch {
		case ch == '"':
			inQuotes = !inQuotes
		case ch == '\\' && i+1 < len(s) && inQuotes:
			i++
			current.WriteByte(s[i])
		case (ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n') && !inQuotes:
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteByte(ch)
		}
	}
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args
}

func parseFloat(s string) (float64, error) {
	var val float64
	_, err := fmt.Sscanf(s, "%f", &val)
	return val, err
}

func parseInt(s string) (int, error) {
	var val int
	_, err := fmt.Sscanf(s, "%d", &val)
	return val, err
}

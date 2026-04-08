package params

import (
	"reflect"
	"testing"
)

func TestShellSplit(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"--endpoint 10.0.0.1 --scope capacity", []string{"--endpoint", "10.0.0.1", "--scope", "capacity"}},
		{`--endpoint "10.0.0.1" --apitoken "my token"`, []string{"--endpoint", "10.0.0.1", "--apitoken", "my token"}},
		{"  --insecure  --scope  hardware  ", []string{"--insecure", "--scope", "hardware"}},
		{"", nil},
	}
	for _, tt := range tests {
		got := shellSplit(tt.input)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("shellSplit(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestParse(t *testing.T) {
	input := "--endpoint 10.0.0.1 --apitoken mytoken --scope capacity --insecure --warning 75 --critical 85"
	p, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if p.Endpoint != "10.0.0.1" {
		t.Errorf("Endpoint = %q, want %q", p.Endpoint, "10.0.0.1")
	}
	if p.APIToken != "mytoken" {
		t.Errorf("APIToken = %q, want %q", p.APIToken, "mytoken")
	}
	if p.Scope != "capacity" {
		t.Errorf("Scope = %q, want %q", p.Scope, "capacity")
	}
	if !p.Insecure {
		t.Error("Insecure = false, want true")
	}
	if p.Warning != 75 {
		t.Errorf("Warning = %v, want 75", p.Warning)
	}
	if p.Critical != 85 {
		t.Errorf("Critical = %v, want 85", p.Critical)
	}
}

func TestParseDefaults(t *testing.T) {
	input := "--endpoint 10.0.0.1 --apitoken tok --scope performance"
	p, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if p.AuthMode != "session" {
		t.Errorf("AuthMode = %q, want %q", p.AuthMode, "session")
	}
	if p.Warning != 80 {
		t.Errorf("Warning = %v, want 80", p.Warning)
	}
	if p.Critical != 90 {
		t.Errorf("Critical = %v, want 90", p.Critical)
	}
	if !p.Insecure {
		t.Error("Insecure = false, want true (default)")
	}
}

func TestParseSecureFlag(t *testing.T) {
	input := "--endpoint 10.0.0.1 --apitoken tok --scope performance --secure"
	p, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if p.Insecure {
		t.Error("Insecure = true, want false after --secure")
	}
}

func TestParseVolumes(t *testing.T) {
	input := "--endpoint 10.0.0.1 --apitoken tok --scope volumes --volumes vol1,vol2,vol3"
	p, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	want := []string{"vol1", "vol2", "vol3"}
	if !reflect.DeepEqual(p.Volumes, want) {
		t.Errorf("Volumes = %v, want %v", p.Volumes, want)
	}
}

func TestParseMissingRequired(t *testing.T) {
	tests := []struct {
		input string
		errSub string
	}{
		{"--apitoken tok --scope capacity", "--endpoint is required"},
		{"--endpoint 10.0.0.1 --scope capacity", "--apitoken is required"},
		{"--endpoint 10.0.0.1 --apitoken tok", "--scope is required"},
	}
	for _, tt := range tests {
		_, err := Parse(tt.input)
		if err == nil {
			t.Errorf("Parse(%q) expected error containing %q", tt.input, tt.errSub)
			continue
		}
		if err.Error() != tt.errSub {
			t.Errorf("Parse(%q) error = %q, want %q", tt.input, err.Error(), tt.errSub)
		}
	}
}

func TestParseInvalidScope(t *testing.T) {
	input := "--endpoint 10.0.0.1 --apitoken tok --scope invalid"
	_, err := Parse(input)
	if err == nil {
		t.Error("expected error for invalid scope")
	}
}

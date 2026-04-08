package sensor

import (
	"encoding/json"
	"testing"
)

func TestNewResult(t *testing.T) {
	r := NewResult("test message")
	if r.Version != 3 {
		t.Errorf("Version = %d, want 3", r.Version)
	}
	if r.Status != "ok" {
		t.Errorf("Status = %q, want %q", r.Status, "ok")
	}
	if r.Message != "test message" {
		t.Errorf("Message = %q, want %q", r.Message, "test message")
	}
}

func TestNewErrorResult(t *testing.T) {
	r := NewErrorResult("something broke")
	if r.Status != "error" {
		t.Errorf("Status = %q, want %q", r.Status, "error")
	}
}

func TestAddChannels(t *testing.T) {
	r := NewResult("test")
	r.AddIntChannel(10, "Test Int", "count", 42)
	r.AddFloatChannel(11, "Test Float", "percent", 3.14159)
	r.AddCustomFloatChannel(12, "Test Custom", "x:1", 5.678)
	r.AddFloatChannelWithLimits(13, "Test Limits", "percent", 85.5, 80, 90)

	if len(r.Channels) != 4 {
		t.Fatalf("got %d channels, want 4", len(r.Channels))
	}

	if r.Channels[0].Value != int64(42) {
		t.Errorf("channel 0 value = %v, want 42", r.Channels[0].Value)
	}
	if r.Channels[1].Value != 3.14 {
		t.Errorf("channel 1 value = %v, want 3.14", r.Channels[1].Value)
	}
	if r.Channels[2].DisplayUnit != "x:1" {
		t.Errorf("channel 2 display_unit = %q, want %q", r.Channels[2].DisplayUnit, "x:1")
	}
	if r.Channels[3].Limits == nil {
		t.Fatal("channel 3 limits is nil")
	}
	if *r.Channels[3].Limits.Warning.Upper != 80 {
		t.Errorf("channel 3 warning upper = %v, want 80", *r.Channels[3].Limits.Warning.Upper)
	}
}

func TestResultJSON(t *testing.T) {
	r := NewResult("test")
	r.AddIntChannel(10, "Total", "size_bytes_disk", 1000)

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if parsed["version"].(float64) != 3 {
		t.Errorf("JSON version = %v, want 3", parsed["version"])
	}
	if parsed["status"].(string) != "ok" {
		t.Errorf("JSON status = %v, want ok", parsed["status"])
	}
}

func TestRoundFloat(t *testing.T) {
	tests := []struct {
		val  float64
		prec int
		want float64
	}{
		{3.14159, 2, 3.14},
		{3.145, 2, 3.15},
		{100.0, 2, 100.0},
		{0.001, 2, 0.0},
	}
	for _, tt := range tests {
		got := roundFloat(tt.val, tt.prec)
		if got != tt.want {
			t.Errorf("roundFloat(%v, %d) = %v, want %v", tt.val, tt.prec, got, tt.want)
		}
	}
}

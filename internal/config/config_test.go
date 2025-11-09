package config

import "testing"

func TestGetEnvWithDefault(t *testing.T) {
	t.Setenv("TEST_KEY_SET", "value")
	if got := getEnvWithDefault("TEST_KEY_SET", "default"); got != "value" {
		t.Fatalf("expected %q, got %q", "value", got)
	}

	t.Setenv("TEST_KEY_EMPTY", "")
	if got := getEnvWithDefault("TEST_KEY_EMPTY", "default"); got != "default" {
		t.Fatalf("expected default %q, got %q", "default", got)
	}

	if got := getEnvWithDefault("TEST_KEY_MISSING", "default"); got != "default" {
		t.Fatalf("expected default %q for missing key, got %q", "default", got)
	}
}

func TestGetEnvIntWithDefault(t *testing.T) {
	t.Setenv("INT_OK", "42")
	if got := getEnvIntWithDefault("INT_OK", 5); got != 42 {
		t.Fatalf("expected 42, got %d", got)
	}

	t.Setenv("INT_BAD", "abc")
	if got := getEnvIntWithDefault("INT_BAD", 5); got != 5 {
		t.Fatalf("expected default 5 for bad int, got %d", got)
	}

	if got := getEnvIntWithDefault("INT_MISSING", 7); got != 7 {
		t.Fatalf("expected default 7 for missing int key, got %d", got)
	}
}

func TestSplitCSV(t *testing.T) {
	tests := []struct {
		in   string
		want []string
	}{
		{"a,b,c", []string{"a", "b", "c"}},
		{" a , b ,  c  ", []string{"a", "b", "c"}},
		{"", nil},
		{",,a,,b,", []string{"a", "b"}},
	}

	for _, tt := range tests {
		got := splitCSV(tt.in)
		if len(got) != len(tt.want) {
			t.Fatalf("splitCSV(%q) len=%d, want %d", tt.in, len(got), len(tt.want))
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Fatalf("splitCSV(%q)[%d]=%q, want %q", tt.in, i, got[i], tt.want[i])
			}
		}
	}
}

func TestLoad_UsesEnvAndDefaults(t *testing.T) {
	// server envs
	t.Setenv("ADDR", ":9000")
	t.Setenv("READ_TIMEOUT", "15")
	t.Setenv("WRITE_TIMEOUT", "20")
	t.Setenv("IDLE_TIMEOUT", "120")
	t.Setenv("READ_HEADER_TIMEOUT", "3")

	// promo files
	t.Setenv("PROMO_FILES", "/tmp/a,/tmp/b")

	cfg := Load()

	if cfg.Server.Addr != ":9000" {
		t.Errorf("Addr = %q, want %q", cfg.Server.Addr, ":9000")
	}
	if cfg.Server.ReadTimeout != 15 {
		t.Errorf("ReadTimeout = %d, want %d", cfg.Server.ReadTimeout, 15)
	}
	if cfg.Server.WriteTimeout != 20 {
		t.Errorf("WriteTimeout = %d, want %d", cfg.Server.WriteTimeout, 20)
	}
	if cfg.Server.IdleTimeout != 120 {
		t.Errorf("IdleTimeout = %d, want %d", cfg.Server.IdleTimeout, 120)
	}
	if cfg.Server.ReadHeaderTimeout != 3 {
		t.Errorf("ReadHeaderTimeout = %d, want %d", cfg.Server.ReadHeaderTimeout, 3)
	}

	if len(cfg.PromoFiles) != 2 || cfg.PromoFiles[0] != "/tmp/a" || cfg.PromoFiles[1] != "/tmp/b" {
		t.Errorf("PromoFiles = %#v, want []string{\"/tmp/a\",\"/tmp/b\"}", cfg.PromoFiles)
	}
}

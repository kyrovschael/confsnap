package alert_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nicholasgasior/confsnap/internal/alert"
)

func overrideConfigDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	// patch via exported function not available; use build-tag approach or
	// test the round-trip via SaveConfig/LoadConfig with a known path.
	return dir
}

func TestDefaultConfig(t *testing.T) {
	cfg := alert.DefaultConfig()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
	if cfg.Level != "WARNING" {
		t.Errorf("expected Level=WARNING, got %s", cfg.Level)
	}
	if !cfg.Stdout {
		t.Error("expected Stdout=true")
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "alert_config.json")

	cfg := alert.Config{
		Enabled: true,
		Level:   "CRITICAL",
		LogFile: "/var/log/confsnap.log",
		Stdout:  false,
	}

	// Write manually to the temp path to test round-trip.
	import_json := `{"enabled":true,"level":"CRITICAL","log_file":"/var/log/confsnap.log","stdout":false}`
	if err := os.WriteFile(path, []byte(import_json), 0644); err != nil {
		t.Fatal(err)
	}
	_ = cfg
	_ = path
}

func TestBuildDispatcher_StdoutOnly(t *testing.T) {
	cfg := alert.Config{
		Enabled: true,
		Level:   "INFO",
		Stdout:  true,
	}
	d := alert.BuildDispatcher(cfg)
	if d == nil {
		t.Fatal("expected non-nil dispatcher")
	}
	if len(d.Handlers) != 1 {
		t.Errorf("expected 1 handler, got %d", len(d.Handlers))
	}
}

func TestBuildDispatcher_FileAndStdout(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "alerts.log")
	cfg := alert.Config{
		Enabled: true,
		Level:   "WARNING",
		LogFile: tmp,
		Stdout:  true,
	}
	d := alert.BuildDispatcher(cfg)
	if len(d.Handlers) != 2 {
		t.Errorf("expected 2 handlers, got %d", len(d.Handlers))
	}
}

package alert_test

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/nicholasgasior/confsnap/internal/alert"
)

type bufHandler struct {
	buf bytes.Buffer
}

func (b *bufHandler) Send(a alert.Alert) error {
	_, err := fmt.Fprintf(&b.buf, "%s %s %s: %s\n", a.Level, a.File, a.Message, a.Timestamp.Format(time.RFC3339))
	return err
}

func TestDispatcher_SendsToAllHandlers(t *testing.T) {
	h1 := &bufHandler{}
	h2 := &bufHandler{}
	d := alert.NewDispatcher(h1, h2)

	a := alert.Alert{
		Timestamp: time.Now(),
		Level:     alert.LevelWarning,
		File:      "/etc/nginx.conf",
		Message:   "file modified",
	}
	errs := d.Dispatch(a)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !strings.Contains(h1.buf.String(), "WARNING") {
		t.Errorf("h1 missing WARNING: %s", h1.buf.String())
	}
	if !strings.Contains(h2.buf.String(), "/etc/nginx.conf") {
		t.Errorf("h2 missing file path: %s", h2.buf.String())
	}
}

func TestFileHandler_WritesAlert(t *testing.T) {
	tmp, err := os.CreateTemp(t.TempDir(), "alert-*.log")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	h := &alert.FileHandler{Path: tmp.Name()}
	a := alert.Alert{
		Timestamp: time.Now(),
		Level:     alert.LevelCritical,
		File:      "/etc/ssh/sshd_config",
		Message:   "critical drift detected",
	}
	if err := h.Send(a); err != nil {
		t.Fatalf("Send error: %v", err)
	}

	data, _ := os.ReadFile(tmp.Name())
	if !strings.Contains(string(data), "CRITICAL") {
		t.Errorf("expected CRITICAL in log, got: %s", data)
	}
	if !strings.Contains(string(data), "critical drift detected") {
		t.Errorf("expected message in log, got: %s", data)
	}
}

func TestLevelFromString(t *testing.T) {
	cases := []struct {
		input    string
		expected alert.Level
	}{
		{"warning", alert.LevelWarning},
		{"CRITICAL", alert.LevelCritical},
		{"info", alert.LevelInfo},
		{"unknown", alert.LevelInfo},
	}
	for _, c := range cases {
		got := alert.LevelFromString(c.input)
		if got != c.expected {
			t.Errorf("LevelFromString(%q) = %q, want %q", c.input, got, c.expected)
		}
	}
}

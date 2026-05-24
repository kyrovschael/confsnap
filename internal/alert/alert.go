package alert

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo    Level = "INFO"
	LevelWarning Level = "WARNING"
	LevelCritical Level = "CRITICAL"
)

// Alert represents a drift alert event.
type Alert struct {
	Timestamp time.Time
	Level     Level
	File      string
	Message   string
}

// Handler defines how an alert is dispatched.
type Handler interface {
	Send(a Alert) error
}

// Dispatcher holds registered handlers and sends alerts to all of them.
type Dispatcher struct {
	Handlers []Handler
}

// NewDispatcher creates a Dispatcher with the given handlers.
func NewDispatcher(handlers ...Handler) *Dispatcher {
	return &Dispatcher{Handlers: handlers}
}

// Dispatch sends the alert to all registered handlers, collecting errors.
func (d *Dispatcher) Dispatch(a Alert) []error {
	var errs []error
	for _, h := range d.Handlers {
		if err := h.Send(a); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// StdoutHandler writes alerts to stdout.
type StdoutHandler struct{}

func (s *StdoutHandler) Send(a Alert) error {
	_, err := fmt.Fprintf(os.Stdout, "[%s] %s %s: %s\n",
		a.Timestamp.Format(time.RFC3339),
		a.Level,
		a.File,
		a.Message,
	)
	return err
}

// FileHandler appends alerts to a log file.
type FileHandler struct {
	Path string
}

func (f *FileHandler) Send(a Alert) error {
	line := fmt.Sprintf("[%s] %s %s: %s\n",
		a.Timestamp.Format(time.RFC3339),
		a.Level,
		a.File,
		a.Message,
	)
	file, err := os.OpenFile(f.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("alert file handler: %w", err)
	}
	defer file.Close()
	_, err = file.WriteString(line)
	return err
}

// LevelFromString parses a Level from a string, defaulting to INFO.
func LevelFromString(s string) Level {
	switch strings.ToUpper(s) {
	case "WARNING":
		return LevelWarning
	case "CRITICAL":
		return LevelCritical
	default:
		return LevelInfo
	}
}

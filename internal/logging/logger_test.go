package logging_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"api-registration-authorization/internal/logging"
	"github.com/rs/zerolog"
)

type example struct{ logger zerolog.Logger }

func (e *example) emit() { e.logger.Info().Str("request_id", "123").Msg("first\nsecond") }

func decode(t *testing.T, buffer *bytes.Buffer) map[string]interface{} {
	t.Helper()
	line := buffer.String()
	if strings.Count(line, "\n") != 1 {
		t.Fatalf("expected one JSON line: %q", line)
	}
	result := map[string]interface{}{}
	if err := json.Unmarshal(buffer.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	return result
}
func TestJSONAndAutomaticSource(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("LOG_LEVEL", "info")
	var buffer bytes.Buffer
	logger, err := logging.New("example", &buffer)
	if err != nil {
		t.Fatal(err)
	}
	logger.Debug().Msg("hidden")
	(&example{logger}).emit()
	entry := decode(t, &buffer)
	for key, want := range map[string]interface{}{
		"service": "example", "environment": "test", "level": "info",
		"message": "first\nsecond", "class": "example", "request_id": "123",
	} {
		if entry[key] != want {
			t.Errorf("%s = %v, want %v", key, entry[key], want)
		}
	}
	if !strings.HasSuffix(entry["function"].(string), "(*example).emit") {
		t.Errorf("wrong function: %v", entry["function"])
	}
	if !strings.Contains(entry["caller"].(string), "logger_test.go:") {
		t.Errorf("wrong caller: %v", entry["caller"])
	}
	timestamp, err := time.Parse(time.RFC3339Nano, entry["timestamp"].(string))
	if err != nil || timestamp.Location() != time.UTC {
		t.Errorf("invalid UTC timestamp: %v", entry["timestamp"])
	}
}
func TestInvalidLevel(t *testing.T) {
	t.Setenv("LOG_LEVEL", "invalid")
	if _, err := logging.New("example", &bytes.Buffer{}); err == nil {
		t.Fatal("expected invalid level error")
	}
}
func TestTemporalContextAndErrors(t *testing.T) {
	t.Setenv("LOG_LEVEL", "debug")
	var buffer bytes.Buffer
	logger, err := logging.New("example", &buffer)
	if err != nil {
		t.Fatal(err)
	}
	parent := logging.NewTemporal(logger)
	child := parent.With("WorkflowID", "workflow-123")
	child.Error("failed", "error", errors.New("failure"))
	entry := decode(t, &buffer)
	if entry["WorkflowID"] != "workflow-123" || entry["error"] != "failure" || entry["level"] != "error" {
		t.Fatal(entry)
	}
	buffer.Reset()
	parent.Info("parent")
	if _, exists := decode(t, &buffer)["WorkflowID"]; exists {
		t.Fatal("child mutated parent")
	}
}

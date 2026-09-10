package logging

import (
	"io"
	stdlog "log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// New emits one JSON object per line. All services use the same field names.
func New(service string, output io.Writer) (zerolog.Logger, error) {
	level := zerolog.InfoLevel
	if value := os.Getenv("LOG_LEVEL"); value != "" {
		parsed, err := zerolog.ParseLevel(value)
		if err != nil {
			return zerolog.Logger{}, err
		}
		level = parsed
	}
	return newLogger(service, output).Level(level), nil
}

func newLogger(service string, output io.Writer) zerolog.Logger {
	return zerolog.New(zerolog.SyncWriter(output)).With().
		Str("service", service).Str("environment", os.Getenv("APP_ENV")).
		Logger().Hook(sourceHook{})
}

// Init must run after loading .env and before starting clients or workers.
func Init(service string) func() {
	logger, err := New(service, os.Stdout)
	if err != nil {
		fallback := newLogger(service, os.Stdout)
		fallback.Fatal().Err(err).Msg("Invalid logging configuration")
	}
	previous := log.Logger
	writer, flags, prefix := stdlog.Writer(), stdlog.Flags(), stdlog.Prefix()
	log.Logger = logger
	stdlog.SetOutput(standardWriter{logger})
	stdlog.SetFlags(0)
	stdlog.SetPrefix("")
	return func() {
		stdlog.SetOutput(writer)
		stdlog.SetFlags(flags)
		stdlog.SetPrefix(prefix)
		log.Logger = previous
	}
}

type standardWriter struct{ logger zerolog.Logger }

func (w standardWriter) Write(p []byte) (int, error) {
	w.logger.Info().Msg(strings.TrimRight(string(p), "\r\n"))
	return len(p), nil
}

type sourceHook struct{}

func (sourceHook) Run(event *zerolog.Event, _ zerolog.Level, _ string) {
	event.Str("timestamp", time.Now().UTC().Format(time.RFC3339Nano))
	// Inspect only enabled log events; skip adapters to report the actual caller.
	var pcs [64]uintptr
	count := runtime.Callers(2, pcs[:])
	frames := runtime.CallersFrames(pcs[:count])
	for {
		frame, more := frames.Next()
		name := frame.Function
		if name != "" && !infrastructure(name) {
			event.Str("caller", frame.File+":"+strconv.Itoa(frame.Line)).Str("function", name)
			if start := strings.Index(name, ".("); start >= 0 {
				if end := strings.Index(name[start+2:], ")"); end >= 0 {
					event.Str("class", strings.TrimPrefix(name[start+2:start+2+end], "*"))
				}
			}
			return
		}
		if !more {
			return
		}
	}
}
func infrastructure(name string) bool {
	return strings.Contains(name, "/internal/logging.") ||
		strings.HasPrefix(name, "github.com/rs/zerolog") ||
		strings.HasPrefix(name, "github.com/gofiber/fiber/") ||
		strings.HasPrefix(name, "go.temporal.io/sdk/") ||
		strings.HasPrefix(name, "gorm.io/") ||
		strings.HasPrefix(name, "log.") ||
		strings.HasPrefix(name, "runtime.")
}

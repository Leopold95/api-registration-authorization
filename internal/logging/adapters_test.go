package logging_test

import (
	"bytes"
	"context"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"api-registration-authorization/internal/logging"
	"github.com/gofiber/fiber/v3"
	fiberlog "github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/rs/zerolog/log"
	gormlog "gorm.io/gorm/logger"
)

func TestGORMLevelsAndNoSQL(t *testing.T) {
	t.Setenv("LOG_LEVEL", "info")
	var buffer bytes.Buffer
	logger, err := logging.New("example", &buffer)
	if err != nil {
		t.Fatal(err)
	}
	adapter := logging.NewGORM(logger)
	sql := func() (string, int64) { t.Fatal("SQL must not be rendered"); return "", 0 }
	adapter.Trace(context.Background(), time.Now(), sql, errors.New("database unavailable"))
	entry := decode(t, &buffer)
	if entry["level"] != "error" || entry["error"] != "database unavailable" {
		t.Fatal(entry)
	}
	buffer.Reset()
	adapter.LogMode(gormlog.Silent).Trace(context.Background(), time.Now(), sql, errors.New("hidden"))
	if buffer.Len() != 0 {
		t.Fatal(buffer.String())
	}
}
func TestFiberAdapter(t *testing.T) {
	t.Setenv("LOG_LEVEL", "info")
	var buffer bytes.Buffer
	logger, err := logging.New("example", &buffer)
	if err != nil {
		t.Fatal(err)
	}
	logging.InstallFiber(logger)
	t.Cleanup(func() { logging.InstallFiber(log.Logger) })
	fiberlog.Errorw("failed", "error", errors.New("failure"))
	entry := decode(t, &buffer)
	if entry["level"] != "error" || entry["error"] != "failure" {
		t.Fatal(entry)
	}
	buffer.Reset()
	fiberlog.SetLevel(fiberlog.LevelDebug)
	fiberlog.Debug("enabled")
	if decode(t, &buffer)["level"] != "debug" {
		t.Fatal("wrong level")
	}
}
func TestHTTPFinalStatusAndRecovery(t *testing.T) {
	for _, test := range []struct {
		name    string
		status  int
		handler fiber.Handler
	}{
		{"ok", 200, func(c fiber.Ctx) error { return c.SendString("ok") }},
		{"error", 401, func(c fiber.Ctx) error { return fiber.ErrUnauthorized }},
		{"panic", 500, func(c fiber.Ctx) error { panic("failure") }},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("LOG_LEVEL", "info")
			var buffer bytes.Buffer
			logger, err := logging.New("example", &buffer)
			if err != nil {
				t.Fatal(err)
			}
			previous := log.Logger
			log.Logger = logger
			defer func() { log.Logger = previous }()
			app := fiber.New()
			app.Use(logging.HTTP())
			app.Use(recover.New())
			app.Get("/test", test.handler)
			response, err := app.Test(httptest.NewRequest("GET", "/test?token=secret", nil))
			if err != nil {
				t.Fatal(err)
			}
			defer response.Body.Close()
			entry := decode(t, &buffer)
			if response.StatusCode != test.status || entry["status"] != float64(test.status) {
				t.Fatal(entry)
			}
			if entry["path"] != "/test" {
				t.Fatal("query string leaked", entry)
			}
		})
	}
}

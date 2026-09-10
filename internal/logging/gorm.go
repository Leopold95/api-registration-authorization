package logging

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
)

type GORM struct {
	logger zerolog.Logger
	level  gormlog.LogLevel
}

func NewGORM(logger zerolog.Logger) *GORM { return &GORM{logger: logger, level: gormlog.Warn} }
func (l *GORM) LogMode(level gormlog.LogLevel) gormlog.Interface {
	copy := *l
	copy.level = level
	return &copy
}
func (l *GORM) Info(_ context.Context, msg string, args ...interface{}) {
	if l.level >= gormlog.Info {
		l.logger.Info().Msgf(msg, args...)
	}
}
func (l *GORM) Warn(_ context.Context, msg string, args ...interface{}) {
	if l.level >= gormlog.Warn {
		l.logger.Warn().Msgf(msg, args...)
	}
}
func (l *GORM) Error(_ context.Context, msg string, args ...interface{}) {
	if l.level >= gormlog.Error {
		l.logger.Error().Msgf(msg, args...)
	}
}
func (l *GORM) Trace(_ context.Context, begin time.Time, _ func() (string, int64), err error) {
	if l.level == gormlog.Silent {
		return
	}
	elapsed := time.Since(begin)
	var event *zerolog.Event
	var message string
	// Do not render SQL/bound parameters: authentication queries contain secrets.
	switch {
	case err != nil && !errors.Is(err, gorm.ErrRecordNotFound) && l.level >= gormlog.Error:
		event, message = l.logger.Error().Err(err), "Database query failed"
	case elapsed > 200*time.Millisecond && l.level >= gormlog.Warn:
		event, message = l.logger.Warn(), "Slow database query"
	case l.level >= gormlog.Info:
		event, message = l.logger.Info(), "Database query completed"
	default:
		return
	}
	event.Float64("duration_ms", float64(elapsed)/float64(time.Millisecond)).Msg(message)
}

var _ gormlog.Interface = (*GORM)(nil)

package logging

import (
	"fmt"

	"github.com/rs/zerolog"
	temporallog "go.temporal.io/sdk/log"
)

// Temporal preserves SDK fields and workflow replay suppression.
type Temporal struct{ logger zerolog.Logger }

func NewTemporal(logger zerolog.Logger) *Temporal { return &Temporal{logger: logger} }
func keyValues(fields []interface{}) map[string]interface{} {
	result := make(map[string]interface{}, (len(fields)+1)/2)
	for i := 0; i < len(fields); i += 2 {
		var value interface{} = "(missing)"
		if i+1 < len(fields) {
			value = fields[i+1]
		}
		if err, ok := value.(error); ok && err != nil {
			value = err.Error()
		}
		result[fmt.Sprint(fields[i])] = value
	}
	return result
}
func (l *Temporal) Debug(msg string, fields ...interface{}) {
	l.logger.Debug().Fields(keyValues(fields)).Msg(msg)
}
func (l *Temporal) Info(msg string, fields ...interface{}) {
	l.logger.Info().Fields(keyValues(fields)).Msg(msg)
}
func (l *Temporal) Warn(msg string, fields ...interface{}) {
	l.logger.Warn().Fields(keyValues(fields)).Msg(msg)
}
func (l *Temporal) Error(msg string, fields ...interface{}) {
	l.logger.Error().Fields(keyValues(fields)).Msg(msg)
}
func (l *Temporal) With(fields ...interface{}) temporallog.Logger {
	return &Temporal{logger: l.logger.With().Fields(keyValues(fields)).Logger()}
}

var _ temporallog.Logger = (*Temporal)(nil)
var _ temporallog.WithLogger = (*Temporal)(nil)

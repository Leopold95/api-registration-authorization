package logging

import (
	"fmt"
	fiberlog "github.com/gofiber/fiber/v3/log"
	"github.com/rs/zerolog"
	"io"
	"sync"
)

type Fiber struct {
	mu     sync.RWMutex
	logger zerolog.Logger
}

func InstallFiber(logger zerolog.Logger) { fiberlog.SetLogger(&Fiber{logger: logger}) }
func (l *Fiber) emit(level zerolog.Level, msg string, fields ...interface{}) {
	logger := l.Logger()
	event := logger.WithLevel(level)
	if level == zerolog.FatalLevel {
		event = logger.Fatal()
	}
	if level == zerolog.PanicLevel {
		event = logger.Panic()
	}
	event.Fields(keyValues(fields)).Msg(msg)
}
func (l *Fiber) Logger() zerolog.Logger                          { l.mu.RLock(); defer l.mu.RUnlock(); return l.logger }
func (l *Fiber) WithContext(_ interface{}) fiberlog.CommonLogger { return l }
func (l *Fiber) SetOutput(writer io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger = l.logger.Output(zerolog.SyncWriter(writer))
}
func (l *Fiber) SetLevel(level fiberlog.Level) {
	levels := []zerolog.Level{zerolog.TraceLevel, zerolog.DebugLevel, zerolog.InfoLevel, zerolog.WarnLevel, zerolog.ErrorLevel, zerolog.FatalLevel, zerolog.PanicLevel}
	if int(level) < 0 || int(level) >= len(levels) {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger = l.logger.Level(levels[level])
}
func (l *Fiber) Trace(args ...interface{}) { l.emit(zerolog.TraceLevel, fmt.Sprint(args...)) }
func (l *Fiber) Tracef(msg string, args ...interface{}) {
	l.emit(zerolog.TraceLevel, fmt.Sprintf(msg, args...))
}
func (l *Fiber) Tracew(msg string, args ...interface{}) { l.emit(zerolog.TraceLevel, msg, args...) }
func (l *Fiber) Debug(args ...interface{})              { l.emit(zerolog.DebugLevel, fmt.Sprint(args...)) }
func (l *Fiber) Debugf(msg string, args ...interface{}) {
	l.emit(zerolog.DebugLevel, fmt.Sprintf(msg, args...))
}
func (l *Fiber) Debugw(msg string, args ...interface{}) { l.emit(zerolog.DebugLevel, msg, args...) }
func (l *Fiber) Info(args ...interface{})               { l.emit(zerolog.InfoLevel, fmt.Sprint(args...)) }
func (l *Fiber) Infof(msg string, args ...interface{}) {
	l.emit(zerolog.InfoLevel, fmt.Sprintf(msg, args...))
}
func (l *Fiber) Infow(msg string, args ...interface{}) { l.emit(zerolog.InfoLevel, msg, args...) }
func (l *Fiber) Warn(args ...interface{})              { l.emit(zerolog.WarnLevel, fmt.Sprint(args...)) }
func (l *Fiber) Warnf(msg string, args ...interface{}) {
	l.emit(zerolog.WarnLevel, fmt.Sprintf(msg, args...))
}
func (l *Fiber) Warnw(msg string, args ...interface{}) { l.emit(zerolog.WarnLevel, msg, args...) }
func (l *Fiber) Error(args ...interface{})             { l.emit(zerolog.ErrorLevel, fmt.Sprint(args...)) }
func (l *Fiber) Errorf(msg string, args ...interface{}) {
	l.emit(zerolog.ErrorLevel, fmt.Sprintf(msg, args...))
}
func (l *Fiber) Errorw(msg string, args ...interface{}) { l.emit(zerolog.ErrorLevel, msg, args...) }
func (l *Fiber) Fatal(args ...interface{})              { l.emit(zerolog.FatalLevel, fmt.Sprint(args...)) }
func (l *Fiber) Fatalf(msg string, args ...interface{}) {
	l.emit(zerolog.FatalLevel, fmt.Sprintf(msg, args...))
}
func (l *Fiber) Fatalw(msg string, args ...interface{}) { l.emit(zerolog.FatalLevel, msg, args...) }
func (l *Fiber) Panic(args ...interface{})              { l.emit(zerolog.PanicLevel, fmt.Sprint(args...)) }
func (l *Fiber) Panicf(msg string, args ...interface{}) {
	l.emit(zerolog.PanicLevel, fmt.Sprintf(msg, args...))
}
func (l *Fiber) Panicw(msg string, args ...interface{}) { l.emit(zerolog.PanicLevel, msg, args...) }

var _ fiberlog.AllLogger[zerolog.Logger] = (*Fiber)(nil)

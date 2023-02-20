package logger

import (
	"fmt"
	"os"
	"time"
)

type Logger interface {
	AddHandler(handler Handler)
	Close() error

	Log(level LogLevel, message string, context ...any)

	Emergency(message string, context ...any)
	Alert(message string, context ...any)
	Critical(message string, context ...any)
	Error(message string, context ...any)
	Warning(message string, context ...any)
	Notice(message string, context ...any)
	Info(message string, context ...any)
	Debug(message string, context ...any)
}

type Monolog struct {
	handlers []Handler
}

func New() Logger {
	return &Monolog{}
}

func (l *Monolog) AddHandler(handler Handler) {
	l.handlers = append(l.handlers, handler)
}

func (l *Monolog) Close() error {
	var err error
	for k := range l.handlers {
		err = l.handlers[k].Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to close %T handler\n", l.handlers[k])
		}
	}
	return nil
}

func (l *Monolog) Log(level LogLevel, message string, context ...any) {
	record := Record{
		level,
		time.Now(),
		message,
		context,
	}

	for _, handler := range l.handlers {
		if handler.Level() >= level {
			handler.Handle(record)

			if !handler.IsBubble() {
				break
			}
		}
	}
}

func (l *Monolog) Emergency(message string, context ...any) {
	l.Log(LogLevelEmergency, message, context...)
}
func (l *Monolog) Alert(message string, context ...any) {
	l.Log(LogLevelAlert, message, context...)
}

func (l *Monolog) Critical(message string, context ...any) {
	l.Log(LogLevelCritical, message, context...)
}

func (l *Monolog) Error(message string, context ...any) {
	l.Log(LogLevelError, message, context...)
}

func (l *Monolog) Warning(message string, context ...any) {
	l.Log(LogLevelWarning, message, context...)
}

func (l *Monolog) Notice(message string, context ...any) {
	l.Log(LogLevelNotice, message, context...)
}

func (l *Monolog) Info(message string, context ...any) {
	l.Log(LogLevelInfo, message, context...)
}

func (l *Monolog) Debug(message string, context ...any) {
	l.Log(LogLevelDebug, message, context...)
}

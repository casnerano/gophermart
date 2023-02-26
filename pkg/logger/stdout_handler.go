package logger

import (
	"fmt"
	"os"
)

type StdOutHandler struct {
	formatter Formatter
	level     LogLevel
	bubble    bool
}

func NewStdOutHandler(formatter Formatter, level LogLevel, bubble bool) Handler {
	return &StdOutHandler{formatter, level, bubble}
}

func (s *StdOutHandler) Handle(record Record) bool {
	_, err := fmt.Fprintln(os.Stdout, string(s.formatter.Format(record)))
	return err == nil
}

func (s *StdOutHandler) Level() LogLevel {
	return s.level
}

func (s *StdOutHandler) IsBubble() bool {
	return s.bubble
}

func (s *StdOutHandler) Close() error {
	return nil
}

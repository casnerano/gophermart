package logger

import (
	"fmt"
)

const layoutDate = "02.01.2006 15:04:05"

type TextFormatter struct{}

func NewTextFormatter() Formatter {
	return &TextFormatter{}
}

func (j TextFormatter) Format(record Record) []byte {
	text := fmt.Sprintf(
		"%s %s %s",
		"["+record.Level.String()+"]",
		record.Date.Format(layoutDate),
		record.Message,
	)

	if len(record.Context) > 0 {
		text = fmt.Sprintf("%s %v", text, record.Context)
	}

	return []byte(text)
}

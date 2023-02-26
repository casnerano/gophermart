package logger

import "encoding/json"

type JSONFormatter struct{}

func NewJSONFormatter() Formatter {
	return &JSONFormatter{}
}

func (j JSONFormatter) Format(record Record) []byte {
	b, _ := json.Marshal(record)
	return b
}

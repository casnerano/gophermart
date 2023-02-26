package logger

import "time"

type Record struct {
	Level   LogLevel  `json:"level"`
	Date    time.Time `json:"date"`
	Message string    `json:"message"`
	Context []any     `json:"context,omitempty"`
}

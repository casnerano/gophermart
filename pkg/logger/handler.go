package logger

type Handler interface {
	Handle(record Record) bool
	Level() LogLevel
	IsBubble() bool
	Close() error
}

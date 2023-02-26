package logger

type Formatter interface {
	Format(record Record) []byte
}

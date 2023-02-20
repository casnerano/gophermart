package logger

type LogLevel int

const (
	LogLevelEmergency LogLevel = iota
	LogLevelAlert
	LogLevelCritical
	LogLevelError
	LogLevelWarning
	LogLevelNotice
	LogLevelInfo
	LogLevelDebug
)

var levelString = map[LogLevel]string{
	LogLevelEmergency: "EMERGENCY",
	LogLevelAlert:     "ALERT",
	LogLevelCritical:  "CRITICAL",
	LogLevelError:     "ERROR",
	LogLevelWarning:   "WARNING",
	LogLevelNotice:    "NOTICE",
	LogLevelInfo:      "INFO",
	LogLevelDebug:     "DEBUG",
}

func (l LogLevel) String() string {
	return levelString[l]
}

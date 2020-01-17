package logger

type LogLevel int

const (
	LogLevelPanic   = LogLevel(1)
	LogLevelFatal   = LogLevel(2)
	LogLevelError   = LogLevel(4)
	LogLevelPrint   = LogLevel(8)
	LogLevelWarning = LogLevel(16)
	LogLevelInfo    = LogLevel(32)
	LogLevelTrace   = LogLevel(64)
	LogLevelDebug   = LogLevel(128)
)

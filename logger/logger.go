package logger

type Logger struct {
	Output    Output
	Formatter Formatter
	Prefixs   []Prefix
}

func (l *Logger) Log(v ...interface{}) {

}
func NewLogger() *Logger {
	return &Logger{}
}
func createLogger(o Output, f Formatter, p ...Prefix) *Logger {
	return &Logger{
		Output:    o,
		Formatter: f,
		Prefixs:   p,
	}
}

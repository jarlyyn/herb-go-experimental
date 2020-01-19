package logger

import (
	"log"
)

type Logger struct {
	Writer
	Formatter Formatter
	Prefixs   []Prefix
}

func (l *Logger) Log(v ...interface{}) {
	var data []byte
	var err error
	defer func() {
		if err != nil {
			log.Fatalln(err.Error())
		}
	}()
	if l.Formatter == nil {
		data, err = DefaultFormatter.Format(v...)
	} else {
		data, err = l.Formatter.Format(v...)
	}
	if err != nil {
		return
	}
	var output = []byte{}
	for k := range l.Prefixs {
		output = append(output, l.Prefixs[k].NewPrefix()...)
	}
	output = append(output, data...)
	_, err = l.Writer.Write(output)

}

func (l *Logger) SetWriter(w Writer) *Logger {
	l.Writer = w
	return l
}
func (l *Logger) SetFormatter(f Formatter) *Logger {
	l.Formatter = f
	return l
}
func (l *Logger) SetPrefixs(p ...Prefix) *Logger {
	l.Prefixs = p
	return l
}
func (l *Logger) AppendPrefixs(p ...Prefix) *Logger {
	l.Prefixs = append(l.Prefixs, p...)
	return l
}
func (l *Logger) Clone() *Logger {
	p := make([]Prefix, len(l.Prefixs))
	copy(p, l.Prefixs)
	return &Logger{
		Writer:    l.Writer,
		Formatter: l.Formatter,
		Prefixs:   p,
	}
}

func (l *Logger) SubLogger() *Logger {
	logger := l.Clone()
	logger.Writer = l
	return logger
}

func NewLogger() *Logger {
	return &Logger{}
}
func createLogger(w Writer, f Formatter, p ...Prefix) *Logger {
	return &Logger{
		Writer:    w,
		Formatter: f,
		Prefixs:   p,
	}
}

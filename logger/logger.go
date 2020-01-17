package logger

type Logger interface {
	Panicln(v ...interface{})
	Fatalln(v ...interface{})
	Errorln(v ...interface{})
	Warningln(v ...interface{})
	Infoln(v ...interface{})
	Traceln(v ...interface{})
	Debugln(v ...interface{})
	Println(v ...interface{})
}
type EmptyLogger struct{}

func (l EmptyLogger) Panicln(v ...interface{})   {}
func (l EmptyLogger) Fatalln(v ...interface{})   {}
func (l EmptyLogger) Errorln(v ...interface{})   {}
func (l EmptyLogger) Warningln(v ...interface{}) {}
func (l EmptyLogger) Infoln(v ...interface{})    {}
func (l EmptyLogger) Traceln(v ...interface{})   {}
func (l EmptyLogger) Debugln(v ...interface{})   {}
func (l EmptyLogger) Println(v ...interface{})   {}

type StartardLogger struct {
	EmptyLogger
	Writer    Writer
	Flag      Flag
	Formatter Formatter
	Prefixs   []Prefix
}

func (l *StartardLogger) Panicln(v ...interface{}) {
	if l.Flag.HasFlag(FlagDisablePanic) {
		return
	}
}

func (l *StartardLogger) Fatalln(v ...interface{}) {
	if l.Flag.HasFlag(FlagDisableFatal) {
		return
	}
}
func (l *StartardLogger) Errorln(v ...interface{}) {
	if l.Flag.HasFlag(FlagDisableError) {
		return
	}
}
func (l *StartardLogger) Println(v ...interface{}) {
	if l.Flag.HasFlag(FlagDisablePrint) {
		return
	}
}
func (l *StartardLogger) Warningln(v ...interface{}) {
	if l.Flag.HasFlag(FlagDisableWarning) {
		return
	}
}
func (l *StartardLogger) Infoln(v ...interface{}) {
	if l.Flag.HasFlag(FlagDisableInfo) {
		return
	}
}
func (l *StartardLogger) Traceln(v ...interface{}) {
	if l.Flag.HasFlag(FlagDisableTrace) {
		return
	}
}
func (l *StartardLogger) Debugln(v ...interface{}) {
	if l.Flag.HasFlag(FlagDisableDebug) {
		return
	}
}

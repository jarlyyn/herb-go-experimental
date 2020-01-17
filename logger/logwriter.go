package logger

import (
	"io"
)

type Writer interface {
	MustWriter(LogLevel,)
}

type StdWriter struct{}

func (w *StdWriter) PanicWriter() io.Writer {
	return Stderr
}
func (w *StdWriter) FatalWriter() io.Writer {
	return Stderr
}
func (w *StdWriter) ErrorWriter() io.Writer {
	return Stderr
}
func (w *StdWriter) PrintWriter() io.Writer {
	return Stdout
}
func (w *StdWriter) WarningWriter() io.Writer {
	return Stderr
}
func (w *StdWriter) InfoWriter() io.Writer {
	return Stdout
}
func (w *StdWriter) TraceWriter() io.Writer {
	return Stdout
}
func (w *StdWriter) DebugWriter() io.Writer {
	return Stdout
}

type LeveledWriter struct {
	StdWriter
	panicWriter   io.Writer
	fatalWriter   io.Writer
	errorWriter   io.Writer
	warningWriter io.Writer
	infoWriter    io.Writer
	traceWriter   io.Writer
	debugWriter   io.Writer
	printWriter   io.Writer
}

func (w *LeveledWriter) PanicWriter() io.Writer {
	if w.panicWriter != nil {
		return w.panicWriter
	}
	return Stderr
}
func (w *LeveledWriter) FatalWriter() io.Writer {
	if w.fatalWriter != nil {
		return w.fatalWriter
	}
	return Stderr
}
func (w *LeveledWriter) ErrorWriter() io.Writer {
	if w.errorWriter != nil {
		return w.errorWriter
	}
	return Stderr
}
func (w *LeveledWriter) PrintWriter() io.Writer {
	if w.printWriter != nil {
		return w.printWriter
	}
	return Stdout
}
func (w *LeveledWriter) WarningWriter() io.Writer {
	if w.warningWriter != nil {
		return w.warningWriter
	}
	return Stderr
}
func (w *LeveledWriter) InfoWriter() io.Writer {
	if w.infoWriter != nil {
		return w.infoWriter
	}
	return Stdout
}
func (w *LeveledWriter) TraceWriter() io.Writer {
	if w.traceWriter != nil {
		return w.traceWriter
	}
	return Stdout
}
func (w *LeveledWriter) DebugWriter() io.Writer {
	if w.debugWriter != nil {
		return w.debugWriter
	}
	return Stdout
}

func (w *LeveledWriter) SetPanicWriter(writer io.Writer) {
	w.panicWriter = writer
}
func (w *LeveledWriter) SetFatalWriter(writer io.Writer) {
	w.fatalWriter = writer
}
func (w *LeveledWriter) SetErrorWriter(writer io.Writer) {
	w.errorWriter = writer
}
func (w *LeveledWriter) SetPrintWriter(writer io.Writer) {
	w.printWriter = writer
}
func (w *LeveledWriter) SetWarningWriter(writer io.Writer) {
	w.warningWriter = writer
}
func (w *LeveledWriter) SetInfoWriter(writer io.Writer) {
	w.infoWriter = writer
}
func (w *LeveledWriter) SetTraceWriter(writer io.Writer) {
	w.traceWriter = writer
}
func (w *LeveledWriter) SetDebugWriter(writer io.Writer) {
	w.debugWriter = writer
}

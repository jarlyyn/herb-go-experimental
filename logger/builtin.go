package logger

var PanicPrefix = FixedPrefix("panic:")
var FatalPrefix = FixedPrefix("fatal:")
var ErrorPrefix = FixedPrefix("error:")
var PrintPrefix = FixedPrefix("print:")
var WarningPrefix = FixedPrefix("warning:")
var InfoPrefix = FixedPrefix("info:")
var TracePrefix = FixedPrefix("trace:")
var DebugPrefix = FixedPrefix("debug:")

var PanicLogger *Logger
var FatalLogger *Logger
var ErrorLogger *Logger
var PrintLogger *Logger
var WarningLogger *Logger
var InfoLogger *Logger
var TraceLogger *Logger
var DebugLogger *Logger

func ResetBuiltinLoggers() {
	PanicLogger = createLogger(Stderr, nil, DefaultTimePrefix, PanicPrefix)
	FatalLogger = createLogger(Stderr, nil, DefaultTimePrefix, FatalPrefix)
	ErrorLogger = createLogger(Stderr, nil, DefaultTimePrefix, ErrorPrefix)
	PrintLogger = createLogger(Stdout, nil, DefaultTimePrefix, PrintPrefix)
	WarningLogger = createLogger(Stdout, nil, DefaultTimePrefix, WarningPrefix)
	InfoLogger = createLogger(Stdout, nil, DefaultTimePrefix, InfoPrefix)
	TraceLogger = createLogger(Null, nil, DefaultTimePrefix, TracePrefix)
	DebugLogger = createLogger(Null, nil, DefaultTimePrefix, DebugPrefix)
}

var Panic = PanicLogger.Log
var Fatal = FatalLogger.Log
var Error = ErrorLogger.Log
var Print = PrintLogger.Log
var Warning = WarningLogger.Log
var Info = InfoLogger.Log
var Trace = TraceLogger.Log
var Debug = DebugLogger.Log

func EnableDevelopmengLoggers() {
	TraceLogger.Writer = Stdout
	DebugLogger.Writer = Stdout
}

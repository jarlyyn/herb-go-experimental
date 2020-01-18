package logger

var PanicLogger = createLogger(Stderr, DefaultFormatter, DefaultTimePrefix)
var FatalLogger = createLogger(Stderr, DefaultFormatter, DefaultTimePrefix)
var ErrorLogger = createLogger(Stderr, DefaultFormatter, DefaultTimePrefix)
var PrintLogger = createLogger(Stdout, DefaultFormatter, DefaultTimePrefix)
var WarningLogger = createLogger(Stdout, DefaultFormatter, DefaultTimePrefix)
var InfoLogger = createLogger(Stdout, DefaultFormatter, DefaultTimePrefix)
var TraceLogger = createLogger(Null, DefaultFormatter, DefaultTimePrefix)
var DebugLogger = createLogger(Null, DefaultFormatter, DefaultTimePrefix)

var Panic = PanicLogger.Log
var Fatal = FatalLogger.Log
var Error = ErrorLogger.Log
var Print = PrintLogger.Log
var Warning = WarningLogger.Log
var Info = InfoLogger.Log
var Trace = TraceLogger.Log
var Debug = DebugLogger.Log

package logger

type Flag int64

const (
	FlagDisablePanic   = Flag(1)
	FlagDisableFatal   = Flag(2)
	FlagDisableError   = Flag(4)
	FlagDisablePrint   = Flag(8)
	FlagDisableWarning = Flag(16)
	FlagDisableInfo    = Flag(32)
	FlagDisableTrace   = Flag(64)
	FlagDisableDebug   = Flag(128)
)

var FlagProductionMode Flag = FlagDisableWarning | FlagDisableInfo | FlagDisableTrace | FlagDisableDebug

func (f Flag) SetFlag(flag Flag) Flag {
	return f | flag
}

func (f Flag) UnsetFlag(falg Flag) Flag {
	return f ^ falg
}

func (f Flag) HasFlag(flag Flag) bool {
	return f&flag != 0
}

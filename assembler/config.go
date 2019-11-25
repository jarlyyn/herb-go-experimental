package assembler

type Config struct {
	Checkers  []TypeChecker
	Unifiers  Unifiers
	Tagname   string
	TagParser func(value string) (*Tag, error)
}

func NewConfig() *Config {
	return &Config{
		TagParser: ParseTag,
	}
}

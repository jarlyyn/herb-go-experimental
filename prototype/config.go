package prototype

type Config struct {
	Version      *Version
	Dependencies map[string]*Version
}

func NewConfig() *Config {
	return &Config{
		Dependencies: map[string]*Version{},
	}
}

package preset

type Driver interface {
	Preset(string, Config) error
}

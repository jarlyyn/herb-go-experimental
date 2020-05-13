package prototype

type Dependency interface {
	ID() string
	Version() *Version
}

type DependencyID string

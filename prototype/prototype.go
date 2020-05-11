package prototype

type Version struct {
	Major int
	Minor int
}
type Prototype struct {
	Version    *Version
	Root       *Component
	Components []Type
}

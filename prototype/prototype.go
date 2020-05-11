package prototype

type Version struct {
	Major int
	Minor int
}
type Prototype struct {
	Version     Version
	Root        *Component
	DataSources *DataSources
	Components  []Type
}

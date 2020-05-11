package prototype

type Source string
type Field string

type Data struct {
	Source Source
	Field  Field
}

type DataSource struct {
	Source Source
	URL    *URL
}

type DataSources []*DataSource

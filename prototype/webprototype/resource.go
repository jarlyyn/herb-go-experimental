package webprototype

type ResourceType string

var ResourceTypeJavascript = ResourceType("javascript")
var ResourceTypeCSS = ResourceType("CSS")

type Resource struct {
	Name    string
	Type    ResourceType
	Content string
	Target  string
}

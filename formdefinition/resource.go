package formdefinition

type ResourceType string

var ResourceTypeJavaScript = ResourceType("javascript")
var ResourceTypeCSS = ResourceType("style")

type Resource struct {
	Name    string
	Type    ResourceType
	Options map[string]string
	Inline  bool
	Target  string
}

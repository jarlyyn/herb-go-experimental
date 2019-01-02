package notificationmanager

type complied func(data interface{}) ([]byte, error)

type Engine interface {
	Compile(template string) (complied, error)
}

const FieldTypeString = "string"
const FieldTypeJSONList = "jsonlist"
const FieldTypeJSONMap = "jsonmap"

type Field struct {
	Name        string
	Type        string
	Enabled     bool
	Description string
	Template    string
}

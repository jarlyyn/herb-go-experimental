package unmarshaler

import (
	"strings"
)

//Tag struct field tag struct
type Tag struct {
	//Name parsed name tag
	Name string
	//Flags field flags
	Flags map[string]string
	//Ignored field ignored  flag
	Ignored bool
}

//NewTag create new tag struct
func NewTag() *Tag {
	return &Tag{
		Flags: map[string]string{},
	}
}

//ParseTag default parse tag func
//Parse field tag with given config
//Return parsed tag and any error if rasised.
func ParseTag(c *Config, value string) (*Tag, error) {
	t := NewTag()
	value = strings.TrimSpace(value)
	if value == "" {
		return t, nil
	}
	if value == "-" {
		t.Ignored = true
		return t, nil
	}
	v := strings.Split(value, ",")
	t.Name = strings.TrimSpace(v[0])
	l := len(v)
	for i := 1; i < l; i++ {
		k := strings.TrimSpace(v[i])
		if k != "" {
			t.Flags[k] = "1"
		}
	}
	return t, nil
}

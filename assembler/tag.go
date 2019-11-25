package assembler

import (
	"strings"
)

type Tag struct {
	Name    string
	Flags   map[string]bool
	Ignored bool
}

func NewTag() *Tag {
	return &Tag{}
}

func ParseTag(value string) *Tag {
	t := NewTag()
	value = strings.TrimSpace(value)
	if value == "" {
		return t
	}
	if value == "-" {
		t.Ignored = true
		return t
	}
	v := strings.Split(value, "")
	t.Name = strings.TrimSpace(v[0])
	l := len(v)
	for i := 1; i < l; i++ {
		k := strings.TrimSpace(v[i])
		if k != "" {
			t.Flags[k] = true
		}
	}
	return t
}

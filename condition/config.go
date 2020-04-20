package condition

import "net/http"

type Config struct {
	Type       string
	Config     func(v interface{}) error
	Not        bool
	Or         bool
	Disabled   bool
	Conditions []*Config
}

func (c *Config) Create(creator ConditionCreator) (Condition, error) {
	var err error
	pc := NewPlainCondition()
	pc.Condition, err = creator.CreateCondition(c.Type, c.Config)
	if err != nil {
		return nil, err
	}
	pc.Not = c.Not
	pc.Or = c.Or
	pc.Disabled = c.Disabled
	for k := range c.Conditions {
		condition, err := c.Conditions[k].Create(creator)
		if err != nil {
			return nil, err
		}
		pc.Conditions = append(pc.Conditions, condition)
	}
	return pc, nil
}

type PlainCondition struct {
	Condition  Condition
	Not        bool
	Or         bool
	Disabled   bool
	Conditions []Condition
}

func NewPlainCondition() *PlainCondition {
	return &PlainCondition{
		Conditions: []Condition{},
	}
}
func (c *PlainCondition) CheckRequest(r *http.Request) (bool, error) {
	var result bool
	var err error
	if c.Disabled {
		return false, nil
	}

	if len(c.Conditions) != 0 {
		conditions := append([]Condition{c.Condition}, c.Conditions...)
		if c.Or {
			result, err = Or(r, conditions...)
		} else {
			result, err = And(r, conditions...)
		}
	} else {
		result, err = c.Condition.CheckRequest(r)
	}
	if err != nil {
		return false, err
	}
	if c.Not {
		result = !result
	}
	return result, nil
}

type ConditionCreator interface {
	CreateCondition(string, func(interface{}) error) (Condition, error)
}

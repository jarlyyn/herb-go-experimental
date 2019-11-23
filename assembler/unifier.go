package assembler

type Unifier interface {
	Unify(a *Assembler, v interface{}) (bool, error)
}

type Unifiers map[interface{}][]Unifier

func (u *Unifiers) Unify(a *Assembler, v interface{}) error {
	tp, err := a.CheckType()
	if err != nil {
		return err
	}
	unifiers, ok := (*u)[tp]
	if ok == false {
		return nil
	}
	for k := range unifiers {
		result, err := unifiers[k].Unify(a, v)
		if err != nil {
			return err
		}
		if result {
			return nil
		}
	}
	return nil
}

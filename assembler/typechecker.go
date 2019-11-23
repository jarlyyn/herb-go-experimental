package assembler

type TypeChecker struct {
	Type      interface{}
	CheckType func(a *Assembler) (bool, error)
}

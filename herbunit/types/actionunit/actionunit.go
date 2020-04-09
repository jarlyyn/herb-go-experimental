package actionunit

type ActionUnit struct {
}

func (u *ActionUnit) InitUnitType() error {
	return nil
}
func (u *ActionUnit) Summary() (interface{}, error) {
	return nil, nil
}
func (u *ActionUnit) PlainSummary() (string, error) {
	return "", nil
}
func (u *ActionUnit) Keyword() string {
	return "action"
}

package memberunit

type MemberUnit struct {
}

func (u *MemberUnit) InitUnitType() error {
	return nil
}
func (u *MemberUnit) Summary() (interface{}, error) {
	return nil, nil
}
func (u *MemberUnit) PlainSummary() (string, error) {
	return "", nil
}
func (u *MemberUnit) Keyword() string {
	return "member"
}

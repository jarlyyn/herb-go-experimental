package cacheunit

type CacheUnit struct {
}

func (u *CacheUnit) InitUnitType() error {
	return nil
}
func (u *CacheUnit) Summary() (interface{}, error) {
	return nil, nil
}
func (u *CacheUnit) PlainSummary() (string, error) {
	return "", nil
}
func (u *CacheUnit) Keyword() string {
	return "cache"
}

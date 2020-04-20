package condition

import "net/http"

type Condition interface {
	CheckRequest(*http.Request) (bool, error)
}

func Not(r *http.Request, c Condition) (bool, error) {
	ok, err := c.CheckRequest(r)
	if err != nil {
		return false, err
	}
	return !ok, nil
}

func And(r *http.Request, c ...Condition) (bool, error) {
	for k := range c {
		ok, err := c[k].CheckRequest(r)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

func Or(r *http.Request, c ...Condition) (bool, error) {
	for k := range c {
		ok, err := c[k].CheckRequest(r)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

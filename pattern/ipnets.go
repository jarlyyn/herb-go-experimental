package pattern

import (
	"net"
	"net/http"
)

func GetRequestIPAddress(r *http.Request) string {

}
func GetRequestIP(r *http.Request) net.IP {

}

type IPNets []*net.IPNet

func (i IPNets) IsEmpty() bool {
	return len(i) == 0
}

func (i IPNets) Match(r *http.Request) (bool, error) {
	if i.IsEmpty() {
		return false, nil
	}
	ip := GetRequestIP(r)
	if ip == nil {
		return false, nil
	}
	for k := range i {
		if i[k].Contains(ip) {
			return true, nil
		}
	}
	return false, nil
}

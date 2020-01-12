package pattern

import (
	"context"
	"net"
	"net/http"
)

func GetRequestIPAddress(r *http.Request) string {
	v := r.Context().Value(ContextKeyIPAddress)
	if v != nil {
		return v.(string)
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	ctx := context.WithValue(r.Context(), ContextKeyIPAddress, ip)
	req := r.WithContext(ctx)
	*r = *req
	return ip

}
func GetRequestIP(r *http.Request) net.IP {
	v := r.Context().Value(ContextKeyIP)
	if v != nil {
		return v.(net.IP)
	}
	ip := net.ParseIP(GetRequestIPAddress(r))
	ctx := context.WithValue(r.Context(), ContextKeyIP, ip)
	req := r.WithContext(ctx)
	*r = *req
	return ip
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

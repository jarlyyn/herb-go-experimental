package healthmonitor

import "net/http"

type Config struct {
	Disabled                 bool
	Name                     string
	Description              string
	URL                      string
	Headers                  http.Header
	StatusCodeOnly           bool
	DaysDomainExpiredWarning int
	DaysDomainExpiredError   int
	DaysCertExpiredWarning   int
	DaysCertExpiredError     int
}

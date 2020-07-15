package token

import (
	"time"
)

type Refresher interface {
	Refresh(ID, *time.Time)
}

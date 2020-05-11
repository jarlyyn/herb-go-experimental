package prototype

import (
	"github.com/herb-go/herb/ui"
)

type ConsumerData struct {
	Prototype    *Prototype
	Translations ui.Collection
}

type Consumer interface {
	Consume(ConsumerData) error
}

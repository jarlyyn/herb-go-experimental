package uniqueid

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"math/rand"
	"sync/atomic"
	"time"
)

//SimpleID simple id driver
type SimpleID struct {
	Current *uint32
	Suff    string
}

//GenerateID generate unique id.
//Return  generated id and any error if rasied.
func (i *SimpleID) GenerateID() (string, error) {
	buf := bytes.NewBuffer(nil)
	ts := time.Now().UnixNano()
	err := binary.Write(buf, binary.BigEndian, ts)
	if err != nil {
		return "", err
	}
	c := atomic.AddUint32(i.Current, 1)
	err = binary.Write(buf, binary.BigEndian, c)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(buf.Bytes()) + i.Suff, nil
}

// NewSimpleID create new simpleid driver
func NewSimpleID() *SimpleID {
	time.Sleep(time.Millisecond)
	var c = rand.Uint32()
	return &SimpleID{
		Current: &c,
		Suff:    "",
	}
}

//SimpleIDFactory simple id driver factory
func SimpleIDFactory(conf Config, prefix string) (Driver, error) {
	i := NewSimpleID()
	conf.Get(prefix+"Suff", &i.Suff)
	return i, nil

}
func init() {
	Register("simpleid", SimpleIDFactory)
}

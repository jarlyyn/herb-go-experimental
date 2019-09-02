package uniqueid

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"sync/atomic"
	"time"
)

type SimpleID struct {
	Current *uint32
	Suff    string
}

func (i *SimpleID) Generate() (string, error) {
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

func NewSimpleID() *SimpleID {
	var c = uint32(0)
	return &SimpleID{
		Current: &c,
		Suff:    "",
	}
}

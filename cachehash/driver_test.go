package cachehash

import (
	"encoding/binary"
	"testing"
)

func TestInt64(t *testing.T) {
	var v int64
	v = v - 1
	var bytes = make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(v))
	if int64(binary.BigEndian.Uint64(bytes)) != -1 {
		t.Fatal(bytes)
	}
}

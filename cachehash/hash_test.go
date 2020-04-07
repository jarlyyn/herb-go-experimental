package cachehash

import (
	"bytes"
	"testing"
)

func TestHash(t *testing.T) {
	var status *Status
	h := NewHash()
	if !h.isEmpty() {
		t.Fatal(h)
	}
	d1 := NewData("key1", 1235, []byte("key1"))
	status = h.set(d1, 1234)
	if status.Changed != true ||
		status.Delta != 4 ||
		status.FirstExpired != 1235 ||
		status.LastExpired != 1235 ||
		status.Size != 4 {
		t.Fatal(status)
	}
	status = h.set(d1, 1234)
	if status.Changed != true ||
		status.Delta != 0 ||
		status.FirstExpired != 1235 ||
		status.LastExpired != 1235 ||
		status.Size != 4 {
		t.Fatal(status)
	}
	d2 := NewData("key2", 1236, []byte("keyd2"))
	status = h.set(d2, 1234)
	if status.Changed != true ||
		status.Delta != 5 ||
		status.FirstExpired != 1235 ||
		status.LastExpired != 1236 ||
		status.Size != 9 {
		t.Fatal(status)
	}
	d3 := NewData("key3", 1000, []byte("keyd2"))
	status = h.set(d3, 1234)
	if status.Changed != true ||
		status.Delta != 0 ||
		status.FirstExpired != 1235 ||
		status.LastExpired != 1236 ||
		status.Size != 9 {
		t.Fatal(status)
	}
	v1 := h.get(d1.Key, 1234)
	if bytes.Compare(v1.Data, d1.Data) != 0 {
		t.Fatal(v1)
	}
	v1e := h.get(d1.Key, 1250)
	if v1e != nil {
		t.Fatal(v1e)
	}
	vne := h.get("notexist", 1234)
	if vne != nil {
		t.Fatal(v1e)
	}
}

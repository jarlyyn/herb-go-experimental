package logger

import "testing"

func TestCsvFormat(t *testing.T) {
	f := &CsvFormatter{}
	b, err := f.Format("hi", "\"123", "o,k")
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `hi,"""123","o,k"` {
		t.Fatal(string(b))
	}

}

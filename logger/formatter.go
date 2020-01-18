package logger

import (
	"bytes"
	"encoding/csv"
	"fmt"
)

type Formatter interface {
	Format(v ...interface{}) ([]byte, error)
}

type CsvFormatter struct {
}

func (f *CsvFormatter) Format(v ...interface{}) ([]byte, error) {
	data := make([]string, len(v))
	for i := range v {
		data[i] = fmt.Sprint(v[i])
	}
	buf := bytes.NewBuffer(nil)
	w := csv.NewWriter(buf)
	err := w.Write(data)
	if err != nil {
		return nil, err
	}
	w.Flush()
	output := bytes.TrimRight(buf.Bytes(), "\n")
	return output, nil
}

var DefaultFormatter = &CsvFormatter{}

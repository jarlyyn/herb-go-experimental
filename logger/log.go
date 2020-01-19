package logger

import (
	"io"
	"os"
	"sync"
)

type Writer interface {
	Open() error
	Close() error
	Write(p []byte) (n int, err error)
}

type NullWriter struct {
}

func (o *NullWriter) Open() error {
	return nil
}
func (o *NullWriter) Close() error {
	return nil
}
func (o *NullWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

type IOWriter struct {
	io.Writer
}

func (o *IOWriter) Open() error {
	return nil
}
func (o *IOWriter) Close() error {
	return nil
}

var Stdout Writer = &IOWriter{
	Writer: os.Stdout,
}

var Stderr Writer = &IOWriter{
	Writer: os.Stderr,
}

var Null Writer = &NullWriter{}

type FileWriter struct {
	lock sync.RWMutex
	Path string
	Mode os.FileMode
	file *os.File
}

func (o *FileWriter) Open() error {
	o.lock.Lock()
	defer o.lock.Unlock()
	file, err := os.OpenFile(o.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, o.Mode)
	if err != nil {
		return err
	}
	o.file = file
	return nil
}
func (o *FileWriter) Close() error {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.file = nil
	return o.file.Close()
}
func (o *FileWriter) Write(p []byte) (n int, err error) {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.file.Write(p)
}

func NewFileWriter(path string, mode os.FileMode) *FileWriter {
	return &FileWriter{
		Path: path,
		Mode: mode,
	}
}

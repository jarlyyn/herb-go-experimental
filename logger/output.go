package logger

import (
	"io"
	"os"
)

type Output interface {
	MustOpen()
	MustClose()
	MustWrite([]byte)
}

type nullOutput struct {
}

func (o *nullOutput) MustOpen() {
}
func (o *nullOutput) MustClose() {
}
func (o *nullOutput) MustWrite(p []byte) {
}

type IOWriterOutput struct {
	IOWriter io.Writer
}

func (o *IOWriterOutput) MustOpen() {
}
func (o *IOWriterOutput) MustClose() {
}

func (o *IOWriterOutput) MustWrite(p []byte) {
	o.IOWriter.Write(p)
}

var Stdout Output = &IOWriterOutput{
	IOWriter: os.Stdout,
}

var Stderr Output = &IOWriterOutput{
	IOWriter: os.Stderr,
}

var Null Output = &nullOutput{}

type FileOutput struct {
	Path string
	Mode os.FileMode
	file *os.File
}

func (o *FileOutput) MustOpen() {
	file, err := os.OpenFile(o.Path, os.O_CREATE|os.O_APPEND, o.Mode)
	if err != nil {
		panic(err)
	}
}
func (o *FileOutput) MustClose() {
	err := o.file.Close()
	if err != nil {
		panic(err)
	}
}
func (o *FileOutput) MustWrite(p []byte) {
	_, err := o.file.Write(p)
	if err != nil {
		panic(err)
	}
}

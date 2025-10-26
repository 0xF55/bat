package main

import (
	"log"
	"os"
)

type oWriter struct {
	name   string // file name
	lines  int    // wroted lines
	stream *os.File
}

func NewWriter() *oWriter {
	writer := &oWriter{}
	writer.lines = 0
	writer.name = OutputFile
	f, err := os.OpenFile(writer.name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open|create %s: %s", writer.name, err.Error())
		return nil
	}
	writer.stream = f
	return writer
}

func (w *oWriter) Write(data string) {
	defer w.stream.Sync()

	data_bytes := []byte(data)
	data_bytes = append(data_bytes, 0x0A) // new line
	w.stream.Write(data_bytes)
	w.lines += 1
}

package main

import (
	"bytes"
	"log"
	"os"
	"sync"
)

// for main writer
type oWriter struct {
	mux    sync.RWMutex
	name   string   // file name
	lines  int      // final wroted lines
	stream *os.File // file
}

// for every bat file
type bWriter struct {
	buffer bytes.Buffer
	lines  int
}

func (bW *bWriter) Write(data string) {
	bW.buffer.WriteString(data + "\n")
	bW.lines++
}

/* These two functions for the main writer */

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

func (w *oWriter) Write(bW *bWriter) {
	w.mux.Lock()
	defer w.mux.Unlock()
	data_bytes := bW.buffer.Bytes()
	w.stream.Write(data_bytes)
	w.lines += bW.lines
}

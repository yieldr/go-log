package log

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

type bufSink struct {
	out    *bufio.Writer
	file   *os.File
	format string
	fields []string
	mux    sync.Mutex
}

func (sink *bufSink) open(name string) (err error) {
	sink.file, err = os.OpenFile(name, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		sink.file, err = os.Create(name)
		if err != nil {
			return errors.New("logging: unable to open or create file")
		}
	}
	return
}

func (sink *bufSink) close() error {
	err := sink.flush()
	if err != nil {
		return err
	}
	return sink.file.Close()
}

func (sink *bufSink) flush() error {
	return sink.out.Flush()
}

func (sink *bufSink) reload() (err error) {
	name := sink.file.Name()
	sink.close()
	err = sink.open(name)
	if err == nil {
		sink.out = bufio.NewWriter(sink.file)
	}
	return
}

func (sink *bufSink) daemon() {
	flush := time.NewTicker(flushInterval)
	reload := time.NewTicker(reloadInterval)
	for {
		select {
		case <-flush.C:
			sink.Flush()
		case <-reload.C:
			sink.Reload()
		}
	}
}

func (sink *bufSink) Log(fields Fields) {
	vals := make([]interface{}, len(sink.fields))
	for i, field := range sink.fields {
		var ok bool
		vals[i], ok = fields[field]
		if !ok {
			vals[i] = "???"
		}
	}

	sink.mux.Lock()
	defer sink.mux.Unlock()
	fmt.Fprintf(sink.out, sink.format, vals...)
}

func (sink *bufSink) Write(b []byte) (int, error) {
	sink.mux.Lock()
	defer sink.mux.Unlock()
	return sink.out.Write(b)
}

// Closes and reopens the output file, in order to momentarily release it's file
// handle. Typically this functionality is combined with a SIGHUP system signal.
// Before reloading, the content of the buffer is flushed.
func (sink *bufSink) Reload() error {
	sink.mux.Lock()
	defer sink.mux.Unlock()
	return sink.reload()
}

// Flushes the buffer to disk.
func (sink *bufSink) Flush() error {
	sink.mux.Lock()
	defer sink.mux.Unlock()
	return sink.flush()
}

// Closes any open file handles used by the Sink.
func (sink *bufSink) Close() error {
	sink.mux.Lock()
	defer sink.mux.Unlock()
	return sink.close()
}

// Returns a new Sink able to buffer output and periodically flush to disk.
func FileSink(name string, format string, fields []string) (*bufSink, error) {
	sink := &bufSink{
		format: format,
		fields: fields,
	}
	err := sink.open(name)
	if err != nil {
		return nil, err
	}
	sink.out = bufio.NewWriter(sink.file)
	go sink.daemon()
	return sink, nil
}

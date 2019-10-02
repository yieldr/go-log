package logstream

import (
	"testing"
	"time"

	"github.com/yieldr/go-log/log"
)

func TestNewLogstream(t *testing.T) {
	l := New(nil, time.Second, log.BasicFormat, log.BasicFields)
	if l == nil {
		t.Fatal("NewLogstream failed")
	}
}

func TestLogStreamLog(t *testing.T) {
	stream := new(StreamMock)

	l := &Logstream{
		format: log.BasicFormat,
		fields: log.BasicFields,
		writer: NewStreamWriter(stream),
	}

	fields := log.Fields{
		"time":     func() interface{} { return "now" },
		"priority": func() interface{} { return "INFO" },
		"message":  func() interface{} { return "foo" },
	}
	l.Log(fields)
	l.Log(fields)
	l.Flush()

	if stream.buf.String() != "now [INFO] foo\nnow [INFO] foo\n" {
		t.Error("Logstream log buffer not matched")
	}
}

func TestLogStreamRun(t *testing.T) {
	stream := new(StreamMock)

	l := &Logstream{
		interval: time.Second * 3,
		format:   log.BasicFormat,
		fields:   log.BasicFields,
		errChan:  make(chan error),
		stopChan: make(chan bool),
		writer:   NewStreamWriter(stream),
	}

	// run in background
	go l.Run()

	// log data
	fields := log.Fields{
		"time":     func() interface{} { return "now" },
		"priority": func() interface{} { return "INFO" },
		"message":  func() interface{} { return "foo" },
	}
	l.Log(fields)

	// data is flushed every 5s
	time.Sleep(time.Second * 5)

	if 0 != l.writer.buf.getSize() {
		t.Error("writer buffer size should be 0.")
	}

	if "now [INFO] foo\n" != stream.buf.String() {
		t.Error("writer buffer content not match")
	}

	// stop
	l.Stop()
}

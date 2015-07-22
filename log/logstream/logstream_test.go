package logstream

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yieldr/go-log/log"
)

func TestNewLogstream(t *testing.T) {
	l := New(nil, time.Second, log.BasicFormat, log.BasicFields)
	assert.NotNil(t, l)
}

func TestLogStreamLog(t *testing.T) {
	stream := new(StreamMock)
	stream.On("Put", mock.Anything).Return(new(StreamResponseMock), nil)

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

	assert.Equal(t, "now [INFO] foo\n", string(l.writer.buffer[0]))
}

func TestLogStreamRun(t *testing.T) {
	stream := new(StreamMock)
	stream.On("Put", mock.Anything).Return(new(StreamResponseMock), nil)

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
	for i := 0; i < len(l.writer.buffer); i++ {
		assert.Nil(t, l.writer.buffer[i])
	}
	assert.Equal(t, "now [INFO] foo\n", stream.buf.String())

	// stop
	l.Stop()
}

package logstream

import (
	"fmt"
	"sync"
	"time"

	"github.com/yieldr/go-log/log"
)

// Logstream is a special case of sink which writes to a stream.
type Logstream struct {
	interval time.Duration
	format   string
	fields   []string

	errChan  chan error
	stopChan chan bool
	mux      sync.Mutex

	stream Stream
	writer *StreamWriter
}

// New returns a new Logstream using the supplied arguments.
func New(stream Stream, interval time.Duration, format string, fields []string) *Logstream {
	return &Logstream{
		interval: interval,
		format:   format,
		fields:   fields,
		stream:   stream,

		errChan:  make(chan error),
		stopChan: make(chan bool),

		writer: NewStreamWriter(stream),
	}
}

// Flush forces data to be written into stream.
func (l *Logstream) Flush() error {
	l.mux.Lock()
	defer l.mux.Unlock()
	return l.flush()
}

// Close a logstream.
func (l *Logstream) Close() error {
	l.mux.Lock()
	defer l.mux.Unlock()
	return l.close()
}

// Reload reloads a logstream.
func (l *Logstream) Reload() error {
	l.mux.Lock()
	defer l.mux.Unlock()
	return l.reload()
}

// Log satisfies the log.Sink interface so it can be supplied as an argument to
// log.New(). It writes the log to the internal buffer, using the format and
// fields.
func (l *Logstream) Log(fields log.Fields) {
	l.mux.Lock()
	defer l.mux.Unlock()

	vals := make([]interface{}, len(l.fields))
	for i, field := range l.fields {
		if fn, ok := fields[field]; ok {
			vals[i] = fn()
		} else {
			vals[i] = "???"
		}
	}

	fmt.Fprintf(l.writer, l.format, vals...)
}

// Write writes p to the internal buffer.
func (l *Logstream) Write(p []byte) (int, error) {
	l.mux.Lock()
	defer l.mux.Unlock()
	return l.writer.Write(p)
}

// Run is usually used as a deamon. All the buffered data is flushed periodically
// until it is stopped.
func (l *Logstream) Run() {
	flush := time.NewTicker(l.interval)
	for {
		select {
		case <-flush.C:
			if err := l.Flush(); err != nil {
				l.errChan <- err
			}
		case <-l.stopChan:
			flush.Stop()
			return
		}
	}
}

// Stop ends the execution of Run.
func (l *Logstream) Stop() {
	l.stopChan <- true
}

// Done returns a channel which will receive when the Run method has returned.
func (l *Logstream) Done() <-chan bool {
	return l.stopChan
}

// Error returns a channel which will receive an error if one was encountered
// during a rotate or flush operation. It's up to the caller to handle the error
// so Logrotate will keep running until Stop is called.
func (l *Logstream) Error() <-chan error {
	return l.errChan
}

// flush buffered data in writer.
func (l *Logstream) flush() error {
	return l.writer.Flush()
}

// close the writer.
func (l *Logstream) close() error {
	return l.writer.Close()
}

// reload the writer.
func (l *Logstream) reload() error {
	if err := l.flush(); err != nil {
		return err
	}
	if err := l.writer.Close(); err != nil {
		return err
	}

	l.writer = NewStreamWriter(l.stream)
	return nil
}

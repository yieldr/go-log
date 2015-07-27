package logstream

import (
	"bytes"
	"errors"

	"github.com/stretchr/testify/mock"
)

// StreamRecord represents a record to be sent to Stream.
type StreamRecord []byte

// Stream defines a remove stream.
type Stream interface {
	Put([]StreamRecord) (StreamResponse, error)
	Close() error
}

// StreamResponse defines a repsonse from a remote stream.
type StreamResponse interface {
	GoString() string
	String() string
}

// StreamWriter is used to write data into a stream.
// Data is buffered at first, data is flushed into the stream when buffer is full.
// It also has a Write() function to implement io.Writer interface.
type StreamWriter struct {
	stream Stream

	maxBufferItems int
	maxBufferSize  int

	buf *recordBuffer
}

// NewStreamWriter creates a new stream writer.
func NewStreamWriter(s Stream) *StreamWriter {
	return &StreamWriter{
		stream:         s,
		maxBufferItems: 500,
		maxBufferSize:  1024 * 1024, //1MB

		buf: newRecordBuffer(500),
	}
}

// Write writes len(p) bytes from p to the underlying data stream.
// It returns the number of bytes written from p (0 <= n <= len(p))
// and any error encountered that caused the write to stop early.
func (s *StreamWriter) Write(p []byte) (n int, err error) {
	n = len(p)

	if s.buf.getItems() >= s.maxBufferItems || s.buf.getSize() >= s.maxBufferSize {
		if err = s.Flush(); err != nil {
			return 0, err
		}
	}

	// Do not just retain or modify p, copy it!
	// See:http://golang.org/pkg/io/#Writer
	data := make([]byte, n)
	copy(data, p)

	if err = s.buf.append(StreamRecord(data)); err != nil {
		return 0, err
	}

	return
}

// Flush buffered data into the stream.
func (s *StreamWriter) Flush() error {
	_, err := s.stream.Put(s.buf.getRecords())
	s.Reset()
	return err
}

// Reset the internal fields in s.
func (s *StreamWriter) Reset() {
	s.buf.reset()
}

// Close the stream in s.
func (s *StreamWriter) Close() error {
	return s.stream.Close()
}

// recordBuffer uses a pre-allocated, fixed-size slice to buffer StreamRecord.
type recordBuffer struct {
	records []StreamRecord
	pos     int
	size    int
}

// newWriterBuffer returns a new initialized recordBuffer.
func newRecordBuffer(maxRecords int) *recordBuffer {
	return &recordBuffer{
		records: make([]StreamRecord, maxRecords),
		pos:     0,
		size:    0,
	}
}

// reset resets the current position and size of records.
func (r *recordBuffer) reset() {
	r.pos = 0
	r.size = 0
}

// append appends r into the r.records.
func (r *recordBuffer) append(s StreamRecord) error {
	if r.pos >= len(r.records) {
		return errors.New("reach the end of buffer.")
	}

	r.records[r.pos] = s
	r.pos++
	r.size += len(s)
	return nil
}

// getSize returns the byte size of r.
func (r *recordBuffer) getSize() int {
	return r.size
}

// getItems returns the number of records in r.
func (r *recordBuffer) getItems() int {
	return r.pos
}

// getRecords returns a new slice of all the stored records in r.
func (r *recordBuffer) getRecords() []StreamRecord {
	return r.records[0:r.pos]
}

// StreamResponseMock is a mock for StreamResponse.
type StreamResponseMock struct {
	StreamResponse
	mock.Mock
}

// StreamMock is a mock for Stream.
type StreamMock struct {
	Stream
	mock.Mock
	buf bytes.Buffer
}

// Put is mocked to return assigned values. It also records
// any input data.
func (s *StreamMock) Put(records []StreamRecord) (StreamResponse, error) {
	for _, r := range records {
		s.buf.Write(r)
	}

	args := s.Called(records)
	return args.Get(0).(StreamResponse), args.Error(1)
}

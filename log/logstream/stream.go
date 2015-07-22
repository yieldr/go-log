package logstream

import (
	"bytes"

	"github.com/stretchr/testify/mock"
)

const (
	maxRecordEntryNum   = 500
	maxRecordBufferSize = 1024 * 1024
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

	buffer     []StreamRecord
	bufferSize int
	bufferIdx  int

	maxBufferItems int
	maxBufferSize  int
}

// NewStreamWriter creates a new stream writer.
func NewStreamWriter(s Stream) *StreamWriter {
	return &StreamWriter{
		stream: s,

		bufferSize: 0,
		bufferIdx:  0,
		buffer:     make([]StreamRecord, maxRecordEntryNum),

		maxBufferItems: maxRecordEntryNum,
		maxBufferSize:  maxRecordBufferSize,
	}
}

// Write writes len(p) bytes from p to the underlying data stream.
// It returns the number of bytes written from p (0 <= n <= len(p))
// and any error encountered that caused the write to stop early.
func (s *StreamWriter) Write(p []byte) (n int, err error) {
	n = len(p)

	if n == 0 {
		return
	}

	if s.bufferSize >= s.maxBufferSize || s.bufferIdx >= s.maxBufferItems {
		err = s.Flush()
	}

	s.buffer[s.bufferIdx] = StreamRecord(p)
	s.bufferSize += n
	s.bufferIdx++

	return
}

// Flush buffered data into the stream.
func (s *StreamWriter) Flush() error {
	_, err := s.stream.Put(s.buffer)
	s.Reset()
	return err
}

// Reset the internal fields in s.
func (s *StreamWriter) Reset() {
	// clean up slice, instead of making a new slice
	for i := 0; i < len(s.buffer); i++ {
		s.buffer[i] = nil
	}
	s.bufferSize = 0
	s.bufferIdx = 0
}

// Close the stream in s.
func (s *StreamWriter) Close() error {
	return s.stream.Close()
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

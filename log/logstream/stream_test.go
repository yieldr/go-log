package logstream

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

func TestStreamWriter(t *testing.T) {

	stream := new(StreamMock)
	stream.On("Put", mock.Anything).Return(new(StreamResponseMock), nil)

	w := &StreamWriter{
		stream: stream,

		bufferSize: 0,
		bufferIdx:  0,
		buffer:     make([]StreamRecord, 2),

		maxBufferItems: 2,
		maxBufferSize:  5,
	}

	w.Write([]byte{1, 2, 3})
	w.Write([]byte{1, 2, 3})
	w.Write([]byte{1, 2, 3})

	assert.Equal(t, []byte{1, 2, 3, 1, 2, 3}, stream.buf.Bytes())
	assert.Equal(t, []StreamRecord{StreamRecord([]byte{1, 2, 3}), nil}, w.buffer)
	assert.Equal(t, 1, w.bufferIdx)
	assert.Equal(t, 3, w.bufferSize)
}

func TestStreamWriterWriteNoError(t *testing.T) {

	stream := new(StreamMock)
	stream.On("Put", mock.Anything).Return(new(StreamResponseMock), nil)

	tests := []struct {
		writer            *StreamWriter
		input             []byte
		expectedStreamBuf []byte
		expectedWriterBuf []StreamRecord
	}{
		// flush is not triggered
		{
			writer: &StreamWriter{
				stream:         stream,
				buffer:         []StreamRecord{nil, nil},
				bufferSize:     0,
				bufferIdx:      0,
				maxBufferItems: 2,
				maxBufferSize:  4,
			},
			input:             []byte{1, 2},
			expectedStreamBuf: nil,
			expectedWriterBuf: []StreamRecord{
				StreamRecord([]byte{1, 2}),
				nil,
			},
		},
		// buffer size exceeds, flush is triggered
		{
			writer: &StreamWriter{
				stream: stream,
				buffer: []StreamRecord{
					StreamRecord([]byte{1, 2, 3, 4, 5}),
					nil,
				},
				bufferSize:     5,
				bufferIdx:      1,
				maxBufferItems: 2,
				maxBufferSize:  4,
			},
			input:             []byte{6, 7},
			expectedStreamBuf: []byte{1, 2, 3, 4, 5},
			expectedWriterBuf: []StreamRecord{
				StreamRecord([]byte{6, 7}),
				nil,
			},
		},
		// buffer items exceeds, flush is triggered
		{
			writer: &StreamWriter{
				stream: stream,
				buffer: []StreamRecord{
					StreamRecord([]byte{1}),
					StreamRecord([]byte{2}),
				},
				bufferSize:     2,
				bufferIdx:      2,
				maxBufferItems: 2,
				maxBufferSize:  4,
			},
			input:             []byte{3},
			expectedStreamBuf: []byte{1, 2},
			expectedWriterBuf: []StreamRecord{
				StreamRecord([]byte{3}),
				nil,
			},
		},
	}

	for _, test := range tests {
		stream.buf.Reset()

		n, err := test.writer.Write(test.input)
		assert.Equal(t, len(test.input), n)
		assert.NoError(t, err)

		assert.Equal(t, test.expectedStreamBuf, stream.buf.Bytes())
		assert.Equal(t, test.expectedWriterBuf, test.writer.buffer)
	}
}

func TestStreamWriterFlushNoError(t *testing.T) {

	stream := new(StreamMock)
	stream.On("Put", mock.Anything).Return(new(StreamResponseMock), nil)

	// init writer
	writer := &StreamWriter{
		stream: stream,
		buffer: []StreamRecord{StreamRecord([]byte{1, 2})},
	}

	writer.Flush()
	assert.Equal(t, []StreamRecord{nil}, writer.buffer)
	assert.Equal(t, 0, writer.bufferSize)
	assert.Equal(t, []byte{1, 2}, stream.buf.Bytes())

	// write more
	writer.buffer = []StreamRecord{StreamRecord([]byte{3})}

	writer.Flush()
	assert.Equal(t, []StreamRecord{nil}, writer.buffer)
	assert.Equal(t, 0, writer.bufferSize)
	assert.Equal(t, []byte{1, 2, 3}, stream.buf.Bytes())
}

func BenchmarkStreamWriter(b *testing.B) {
	stream := new(StreamMock)
	stream.On("Put", mock.Anything).Return(new(StreamResponseMock), nil)

	w := NewStreamWriter(stream)
	input := []byte{
		1, 2, 3, 4, 5, 6, 7, 8, 9,
	}

	for i := 0; i < b.N; i++ {
		w.Write(input)
	}
}

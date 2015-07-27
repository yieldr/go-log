package logstream

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStreamWriterWriteNoError(t *testing.T) {

	stream := new(StreamMock)
	stream.On("Put", mock.Anything).Return(new(StreamResponseMock), nil)

	tests := []struct {
		writer   *StreamWriter
		writes   int
		input    []byte
		expected []byte
	}{
		// flush is not triggered
		{
			writer: &StreamWriter{
				stream:         stream,
				buf:            newRecordBuffer(500),
				maxBufferItems: 2,
				maxBufferSize:  4,
			},
			writes:   1,
			input:    []byte{1, 2},
			expected: nil,
		},
		// buffer size exceeds, flush is triggered
		{
			writer: &StreamWriter{
				stream:         stream,
				buf:            newRecordBuffer(500),
				maxBufferItems: 2,
				maxBufferSize:  4,
			},
			writes:   2,
			input:    []byte{1, 2, 3, 4, 5},
			expected: []byte{1, 2, 3, 4, 5},
		},
		// buffer items exceeds, flush is triggered
		{
			writer: &StreamWriter{
				stream:         stream,
				buf:            newRecordBuffer(500),
				maxBufferItems: 2,
				maxBufferSize:  4,
			},
			writes:   3,
			input:    []byte{1, 2},
			expected: []byte{1, 2, 1, 2},
		},
	}

	for _, test := range tests {
		stream.buf.Reset()

		for i := 0; i < test.writes; i++ {
			_, err := test.writer.Write(test.input)
			assert.NoError(t, err)
		}

		assert.Equal(t, test.expected, stream.buf.Bytes())
	}
}

func TestStreamWriterFlushNoError(t *testing.T) {

	stream := new(StreamMock)
	stream.On("Put", mock.Anything).Return(new(StreamResponseMock), nil)

	// init writer
	writer := &StreamWriter{
		stream: stream,
		buf:    newRecordBuffer(500),
	}

	writer.Flush()
	assert.Equal(t, 0, writer.buf.getItems())
	assert.Equal(t, 0, writer.buf.getSize())
	assert.Equal(t, []byte(nil), stream.buf.Bytes())

	// write more
	writer.Write([]byte{1, 2, 3})

	writer.Flush()
	assert.Equal(t, 0, writer.buf.getItems())
	assert.Equal(t, 0, writer.buf.getSize())
	assert.Equal(t, []byte{1, 2, 3}, stream.buf.Bytes())
}

func BenchmarkStreamWriter(b *testing.B) {

	stream := new(StreamMock)
	stream.On("Put", mock.Anything).Return(new(StreamResponseMock), nil)

	w := NewStreamWriter(stream)
	p := []byte{
		1, 2, 3, 4, 5,
		1, 2, 3, 4, 5,
		1, 2, 3, 4, 5,
		1, 2, 3, 4, 5,
	}

	for i := 0; i < b.N; i++ {
		_, err := w.Write(p)
		if err != nil {
			b.FailNow()
		}
	}

	b.ReportAllocs()
}

package logstream

import "testing"

func assertByteSliceEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestStreamWriterWriteNoError(t *testing.T) {

	stream := new(StreamMock)

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
			expected: []byte{},
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
			if err != nil {
				t.Error(err)
			}
		}

		if !assertByteSliceEqual(test.expected, stream.buf.Bytes()) {
			t.Error("buffer data not match.")
		}
	}
}

func TestStreamWriterFlushNoError(t *testing.T) {

	stream := new(StreamMock)

	// init writer
	writer := &StreamWriter{
		stream: stream,
		buf:    newRecordBuffer(500),
	}

	writer.Flush()
	if 0 != writer.buf.getItems() {
		t.Error("writer buffer items should be 0.")
	}
	if 0 != writer.buf.getSize() {
		t.Error("writer buffer size should be 0.")
	}
	if []byte(nil) != stream.buf.Bytes() {
		t.Error("writer buffer should be nil.")
	}

	// write more
	writer.Write([]byte{1, 2, 3})

	writer.Flush()
	if 0 != writer.buf.getItems() {
		t.Error("writer buffer items should be 0.")
	}
	if 0 != writer.buf.getSize() {
		t.Error("writer buffer size should be 0.")
	}
	if !assertByteSliceEqual([]byte{1, 2, 3}, stream.buf.Bytes()) {
		t.Error("writer buffer not match,")
	}
}

func BenchmarkStreamWriter(b *testing.B) {

	stream := new(StreamMock)

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

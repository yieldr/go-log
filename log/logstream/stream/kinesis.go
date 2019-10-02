package stream

import (
	"errors"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesis"

	"github.com/yieldr/go-log/log/logstream"
)

// Kinesis implements Stream interface and wraps a kinesis client.
type Kinesis struct {
	name   string
	stream *kinesis.Kinesis
}

// New created a new Kinesis stream with given name and config.
func New(name string, c aws.Config) logstream.Stream {
	return &Kinesis{
		name:   name,
		stream: kinesis.New(&c),
	}
}

// Put records into a remote kinesis stream.
func (k *Kinesis) Put(records []logstream.StreamRecord) (logstream.StreamResponse, error) {
	if len(records) == 0 {
		return nil, errors.New("empty records for kinesis.")
	}

	entries := make([]*kinesis.PutRecordsRequestEntry, len(records))
	for i, record := range records {
		entries[i] = &kinesis.PutRecordsRequestEntry{
			Data:         []byte(record),
			PartitionKey: aws.String(k.getPartitionKey(i)),
		}
	}

	params := &kinesis.PutRecordsInput{
		Records:    entries,
		StreamName: aws.String(k.name),
	}

	return k.stream.PutRecords(params)
}

// Close.
func (k *Kinesis) Close() error {
	return nil
}

// getPartitionKey generates a random string, based on i(optional).
func (k *Kinesis) getPartitionKey(i int) string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

package logstream

import (
	"errors"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

// Kinesis implements Stream interface and wraps a kinesis client.
// TODO: need aws.Config.
type Kinesis struct {
	streamName string
	stream     kinesis.Kinesis
}

// Put records into a remote kinesis stream.
func (k *Kinesis) Put(records []StreamRecord) (StreamResponse, error) {
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
		StreamName: aws.String(k.streamName),
	}

	return k.stream.PutRecords(params)
}

// Close.
// TODO: do we close the connection to kinesis?
func (k *Kinesis) Close() error {
	return nil
}

// getPartitionKey generates a random string, based on i(optional).
func (k *Kinesis) getPartitionKey(i int) string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

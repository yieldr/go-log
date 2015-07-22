package stream

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/yieldr/go-log/log/logstream"
)

// Kinesis implements Stream interface and wraps a kinesis client.
type Kinesis struct {
	name   string
	client *kinesis.Kinesis
}

// New a kinesis stream with given name and config.
func New(name string, config *aws.Config) *Kinesis {
	return &Kinesis{
		name:   name,
		client: kinesis.New(config),
	}
}

// Put records into a remote kinesis client.
func (k *Kinesis) Put(records []logstream.StreamRecord) (logstream.StreamResponse, error) {

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

	return k.client.PutRecords(params)
}

// Close.
// TODO: do we close the connection to kinesis?
func (k *Kinesis) Close() error {
	return nil
}

// getPartitionKey generates a random string, based on i(optional).
func (k *Kinesis) getPartitionKey(i int) string {
	return strconv.FormatInt(time.Now().Unix()+(int64)(i), 10)
}

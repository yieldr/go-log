package functional

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/yieldr/go-log/log/logstream"
	"github.com/yieldr/go-log/log/logstream/stream"
)

func TestKinesisStream(t *testing.T) {

	fmt.Println("Test Kinesis Stream.")

	// your own credentials
	streamName := ""
	region := ""
	accessKeyId := ""
	secretAccessKey := ""

	config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKeyId, secretAccessKey, ""),
		Region:      region,
	}

	demo := stream.New(streamName, config)

	for i := 0; i < 10; i++ {

		records := []logstream.StreamRecord{
			logstream.StreamRecord([]byte(time.Now().String() + "-1")),
			logstream.StreamRecord([]byte(time.Now().String() + "-2")),
			logstream.StreamRecord([]byte(time.Now().String() + "-3")),
		}

		resp, err := demo.Put(records)

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(resp)
	}
}

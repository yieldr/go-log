package stream

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/yieldr/go-log/log/logstream"
)

func TestPut(t *testing.T) {

	region := "eu-west-1"
	accessKeyId := "AKIAI6E2ROMVG4TX3LRQ"
	secretAccessKey := "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"

	config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKeyId, secretAccessKey, ""),
		Region:      region,
	}

	demo := New("demo", config)

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

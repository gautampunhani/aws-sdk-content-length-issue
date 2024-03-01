package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/logging"
)

func main() {
	err := DownloadCarrierDocument("tcga-2-open", "00db5a87-9d10-48cb-b1e4-928065681023/TCGA-31-1951-01A-01D-0651-02_BioSizing.tsv", "us-east-1")
	if err != nil {
		fmt.Printf("Error %s", err.Error())
		return
	}
}

func DownloadCarrierDocument(bucket string, objectKey string, region string) error {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
	}

	todo := context.TODO()
	client, err := s3Client(todo, region)
	if err != nil {
		return err
	}

	result, err := client.GetObject(todo, input)
	if err != nil {
		return err
	}

	fmt.Printf("\n\nContent Length %d", result.ContentLength)
	return nil
}

func s3Client(ctx context.Context, region string) (*s3.Client, error) {
	var cfg aws.Config
	var err error

	cfg, err = config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(aws.AnonymousCredentials{}),
		config.WithClientLogMode(aws.LogRequest|aws.LogResponse|aws.LogResponseEventMessage),
		config.WithLogger(&loggerForAWS{}),
	)

	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(cfg), nil
}

type loggerForAWS struct {
}

func (log *loggerForAWS) Logf(classification logging.Classification, format string, v ...interface{}) {
	fmt.Println("Logger Message : ")
	fmt.Printf(format, v)
}

package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type RemoteStorage struct {
	client *s3.Client
	bucket string
}

func NewRemoteStorage() (*RemoteStorage, error) {
	endpoint := os.Getenv("AGREPL_R2_ENDPOINT")
	accessKey := os.Getenv("AGREPL_R2_ACCESS_KEY_ID")
	secretKey := os.Getenv("AGREPL_R2_SECRET_ACCESS_KEY")
	bucket := os.Getenv("AGREPL_R2_BUCKET")

	if endpoint == "" || accessKey == "" || secretKey == "" || bucket == "" {
		return nil, fmt.Errorf("remote storage environment variables not set (AGREPL_R2_ENDPOINT, AGREPL_R2_ACCESS_KEY_ID, AGREPL_R2_SECRET_ACCESS_KEY, AGREPL_R2_BUCKET)")
	}

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               endpoint,
			HostnameImmutable: true,
			SigningRegion:     "auto",
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load SDK config: %w", err)
	}

	return &RemoteStorage{
		client: s3.NewFromConfig(cfg),
		bucket: bucket,
	}, nil
}

func (rs *RemoteStorage) Push(ctx context.Context, runID string, data []byte) error {
	key := fmt.Sprintf("runs/%s.json", runID)
	_, err := rs.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(rs.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return fmt.Errorf("failed to upload run to R2: %w", err)
	}
	return nil
}

func (rs *RemoteStorage) Pull(ctx context.Context, runID string) ([]byte, error) {
	key := fmt.Sprintf("runs/%s.json", runID)
	output, err := rs.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(rs.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download run from R2: %w", err)
	}
	defer output.Body.Close()

	data, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read run data: %w", err)
	}
	return data, nil
}

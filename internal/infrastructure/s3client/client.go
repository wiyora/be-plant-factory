package s3client

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rizalarfiyan/be-plant-factory/internal/config"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
)

type S3Client struct {
	client *s3.Client
	log    zerolog.Logger
	bucket string
}

func New(i do.Injector) (*S3Client, error) {
	rawLog := do.MustInvoke[*zerolog.Logger](i)
	cfg := do.MustInvoke[*config.Config](i)

	log := logger.WithLayer(rawLog, logger.LayerS3)
	log.Info().Msg("Connecting to S3")

	sdkConfig, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithRegion("us-east-1"),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.S3.AccessKey,
			cfg.S3.SecretKey,
			"",
		)),
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to load s3 config")
		return nil, fmt.Errorf("load s3 config: %w", err)
	}

	client := s3.NewFromConfig(sdkConfig, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.S3.Endpoint)
		o.UsePathStyle = true
	})

	s3Client := &S3Client{
		client: client,
		log:    log,
		bucket: cfg.S3.BucketName,
	}

	if err := s3Client.ensureBucket(context.Background()); err != nil {
		log.Error().Err(err).Msg("failed to ensure s3 bucket")
		return nil, err
	}

	log.Info().Msg("successfully connected to s3")
	return s3Client, nil
}

func (s *S3Client) ensureBucket(ctx context.Context) error {
	_, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(s.bucket),
	})
	if err == nil {
		return nil
	}

	s.log.Info().Str("bucket", s.bucket).Msg("bucket not found, creating")
	_, err = s.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(s.bucket),
	})

	if err != nil {
		s.log.Error().Err(err).Str("bucket", s.bucket).Msg("failed to create bucket")
		return fmt.Errorf("create bucket: %w", err)
	}

	s.log.Info().Str("bucket", s.bucket).Msg("bucket created successfully")
	return nil
}

func (s *S3Client) Client() *s3.Client {
	return s.client
}

func (s *S3Client) Bucket() string {
	return s.bucket
}

func (s *S3Client) HealthCheckWithContext(ctx context.Context) error {
	_, err := s.client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		s.log.Error().Err(err).Msg("s3 health check failed")
		return fmt.Errorf("s3 health check failed: %w", err)
	}

	return nil
}

func (s *S3Client) Shutdown(ctx context.Context) error {
	return nil
}

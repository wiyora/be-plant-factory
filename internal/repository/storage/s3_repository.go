package storage

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/s3client"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/helper"
	"github.com/samber/do/v2"
)

type S3Repository interface {
	GeneratePresignedUpload(ctx context.Context, req entity.StorageGeneratePresignedUpload) (StoragePresignedUpload, error)
	CleanupTemporary(ctx context.Context) error
	BatchDelete(ctx context.Context, objects []types.ObjectIdentifier) error
}

type s3Repository struct {
	client *s3.Client
	bucket string
}

func NewS3Repository(i do.Injector) (S3Repository, error) {
	client := do.MustInvoke[*s3client.S3Client](i)

	return s3Repository{
		client: client.Client(),
		bucket: client.Bucket(),
	}, nil
}

func (r s3Repository) GeneratePresignedUpload(ctx context.Context, req entity.StorageGeneratePresignedUpload) (StoragePresignedUpload, error) {
	log := logger.WithLayerCtx(ctx, logger.LayerS3Repository)

	config, _ := req.Type.Config()
	presigner := s3.NewPresignClient(r.client)

	output, err := presigner.PresignPostObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucket),
		Key:         aws.String(req.Type.NewTempFile(req.Extension)),
		ContentType: aws.String(req.MimeType),
	}, func(o *s3.PresignPostOptions) {
		o.Expires = config.MaxPresigned
		o.Conditions = []interface{}{
			[]interface{}{"content-length-range", config.MinSize, config.MaxSize},
			map[string]string{"content-type": req.MimeType},
		}
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to generate presigned upload URL")
		return StoragePresignedUpload{}, err
	}

	return StoragePresignedUpload{
		URL:    output.URL,
		Fields: output.Values,
	}, nil
}

func (r s3Repository) CleanupTemporary(ctx context.Context) error {
	log := logger.WithLayerCtx(ctx, logger.LayerS3Repository)
	log.Info().Msg("cleanup temporary files")

	for _, st := range entity.StorageTypes {
		cfg, _ := st.Config()
		prefix := st.TempPath()
		log.Info().Str("path", prefix).Msg("scanning expired files")

		cutoffTime := time.Now().Add(-cfg.MaxAge)
		expiredObjects, err := r.collectExpiredObjects(ctx, prefix, cutoffTime)
		if err != nil {
			log.Error().Err(err).Str("path", prefix).Msg("failed to scan expired files")
			return err
		}

		log.Info().Str("path", prefix).Int("expired_count", len(expiredObjects)).Msg("found expired files")

		if err := r.BatchDelete(ctx, expiredObjects); err != nil {
			log.Error().Err(err).Str("path", prefix).Msg("failed to clean up path")
			return err
		}

		log.Info().Str("path", prefix).Int("deleted_count", len(expiredObjects)).Msg("successfully cleaned up expired files")
	}

	return nil
}

func (r s3Repository) BatchDelete(ctx context.Context, objects []types.ObjectIdentifier) error {
	if len(objects) == 0 {
		return nil
	}

	log := logger.WithLayerCtx(ctx, logger.LayerS3Repository)
	batches := helper.ChunkSlice(objects, 1000)

	for _, batch := range batches {
		input := &s3.DeleteObjectsInput{
			Bucket: aws.String(r.bucket),
			Delete: &types.Delete{
				Objects: batch,
				Quiet:   aws.Bool(true),
			},
		}

		if _, err := r.client.DeleteObjects(ctx, input); err != nil {
			log.Error().Err(err).Int("batch_size", len(batch)).Msg("failed to delete objects")
			return err
		}
	}

	return nil
}

func (r s3Repository) collectExpiredObjects(ctx context.Context, prefix string, cutoffTime time.Time) ([]types.ObjectIdentifier, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(r.bucket),
		Prefix: aws.String(prefix),
	}

	paginator := s3.NewListObjectsV2Paginator(r.client, input)
	var expiredObjects []types.ObjectIdentifier

	log := logger.WithLayerCtx(ctx, logger.LayerS3Repository)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			log.Error().Err(err).Str("prefix", prefix).Msg("failed to list objects")
			return nil, err
		}

		for _, object := range page.Contents {
			if *object.Key == prefix || *object.Key == prefix+"/" {
				continue
			}

			if object.LastModified.Before(cutoffTime) {
				expiredObjects = append(expiredObjects, types.ObjectIdentifier{
					Key: object.Key,
				})
			}
		}
	}

	return expiredObjects, nil
}

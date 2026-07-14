package storage

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	domainError "github.com/rizalarfiyan/be-plant-factory/internal/domain/error"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/s3client"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/helper"
	"github.com/samber/do/v2"
)

type S3Repository interface {
	GeneratePresignedUpload(ctx context.Context, req entity.StorageGeneratePresignedUpload) (StoragePresignedUpload, error)
	CleanupTemporary(ctx context.Context) error
	BatchDelete(ctx context.Context, objects []types.ObjectIdentifier) error
	Move(ctx context.Context, srcKey, dstKey string) error
	Diff(ctx context.Context, diff entity.StorageDiff) error
	Diffs(ctx context.Context, diff entity.StorageDiffs) error
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

func (r s3Repository) Move(ctx context.Context, srcKey, dstKey string) error {
	log := logger.WithLayerCtx(ctx, logger.LayerS3Repository)

	copyInput := &s3.CopyObjectInput{
		Bucket:     aws.String(r.bucket),
		CopySource: aws.String(r.bucket + "/" + srcKey),
		Key:        aws.String(dstKey),
	}

	if _, err := r.client.CopyObject(ctx, copyInput); err != nil {
		log.Error().Err(err).Str("src", srcKey).Str("dst", dstKey).Msg("failed to copy object")
		return err
	}

	deleteInput := &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(srcKey),
	}

	if _, err := r.client.DeleteObject(ctx, deleteInput); err != nil {
		log.Error().Err(err).Str("src", srcKey).Msg("failed to delete source object after copy")
		return err
	}

	return nil
}

func (r s3Repository) Diff(ctx context.Context, diff entity.StorageDiff) error {
	log := logger.WithLayerCtx(ctx, logger.LayerS3Repository)

	if !diff.IsValid {
		return nil
	}

	if !helper.IsEmptyStruct(diff.Added) {
		if err := r.Move(ctx, diff.Added.From, diff.Added.To); err != nil {
			log.Error().Err(err).Interface("added", diff.Added).Msg("failed to process diff move operation")
			return err
		}
	}

	if diff.Removed != "" {
		obj := types.ObjectIdentifier{Key: aws.String(diff.Removed)}
		if err := r.BatchDelete(ctx, []types.ObjectIdentifier{obj}); err != nil {
			log.Error().Err(err).Str("removed", diff.Removed).Msg("failed to process diff remove operation")
			return err
		}
	}

	return nil
}

func (r s3Repository) Diffs(ctx context.Context, diff entity.StorageDiffs) error {
	log := logger.WithLayerCtx(ctx, logger.LayerS3Repository)

	if !diff.IsValid {
		return nil
	}

	var errCount int

	for _, moveOp := range diff.Added {
		if helper.IsEmptyStruct(diff.Added) {
			continue
		}

		if err := r.Move(ctx, moveOp.From, moveOp.To); err != nil {
			log.Error().Err(err).Interface("move_op", moveOp).Msg("failed to move object in bulk diff")
			errCount++
		}
	}

	if len(diff.Removed) > 0 {
		var objects []types.ObjectIdentifier
		for _, key := range diff.Removed {
			if key != "" {
				objects = append(objects, types.ObjectIdentifier{Key: aws.String(key)})
			}
		}

		if len(objects) > 0 {
			if err := r.BatchDelete(ctx, objects); err != nil {
				log.Error().Err(err).Int("count", len(objects)).Msg("failed to delete objects in bulk diff")
				errCount++
			}
		}
	}

	if errCount > 0 {
		log.Error().Int("error_count", errCount).Msg("bulk storage diff completed with errors")
		return domainError.ErrStorageDiffInvalid
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

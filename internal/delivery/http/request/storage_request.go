package request

import (
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	domainError "github.com/rizalarfiyan/be-plant-factory/internal/domain/error"
)

type StoragePresignedUploadRequest struct {
	MimeType  string             `json:"mime-type" validate:"required" example:"image/png"`
	FileSize  int64              `json:"file-size" validate:"required,gt=0" example:"1048576"`
	Extension string             `json:"extensions" validate:"required" example:".png"`
	Type      entity.StorageType `json:"type" validate:"required,enum" example:"avatar"`
}

func (r StoragePresignedUploadRequest) Validate() error {
	config, _ := r.Type.Config()

	if !config.MimeType.Contains(r.MimeType) {
		return domainError.NewManualValidation("mime-type", "INVALID")
	}

	if !config.Extension.Contains(r.Extension) {
		return domainError.NewManualValidation("extensions", "INVALID")
	}

	if r.FileSize > config.MaxSize {
		return domainError.NewManualValidation("file-size", "MAX_FILE_SIZE_EXCEEDED")
	}

	if r.FileSize < config.MinSize {
		return domainError.NewManualValidation("file-size", "MIN_FILE_SIZE_EXCEEDED")
	}

	return nil
}

func (r StoragePresignedUploadRequest) ToEntity() entity.StorageGeneratePresignedUpload {
	return entity.StorageGeneratePresignedUpload{
		MimeType:  r.MimeType,
		FileSize:  r.FileSize,
		Extension: r.Extension,
		Type:      r.Type,
	}
}

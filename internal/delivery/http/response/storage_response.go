package response

import "github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"

type StoragePresignedUploadResponse struct {
	URL    string            `json:"url"`
	Fields map[string]string `json:"fields"`
}

func NewStoragePresignedUploadResponse(e entity.StoragePresignedUpload) StoragePresignedUploadResponse {
	return StoragePresignedUploadResponse{
		URL:    e.URL,
		Fields: e.Fields,
	}
}

package storage

import "github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"

type StoragePresignedUpload struct {
	URL    string
	Fields map[string]string
}

func (p StoragePresignedUpload) ToEntity() entity.StoragePresignedUpload {
	return entity.StoragePresignedUpload{
		URL:    p.URL,
		Fields: p.Fields,
	}
}

package entity

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/google/uuid"
)

type StorageType string

const (
	StorageTypeAvatar     StorageType = "avatar"
	StorageTypeTenantLogo StorageType = "tenant-logo"
)

type StorageConfig struct {
	Path         string
	MimeType     mapset.Set[string]
	Extension    mapset.Set[string]
	MinSize      int64
	MaxSize      int64
	MaxAge       time.Duration
	MaxPresigned time.Duration
}

type StorageDiff struct {
	Added   []string
	Removed []string
	Current []string
}

var (
	storageConfig = map[StorageType]StorageConfig{
		StorageTypeAvatar: {
			Path:         "user/avatar",
			MimeType:     mapset.NewSet("image/jpeg", "image/png"),
			Extension:    mapset.NewSet(".jpg", ".jpeg", ".png"),
			MinSize:      10 * 1024,       // 10 KB
			MaxSize:      750 * 1024,      // 750 KB
			MaxAge:       24 * time.Hour,  // 1 day
			MaxPresigned: 7 * time.Minute, // 7 minutes
		},
		StorageTypeTenantLogo: {
			Path:         "tenant/logo",
			MimeType:     mapset.NewSet("image/jpeg", "image/png"),
			Extension:    mapset.NewSet(".jpg", ".jpeg", ".png"),
			MinSize:      10 * 1024,          // 10 KB
			MaxSize:      750 * 1024,         // 750 KB
			MaxAge:       3 * 24 * time.Hour, // 3 days
			MaxPresigned: 7 * time.Minute,    // 7 minutes
		},
	}

	StorageTypes = []StorageType{
		StorageTypeAvatar,
		StorageTypeTenantLogo,
	}
)

func (s StorageType) Config() (StorageConfig, bool) {
	cfg, ok := storageConfig[s]
	return cfg, ok
}

func (s StorageType) Valid() bool {
	_, ok := storageConfig[s]
	return ok
}

func (s StorageType) Path() string {
	if cfg, ok := s.Config(); ok {
		return cfg.Path
	}
	return ""
}

func (s StorageType) TempPath() string {
	path := s.Path()
	if path == "" {
		return ""
	}
	return fmt.Sprintf("temp/%s", path)
}

func (s StorageType) IsValidFile(files ...string) bool {
	cfg, ok := s.Config()
	if !ok || len(files) == 0 {
		return false
	}

	basePath := s.Path()
	tempPath := s.TempPath()

	for _, file := range files {
		dir := filepath.Dir(file)
		if dir != basePath && dir != tempPath {
			return false
		}

		ext := filepath.Ext(file)
		if !cfg.Extension.Contains(ext) {
			return false
		}

		baseName := filepath.Base(file)
		uuidPart := strings.TrimSuffix(baseName, ext)

		if _, err := uuid.Parse(uuidPart); err != nil {
			return false
		}
	}

	return true
}

func (s StorageType) Diff(oldFiles []string, newFiles []string) StorageDiff {
	oldSet := mapset.NewSet(oldFiles...)
	newSet := mapset.NewSet(newFiles...)

	// TODO: check added and move from temp to original path
	return StorageDiff{
		Added:   newSet.Difference(oldSet).ToSlice(),
		Removed: oldSet.Difference(newSet).ToSlice(),
		Current: oldSet.Intersect(newSet).ToSlice(),
	}
}

func (s StorageType) NewFile(extension string) string {
	if cfg, ok := s.Config(); ok {
		return fmt.Sprintf("%s/%s%s", cfg.Path, uuid.Must(uuid.NewV7()).String(), extension)
	}

	return ""
}

func (s StorageType) NewTempFile(extension string) string {
	newFile := s.NewFile(extension)
	if newFile != "" {
		return fmt.Sprintf("temp/%s", newFile)
	}

	return ""
}

type StorageGeneratePresignedUpload struct {
	MimeType  string
	FileSize  int64
	Extension string
	Type      StorageType
}

type StoragePresignedUpload struct {
	URL    string
	Fields map[string]string
}

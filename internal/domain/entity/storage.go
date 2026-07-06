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

type StorageMove struct {
	From string
	To   string
}

type StorageDiff struct {
	Added   StorageMove
	Removed string
	Result  string
	IsValid bool
}

type StorageDiffs struct {
	Added   []StorageMove
	Removed []string
	Result  []string
	IsValid bool
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

	tempPrefix = "temp/"
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
	return tempPrefix + path
}

func (s StorageType) IsValidFile(file string) bool {
	if file == "" {
		return false
	}

	cfg, ok := s.Config()
	if !ok {
		return false
	}

	ext := filepath.Ext(file)
	if !cfg.Extension.Contains(ext) {
		return false
	}

	baseName := filepath.Base(file)
	uuidPart := baseName[:len(baseName)-len(ext)]
	if _, err := uuid.Parse(uuidPart); err != nil {
		return false
	}

	dir := filepath.Dir(file)
	return dir == cfg.Path || dir == (tempPrefix+cfg.Path)
}

func (s StorageType) IsValidFiles(files ...string) bool {
	if len(files) == 0 {
		return false
	}

	cfg, ok := s.Config()
	if !ok {
		return false
	}

	basePath := cfg.Path
	tempPath := tempPrefix + basePath

	for _, file := range files {
		if file == "" {
			return false
		}

		ext := filepath.Ext(file)
		if !cfg.Extension.Contains(ext) {
			return false
		}

		baseName := filepath.Base(file)
		uuidPart := baseName[:len(baseName)-len(ext)]
		if _, err := uuid.Parse(uuidPart); err != nil {
			return false
		}

		dir := filepath.Dir(file)
		if dir != basePath && dir != tempPath {
			return false
		}
	}

	return true
}

func (s StorageType) Diff(oldFile string, newFile string) StorageDiff {
	if oldFile != "" && !s.IsValidFile(oldFile) {
		return StorageDiff{
			IsValid: false,
		}
	}

	if newFile != "" && !s.IsValidFile(newFile) {
		return StorageDiff{
			IsValid: false,
		}
	}

	if newFile == "" || newFile == oldFile {
		return StorageDiff{
			Result:  oldFile,
			IsValid: true,
		}
	}

	basePathPrefix := s.Path()
	if strings.HasPrefix(newFile, basePathPrefix) {
		return StorageDiff{
			Result:  oldFile,
			IsValid: false,
		}
	}

	tempPathPrefix := tempPrefix + basePathPrefix
	if strings.HasPrefix(newFile, tempPathPrefix) {
		finalDestination := newFile[len(tempPrefix):]
		return StorageDiff{
			Added: StorageMove{
				From: newFile,
				To:   finalDestination,
			},
			Removed: oldFile,
			Result:  finalDestination,
			IsValid: true,
		}
	}

	return StorageDiff{
		IsValid: false,
	}
}

func (s StorageType) Diffs(oldFiles []string, newFiles []string) StorageDiffs {
	diffs := StorageDiffs{
		Added:   []StorageMove{},
		Removed: []string{},
		Result:  []string{},
		IsValid: true,
	}

	oldFileSet := mapset.NewSet[string]()
	validOldCount := 0
	for _, file := range oldFiles {
		if file != "" {
			if !s.IsValidFile(file) {
				return StorageDiffs{
					IsValid: false,
				}
			}
			oldFileSet.Add(file)
			validOldCount++
		}
	}

	newFileSet := mapset.NewSet[string]()
	validNewCount := 0
	for _, file := range newFiles {
		if file != "" {
			if !s.IsValidFile(file) {
				return StorageDiffs{
					IsValid: false,
				}
			}
			newFileSet.Add(file)
			validNewCount++
		}
	}

	if validNewCount == 0 {
		if validOldCount > 0 {
			diffs.Result = make([]string, 0, validOldCount)
			oldFileSet.Each(func(file string) bool {
				diffs.Result = append(diffs.Result, file)
				return false
			})
		}
		return diffs
	}

	basePathPrefix := s.Path()
	tempPathPrefix := tempPrefix + basePathPrefix

	diffs.Added = make([]StorageMove, 0, validNewCount)
	diffs.Result = make([]string, 0, validNewCount)

	var iterErr bool
	newFileSet.Each(func(newFile string) bool {
		if strings.HasPrefix(newFile, basePathPrefix) {
			if oldFileSet.Contains(newFile) {
				diffs.Result = append(diffs.Result, newFile)
				oldFileSet.Remove(newFile)
				return false
			}

			iterErr = true
			return true
		}

		if strings.HasPrefix(newFile, tempPathPrefix) {
			finalDestination := newFile[len(tempPrefix):]
			diffs.Added = append(diffs.Added, StorageMove{
				From: newFile,
				To:   finalDestination,
			})

			diffs.Result = append(diffs.Result, finalDestination)
			return false
		}

		iterErr = true
		return true
	})

	if iterErr {
		return StorageDiffs{
			IsValid: false,
		}
	}

	if oldFileSet.Cardinality() > 0 {
		diffs.Removed = make([]string, 0, oldFileSet.Cardinality())
		oldFileSet.Each(func(oldFile string) bool {
			diffs.Removed = append(diffs.Removed, oldFile)
			return false
		})
	}

	return diffs
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
		return tempPrefix + newFile
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

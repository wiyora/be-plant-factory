package logger

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/rizalarfiyan/be-plant-factory/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func New(env config.AppEnv) (*zerolog.Logger, error) {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339

	output := zerolog.ConsoleWriter{
		Out:          os.Stderr,
		TimeFormat:   time.RFC3339,
		FormatCaller: formatCaller(),
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("%v", i)
		},
	}

	if !env.IsServerEnv() {
		output.Out = os.Stdout
	}

	logger := zerolog.New(output).With().Timestamp().Caller().Logger()
	log.Logger = logger

	return &logger, nil
}

func formatCaller() zerolog.Formatter {
	baseDir, err := filepath.Abs("./")
	isValidBaseDir := err == nil

	return func(i interface{}) string {
		if !isValidBaseDir {
			return fmt.Sprintf("%v", i)
		}

		fullPath, ok := i.(string)
		if !ok {
			return fmt.Sprintf("%v", i)
		}

		if !strings.HasPrefix(fullPath, baseDir) || !isValidBaseDir {
			return path.Base(fullPath)
		}

		currentFile, err := filepath.Rel(baseDir, fullPath)
		if err != nil {
			return path.Base(fullPath)
		}

		return currentFile
	}
}

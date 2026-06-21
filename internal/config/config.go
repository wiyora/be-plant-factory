package config

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/creasty/defaults"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	dotEnv "github.com/knadh/koanf/parsers/dotenv"
	envProvider "github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
)

type Config struct {
	App      AppConfig      `env:"APP"`
	HTTP     HTTPConfig     `env:"HTTP"`
	Fiber    FiberConfig    `env:"FIBER"`
	CORS     CORSConfig     `env:"CORS"`
	Log      LogConfig      `env:"LOG"`
	Database DatabaseConfig `env:"DB"`
	Redis    RedisConfig    `env:"REDIS"`
	JWT      JWTConfig      `env:"JWT"`
	Auth     AuthConfig     `env:"AUTH"`
	Swagger  SwaggerConfig  `env:"SWAGGER"`
	Frontend FrontendConfig `env:"FRONTEND"`
}

type Loader interface {
	Load() (*Config, error)
}

type loader struct {
	k        *koanf.Koanf
	filePath string
	paths    map[string]string
	log      zerolog.Logger
}

func NewConfig(filePath string, log zerolog.Logger) Loader {
	return &loader{
		k:        koanf.New("."),
		filePath: filePath,
		paths:    make(map[string]string),
		log:      log,
	}
}

func (l *loader) Load() (*Config, error) {
	l.log.Info().Msg("Loading configuration")

	var cfg Config
	if err := defaults.Set(&cfg); err != nil {
		l.log.Error().Err(err).Msg("failed to set default values for config")
		return nil, err
	}

	l.getEnvPaths(reflect.TypeOf(cfg), "", "", l.paths)

	env := MustLoadAppEnv()
	if env.IsServerEnv() {
		if err := l.loadFromEnv(); err != nil {
			l.log.Error().Err(err).Msg("failed to load config from environment variables")
			return nil, err
		}
	} else {
		if err := l.loadFromFile(); err != nil {
			l.log.Error().Err(err).Msg("failed to load config from file")
			return nil, err
		}
	}

	if err := l.k.UnmarshalWithConf("", &cfg, koanf.UnmarshalConf{
		Tag: "env",
	}); err != nil {
		l.log.Error().Err(err).Msg("failed to unmarshal config")
		return nil, err
	}

	if err := l.validate(&cfg); err != nil {
		l.log.Error().Err(err).Msg("config validation failed")
		return nil, err
	}

	// Initialize any additional fields
	cfg.Frontend.Init()

	l.log.Info().Msg("Configuration loaded successfully")
	return &cfg, nil
}

func (l *loader) getEnvPaths(t reflect.Type, envPath, structPath string, paths map[string]string) {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("env")
		if tag == "" || tag == "-" {
			continue
		}

		newEnvPath := tag
		if envPath != "" {
			newEnvPath = envPath + "_" + tag
		}

		newStructPath := tag
		if structPath != "" {
			newStructPath = structPath + "." + tag
		}

		paths[newEnvPath] = newStructPath

		if field.Type.Kind() == reflect.Struct {
			l.getEnvPaths(field.Type, newEnvPath, newStructPath, paths)
		} else if field.Type.Kind() == reflect.Pointer && field.Type.Elem().Kind() == reflect.Struct {
			l.getEnvPaths(field.Type.Elem(), newEnvPath, newStructPath, paths)
		}
	}
}

func (l *loader) validate(cfg *Config) error {
	v := validator.New()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("env"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)
	trans, found := uni.GetTranslator("en")
	if !found {
		return fmt.Errorf("translator not found for locale 'en'")
	}

	err := enTranslations.RegisterDefaultTranslations(v, trans)
	if err != nil {
		return fmt.Errorf("failed to register translations: %v", err)
	}

	if err := v.Struct(cfg); err != nil {
		var errs validator.ValidationErrors
		if ok := errors.As(err, &errs); !ok {
			return err
		}

		for _, e := range errs {
			keys := strings.Split(e.Namespace(), ".")
			if len(keys) > 1 {
				keys = keys[1:]
			}

			key := strings.Join(keys, "_")
			errMsg := e.Translate(trans)
			errMsg = strings.Replace(errMsg, e.Field(), key, 1)
			l.log.Error().Msg(errMsg)
		}

		return fmt.Errorf("validation failed with %d error(s)", len(errs))
	}

	return nil
}

func (l *loader) loadFromEnv() error {
	return l.k.Load(envProvider.Provider(".", envProvider.Opt{
		TransformFunc: func(k, v string) (string, any) {
			if p, ok := l.paths[k]; ok {
				return p, v
			}

			return "", nil
		},
	}), nil)
}

func (l *loader) loadFromFile() error {
	envConfig := dotEnv.ParserEnv("", ".", func(s string) string {
		if p, ok := l.paths[s]; ok {
			return p
		}

		return strings.ReplaceAll(s, "_", ".")
	})

	return l.k.Load(file.Provider(l.filePath), envConfig)
}

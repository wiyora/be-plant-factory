package middleware

import (
	"fmt"
	"strings"

	swaggerUI "github.com/gofiber/contrib/v3/swaggerui"
	fiberZerolog "github.com/gofiber/contrib/v3/zerolog"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/basicauth"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/helper"
	"github.com/rs/zerolog"
)

func (m *middleware) Register(app *fiber.App) {
	app.Use(requestid.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     m.conf.CORS.AllowOrigins,
		AllowMethods:     m.conf.CORS.AllowMethods,
		AllowHeaders:     m.conf.CORS.AllowHeaders,
		AllowCredentials: m.conf.CORS.AllowCredentials,
		ExposeHeaders:    m.conf.CORS.ExposeHeaders,
	}))

	app.Use(compress.New(compress.Config{
		Level: compress.Level(m.conf.Fiber.CompressLevel),
	}))

	app.Use(helmet.New())

	app.Use(recover.New(recover.Config{
		EnableStackTrace: !m.conf.App.Env.IsServerEnv(),
	}))

	app.Use(fiberZerolog.New(fiberZerolog.Config{
		FieldsSnakeCase: true,
		Fields: []string{
			fiberZerolog.FieldIP,
			fiberZerolog.FieldMethod,
			fiberZerolog.FieldPath,
			fiberZerolog.FieldURL,
			fiberZerolog.FieldMethod,
			fiberZerolog.FieldLatency,
			fiberZerolog.FieldStatus,
			fiberZerolog.FieldBody,
			fiberZerolog.FieldError,
			fiberZerolog.FieldRequestID,
		},
		GetLogger: func(c fiber.Ctx) zerolog.Logger {
			logCtx := m.log.With()

			logCtx = logCtx.Str("layer", string(logger.LayerHttp))

			if user, ok := helper.GetLocalUser(c); ok {
				logCtx = logCtx.Str("user_id", user.ID.String())
			}

			if sessionID, ok := helper.GetLocalSessionID(c); ok {
				logCtx = logCtx.Str("session_id", sessionID.String())
			}

			return logCtx.Logger()
		},
		Next: func(c fiber.Ctx) bool {
			path := c.Path()
			return path == "/health" || strings.HasPrefix(path, "/swagger")
		},
	}))

	app.Use(func(c fiber.Ctx) error {
		reqID := c.GetRespHeader(fiber.HeaderXRequestID)
		reqLogger := m.log.With().Str("request_id", reqID).Logger()
		ctx := reqLogger.WithContext(c.Context())
		c.SetContext(ctx)
		return c.Next()
	})

	if m.conf.Swagger.Username != "" || m.conf.Swagger.Password != "" {
		app.Get("/swagger/*", basicauth.New(basicauth.Config{
			Users: map[string]string{
				m.conf.Swagger.Username: helper.BasicAuthPassword(m.conf.Swagger.Password),
			},
		}))
	}

	app.Use(swaggerUI.New(swaggerUI.Config{
		BasePath: "/",
		FilePath: "./swagger/swagger.json",
		Path:     "swagger",
		Title:    fmt.Sprintf("%s API Docs", m.conf.App.Name),
	}))
}

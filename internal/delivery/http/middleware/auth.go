package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-plant-factory/internal/delivery/http/response"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/constant"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/helper"
)

func (m *middleware) Auth() fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := c.Context()
		log := logger.WithLayerCtx(ctx, logger.LayerMiddleware)

		authToken := c.Cookies(constant.CookieAccessToken)
		if authToken == "" {
			return response.New(c, code.Unauthorized)
		}

		claims, err := m.jwt.ParseAccessToken(authToken)
		if err != nil || helper.IsEmptyStruct(claims) {
			log.Error().Err(err).Msg("failed to parse access token")
			return response.New(c, code.Unauthorized)
		}

		accessTokenID, err := uuid.Parse(claims.ID)
		if err != nil {
			log.Error().Err(err).Msg("invalid token claims ID")
			return response.New(c, code.Unauthorized)
		}

		isBlacklisted := m.token.IsBlacklisted(ctx, accessTokenID)
		if isBlacklisted {
			log.Error().Msg("token is blacklisted")
			return response.New(c, code.Unauthorized)
		}

		user, err := m.cache.GetById(ctx, claims.UserID)
		if err != nil {
			log.Error().Err(err).Msg("failed to get user from cache")
			return response.New(c, code.Unauthorized)
		}

		if helper.IsEmptyStruct(user) {
			log.Error().Msg("user not found in cache")
			return response.New(c, code.Unauthorized)
		}

		// store into locals
		c.Locals(constant.KeyAuthUser, entity.AuthContext{
			ID:          user.ID,
			Name:        user.Name,
			Email:       user.Email,
			Avatar:      user.Avatar,
			CurrentStep: user.CurrentStep,
		})
		c.Locals(constant.KeyAuthSessionID, claims.SessionID)
		c.Locals(constant.KeyAuthAccessToken, entity.AccessTokenContext{
			ID:        accessTokenID,
			ExpiredAt: claims.ExpiresAt.Time,
		})

		// store into context
		reqLogger := m.log.With().Str("user_id", claims.UserID.String()).Str("session_id", claims.SessionID.String()).Logger()
		newCtx := reqLogger.WithContext(ctx)
		c.SetContext(newCtx)

		// store into local and context
		fiber.StoreInContext(c, constant.KeyAuthUserID, user.ID)

		return c.Next()
	}
}

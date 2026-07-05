package handler

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/rizalarfiyan/be-plant-factory/internal/config"
	"github.com/rizalarfiyan/be-plant-factory/internal/delivery/http/response"
	domainError "github.com/rizalarfiyan/be-plant-factory/internal/domain/error"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/constant"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/helper"
	"github.com/rizalarfiyan/be-plant-factory/internal/usecase/auth"
	"github.com/samber/do/v2"
)

type AuthHandler interface {
	GoogleLogin(c fiber.Ctx) error
	GoogleCallback(c fiber.Ctx) error
	Me(c fiber.Ctx) error
	Logout(c fiber.Ctx) error
	RefreshToken(c fiber.Ctx) error
}

type authHandler struct {
	auth auth.AuthUseCase `do:""`
	conf *config.Config   `do:""`
}

func NewAuthHandler(i do.Injector) (AuthHandler, error) {
	return do.InvokeStruct[authHandler](i)
}

// GoogleLogin godoc
//
//	@Summary		Google Login
//	@Description	Google Login of the application
//	@ID				google-login
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			redirect_url	query		string																false	"URL to redirect to after successful login (optional, overrides default success redirect URL)"
//	@Success		200				{object}	response.BaseSwaggerResponse{data=response.AuthGoogleLoginResponse}	"Google login URL retrieved successfully. Available code (SUCCESS)"
//	@Failure		400				{object}	response.BaseSwaggerEmptyResponse{}									"Invalid redirect URL provided. Available code (INVALID_REDIRECT_URL)"
//	@Failure		500				{object}	response.BaseSwaggerEmptyResponse{}									"Failed to retrieve Google login URL. Available code (INTERNAL_SERVER_ERROR)"
//	@Router			/auth/google [get]
func (h authHandler) GoogleLogin(c fiber.Ctx) error {
	redirectUrl := c.Query("redirect_url")

	req := auth.GoogleRequest{
		BaseUrl: h.conf.Frontend.DefaultAllowedUrl(),
	}

	if redirectUrl != "" {
		u, err := url.Parse(redirectUrl)
		if err != nil {
			return response.New(c, code.InvalidRedirectURL)
		}

		baseUrl := fmt.Sprintf("%s://%s", u.Scheme, u.Host)
		if !h.conf.Frontend.IsValidAllowedUrl(baseUrl) {
			return response.New(c, code.InvalidRedirectURL)
		}

		req.BaseUrl = baseUrl
		req.RedirectUri = u.RequestURI()
	}

	authUrl, err := h.auth.GoogleLoginURL(c.Context(), req)
	if err != nil {
		return err
	}

	res := response.NewAuthGoogleLoginResponse(authUrl)
	return response.New(c, code.OK, response.WithData(res))
}

// GoogleCallback godoc
//
//	@Summary		Google Callback
//	@Description	Google Callback for the application
//	@ID				google-callback
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			code	query		string	true	"Authorization code from Google"
//	@Param			state	query		string	true	"State parameter for CSRF protection"
//	@Success		302		{string}	string	"Redirects to the success URL with authentication cookies set"
//	@Failure		302		{string}	string	"Redirects to the failure URL with error code as query parameter. Available code (INTERNAL_SERVER_ERROR, INVALID_STATE, PROVIDER_ERROR, TOKEN_EXCHANGE_FAILED, INVALID_CALLBACK, FAILED_CALLBACK)"
//	@Router			/auth/google/callback [get]
func (h authHandler) GoogleCallback(c fiber.Ctx) error {
	authCode := c.Query("code")
	state := c.Query("state")
	if authCode == "" || state == "" {
		return h.redirectToFailure(c, h.conf.Frontend.DefaultAllowedUrl(), code.InvalidCallback)
	}

	req := auth.CallbackRequest{
		State:      state,
		Code:       authCode,
		DeviceName: h.auth.GetDeviceName(c.Get(fiber.HeaderUserAgent)),
		IPAddress:  c.IP(),
	}

	res, err := h.auth.HandleGoogleCallback(c.Context(), req)
	if err != nil {
		if appErr, ok := domainError.IsAppError(err); ok {
			return h.redirectToFailure(c, res.BaseURL, appErr.Code)
		}

		return h.redirectToFailure(c, res.BaseURL, code.FailedCallback)
	}

	h.setCookies(c, res.TokenResponse)
	return h.redirectToSuccess(c, res.BaseURL, res.RedirectUri, res.IsCompleted)
}

func (h authHandler) redirectToFailure(c fiber.Ctx, baseURL string, code code.AppCode) error {
	u, _ := url.Parse(baseURL)
	u = u.JoinPath("login")
	params := url.Values{}
	params.Add("error_code", code.Code)
	u.RawQuery = params.Encode()
	return c.Redirect().To(u.String())
}

func (h authHandler) redirectToSuccess(c fiber.Ctx, baseURL, redirectURI string, isCompleted bool) error {
	u, _ := url.Parse(baseURL)

	target := "dashboard"
	if !isCompleted {
		target = "getting-started"
	} else if redirectURI != "" {
		target = redirectURI
	}

	if path, query, found := strings.Cut(target, "?"); found {
		u.Path, _ = url.JoinPath(u.Path, path)
		u.RawQuery = query
	} else {
		u.Path, _ = url.JoinPath(u.Path, target)
	}

	return c.Redirect().To(u.String())
}

func (h authHandler) setCookieHelper(c fiber.Ctx, name, value, path string, httpOnly bool, expires time.Time) {
	c.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		HTTPOnly: httpOnly,
		Secure:   true,
		SameSite: fiber.CookieSameSiteLaxMode,
		Expires:  expires,
		Domain:   ".plant-factory.com",
	})
}

func (h authHandler) setCookies(c fiber.Ctx, req auth.TokenResponse) {
	h.setCookieHelper(c, constant.CookieAccessToken, req.AccessToken, "/", true, req.AccessTokenExpiresAt)
	h.setCookieHelper(c, constant.CookieRefreshToken, req.RefreshToken, "/auth/refresh", true, req.RefreshTokenExpiresAt)
	h.setCookieHelper(c, constant.CookieIsLoggedIn, "true", "/", false, req.RefreshTokenExpiresAt)
}

func (h authHandler) clearCookies(c fiber.Ctx) {
	expires := time.Unix(0, 0)
	h.setCookieHelper(c, constant.CookieAccessToken, "", "/", true, expires)
	h.setCookieHelper(c, constant.CookieRefreshToken, "", "/auth/refresh", true, expires)
	h.setCookieHelper(c, constant.CookieIsLoggedIn, "", "/", false, expires)
}

// Me godoc
//
//	@Summary		Get Current User
//	@Description	Get the currently authenticated user's information
//	@ID				get-current-user
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Security		CookieAccessToken
//	@Success		200	{object}	response.BaseSwaggerResponse{data=response.AuthMeResponse}	"Current user retrieved successfully. Available code (SUCCESS)"
//	@Failure		401	{object}	response.BaseSwaggerEmptyResponse{}							"Unauthorized - user is not authenticated. Available code (UNAUTHORIZED)"
//	@Failure		422	{object}	response.BaseSwaggerEmptyResponse{}							"User not found. Available code (USER_NOT_FOUND)"
//	@Failure		500	{object}	response.BaseSwaggerEmptyResponse{}							"Internal Server Error. Available code (INTERNAL_SERVER_ERROR)"
//	@Router			/auth/me [get]
func (h authHandler) Me(c fiber.Ctx) error {
	user, ok := helper.GetLocalUser(c)
	if !ok {
		return response.New(c, code.Unauthorized)
	}

	res := response.NewAuthMeResponse(user)
	return response.New(c, code.OK, response.WithData(res))
}

// Logout godoc
//
//	@Summary		Logout
//	@Description	Logout the currently authenticated user by invalidating their session and tokens
//	@ID				logout
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Security		CookieAccessToken
//	@Success		200	{object}	response.BaseSwaggerEmptyResponse{}	"User logged out successfully. Available code (SUCCESS)"
//	@Failure		401	{object}	response.BaseSwaggerEmptyResponse{}	"Unauthorized - user is not authenticated. Available code (UNAUTHORIZED)"
//	@Failure		500	{object}	response.BaseSwaggerEmptyResponse{}	"Internal Server Error. Available code (INTERNAL_SERVER_ERROR)"
//	@Router			/auth/logout [post]
func (h authHandler) Logout(c fiber.Ctx) error {
	user, ok := helper.GetLocalUser(c)
	if !ok {
		return response.New(c, code.Unauthorized)
	}

	sessionID, ok := helper.GetLocalSessionID(c)
	if !ok {
		return response.New(c, code.Unauthorized)
	}

	accessToken, ok := helper.GetLocalAccessToken(c)
	if !ok {
		return response.New(c, code.Unauthorized)
	}

	req := auth.LogoutRequest{
		UserID:        user.ID,
		SessionID:     sessionID,
		AccessTokenID: accessToken.ID,
		ExpiredAt:     accessToken.ExpiredAt,
	}

	if err := h.auth.Logout(c.Context(), req); err != nil {
		return err
	}

	h.clearCookies(c)
	return response.New(c, code.OK)
}

// RefreshToken godoc
//
//	@Summary		Refresh Access Token
//	@Description	Refresh the access token using a valid refresh token
//	@ID				refresh-access-token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Security		CookieRefreshToken
//	@Success		200	{object}	response.BaseSwaggerEmptyResponse{}	"Access token refreshed successfully. Available code (SUCCESS)"
//	@Failure		401	{object}	response.BaseSwaggerEmptyResponse{}	"Unauthorized - invalid or expired refresh token. Available code (UNAUTHORIZED, INVALID_REFRESH_TOKEN, REFRESH_TOKEN_NOT_ELIGIBLE[remaining_seconds])"
//	@Failure		500	{object}	response.BaseSwaggerEmptyResponse{}	"Internal Server Error. Available code (INTERNAL_SERVER_ERROR)"
//	@Router			/auth/refresh [post]
func (h authHandler) RefreshToken(c fiber.Ctx) error {
	refreshToken := c.Cookies(constant.CookieRefreshToken)
	if refreshToken == "" || len(refreshToken) != 64 {
		return response.New(c, code.Unauthorized)
	}

	res, err := h.auth.RefreshToken(c.Context(), refreshToken)
	if err != nil {
		h.clearCookies(c)
		return err
	}

	h.setCookies(c, res)
	return response.New(c, code.OK)
}

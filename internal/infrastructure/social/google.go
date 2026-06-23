package social

import (
	"context"
	"time"

	"github.com/bytedance/sonic"
	"github.com/rizalarfiyan/be-plant-factory/internal/config"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type googleProvider struct {
	oauthConf *oauth2.Config
}

type googleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func NewGoogleProvider(cfg *config.Config) SocialAuthProvider {
	return &googleProvider{
		oauthConf: &oauth2.Config{
			ClientID:     cfg.Auth.Google.ClientID,
			ClientSecret: cfg.Auth.Google.ClientSecret,
			RedirectURL:  cfg.Auth.Google.RedirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

func (g *googleProvider) GetLoginURL(state string) string {
	return g.oauthConf.AuthCodeURL(state)
}

func (g *googleProvider) GetUserInfo(ctx context.Context, code string) (entity.SocialUser, error) {
	token, err := g.oauthConf.Exchange(ctx, code)
	if err != nil {
		return entity.SocialUser{}, err
	}

	reqCtx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	client := g.oauthConf.Client(reqCtx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return entity.SocialUser{}, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var user googleUser
	if err := sonic.ConfigStd.NewDecoder(resp.Body).Decode(&user); err != nil {
		return entity.SocialUser{}, err
	}

	return entity.SocialUser{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.GivenName,
		LastName:  user.FamilyName,
		Picture:   user.Picture,
		Provider:  entity.ProviderGoogle,
	}, nil
}

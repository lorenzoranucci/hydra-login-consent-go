package social_login_provider

import (
	"fmt"
	"net/url"

	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/domain"
)

type Google struct {
	id string
	googleClientID      int
	googleLoginCallbackUrl   string
	googleLoginEndpoint string
}

func (f Google) GetLoginURL(loginChallenge string) (*url.URL, error) {
	panic("implement me")
}

func (f Google) GetLoginEndpoint() url.URL {
	panic("implement me")
}

func NewGoogle(
	id string,
	googleClientID int,
	googleLoginCallbackEndpoint string,
	googleLoginEndpoint string,
) *Google {
	return &Google{
		id:                       id,
		googleClientID:         googleClientID,
		googleLoginCallbackUrl: googleLoginCallbackEndpoint,
		googleLoginEndpoint:    googleLoginEndpoint,
	}
}

func (f Google) GetID() string {
	return f.id
}

func (f Google) GetUserByToken(token string) (domain.SocialLoginProviderUser, error) {
	panic("implement me")
}

func (f Google) GetLogisnEndpoint(loginChallenge string) (*url.URL, error) {
	return url.Parse(
		fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&state=%s"),
	)
}

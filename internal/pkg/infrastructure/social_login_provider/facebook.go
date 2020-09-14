package social_login_provider

import (
	"fmt"
	"net/url"

	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/domain"
)

type Facebook struct {
	id string
	facebookClientID      int
	facebookLoginCallbackUrl   string
	facebookLoginEndpoint string
}

func (f Facebook) GetLoginURL(loginChallenge string) (*url.URL, error) {
	panic("implement me")
}

func (f Facebook) GetLoginEndpoint() url.URL {
	panic("implement me")
}

func NewFacebook(
	id string,
	facebookClientID int,
	facebookLoginCallbackEndpoint string,
	facebookLoginEndpoint string,
) *Facebook {
	return &Facebook{
		id:                       id,
		facebookClientID:         facebookClientID,
		facebookLoginCallbackUrl: facebookLoginCallbackEndpoint,
		facebookLoginEndpoint:    facebookLoginEndpoint,
	}
}

func (f Facebook) GetID() string {
	return f.id
}

func (f Facebook) GetUserByToken(token string) (domain.SocialLoginProviderUser, error) {
	panic("implement me")
}

func (f Facebook) GetLogisnEndpoint(loginChallenge string) (*url.URL, error) {
	return url.Parse(
		fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&state=%s"),
	)
}

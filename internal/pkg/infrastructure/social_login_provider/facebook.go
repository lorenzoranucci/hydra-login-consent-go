package social_login_provider

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/domain"
)

type Facebook struct {
	id           string
	clientID     int
	clientSecret     string
	redirectURI  string
	authEndpoint string // https://www.facebook.com/v8.0/dialog/oauth
	tokenEndpoint string // https://graph.facebook.com/v8.0/oauth/access_token
	verifyTokenEndpoint string // https://graph.facebook.com/v8.0/me
}

func NewFacebook(
	id string,
	clientID int,
	clientSecret string,
	redirectURI string,
	authEndpoint string,
	tokenEndpoint string,
	verifyTokenEndpoint string,
) *Facebook {
	return &Facebook{
		id: id,
		clientID: clientID,
		clientSecret: clientSecret,
		redirectURI: redirectURI,
		authEndpoint: authEndpoint,
		tokenEndpoint: tokenEndpoint,
		verifyTokenEndpoint: verifyTokenEndpoint,
	}
}

type FacebookAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func (f *Facebook) GetID() string {
	return f.id
}

func (f *Facebook) GetUserByToken(code string) (*domain.SocialLoginProviderUser, error) {
	accessTokenResponseHTTP, err := http.Get(
		fmt.Sprintf(
			"%s?client_id=%d&client_secret=%s&redirect_uri=%s&code=%s",
			f.tokenEndpoint,
			f.clientID,
			f.clientSecret,
			f.redirectURI,
			code,
		),
	)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(accessTokenResponseHTTP.Body)
	scanner.Scan()
	accessTokenResponseBody :=scanner.Bytes()

	accessTokenResponse := &FacebookAccessTokenResponse{}
	err = json.Unmarshal(accessTokenResponseBody, accessTokenResponse)
	if err != nil {
		return nil, err
	}

	facebookUserResponseHTTP, err := http.Get(
		fmt.Sprintf(
			"%s?fields=id,email,first_name,last_name&access_token=%s",
			f.verifyTokenEndpoint,
			accessTokenResponse.AccessToken,
		),
	)

	if err != nil {
		return nil, err
	}

	scanner = bufio.NewScanner(facebookUserResponseHTTP.Body)
	scanner.Scan()
	facebookUserResponseBody :=scanner.Bytes()

	facebookUser := &domain.SocialLoginProviderUser{}
	err = json.Unmarshal(facebookUserResponseBody, facebookUser)
	if err != nil {
		return nil, err
	}

	return facebookUser, nil
}

func (f *Facebook) GetLoginEndpoint(state string) (*url.URL, error) {
	return url.Parse(
		fmt.Sprintf(
			"%s?client_id=%d&redirect_uri=%s&state=%s",
			f.authEndpoint,
			f.clientID,
			f.redirectURI,
			state,
		),
	)
}

package social_login_provider

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/domain"
)

type Google struct {
	id           string
	clientID     int
	clientSecret     string
	redirectURI  string
	authEndpoint string // https://www.google.com/v8.0/dialog/oauth
	tokenEndpoint string // https://graph.google.com/v8.0/oauth/access_token
	verifyTokenEndpoint string // https://graph.google.com/v8.0/me
}

func NewGoogle(
	id string,
	clientID int,
	clientSecret string,
	redirectURI string,
	authEndpoint string,
	tokenEndpoint string,
	verifyTokenEndpoint string,
) *Google {
	return &Google{
		id: id,
		clientID: clientID,
		clientSecret: clientSecret,
		redirectURI: redirectURI,
		authEndpoint: authEndpoint,
		tokenEndpoint: tokenEndpoint,
		verifyTokenEndpoint: verifyTokenEndpoint,
	}
}

type GoogleAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func (g *Google) GetID() string {
	return g.id
}

func (g *Google) GetUserByToken(code string) (*domain.SocialLoginProviderUser, error) {
	accessTokenResponseHTTP, err := http.Get(
		fmt.Sprintf(
			"%s?client_id=%d&client_secret=%s&redirect_uri=%s&code=%s",
			g.tokenEndpoint,
			g.clientID,
			g.clientSecret,
			g.redirectURI,
			code,
		),
	)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(accessTokenResponseHTTP.Body)
	scanner.Scan()
	accessTokenResponseBody :=scanner.Bytes()

	accessTokenResponse := &GoogleAccessTokenResponse{}
	err = json.Unmarshal(accessTokenResponseBody, accessTokenResponse)
	if err != nil {
		return nil, err
	}

	googleUserResponseHTTP, err := http.Get(
		fmt.Sprintf(
			"%s?fields=id,email,first_name,last_name&access_token=%s",
			g.verifyTokenEndpoint,
			accessTokenResponse.AccessToken,
		),
	)

	if err != nil {
		return nil, err
	}

	scanner = bufio.NewScanner(googleUserResponseHTTP.Body)
	scanner.Scan()
	googleUserResponseBody :=scanner.Bytes()

	googleUser := &domain.SocialLoginProviderUser{}
	err = json.Unmarshal(googleUserResponseBody, googleUser)
	if err != nil {
		return nil, err
	}

	return googleUser, nil
}

func (g *Google) GetLoginEndpoint(loginChallenge string) (*url.URL, error) {
	return url.Parse(
		fmt.Sprintf(
			"%s?client_id=%s&redirect_uri=%s&state=%s",
			g.authEndpoint,
			g.clientID,
			g.redirectURI,
			loginChallenge,
		),
	)
}

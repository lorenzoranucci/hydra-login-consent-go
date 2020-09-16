package social_login_provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/domain"
)

type Facebook struct {
	id           string
	clientID     string
	clientSecret     string
	redirectURI  string
	authEndpoint string // https://www.facebook.com/v8.0/dialog/oauth
	tokenEndpoint string // https://graph.facebook.com/v8.0/oauth/access_token
	verifyTokenEndpoint string // https://graph.facebook.com/v8.0/me
}

func NewFacebook(
	id string,
	clientID string,
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

func (f *Facebook) GetID() string {
	return f.id
}

func (f *Facebook) GetUserByToken(code string) (*domain.SocialLoginProviderUser, error) {
	accessTokenGetURL := fmt.Sprintf(
		"%s?client_id=%s&client_secret=%s&redirect_uri=%s&code=%s",
		f.tokenEndpoint,
		f.clientID,
		f.clientSecret,
		f.redirectURI,
		code,
	)
	accessTokenResponseHTTP, err := http.Get(
		accessTokenGetURL,
	)
	if err != nil {
		return nil, err
	}

	accessTokenResponseBody, err := ioutil.ReadAll(accessTokenResponseHTTP.Body)
	if err != nil {
		return nil, err
	}

	accessTokenResponse := &FacebookAccessTokenResponse{}
	err = json.Unmarshal(accessTokenResponseBody, accessTokenResponse)
	if err != nil {
		return nil, err
	}

	if accessTokenResponse.AccessToken == nil {
		return nil, fmt.Errorf("invalid access token")
	}

	tokenVerifyURL := fmt.Sprintf(
		"%s?fields=id,email,first_name,last_name&access_token=%s",
		f.verifyTokenEndpoint,
		*accessTokenResponse.AccessToken,
	)
	facebookUserResponseHTTP, err := http.Get(
		tokenVerifyURL,
	)

	if err != nil {
		return nil, err
	}

	facebookUserResponseBody, err := ioutil.ReadAll(facebookUserResponseHTTP.Body)
	if err != nil {
		return nil, err
	}

	facebookUser := &FacebookUser{}
	err = json.Unmarshal(facebookUserResponseBody, facebookUser)
	if err != nil {
		return nil, err
	}

	socialLoginProviderUser, err := domain.NewSocialLoginProviderUser(
		facebookUser.ID,
		facebookUser.Email,
		facebookUser.FirstName,
		facebookUser.LastName,
		f.id,
		*accessTokenResponse.AccessToken,
	)
	if err != nil {
		return nil, err
	}

	return socialLoginProviderUser, nil
}

func (f *Facebook) GetLoginEndpoint(loginChallenge string) (*url.URL, error) {
	return url.Parse(
		fmt.Sprintf(
			"%s?client_id=%s&redirect_uri=%s&state=%s&scope=email&auth_type=rerequest",
			f.authEndpoint,
			f.clientID,
			f.redirectURI,
			loginChallenge,
		),
	)
}

type FacebookAccessTokenResponse struct {
	AccessToken *string `json:"access_token,omitempty"`
}

type FacebookUser struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

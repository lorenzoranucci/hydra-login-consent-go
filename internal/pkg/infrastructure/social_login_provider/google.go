package social_login_provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/domain"
)

type Google struct {
	id                  string
	clientID            string
	clientSecret        string
	redirectURI         string
	authEndpoint        string // https://www.google.com/v8.0/dialog/oauth
	tokenEndpoint       string // https://graph.google.com/v8.0/oauth/access_token
	verifyTokenEndpoint string // https://graph.google.com/v8.0/me
}

func NewGoogle(
	id string,
	clientID string,
	clientSecret string,
	redirectURI string,
	authEndpoint string,
	tokenEndpoint string,
	verifyTokenEndpoint string,
) *Google {
	return &Google{
		id:                  id,
		clientID:            clientID,
		clientSecret:        clientSecret,
		redirectURI:         redirectURI,
		authEndpoint:        authEndpoint,
		tokenEndpoint:       tokenEndpoint,
		verifyTokenEndpoint: verifyTokenEndpoint,
	}
}

func (g *Google) GetID() string {
	return g.id
}

func (g *Google) GetUserByToken(code string) (*domain.SocialLoginProviderUser, error) {
	accessTokenBody := fmt.Sprintf(
		"client_id=%s&client_secret=%s&redirect_uri=%s&code=%s&grant_type=authorization_code",
		g.clientID,
		g.clientSecret,
		g.redirectURI,
		code,
	)
	accessTokenResponseHTTP, err := http.Post(
		g.tokenEndpoint,
		"application/x-www-form-urlencoded",
		bytes.NewReader([]byte(accessTokenBody)),
	)
	if err != nil {
		return nil, err
	}

	accessTokenResponseBody, err := ioutil.ReadAll(accessTokenResponseHTTP.Body)
	if err != nil {
		return nil, err
	}

	accessTokenResponse := &GoogleAccessTokenResponse{}
	err = json.Unmarshal(accessTokenResponseBody, accessTokenResponse)
	if err != nil {
		return nil, err
	}

	if accessTokenResponse.AccessToken == nil {
		return nil, fmt.Errorf("invalid access token")
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", g.verifyTokenEndpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *accessTokenResponse.AccessToken))
	googleUserResponseHTTP, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	googleUserResponseBody, err := ioutil.ReadAll(googleUserResponseHTTP.Body)
	if err != nil {
		return nil, err
	}

	googleUser := &GoogleUser{}
	err = json.Unmarshal(googleUserResponseBody, googleUser)
	if err != nil {
		return nil, err
	}

	socialLoginProviderUser, err := g.getDomainUserFromGoogleUser(googleUser, accessTokenResponse)
	if err != nil {
		return nil, err
	}

	return socialLoginProviderUser, nil
}

func (g *Google) GetLoginEndpoint(loginChallenge string) (*url.URL, error) {
	return url.Parse(
		fmt.Sprintf(
			"%s?client_id=%s&redirect_uri=%s&state=%s&"+
				"scope=https://www.googleapis.com/auth/userinfo.profile+"+
				"https://www.googleapis.com/auth/userinfo.email&"+
				"access_type=offline&response_type=code",
			g.authEndpoint,
			g.clientID,
			g.redirectURI,
			loginChallenge,
		),
	)
}

type GoogleAccessTokenResponse struct {
	AccessToken *string `json:"access_token,omitempty"`
}

type GoogleUser struct {
	ID             string                   `json:"resourceName"`
	Names          GoogleUserNames          `json:"names"`
	EmailAddresses GoogleUserEmailAddresses `json:"emailAddresses"`
}

type GoogleUserNames []struct {
	Metadata  GoogleMetadata
	FirstName string `json:"givenName"`
	LastName  string `json:"familyName"`
}

type GoogleUserEmailAddresses []struct {
	Metadata GoogleMetadata
	Value    string `json:"value"`
}

type GoogleMetadata struct {
	Primary bool `json:"primary"`
}

func (g *Google) getDomainUserFromGoogleUser(
	googleUser *GoogleUser,
	accessTokenResponse *GoogleAccessTokenResponse,
) (*domain.SocialLoginProviderUser, error) {
	emailAddress := ""
	if len(googleUser.EmailAddresses) > 0 {
		emailAddress = googleUser.EmailAddresses[0].Value
		if len(googleUser.EmailAddresses) > 1 {
			for _, googleUserEmailAddress := range googleUser.EmailAddresses {
				if googleUserEmailAddress.Metadata.Primary {
					emailAddress = googleUserEmailAddress.Value
					break
				}
			}
		}
	}

	firstName := ""
	lastName := ""
	if len(googleUser.Names) > 0 {
		firstName = googleUser.Names[0].FirstName
		lastName = googleUser.Names[0].LastName
		if len(googleUser.Names) > 1 {
			for _, googleUserName := range googleUser.Names {
				if googleUserName.Metadata.Primary {
					firstName = googleUserName.FirstName
					lastName = googleUserName.LastName
					break
				}
			}
		}
	}
	return domain.NewSocialLoginProviderUser(
		googleUser.ID,
		emailAddress,
		firstName,
		lastName,
		g.id,
		*accessTokenResponse.AccessToken,
	)
}

package domain

import (
	"fmt"
	"net/url"
)

type SocialLoginProvider interface {
	GetID() string
	GetUserByToken(token string) (*SocialLoginProviderUser, error)
	GetLoginEndpoint(loginChallenge string) (*url.URL, error)
}

type SocialLoginProviderUser struct {
	id        string
	email     string
	firstName string
	lastName  string

	socialLoginProviderID string
	accessToken string
}

func (s *SocialLoginProviderUser) Id() string {
	return s.id
}

func (s *SocialLoginProviderUser) Email() string {
	return s.email
}

func (s *SocialLoginProviderUser) FirstName() string {
	return s.firstName
}

func (s *SocialLoginProviderUser) LastName() string {
	return s.lastName
}

func (s *SocialLoginProviderUser) SocialLoginProviderID() string {
	return s.socialLoginProviderID
}

func (s *SocialLoginProviderUser) AccessToken() string {
	return s.accessToken
}

func NewSocialLoginProviderUser(
	id string,
	email string,
	firstName string,
	lastName string,
	socialLoginProviderId string,
	accessToken string,
) (*SocialLoginProviderUser, error) {
	if id == "" {
		return nil, fmt.Errorf("missing social user id")
	}

	if email == "" {
		return nil, fmt.Errorf("missing social user email")
	}

	if firstName == "" {
		return nil, fmt.Errorf("missing social user first name")
	}

	if lastName == "" {
		return nil, fmt.Errorf("missing social user last name")
	}

	if socialLoginProviderId == "" {
		return nil, fmt.Errorf("missing social login provider id")
	}

	if accessToken == "" {
		return nil, fmt.Errorf("missing access token")
	}

	return &SocialLoginProviderUser{
		id: id,
		email: email,
		firstName: firstName,
		lastName: lastName,
		socialLoginProviderID: socialLoginProviderId,
		accessToken: accessToken,
	}, nil
}

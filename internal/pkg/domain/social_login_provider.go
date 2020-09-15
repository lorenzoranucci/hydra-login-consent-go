package domain

import "net/url"

type SocialLoginProviderUser struct {
	ID        string
	Email     string
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type SocialLoginProvider interface {
	GetID() string
	GetUserByToken(token string) (*SocialLoginProviderUser, error)
	GetLoginEndpoint(loginChallenge string) (*url.URL, error)
}

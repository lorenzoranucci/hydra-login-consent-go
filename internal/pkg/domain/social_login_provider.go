package domain

import "net/url"

type SocialLoginProviderUser struct {
	ID string
	Email string
	Name string
	SecondName string
}

type SocialLoginProvider interface {
	GetID() string
	GetUserByToken(token string) (SocialLoginProviderUser, error)
	GetLoginURL(loginChallenge string) (*url.URL, error)
}

package domain

import (
	"fmt"

	"github.com/google/uuid"
)

type User struct {
	UUID                     uuid.UUID
	Email                    string
	Name                     string
	SecondName               string
	AutoLoginToken           string
	SocialLoginProviderUsers UserSocialLoginProviderUsers
}

type UserSocialLoginProviderUsers map[SocialLoginProvider]SocialLoginProviderUser

func CreateUserFromSocialLoginProviderUser(
	socialLoginProviderUser SocialLoginProviderUser,
	socialLoginProvider SocialLoginProvider,
) *User {
	return &User{
		UUID:       uuid.New(),
		Email:      socialLoginProviderUser.Email,
		Name:       socialLoginProviderUser.Name,
		SecondName: socialLoginProviderUser.SecondName,
		SocialLoginProviderUsers: UserSocialLoginProviderUsers{
			socialLoginProvider: socialLoginProviderUser,
		},
	}
}

func (u *User) AddSocialLoginProviderUser(
	socialLoginProviderUser SocialLoginProviderUser,
	socialLoginProvider SocialLoginProvider,
) error {
	_, found := u.SocialLoginProviderUsers[socialLoginProvider]
	if found {
		return fmt.Errorf(
			"user '%s' already has an social user associated for provider %s",
			u.UUID.String(),
			socialLoginProvider.GetID(),
		)
	}

	if u.Email != socialLoginProviderUser.Email {
		return fmt.Errorf(
			"user '%s' email '%s' is not equal to social '%s' user '%s' email '%s'",
			u.UUID.String(),
			u.Email,
			socialLoginProvider.GetID(),
			socialLoginProviderUser.ID,
			socialLoginProviderUser.Email,
		)
	}
	u.SocialLoginProviderUsers[socialLoginProvider] = socialLoginProviderUser
	return nil
}

type UserRepository interface {
	FindByEmail(email string) (*User, bool, error)
	FindByEmailAndPassword(email string) (*User, bool, error)
	FindBySocialLoginProviderUserID(
		socialLoginProviderUserID string,
		socialLoginProvider SocialLoginProvider,
	) (*User, bool, error)
	FindByAutoLoginToken(autoLoginToken string) (*User, bool, error)
}

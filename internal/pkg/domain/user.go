package domain

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             int
	UUID           uuid.UUID
	Email          string
	Password       string
	FirstName      string
	LastName       string
	AutoLoginToken string
	Roles          []string

	SocialLoginProviderUsers UserSocialLoginProviderUsers
}

type UserSocialLoginProviderUsers map[SocialLoginProvider]SocialLoginProviderUser

func CreateUserFromSocialLoginProviderUser(
	socialLoginProviderUser SocialLoginProviderUser,
	socialLoginProvider SocialLoginProvider,
) *User {
	return &User{
		UUID:           uuid.New(),
		Email:          socialLoginProviderUser.Email,
		Password:       password(15),
		FirstName:      socialLoginProviderUser.FirstName,
		LastName:       socialLoginProviderUser.LastName,
		AutoLoginToken: uuid.New().String(),
		Roles:          []string{},
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
	FindByEmailAndPassword(email string, password string) (*User, bool, error)
	FindBySocialLoginProviderUserID(
		socialLoginProviderUserID string,
		socialLoginProvider SocialLoginProvider,
	) (*User, bool, error)
	FindByAutoLoginToken(autoLoginToken string) (*User, bool, error)
	Persist(user User) error
}

func password(length int) string {
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

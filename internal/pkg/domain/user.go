package domain

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type User struct {
	id             int
	uuid           uuid.UUID
	email          string
	password       string
	firstName      string
	lastName       string
	autoLoginToken string
	roles          []string

	socialLoginProviderUsers UserSocialLoginProviderUsers
}

func NewUser(
	uuid uuid.UUID,
	email string,
	password string,
	firstName string,
	lastName string,
	autoLoginToken string,
	roles []string,
	socialLoginProviderUsers UserSocialLoginProviderUsers,
) *User {
	return &User{
		uuid:                     uuid,
		email:                    email,
		password:                 password,
		firstName:                firstName,
		lastName:                 lastName,
		autoLoginToken:           autoLoginToken,
		roles:                    roles,
		socialLoginProviderUsers: socialLoginProviderUsers,
	}
}

func (u *User) Id() int {
	return u.id
}

func (u *User) Uuid() uuid.UUID {
	return u.uuid
}

func (u *User) Email() string {
	return u.email
}

func (u *User) Password() string {
	return u.password
}

func (u *User) FirstName() string {
	return u.firstName
}

func (u *User) LastName() string {
	return u.lastName
}

func (u *User) AutoLoginToken() string {
	return u.autoLoginToken
}

func (u *User) Roles() []string {
	return u.roles
}

func (u *User) SocialLoginProviderUsers() UserSocialLoginProviderUsers {
	return u.socialLoginProviderUsers
}

type UserSocialLoginProviderUsers map[string]*SocialLoginProviderUser

func CreateUserFromSocialLoginProviderUser(
	socialLoginProviderUser *SocialLoginProviderUser,
) (*User, error) {
	return NewUser(
		uuid.New(),
		socialLoginProviderUser.Email(),
		password(15),
		socialLoginProviderUser.firstName,
		socialLoginProviderUser.lastName,
		generateAutoLoginToken(),
		[]string{},
		UserSocialLoginProviderUsers{
			socialLoginProviderUser.SocialLoginProviderID(): socialLoginProviderUser,
		},
	), nil
}

func (u *User) AddSocialLoginProviderUser(
	socialLoginProviderUser *SocialLoginProviderUser,
) error {
	_, found := u.socialLoginProviderUsers[socialLoginProviderUser.socialLoginProviderID]
	if found {
		return fmt.Errorf(
			"user '%s' already has an social user associated for provider %s",
			u.uuid.String(),
			socialLoginProviderUser.socialLoginProviderID,
		)
	}

	if u.email != socialLoginProviderUser.email {
		return fmt.Errorf(
			"user '%s' email '%s' is not equal to social '%s' user '%s' email '%s'",
			u.uuid.String(),
			u.email,
			socialLoginProviderUser.socialLoginProviderID,
			socialLoginProviderUser.id,
			socialLoginProviderUser.email,
		)
	}
	u.socialLoginProviderUsers[socialLoginProviderUser.socialLoginProviderID] = socialLoginProviderUser
	return nil
}

type UserRepository interface {
	FindByEmail(email string) (*User, bool, error)
	FindByEmailAndPassword(email string, password string) (*User, bool, error)
	FindBySocialLoginProviderUser(
		socialLoginProviderUser *SocialLoginProviderUser,
	) (*User, bool, error)
	FindByAutoLoginToken(autoLoginToken string) (*User, bool, error)
	Persist(user User) error
}

/** todo tbd*/
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

// todo tbd
func generateAutoLoginToken() string {
	return uuid.New().String()
}

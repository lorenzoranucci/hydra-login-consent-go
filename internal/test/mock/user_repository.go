package mock

import (
"github.com/google/uuid"
"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/domain"
)

var usersByEmail = map[string]*domain.User{
	"foo@bar.com": {
		ID:                       1,
		UUID:                     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		Email:                    "foo@bar.com",
		FirstName:                "Foo",
		LastName:                 "Bar",
		AutoLoginToken:           "12345",
		Roles:                    []string{"form_farmer"},
		SocialLoginProviderUsers: nil,
	},
}

var usersByEmailAndPassword = map[string]*domain.User{
	"foo@bar.com1234": {
		ID:                       1,
		UUID:                     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		Email:                    "foo@bar.com",
		Password:                 "1234",
		FirstName:                "Foo",
		LastName:                 "Bar",
		AutoLoginToken:           "12345",
		Roles:                    []string{"form_farmer"},
		SocialLoginProviderUsers: nil,
	},
}

var usersBySocial = map[string]*domain.User{
}

var usersByAltk = map[string]*domain.User{
	"12345": {
		ID:                       1,
		UUID:                     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		Email:                    "foo@bar.com",
		Password:                 "1234",
		FirstName:                "Foo",
		LastName:                 "Bar",
		AutoLoginToken:           "12345",
		Roles:                    []string{"form_farmer"},
		SocialLoginProviderUsers: nil,
	},
}

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (u *UserRepository) FindByEmail(email string) (*domain.User, bool, error) {
	user, found :=  usersByEmail[email]
	return user, found, nil
}

func (u *UserRepository) FindByEmailAndPassword(email string, password string) (*domain.User, bool, error) {
	user, found :=  usersByEmailAndPassword[email+password]
	return user, found, nil
}

func (u *UserRepository) FindBySocialLoginProviderUserID(
	socialLoginProviderUserID string,
	socialLoginProvider domain.SocialLoginProvider,
) (*domain.User, bool, error) {
	user, found :=  usersByEmailAndPassword[socialLoginProviderUserID+socialLoginProvider.GetID()]
	return user, found, nil
}

func (u *UserRepository) FindByAutoLoginToken(autoLoginToken string) (*domain.User, bool, error) {
	user, found :=  usersByAltk[autoLoginToken]
	return user, found, nil
}

func (u *UserRepository) Persist(user domain.User) error {
	usersByAltk[user.AutoLoginToken] = &user
	usersByEmail[user.Email] = &user
	usersByEmailAndPassword[user.Email+user.Password] = &user
	for social, socialUser := range user.SocialLoginProviderUsers {
		usersBySocial[socialUser.ID+social.GetID()] = &user
	}
	return nil
}


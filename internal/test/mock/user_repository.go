package mock

import (
	"github.com/google/uuid"
	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/domain"
)

var fooBar = domain.NewUser(
	uuid.MustParse("00000000-0000-0000-0000-000000000001"),
	"foo@bar.com",
	"1234",
	"Foo",
	"Bar",
	"12345",
	[]string{"form_farmer"},
	nil,
)

var usersByEmail = map[string]*domain.User{
	"foo@bar.com": fooBar,
}

var usersByEmailAndPassword = map[string]*domain.User{
	"foo@bar.com1234": fooBar,
}

var usersBySocial = map[string]*domain.User{
}

var usersByAltk = map[string]*domain.User{
	"12345": fooBar,
}

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (u *UserRepository) FindByEmail(email string) (*domain.User, bool, error) {
	user, found := usersByEmail[email]
	return user, found, nil
}

func (u *UserRepository) FindByEmailAndPassword(email string, password string) (*domain.User, bool, error) {
	user, found := usersByEmailAndPassword[email+password]
	return user, found, nil
}

func (u *UserRepository) FindBySocialLoginProviderUser(
	socialLoginProviderUser *domain.SocialLoginProviderUser,
) (*domain.User, bool, error) {
	user, found := usersBySocial[socialLoginProviderUser.Id()+socialLoginProviderUser.SocialLoginProviderID()]
	return user, found, nil
}

func (u *UserRepository) FindByAutoLoginToken(autoLoginToken string) (*domain.User, bool, error) {
	user, found := usersByAltk[autoLoginToken]
	return user, found, nil
}

func (u *UserRepository) Persist(user domain.User) error {
	usersByAltk[user.AutoLoginToken()] = &user
	usersByEmail[user.Email()] = &user
	usersByEmailAndPassword[user.Email()+user.Password()] = &user
	for _, socialUser := range user.SocialLoginProviderUsers() {
		usersBySocial[socialUser.Id()+socialUser.SocialLoginProviderID()] = &user
	}
	return nil
}

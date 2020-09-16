package mysql

import (
	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/domain"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (u *UserRepository) FindByEmail(email string) (*domain.User, bool, error) {
	panic("implement me")
}

func (u *UserRepository) FindByEmailAndPassword(email string, password string) (*domain.User, bool, error) {
	panic("implement me")
}

func (u *UserRepository) FindBySocialLoginProviderUser(
	socialLoginProviderUser *domain.SocialLoginProviderUser,
) (*domain.User, bool, error) {
	panic("implement me")
}

func (u *UserRepository) FindByAutoLoginToken(autoLoginToken string) (*domain.User, bool, error) {
	panic("implement me")
}

func (u *UserRepository) Persist(user domain.User) error {
	panic("implement me")
}

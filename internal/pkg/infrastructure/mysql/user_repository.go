package mysql

import "github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/domain"

type UserRepository struct {
	
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (u UserRepository) FindByEmail(email string) (*domain.User, bool, error) {
	panic("implement me")
}

func (u UserRepository) FindByEmailAndPassword(email string) (*domain.User, bool, error) {
	panic("implement me")
}

func (u UserRepository) FindBySocialLoginProviderUserID(socialLoginProviderUserID string, socialLoginProvider domain.SocialLoginProvider) (*domain.User, bool, error) {
	panic("implement me")
}

func (u UserRepository) FindByAutoLoginToken(autoLoginToken string) (*domain.User, bool, error) {
	panic("implement me")
}


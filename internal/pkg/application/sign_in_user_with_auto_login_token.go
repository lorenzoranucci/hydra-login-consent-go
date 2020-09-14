package application

import "github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/domain"

type SignInUserWithAutoLoginTokenRequest struct {
	AutoLoginToken string
}

type SignInWithAutoLoginTokenServiceInterface interface {
	Execute(
		signInRequest SignInUserWithAutoLoginTokenRequest,
	) (*domain.User, bool, error)
}

type SignInWithAutoLoginTokenService struct {
	userRepository      domain.UserRepository
}

func NewSignInWithAutoLoginTokenService(userRepository domain.UserRepository) *SignInWithAutoLoginTokenService {
	return &SignInWithAutoLoginTokenService{userRepository: userRepository}
}

func (s SignInWithAutoLoginTokenService) Execute(
	signInRequest SignInUserWithAutoLoginTokenRequest,
) (*domain.User, bool, error) {
	return s.userRepository.FindByAutoLoginToken(signInRequest.AutoLoginToken)
}

package application

import "github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/domain"

type SignInUserWithEmailAndPasswordRequest struct {
	Email string
	Password string
}

type SignInUserWithEmailAndPasswordServiceInterface interface {
	Execute(
		signInRequest SignInUserWithEmailAndPasswordRequest,
	) (*domain.User, bool, error)
}

type SignInUserWithEmailAndPasswordService struct {
	userRepository      domain.UserRepository
}

func NewSignInUserWithEmailAndPasswordService(userRepository domain.UserRepository) *SignInUserWithEmailAndPasswordService {
	return &SignInUserWithEmailAndPasswordService{userRepository: userRepository}
}

func (s SignInUserWithEmailAndPasswordService) Execute(
	signInRequest SignInUserWithEmailAndPasswordRequest,
) (*domain.User, bool, error) {
	return s.userRepository.FindByEmailAndPassword(signInRequest.Password)
}

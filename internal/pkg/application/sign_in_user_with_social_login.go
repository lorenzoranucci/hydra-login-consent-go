package application

import (
	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/domain"
)

type SignInUserWithSocialLoginRequest struct {
	SocialLoginProviderToken string
	SocialLoginProvider      domain.SocialLoginProvider
}

type SignInUserWithSocialLoginServiceInterface interface {
	Execute(
		signInRequest SignInUserWithSocialLoginRequest,
	) (*domain.User, error)
}

type SignInWithSocialLoginService struct {
	userRepository domain.UserRepository
}

func NewSignInUserWithSocialLoginService(userRepository domain.UserRepository) *SignInWithSocialLoginService {
	return &SignInWithSocialLoginService{userRepository: userRepository}
}

func (s SignInWithSocialLoginService) Execute(
	signInRequest SignInUserWithSocialLoginRequest,
) (*domain.User, error) {
	socialLoginProviderUser, err := signInRequest.SocialLoginProvider.GetUserByToken(
		signInRequest.SocialLoginProviderToken,
	)
	if err != nil {
		return nil, err
	}

	user, found, err := s.userRepository.FindBySocialLoginProviderUserID(
		socialLoginProviderUser.ID,
		signInRequest.SocialLoginProvider,
	)
	if err != nil {
		return nil, err
	}

	if !found {
		user, found, err = s.userRepository.FindByEmail(
			socialLoginProviderUser.Email,
		)
		if err != nil {
			return nil, err
		}

		if found {
			err = user.AddSocialLoginProviderUser(socialLoginProviderUser, signInRequest.SocialLoginProvider)
			if err != nil {
				return nil, err
			}
		} else {
			user = domain.CreateUserFromSocialLoginProviderUser(
				socialLoginProviderUser,
				signInRequest.SocialLoginProvider,
			)
		}
	}

	return user, nil
}

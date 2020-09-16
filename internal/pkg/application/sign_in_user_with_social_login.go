package application

import (
	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/domain"
)

type SignInUserWithSocialLoginRequest struct {
	SocialLoginProviderUser      *domain.SocialLoginProviderUser
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
	domainUserBySocialLoginProviderUserID, foundBySocialLoginProviderUserID, err := s.userRepository.
		FindBySocialLoginProviderUser(
			signInRequest.SocialLoginProviderUser,
		)
	if err != nil {
		return nil, err
	}

	if !foundBySocialLoginProviderUserID {
		domainUserByEmail, foundByEmail, err := s.userRepository.FindByEmail(
			signInRequest.SocialLoginProviderUser.Email(),
		)
		if err != nil {
			return nil, err
		}

		if foundByEmail {
			return s.addSocialLoginProviderUserToDomainUser(signInRequest, domainUserByEmail, signInRequest.SocialLoginProviderUser)
		} else {
			return s.createANewUserFromSocialLoginProviderUser(signInRequest.SocialLoginProviderUser, err)
		}
	}

	return domainUserBySocialLoginProviderUserID, nil
}

func (s SignInWithSocialLoginService) addSocialLoginProviderUserToDomainUser(signInRequest SignInUserWithSocialLoginRequest, domainUserByEmail *domain.User, socialLoginProviderUser *domain.SocialLoginProviderUser) (*domain.User, error) {
	err := domainUserByEmail.AddSocialLoginProviderUser(signInRequest.SocialLoginProviderUser)
	if err != nil {
		return nil, err
	}
	err = s.userRepository.Persist(*domainUserByEmail)
	if err != nil {
		return nil, err
	}

	return domainUserByEmail, nil
}

func (s SignInWithSocialLoginService) createANewUserFromSocialLoginProviderUser(
	socialLoginProviderUser *domain.SocialLoginProviderUser,
	err error,
) (*domain.User, error) {
	newDomainUser, err := domain.CreateUserFromSocialLoginProviderUser(
		socialLoginProviderUser,
	)

	if err != nil {
		return nil, err
	}

	err = s.userRepository.Persist(*newDomainUser)
	if err != nil {
		return nil, err
	}

	return newDomainUser, nil
}

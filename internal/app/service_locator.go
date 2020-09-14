package app

import (
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/application"
	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/domain"
	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/infrastructure/hydra"
	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/infrastructure/mysql"
	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/infrastructure/social_login_provider"
)

var serviceLocator *ServiceLocator

type ServiceLocator struct {
	userRepository domain.UserRepository

	signInUserWithAutoLoginTokenService   application.SignInWithAutoLoginTokenServiceInterface
	signInUserWithEmailAndPasswordService application.SignInUserWithEmailAndPasswordServiceInterface
	signInUserWithSocialLoginService      application.SignInUserWithSocialLoginServiceInterface

	hydraClient                 hydra.HydraClientInterface
	socialLoginProviderFactory  social_login_provider.FactoryInterface
	facebookSocialLoginProvider domain.SocialLoginProvider
	googleSocialLoginProvider   domain.SocialLoginProvider

	FacebookSocialLoginProviderID string
	GoogleSocialLoginProviderID   string
}

func GetServiceLocator() *ServiceLocator {
	if serviceLocator == nil {
		serviceLocator = &ServiceLocator{}
	}
	return serviceLocator
}

func ResetServiceLocator() *ServiceLocator {
	serviceLocator = &ServiceLocator{}
	return serviceLocator
}

func (sl *ServiceLocator) UserRepository() domain.UserRepository {
	return sl.mysqlUserRepository()
}

func (sl *ServiceLocator) mysqlUserRepository() mysql.UserRepository {
	_, found := sl.userRepository.(mysql.UserRepository)
	if !found {
		sl.userRepository = *mysql.NewUserRepository()
	}

	return sl.userRepository.(mysql.UserRepository)
}

func (sl *ServiceLocator) CreateSignInUserWithEmailAndPasswordService() application.SignInUserWithEmailAndPasswordServiceInterface {
	if sl.signInUserWithEmailAndPasswordService == nil {
		sl.signInUserWithEmailAndPasswordService = application.NewSignInUserWithEmailAndPasswordService(
			sl.UserRepository(),
		)
	}

	return sl.signInUserWithEmailAndPasswordService
}

func (sl *ServiceLocator) CreateSignInUserWithAutoLoginTokenService() application.SignInWithAutoLoginTokenServiceInterface {
	if sl.signInUserWithAutoLoginTokenService == nil {
		sl.signInUserWithAutoLoginTokenService = application.NewSignInWithAutoLoginTokenService(
			sl.UserRepository(),
		)
	}

	return sl.signInUserWithAutoLoginTokenService
}

func (sl *ServiceLocator) CreateSignInUserWithSocialLoginService() application.SignInUserWithSocialLoginServiceInterface {
	if sl.signInUserWithSocialLoginService == nil {
		sl.signInUserWithSocialLoginService = application.NewSignInUserWithSocialLoginService(
			sl.UserRepository(),
		)
	}

	return sl.signInUserWithSocialLoginService
}

func (sl *ServiceLocator) CreateHydraClient() hydra.HydraClientInterface {
	if sl.hydraClient == nil {
		sl.hydraClient = hydra.NewHydraClientStruct(
			sl.hydraAdminURL(),
		)
	}

	return sl.hydraClient
}

func (sl *ServiceLocator) CreateSocialLoginProviderFactory() social_login_provider.FactoryInterface {
	if sl.socialLoginProviderFactory == nil {
		sl.socialLoginProviderFactory = social_login_provider.NewFactory(
			[]domain.SocialLoginProvider{
				sl.CreateFacebookSocialLoginProvider(),
				sl.CreateGoogleSocialLoginProvider(),
			},
		)
	}

	return sl.socialLoginProviderFactory
}

func (sl *ServiceLocator) CreateFacebookSocialLoginProvider() domain.SocialLoginProvider {
	if sl.facebookSocialLoginProvider == nil {
		clientID, err := strconv.Atoi(os.Getenv("FACEBOOK_CLIENT_ID"))
		if err != nil {
			panic(err)
		}

		sl.facebookSocialLoginProvider = social_login_provider.NewFacebook(
			sl.FacebookID(),
			clientID,
			os.Getenv("FACEBOOK_LOGIN_CALLBACK_ENDPOINT"),
			os.Getenv("FACEBOOK_LOGIN_ENDPOINT"),
		)
	}

	return sl.facebookSocialLoginProvider
}

func (sl *ServiceLocator) CreateGoogleSocialLoginProvider() domain.SocialLoginProvider {
	if sl.googleSocialLoginProvider == nil {
		clientID, err := strconv.Atoi(os.Getenv("GOOGLE_CLIENT_ID"))
		if err != nil {
			panic(err)
		}

		sl.googleSocialLoginProvider = social_login_provider.NewGoogle(
			sl.GoogleID(),
			clientID,
			os.Getenv("GOOGLE_LOGIN_CALLBACK_ENDPOINT"),
			os.Getenv("GOOGLE_LOGIN_ENDPOINT"),
		)
	}

	return sl.googleSocialLoginProvider
}

func (sl *ServiceLocator) FacebookID() string {
	return os.Getenv("FACEBOOK_ID")
}

func (sl *ServiceLocator) GoogleID() string {
	return os.Getenv("GOOGLE_ID")
}

func (sl *ServiceLocator) hydraAdminURL() *url.URL {
	hydraAdminURLString := os.Getenv("HYDRA_ADMIN_URL")
	if hydraAdminURLString == "" {
		panic(fmt.Errorf("HYDRA_ADMIN_URL not found"))
	}

	hydraAdminURL, err := url.Parse(hydraAdminURLString)
	if err != nil {
		panic(err)
	}
	return hydraAdminURL
}

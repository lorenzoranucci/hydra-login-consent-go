package social_login_provider

import (
	"fmt"

	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/domain"
)

type FactoryInterface interface {
	GetSocialLoginProviderByID(socialLoginProviderID string) (domain.SocialLoginProvider, error)
}

type Factory struct {
	enabledSocialLoginProviders []domain.SocialLoginProvider
}

func NewFactory(enabledSocialLoginProviders []domain.SocialLoginProvider) *Factory {
	return &Factory{enabledSocialLoginProviders: enabledSocialLoginProviders}
}

func (s Factory) GetSocialLoginProviderByID(socialLoginProviderID string) (domain.SocialLoginProvider, error) {
	for _, socialLoginProvider := range s.enabledSocialLoginProviders {
		if socialLoginProvider.GetID() == socialLoginProviderID {
			return socialLoginProvider, nil
		}
	}

	return nil, fmt.Errorf("unsupported social login provider '%s'", socialLoginProviderID)
}





package handler

import (
	"fmt"
	"net/http"

	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/application"
	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/infrastructure/hydra"
	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/infrastructure/social_login_provider"
)

type FacebookSocialLoginCallbackHandler struct {
	signInUserWithSocialLoginService application.SignInUserWithSocialLoginServiceInterface

	facebookSocialLoginProvider *social_login_provider.Facebook
	hydraClient                 hydra.HydraClientInterface
}

func NewFacebookSocialLoginCallbackHandler(
	signInUserWithSocialLoginService application.SignInUserWithSocialLoginServiceInterface,
	facebookSocialLoginProvider *social_login_provider.Facebook,
	hydraClient hydra.HydraClientInterface,
) *FacebookSocialLoginCallbackHandler {
	return &FacebookSocialLoginCallbackHandler{signInUserWithSocialLoginService: signInUserWithSocialLoginService, facebookSocialLoginProvider: facebookSocialLoginProvider, hydraClient: hydraClient}
}

func (h *FacebookSocialLoginCallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.handleGet(w, r)
	}
}

func (h *FacebookSocialLoginCallbackHandler) handleGet(w http.ResponseWriter, r *http.Request, ) {
	code := r.URL.Query().Get("code")
	loginChallenge := r.URL.Query().Get("state")

	_, err := h.hydraClient.GetLoginRequest(loginChallenge)
	if err != nil {
		fail(w, fmt.Errorf("invalid login_challenge for facebook social login"), 400)
		return
	}

	user, err := h.signInUserWithSocialLoginService.Execute(
		application.SignInUserWithSocialLoginRequest{
			SocialLoginProviderToken: code,
			SocialLoginProvider:      h.facebookSocialLoginProvider,
		},
	)

	if err != nil {
		err = h.hydraClient.RejectLoginAndRedirect(
			w,
			r,
			loginChallenge,
			fmt.Errorf("error social login with facebook and loginChallenge %s", loginChallenge),
		)

		if err != nil {
			fail(w, fmt.Errorf("cannot reject login challenge"), 500)
		}
		return
	}

	err = h.hydraClient.AcceptLoginAndRedirect(
		w,
		r,
		loginChallenge,
		user.Email,
		nil,
		nil,
	)
	if err != nil {
		fail(w, fmt.Errorf("cannot accept login challenge"), 500)
	}
}

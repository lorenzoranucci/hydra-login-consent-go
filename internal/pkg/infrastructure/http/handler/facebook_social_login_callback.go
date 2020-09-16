package handler

import (
	"fmt"
	"net/http"

	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/application"
	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/domain"
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

	err := h.validateLoginChallangeAndFailOnError(w, loginChallenge)
	if err != nil {
		return
	}

	facebookUser, err := h.facebookSocialLoginProvider.GetUserByToken(code)
	if err != nil {
		h.redirectToFacebookLoginEndpointAndFailOnError(w, r, loginChallenge)
		return
	}

	user, err := h.signInUserWithSocialLoginService.Execute(
		application.SignInUserWithSocialLoginRequest{
			SocialLoginProviderUser: facebookUser,
		},
	)

	if err != nil {
		h.rejectLoginChallengeAndFailOnError(w, r, err, loginChallenge)
		return
	}

	h.acceptLoginChallengeAndFailOnError(w, r, err, loginChallenge, user)
}

func (h *FacebookSocialLoginCallbackHandler) acceptLoginChallengeAndFailOnError(
	w http.ResponseWriter,
	r *http.Request,
	err error,
	loginChallenge string,
	user *domain.User,
) {
	err = h.hydraClient.AcceptLoginAndRedirect(
		w,
		r,
		loginChallenge,
		user.Email(),
	)
	if err != nil {
		fail(w, fmt.Errorf("cannot accept login challenge"), 500)
	}
}

func (h *FacebookSocialLoginCallbackHandler) rejectLoginChallengeAndFailOnError(
	w http.ResponseWriter,
	r *http.Request,
	err error,
	loginChallenge string,
) {
	err = h.hydraClient.RejectLoginAndRedirect(
		w,
		r,
		loginChallenge,
		fmt.Errorf("error social login with facebook and loginChallenge %s", loginChallenge),
	)

	if err != nil {
		fail(w, fmt.Errorf("cannot reject login challenge"), 500)
	}
}

func (h *FacebookSocialLoginCallbackHandler) redirectToFacebookLoginEndpointAndFailOnError(
	w http.ResponseWriter,
	r *http.Request,
	loginChallenge string,
) {
	endpoint, err := h.facebookSocialLoginProvider.GetLoginEndpoint(loginChallenge)
	if err != nil {
		fail(w, fmt.Errorf("cannot redirect to facebook login"), 500)
		return
	}
	http.Redirect(w, r, endpoint.String(), 307)
}

func (h *FacebookSocialLoginCallbackHandler) validateLoginChallangeAndFailOnError(w http.ResponseWriter, loginChallenge string) error {
	_, err := h.hydraClient.GetLoginRequest(loginChallenge)
	if err != nil {
		fail(w, fmt.Errorf("invalid login_challenge for facebook social login"), 400)
		return err
	}
	return nil
}

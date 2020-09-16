package handler

import (
	"fmt"
	"net/http"

	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/application"
	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/domain"
	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/infrastructure/hydra"
	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/infrastructure/social_login_provider"
)

type GoogleSocialLoginCallbackHandler struct {
	signInUserWithSocialLoginService application.SignInUserWithSocialLoginServiceInterface

	googleSocialLoginProvider *social_login_provider.Google
	hydraClient               hydra.HydraClientInterface
}

func NewGoogleSocialLoginCallbackHandler(
	signInUserWithSocialLoginService application.SignInUserWithSocialLoginServiceInterface,
	googleSocialLoginProvider *social_login_provider.Google,
	hydraClient hydra.HydraClientInterface,
) *GoogleSocialLoginCallbackHandler {
	return &GoogleSocialLoginCallbackHandler{signInUserWithSocialLoginService: signInUserWithSocialLoginService, googleSocialLoginProvider: googleSocialLoginProvider, hydraClient: hydraClient}
}

func (h *GoogleSocialLoginCallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.handleGet(w, r)
	}
}

func (h *GoogleSocialLoginCallbackHandler) handleGet(w http.ResponseWriter, r *http.Request, ) {
	code := r.URL.Query().Get("code")
	loginChallenge := r.URL.Query().Get("state")

	err := h.validateLoginChallangeAndFailOnError(w, loginChallenge)
	if err != nil {
		return
	}

	googleUser, err := h.googleSocialLoginProvider.GetUserByToken(code)
	if err != nil {
		h.redirectToGoogleLoginEndpointAndFailOnError(w, r, loginChallenge)
		return
	}

	user, err := h.signInUserWithSocialLoginService.Execute(
		application.SignInUserWithSocialLoginRequest{
			SocialLoginProviderUser: googleUser,
		},
	)

	if err != nil {
		h.rejectLoginChallengeAndFailOnError(w, r, err, loginChallenge)
		return
	}

	h.acceptLoginChallengeAndFailOnError(w, r, err, loginChallenge, user)
}

func (h *GoogleSocialLoginCallbackHandler) acceptLoginChallengeAndFailOnError(
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

func (h *GoogleSocialLoginCallbackHandler) rejectLoginChallengeAndFailOnError(
	w http.ResponseWriter,
	r *http.Request,
	err error,
	loginChallenge string,
) {
	err = h.hydraClient.RejectLoginAndRedirect(
		w,
		r,
		loginChallenge,
		fmt.Errorf("error social login with google and loginChallenge %s", loginChallenge),
	)

	if err != nil {
		fail(w, fmt.Errorf("cannot reject login challenge"), 500)
	}
}

func (h *GoogleSocialLoginCallbackHandler) redirectToGoogleLoginEndpointAndFailOnError(
	w http.ResponseWriter,
	r *http.Request,
	loginChallenge string,
) {
	endpoint, err := h.googleSocialLoginProvider.GetLoginEndpoint(loginChallenge)
	if err != nil {
		fail(w, fmt.Errorf("cannot redirect to google login"), 500)
		return
	}
	http.Redirect(w, r, endpoint.String(), 307)
}

func (h *GoogleSocialLoginCallbackHandler) validateLoginChallangeAndFailOnError(w http.ResponseWriter, loginChallenge string) error {
	_, err := h.hydraClient.GetLoginRequest(loginChallenge)
	if err != nil {
		fail(w, fmt.Errorf("invalid login_challenge for google social login"), 400)
		return err
	}
	return nil
}

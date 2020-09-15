package handler

import (
	"fmt"
	"net/http"

	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/domain"
	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/infrastructure/hydra"
)

type HydraConsentHandler struct {
	userRepository domain.UserRepository
	hydraClient    hydra.HydraClientInterface
}

func NewHydraConsentHandler(
	userRepository domain.UserRepository,
	hydraClient hydra.HydraClientInterface,
) *HydraConsentHandler {
	return &HydraConsentHandler{userRepository: userRepository, hydraClient: hydraClient}
}

func (h *HydraConsentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.handleConsentGet(w, r)
	}

	/* todo implement for third party clients consent handling
	since, for now, we have first party clients only, we can always silently skip consent step
	if r.Method == "POST" {
		h.handleConsentPost(w, r)
	}*/
}

func (h *HydraConsentHandler) handleConsentGet(w http.ResponseWriter, r *http.Request, ) {
	consentChallenge := r.URL.Query().Get("consent_challenge")

	consentRequest, err := h.hydraClient.GetConsentRequest(consentChallenge)
	if err != nil {
		fail(w, fmt.Errorf("invalid consent_challenge"), 400)
		return
	}

	user, found, err := h.userRepository.FindByEmail(consentRequest.UserID)

	if err != nil {
		h.rejectConsentAndRedirect(
			w,
			r,
			consentChallenge,
			fmt.Errorf("internal error finding user with given consent UserID %s", consentRequest.UserID),
		)
		return
	}

	if !found {
		err = h.rejectConsentAndRedirect(
			w,
			r,
			consentChallenge,
			fmt.Errorf("cannot find user with given consent UserID %s", consentRequest.UserID),
		)
		if err != nil {
			fail(w, fmt.Errorf("cannot reject consent challenge"), 500)
		}

		return
	}

	err = h.acceptConsentAndRedirect(
		w,
		r,
		*consentRequest,
		consentChallenge,
		*user,
	)
	if err != nil {
		fail(w, fmt.Errorf("cannot accept consent challenge %s", consentChallenge), 500)
	}
}

func (h *HydraConsentHandler) acceptConsentAndRedirect(
	w http.ResponseWriter,
	r *http.Request,
	consentRequest hydra.ConsentRequest,
	consentChallenge string,
	user domain.User,
) error {
	consentAccepted, err := h.hydraClient.AcceptConsentRequest(
		consentChallenge,
		user,
		consentRequest.RequestedScope,
		consentRequest.RequestedAccessTokenAudience,
		true,
	)

	if err != nil {
		return err
	}

	http.Redirect(w, r, consentAccepted.RedirectTo, 301)
	return nil
}

func (h *HydraConsentHandler) rejectConsentAndRedirect(
	w http.ResponseWriter,
	r *http.Request,
	consentChallenge string,
	consentError error,
) error {
	consentRejected, err := h.hydraClient.RejectConsentRequest(consentChallenge, consentError)

	if err != nil {
		return err
	}

	http.Redirect(w, r, consentRejected.RedirectTo, 301)
	return nil
}

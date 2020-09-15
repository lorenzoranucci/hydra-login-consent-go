package hydra

import (
	"net/http"

	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
)

type LoginRequest struct {
	LoginRequestSkipForUser *LoginRequestSkipForUser
	LoginRequestState       *LoginRequestState
}

type LoginRequestSkipForUser struct {
	UserID            string
}

type LoginAccepted struct {
	RedirectTo string
}

type LoginRejected struct {
	RedirectTo string
}

type LoginRequestState struct {
	ALTK *string
	SocialLoginProviderID *string
}

func (h *HydraClientStruct) GetLoginRequest(loginChallenge string) (*LoginRequest, error) {
	getLoginRequestOk, err := (*h.getAdmin()).GetLoginRequest(
		&admin.GetLoginRequestParams{
			LoginChallenge: loginChallenge,
			Context:        h.getDefaultContext(),
		})

	if err != nil {
		return nil, err
	}

	loginRequest := &LoginRequest{
		LoginRequestSkipForUser: getLoginRequestSkipForUser(getLoginRequestOk),
		LoginRequestState:       getLoginRequestState(getLoginRequestOk),
	}

	return loginRequest, nil
}

func (h *HydraClientStruct) AcceptLoginRequest(
	loginChallenge string,
	userID string,
	remember *bool,
	rememberFor *int64,
) (*LoginAccepted, error) {
	body := &models.AcceptLoginRequest{
		Subject: &userID,
	}

	if remember != nil && rememberFor != nil {
		body.Remember = *remember
		body.RememberFor = *rememberFor
	}

	accept, err := (*h.getAdmin()).AcceptLoginRequest(
		&admin.AcceptLoginRequestParams{
			Body:           body,
			LoginChallenge: loginChallenge,
			Context:        h.getDefaultContext(),
		},
	)

	if err != nil {
		return nil, err
	}

	return &LoginAccepted{RedirectTo: *accept.GetPayload().RedirectTo}, nil
}

func (h *HydraClientStruct) RejectLoginRequest(
	loginChallenge string,
	loginError error,
) (*LoginRejected, error) {
	body := &models.RejectRequest{
		Error: loginError.Error(),
	}

	reject, err := (*h.getAdmin()).RejectLoginRequest(
		&admin.RejectLoginRequestParams{
			Body:           body,
			LoginChallenge: loginChallenge,
			Context:        h.getDefaultContext(),
		},
	)

	if err != nil {
		return nil, err
	}

	return &LoginRejected{RedirectTo: *reject.GetPayload().RedirectTo}, nil
}

func (h *HydraClientStruct) AcceptLoginAndRedirect(
	w http.ResponseWriter,
	r *http.Request,
	loginChallenge string,
	userID string,
	remember *bool,
	rememberFor *int64,
) error {
	loginAccepted, err := h.AcceptLoginRequest(
		loginChallenge,
		userID,
		remember,
		rememberFor,
	)

	if err != nil {
		return err
	}

	http.Redirect(w, r, loginAccepted.RedirectTo, 301)
	return nil
}

func (h *HydraClientStruct) RejectLoginAndRedirect(
	w http.ResponseWriter,
	r *http.Request,
	loginChallenge string,
	loginError error,
) error {
	loginRejected, err :=h.RejectLoginRequest(loginChallenge, loginError)

	if err != nil {
		return err
	}

	http.Redirect(w, r, loginRejected.RedirectTo, 301)
	return nil
}

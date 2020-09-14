package hydra

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	httptransport "github.com/go-openapi/runtime/client"
	hydra "github.com/ory/hydra-client-go/client"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"golang.org/x/oauth2"
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

type HydraClientInterface interface {
	GetLoginRequest(loginChallenge string) (*LoginRequest, error)
	AcceptLoginRequest(
		loginChallenge string,
		userID string,
		remember *bool,
		rememberFor *int64,
	) (*LoginAccepted, error)

	RejectLoginRequest(
		loginChallenge string,
		loginError error,
	) (*LoginRejected, error)
}

type HydraClientStruct struct {
	hydraAdminURL *url.URL
}

func NewHydraClientStruct(hydraAdminURL *url.URL) *HydraClientStruct {
	return &HydraClientStruct{hydraAdminURL: hydraAdminURL}
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

func (h *HydraClientStruct) AcceptConsentRequest(
	consentChallenge string,
	userID string,
	remember *bool,
	rememberFor *int64,
) (*LoginAccepted, error) {
	accept, err := (*h.getAdmin()).AcceptConsentRequest(
		&admin.AcceptConsentRequestParams{
			Body: &models.AcceptConsentRequest{
				GrantAccessTokenAudience: consentRequest.Payload.RequestedAccessTokenAudience, // todo investigate if we can open to all audiences
				GrantScope:               []string{"openid", "offline"},                       // todo investigate how to fill this
				Remember:                 true,
				Session: getConsentRequestSession(consentRequest),
			},
			ConsentChallenge: *consentRequest.Payload.Challenge,
			Context:          DefaultContext,
		})
}

func (h *HydraClientStruct) getDefaultContext() context.Context {
	return context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}})
}

func (h *HydraClientStruct) getAdmin() *admin.ClientService {
	skipTlsClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: time.Second * 1000,
	}
	transport := httptransport.NewWithClient(
		h.hydraAdminURL.Host,
		h.hydraAdminURL.Path,
		[]string{h.hydraAdminURL.Scheme},
		skipTlsClient,
	) // todo fix skip tls
	return &hydra.New(transport, nil).Admin
}

func getLoginRequestSkipForUser(getLoginRequestOk *admin.GetLoginRequestOK) *LoginRequestSkipForUser {
	if getLoginRequestOk.GetPayload().Skip != nil &&
		*getLoginRequestOk.GetPayload().Skip &&
		getLoginRequestOk.GetPayload().Subject != nil {
		return &LoginRequestSkipForUser{UserID: *getLoginRequestOk.GetPayload().Subject}
	}

	return nil
}

func getLoginRequestState(getLoginRequestOk *admin.GetLoginRequestOK) *LoginRequestState {
	if getLoginRequestOk.Payload.RequestURL == nil {
		return nil
	}

	requestURL, err := url.Parse(*getLoginRequestOk.Payload.RequestURL)
	if err != nil {
		return nil
	}

	state := &LoginRequestState{}
	stateJson := requestURL.Query().Get("state")
	err = json.Unmarshal([]byte(stateJson), state)
	if err != nil {
		return nil
	}

	return state
}


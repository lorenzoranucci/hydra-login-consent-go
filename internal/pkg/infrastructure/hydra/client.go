package hydra

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/domain"
	hydra "github.com/ory/hydra-client-go/client"
	"github.com/ory/hydra-client-go/client/admin"
	"golang.org/x/oauth2"
)

type HydraClientInterface interface {
	GetLoginRequest(loginChallenge string) (*LoginRequest, error)

	AcceptLoginAndRedirect(
		w http.ResponseWriter,
		r *http.Request,
		loginChallenge string,
		userID string,
		remember *bool,
		rememberFor *int64,
	) error

	RejectLoginAndRedirect(
		w http.ResponseWriter,
		r *http.Request,
		loginChallenge string,
		loginError error,
	) error

	GetConsentRequest(consentChallenge string) (*ConsentRequest, error)

	AcceptConsentRequest(
		consentChallenge string,
		user domain.User,
		grantScopes []string,
		grantAudience []string,
		remember bool,
	) (*ConsentAccepted, error)

	RejectConsentRequest(
		consentChallenge string,
		consentError error,
	) (*ConsentRejected, error)
}

type HydraClientStruct struct {
	hydraAdminURL *url.URL
}

func NewHydraClientStruct(hydraAdminURL *url.URL) *HydraClientStruct {
	return &HydraClientStruct{hydraAdminURL: hydraAdminURL}
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

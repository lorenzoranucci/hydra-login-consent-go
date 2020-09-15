package hydra

import (
	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/domain"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
)

type ConsentRequest struct {
	UserID                       string
	RequestedAccessTokenAudience []string
	RequestedScope               []string
}

type ConsentAccepted struct {
	RedirectTo string
}

type ConsentRejected struct {
	RedirectTo string
}

func (h *HydraClientStruct) GetConsentRequest(consentChallenge string) (*ConsentRequest, error) {
	getConsentRequestOk, err := (*h.getAdmin()).GetConsentRequest(
		&admin.GetConsentRequestParams{
			ConsentChallenge: consentChallenge,
			Context:        h.getDefaultContext(),
		})

	if err != nil {
		return nil, err
	}

	consentRequest := &ConsentRequest{
		UserID:                       getConsentRequestOk.GetPayload().Subject,
		RequestedAccessTokenAudience: getConsentRequestOk.GetPayload().RequestedAccessTokenAudience,
		RequestedScope:               getConsentRequestOk.GetPayload().RequestedScope,
	}

	return consentRequest, nil
}

func (h *HydraClientStruct) AcceptConsentRequest(
	consentChallenge string,
	user domain.User,
	grantScopes []string,
	grantAudience []string,
	remember bool,
) (*ConsentAccepted, error) {
	accept, err := (*h.getAdmin()).AcceptConsentRequest(
		&admin.AcceptConsentRequestParams{
			Body: &models.AcceptConsentRequest{
				GrantAccessTokenAudience: grantAudience,
				GrantScope:               grantScopes,
				Remember:                 remember,
				Session:                  &models.ConsentRequestSession{
					IDToken:     getIDTokenFromUser(user),
					AccessToken: getAccessTokenFromUser(user),
				},
			},
			ConsentChallenge: consentChallenge,
			Context:          h.getDefaultContext(),
		})

	if err != nil {
		return nil, err
	}

	return &ConsentAccepted{
		RedirectTo: *accept.GetPayload().RedirectTo,
	}, nil
}

func (h *HydraClientStruct) RejectConsentRequest(
	consentChallenge string,
	consentError error,
) (*ConsentRejected, error) {
	body := &models.RejectRequest{
		Error: consentError.Error(),
	}

	reject, err := (*h.getAdmin()).RejectLoginRequest(
		&admin.RejectLoginRequestParams{
			Body:           body,
			LoginChallenge: consentChallenge,
			Context:        h.getDefaultContext(),
		},
	)

	if err != nil {
		return nil, err
	}

	return &ConsentRejected{RedirectTo: *reject.GetPayload().RedirectTo}, nil
}

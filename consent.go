package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	hydra "github.com/ory/hydra-client-go/client"
	admin2 "github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
)

var consentGetTemplate = template.Must(template.New("").Parse(`<html>
<head>
    <title></title>
</head>
<body>
<h1>An application requests access to your data!</h1>
<form action="" method="POST">
    <input type="hidden" name="challenge" value="{{ .ConsentChallenge }}"><input
        type="hidden" name="_csrf" value="{{ .CsrfToken }}">
    <p>Hi {{ .Subject }}, application <strong>{{ .ClientName }}</strong> wants access resources on your behalf and to:
    </p>

    <input type="checkbox" id="openid" value="openid" name="grant_scope">
    <label for="openid">openid</label><br>
    
    <input type="checkbox" id="offline" value="offline" name="grant_scope">
    <label for="offline">offline</label><br>
    
    <p>Do you want to be asked next time when this application wants to access your data? The application will
        not be able to ask for more permissions without your consent.</p>
    <ul></ul>
    <p>
        <input type="checkbox" id="remember" name="remember" value="1">
        <label for="remember">Do not ask me again</label>
    </p>
    <p>
        <input type="submit" id="accept" name="submit" value="Allow access">
        <input type="submit" id="reject" name="submit" value="Deny access">
    </p>
</form>
</body>
</html>`))

func handleConsentGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	consentChallenge := r.URL.Query().Get("consent_challenge")
	if consentChallenge == "" {
		handleError(w, fmt.Errorf("Expected a consent challenge to be set but received none."))
	}

	admin := getHydraAdmin()
	consentRequest, err := admin.Admin.GetConsentRequest(
		&admin2.GetConsentRequestParams{
			ConsentChallenge: consentChallenge,
			Context:          DefaultContext,
		},
	)

	if err != nil {
		handleError(w, err)
		return
	}

	if consentRequest.GetPayload().Skip {
		acceptSkippedConsent(w, r, admin, consentRequest)
		return
	}

	if consentRequest.GetPayload().Client.ClientID == "client-frontend.localhost" {
		acceptInternalClients(w, r, admin, consentRequest)
		return
	}

	_ = consentGetTemplate.Execute(w, struct {
		CsrfToken        string
		ConsentChallenge string
		Subject          string
		ClientName       string
	}{
		CsrfToken:        "change me", // todo investigate how to change this
		ConsentChallenge: consentChallenge,
		Subject:          consentRequest.Payload.Subject,
		ClientName: fmt.Sprintf(
			"%s/%s",
			consentRequest.Payload.Client.Owner,
			consentRequest.Payload.Client.ClientName,
		),
	})
}

func handleConsentPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		handleError(w, err)
		return
	}

	submit := r.Form.Get("submit")
	grantScope := r.Form["grant_scope"]
	remember, _ := strconv.ParseBool(r.Form.Get("remember"))

	admin := getHydraAdmin()
	if submit == "Deny access" {
		reject, err := admin.Admin.RejectConsentRequest(
			&admin2.RejectConsentRequestParams{
				Body: &models.RejectRequest{
					Error:            "access_denied",
					ErrorDescription: "The resource owner denied the request",
				},
				ConsentChallenge: r.Form.Get("consentChallenge"),
				Context:          DefaultContext,
			})

		if err != nil {
			handleError(w, err)
			return
		}

		http.Redirect(w, r, *reject.GetPayload().RedirectTo, 301) // todo check code
		return
	}

	consentRequest, err := admin.Admin.GetConsentRequest(
		&admin2.GetConsentRequestParams{
			ConsentChallenge: r.Form.Get("challenge"),
			Context:          DefaultContext,
		},
	)
	if err != nil {
		handleError(w, err)
		return
	}

	acceptExplicitConsent(w, r, admin, consentRequest, grantScope, remember)
}

func acceptSkippedConsent(w http.ResponseWriter, r *http.Request, admin *hydra.OryHydra, consentRequest *admin2.GetConsentRequestOK) {
	accept, err := admin.Admin.AcceptConsentRequest(
		&admin2.AcceptConsentRequestParams{
			Body: &models.AcceptConsentRequest{
				GrantAccessTokenAudience: consentRequest.Payload.RequestedAccessTokenAudience, // todo investigate if we can open to all audiences
				GrantScope:               consentRequest.Payload.RequestedScope,               // todo investigate how to fill this
				Session: getConsentRequestSession(consentRequest),                                                                             // todo investigate how to fill this
			},
			ConsentChallenge: *consentRequest.Payload.Challenge,
			Context:          DefaultContext,
		})

	if err != nil {
		handleError(w, err)
		return
	}

	http.Redirect(w, r, *accept.GetPayload().RedirectTo, 301) // todo check code
}

func acceptInternalClients(
	w http.ResponseWriter,
	r *http.Request,
	admin *hydra.OryHydra,
	consentRequest *admin2.GetConsentRequestOK,
) {
	accept, err := admin.Admin.AcceptConsentRequest(
		&admin2.AcceptConsentRequestParams{
			Body: &models.AcceptConsentRequest{
				GrantAccessTokenAudience: consentRequest.Payload.RequestedAccessTokenAudience, // todo investigate if we can open to all audiences
				GrantScope:               []string{"openid", "offline"},                       // todo investigate how to fill this
				Remember:                 true,
				Session: getConsentRequestSession(consentRequest),
			},
			ConsentChallenge: *consentRequest.Payload.Challenge,
			Context:          DefaultContext,
		})

	if err != nil {
		handleError(w, err)
		return
	}

	http.Redirect(w, r, *accept.GetPayload().RedirectTo, 301) // todo check code
	return
}

func acceptExplicitConsent(
	w http.ResponseWriter,
	r *http.Request,
	admin *hydra.OryHydra,
	consentRequest *admin2.GetConsentRequestOK,
	scopes []string,
	remember bool,
) {
	accept, err := admin.Admin.AcceptConsentRequest(
		&admin2.AcceptConsentRequestParams{
			Body: &models.AcceptConsentRequest{
				GrantAccessTokenAudience: consentRequest.Payload.RequestedAccessTokenAudience, // todo investigate if we can open to all audiences
				GrantScope:               scopes,                                              // todo investigate how to fill this
				Remember:                 remember,
				Session: getConsentRequestSession(consentRequest),
			},
			ConsentChallenge: *consentRequest.Payload.Challenge,
			Context:          DefaultContext,
		})

	if err != nil {
		handleError(w, err)
		return
	}

	http.Redirect(w, r, *accept.GetPayload().RedirectTo, 301) // todo check code
	return
}

func getConsentRequestSession(consentRequest *admin2.GetConsentRequestOK) *models.ConsentRequestSession {
	return &models.ConsentRequestSession{
		IDToken:     getIDTokenFromSubject(consentRequest.Payload.Subject),
		AccessToken: getAccessTokenFromSubject(consentRequest.Payload.Subject),
	}
}

func getIDTokenFromSubject(subject string) IDToken {
	//todo findUserFromEmail
	return IDToken{
		Name:  "Foo Bar",
		Email: subject,
	}
}

func getAccessTokenFromSubject(subject string) AccessToken {
	//todo findUserFromEmail
	return AccessToken{
		Mercure:  Mercure{
			Subscribe: []string{"users/8f068dce-c664-11e8-a21b-02286126d5c8"},
		},
		User: User{
			Email: subject,
			Id:    1,
			Roles: []string{
				"admin",
			},
			UUID: "8f068dce-c664-11e8-a21b-02286126d5c8",
		},
	}
}

type IDToken struct {
	Name  string
	Email string
}

type AccessToken struct {
	Mercure Mercure
	User User
}

type Mercure struct {
	Subscribe []string
}

type User struct {
	Email string
	Id int
	Roles []string
	UUID string
}

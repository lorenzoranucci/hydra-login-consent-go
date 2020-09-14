package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/url"
	"strconv"

	"github.com/julienschmidt/httprouter"
	admin2 "github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
)

type State struct {
	RedirectURL string
	Altk        string
}

func handleLoginGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	loginChallenge := r.URL.Query().Get("login_challenge")
	admin := getHydraAdmin()

	getLoginRequestOk, err := admin.Admin.GetLoginRequest(
		&admin2.GetLoginRequestParams{
			LoginChallenge: loginChallenge,
			Context:        DefaultContext,
		})

	if err != nil {
		handleError(w, err)
		return
	}

	requestURL, err := url.Parse(*getLoginRequestOk.Payload.RequestURL)
	if requestURL != nil {
		state := &State{}
		stateJson := requestURL.Query().Get("state")
		err = json.Unmarshal([]byte(stateJson), state)
		if state.Altk != "" {
			//todo auth user by altk
			userIDFromAltk := "foo_altk@ppro.it"

			accept, err := admin.Admin.AcceptLoginRequest(
				&admin2.AcceptLoginRequestParams{
					Body: &models.AcceptLoginRequest{
						Subject: &userIDFromAltk,
					},
					LoginChallenge: loginChallenge,
					Context:        DefaultContext,
				},
			)

			if err != nil {
				handleError(w, err)
				return
			}

			http.Redirect(w, r, *accept.GetPayload().RedirectTo, 301) // todo check code
			return
		}
	}

	if err != nil {
		handleError(w, err)
		return
	}

	if *getLoginRequestOk.Payload.Skip {
		accept, err := admin.Admin.AcceptLoginRequest(
			&admin2.AcceptLoginRequestParams{
				Body: &models.AcceptLoginRequest{
					Subject: getLoginRequestOk.Payload.Subject,
				},
				LoginChallenge: loginChallenge,
				Context:        DefaultContext,
			},
		)

		if err != nil {
			handleError(w, err)
			return
		}

		http.Redirect(w, r, *accept.GetPayload().RedirectTo, 301) // todo check code
		return
	}

	renderLoginGetTemplate(w, loginChallenge, "")
}

func handleLoginPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		handleError(w, err)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")
	remember, _ := strconv.ParseBool(r.Form.Get("remember"))

	loginChallenge := r.Form.Get("challenge")
	if email == "" && password == "" {
		loginError := "The username / password combination is not correct"
		renderLoginGetTemplate(w, loginChallenge, loginError)

		return
	}

	admin := getHydraAdmin()

	accept, err := admin.Admin.AcceptLoginRequest(
		&admin2.AcceptLoginRequestParams{
			Body: &models.AcceptLoginRequest{
				Subject:     &email,
				Remember:    remember,
				RememberFor: 3600,
			},
			LoginChallenge: loginChallenge,
			Context:        DefaultContext,
		},
	)

	if err != nil {
		handleError(w, err)
		return
	}

	http.Redirect(w, r, *accept.GetPayload().RedirectTo, 301) // todo check code
	return
}

var loginGetTemplate = template.Must(template.New("").Parse(`<html>
<head></head>
<body>
<h1>Login</h1>
<h2>{{ .Error }}</h2>
<form action="/login" method="POST">
    <input type="hidden" name="_csrf" value="{{ .CsrfToken }}">
    <input type="hidden" name="challenge" value="{{ .LoginChallenge }}">
    <table style="">
        <tbody>
        <tr>
            <td><input type="email" id="email" name="email" placeholder="email@foobar.com"></td>
            <td>(it's "foo@bar.com")</td>
        </tr>
        <tr>
            <td><input type="password" id="password" name="password"></td>
            <td>(it's "foobar")</td>
        </tr>
        </tbody>
    </table>
    <input type="checkbox" id="remember" name="remember" value="1"><label for="remember">Remember me</label><br><input
        type="submit" id="accept" value="Log in"></form>

	<br/>
	<a href="{{ .FacebookLoginEndpoint }}?client_id={{ .FacebookClientID }}&redirect_uri={{ .FacebookRedirectURI }}&state={{ .LoginChallenge}}">Login with Facebook</a>

</body>
</html>`))
func renderLoginGetTemplate(w http.ResponseWriter, loginChallenge string, loginError string) {
	_ = loginGetTemplate.Execute(w, struct {
		CsrfToken           string
		LoginChallenge      string
		Error               string
		FacebookClientID    int
		FacebookRedirectURI string
		FacebookLoginEndpoint string
	}{
		CsrfToken:           "change me",
		LoginChallenge:      loginChallenge,
		Error:               loginError,
		FacebookClientID:    FacebookClientID,
		FacebookRedirectURI: FacebookRedirectURI,
		FacebookLoginEndpoint: FacebookLoginEndpoint,
	})
}


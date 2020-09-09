package main

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	admin2 "github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
)

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
</body>
</html>`))

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

	_ = loginGetTemplate.Execute(w, struct {
		CsrfToken      string
		LoginChallenge string
		Error string
	}{
		CsrfToken:      "change me",
		LoginChallenge: loginChallenge,
		Error: "",
	})
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

	if email == "" && password == "" {
		_ = loginGetTemplate.Execute(w, struct {
			CsrfToken      string
			LoginChallenge string
			Error string
		}{
			CsrfToken:      "change me",
			LoginChallenge: r.Form.Get("challenge"),
			Error: "The username / password combination is not correct",
		})

		return
	}

	admin := getHydraAdmin()

	accept, err := admin.Admin.AcceptLoginRequest(
		&admin2.AcceptLoginRequestParams{
			Body: &models.AcceptLoginRequest{
				Subject: &email,
				Remember: remember,
				RememberFor: 3600,
			},
			LoginChallenge: r.Form.Get("challenge"),
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

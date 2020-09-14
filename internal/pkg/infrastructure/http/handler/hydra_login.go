package handler

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/application"
	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/infrastructure/hydra"
	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/infrastructure/social_login_provider"
)

type HydraLoginHandler struct {
	signInUserWithAutoLoginTokenService application.SignInWithAutoLoginTokenServiceInterface
	signInWithEmailAndPasswordService   application.SignInUserWithEmailAndPasswordServiceInterface

	hydraClient                hydra.HydraClientInterface
	socialLoginProviderFactory social_login_provider.FactoryInterface

	FacebookSocialLoginProviderID string
	GoogleSocialLoginProviderID string
}

func NewHydraLoginHandler(
	signInUserWithAutoLoginTokenService application.SignInWithAutoLoginTokenServiceInterface,
	signInWithEmailAndPasswordService application.SignInUserWithEmailAndPasswordServiceInterface,
	hydraClient hydra.HydraClientInterface,
	socialLoginProviderFactory social_login_provider.FactoryInterface,
	facebookSocialLoginProviderID string,
	googleSocialLoginProviderID string,
) *HydraLoginHandler {
	return &HydraLoginHandler{
		signInUserWithAutoLoginTokenService: signInUserWithAutoLoginTokenService,
		signInWithEmailAndPasswordService: signInWithEmailAndPasswordService,
		hydraClient: hydraClient,
		socialLoginProviderFactory: socialLoginProviderFactory,
		FacebookSocialLoginProviderID: facebookSocialLoginProviderID,
		GoogleSocialLoginProviderID: googleSocialLoginProviderID,
	}
}

func (h *HydraLoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.handleLoginGet(w, r)
	}

	if r.Method == "POST" {
		h.handleLoginPost(w, r)
	}
}

func (h *HydraLoginHandler) handleLoginGet(w http.ResponseWriter, r *http.Request, ) {
	loginChallenge := r.URL.Query().Get("login_challenge")

	loginRequest, err := h.hydraClient.GetLoginRequest(loginChallenge)
	if err != nil {
		fail(w, fmt.Errorf("invalid login_challenge"), 400)
		return
	}

	if loginRequest.LoginRequestSkipForUser != nil {
		h.acceptLoginAndRedirect(
			w,
			r,
			loginChallenge,
			loginRequest.LoginRequestSkipForUser.UserID,
			nil,
			nil,
		)
		return
	}

	if loginRequest.LoginRequestState != nil && loginRequest.LoginRequestState.ALTK != nil {
		user, found, err := h.signInUserWithAutoLoginTokenService.Execute(
			application.SignInUserWithAutoLoginTokenRequest{AutoLoginToken: *loginRequest.LoginRequestState.ALTK},
		)
		if err != nil {
			h.rejectLoginAndRedirect(
				w,
				r,
				loginChallenge,
				fmt.Errorf("internal error finding user with given auto login token"),
			)
			return
		}

		if !found {
			h.rejectLoginAndRedirect(
				w,
				r,
				loginChallenge,
				fmt.Errorf("cannot find user with given auto login token"),
			)
			return
		}

		remember := true
		var rememberFor int64 = 0
		h.acceptLoginAndRedirect(w, r, loginChallenge, user.UUID.String(), &remember, &rememberFor)
		return
	}

	if loginRequest.LoginRequestState != nil && loginRequest.LoginRequestState.SocialLoginProviderID != nil {
		socialLoginProvider, err := h.socialLoginProviderFactory.GetSocialLoginProviderByID(
			*loginRequest.LoginRequestState.SocialLoginProviderID,
		)

		if err != nil {
			h.rejectLoginAndRedirect(
				w,
				r,
				loginChallenge,
				fmt.Errorf("invalid social login provider"),
			)
			return
		}

		endpoint, err := socialLoginProvider.GetLoginURL(loginChallenge)
		if err != nil {
			h.rejectLoginAndRedirect(
				w,
				r,
				loginChallenge,
				fmt.Errorf(
					"invalid social login endpoint for `%s` and login challenge '%s'",
					socialLoginProvider.GetID(),
					loginChallenge,
				),
			)
			return
		}
		http.Redirect(w, r, endpoint.String(), 301)
	}

	h.renderLoginGetTemplate(w, loginChallenge, "")
}

func (h *HydraLoginHandler) handleLoginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fail(w, err, 400)
		return
	}

	loginChallenge := r.Form.Get("login_challenge")

	loginRequest, err := h.hydraClient.GetLoginRequest(loginChallenge)
	if err != nil {
		fail(w, fmt.Errorf("invalid login_challenge"), 400)
		return
	}

	if loginRequest.LoginRequestSkipForUser != nil {
		h.acceptLoginAndRedirect(
			w,
			r,
			loginChallenge,
			loginRequest.LoginRequestSkipForUser.UserID,
			nil,
			nil,
		)
		return
	}

	user, found, err := h.signInWithEmailAndPasswordService.Execute(
		application.SignInUserWithEmailAndPasswordRequest{
			Email:    r.Form.Get("email"),
			Password: r.Form.Get("password"),
		},
	)
	if err != nil {
		h.renderLoginGetTemplate(
			w,
			loginChallenge,
			"internal error finding user with given email and password",
		)
		return
	}

	if !found {
		h.renderLoginGetTemplate(
			w,
			loginChallenge,
			"cannot find user with given email and password",
		)
		return
	}

	remember := true
	var rememberFor int64 = 0
	h.acceptLoginAndRedirect(w, r, loginChallenge, user.UUID.String(), &remember, &rememberFor)
}

func (h *HydraLoginHandler) acceptLoginAndRedirect(
	w http.ResponseWriter,
	r *http.Request,
	loginChallenge string,
	userID string,
	remember *bool,
	rememberFor *int64,
) {
	loginAccepted, err := h.hydraClient.AcceptLoginRequest(
		loginChallenge,
		userID,
		remember,
		rememberFor,
	)

	if err != nil {
		fail(w, fmt.Errorf("cannot accept login challenge"), 500)
		return
	}

	http.Redirect(w, r, loginAccepted.RedirectTo, 301)
	return
}

func (h *HydraLoginHandler) rejectLoginAndRedirect(
	w http.ResponseWriter,
	r *http.Request,
	loginChallenge string,
	loginError error,
) {
	loginRejected, err :=h.hydraClient.RejectLoginRequest(loginChallenge, loginError)

	if err != nil {
		fail(w, fmt.Errorf("cannot reject login challenge"), 500)
		return
	}

	http.Redirect(w, r, loginRejected.RedirectTo, 301)
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
	<a href="{{ .FacebookLoginURL }}">Login with Facebook</a>
	<a href="{{ .GoogleLoginURL }}">Login with Google</a>

</body>
</html>`))


func (h *HydraLoginHandler) renderLoginGetTemplate(
	w http.ResponseWriter,
	loginChallenge string,
	loginError string,
) {
	facebook, err := h.socialLoginProviderFactory.GetSocialLoginProviderByID(h.FacebookSocialLoginProviderID)
	if err != nil {
		fail(w, fmt.Errorf("cannot reject login challenge"), 500)
		return
	}
	google, err := h.socialLoginProviderFactory.GetSocialLoginProviderByID(h.GoogleSocialLoginProviderID)
	if err != nil {
		fail(w, fmt.Errorf("cannot reject login challenge"), 500)
		return
	}

	facebookLoginURL, err := facebook.GetLoginURL(loginChallenge)
	if err != nil {
		fail(w, fmt.Errorf("cannot reject login challenge"), 500)
		return
	}

	googleLoginURL, err := google.GetLoginURL(loginChallenge)
	if err != nil {
		fail(w, fmt.Errorf("cannot reject login challenge"), 500)
		return
	}

	_ = loginGetTemplate.Execute(w, struct {
		CsrfToken        string
		LoginChallenge   string
		Error            string
		FacebookLoginURL string
		GoogleLoginURL   string
	}{
		CsrfToken:        "change me",
		LoginChallenge:   loginChallenge,
		Error:            loginError,
		FacebookLoginURL: facebookLoginURL.String(),
		GoogleLoginURL:   googleLoginURL.String(),
	})
}

package http

import (
	"fmt"
	"net/http"

	"github.com/lorenzoranucci/hydra-login-consent-go/internal/app"
	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/infrastructure/http/handler"
)

type Server struct {
	m    *http.ServeMux
	port int
}

func NewServer(
	port int,
) *Server {
	srv := &Server{
		m:    http.NewServeMux(),
		port: port,
	}

	srv.m.Handle("/favicon.ico", http.NotFoundHandler())

	serviceLocator := app.GetServiceLocator()

	srv.m.Handle(
		"/login",
		handler.NewHydraLoginHandler(
			serviceLocator.CreateSignInUserWithAutoLoginTokenService(),
			serviceLocator.CreateSignInUserWithEmailAndPasswordService(),
			serviceLocator.CreateHydraClient(),
			serviceLocator.CreateSocialLoginProviderFactory(),
			serviceLocator.FacebookID(),
			serviceLocator.GoogleID(),
		),
	)
	srv.m.Handle(
		"/consent",
		handler.NewHydraConsentHandler(
			serviceLocator.UserRepository(),
			serviceLocator.CreateHydraClient(),
		),
	)

	srv.m.Handle(
		"/login/social/facebook",
		handler.NewFacebookSocialLoginCallbackHandler(
			serviceLocator.CreateSignInUserWithSocialLoginService(),
			serviceLocator.CreateFacebookSocialLoginProvider(),
			serviceLocator.CreateHydraClient(),
		),
	)

	srv.m.Handle(
		"/login/social/google",
		handler.NewGoogleSocialLoginCallbackHandler(
			serviceLocator.CreateSignInUserWithSocialLoginService(),
			serviceLocator.CreateGoogleSocialLoginProvider(),
			serviceLocator.CreateHydraClient(),
		),
	)

	return srv
}

func (s *Server) Run() {
	err := http.ListenAndServe(
		fmt.Sprintf(":%d", s.port),
		s.m,
	)

	fmt.Print(err)
	fmt.Println("Server terminated")
}

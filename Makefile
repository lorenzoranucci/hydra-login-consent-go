APP_VERSION ?= dev

_HYDRA_ADMIN_URL=https://ory-hydra-example--hydra:9001
_FACEBOOK_ID=facebook
_FACEBOOK_CLIENT_ID=236838367690513
_FACEBOOK_CLIENT_SECRET=b6f57361a9e222c0a3697918150dcb93
_FACEBOOK_REDIRECT_URI=http://localhost:9020/login/social/facebook
_FACEBOOK_AUTH_ENDPOINT=https://www.facebook.com/v8.0/dialog/oauth
_FACEBOOK_TOKEN_ENDPOINT=https://graph.facebook.com/v8.0/oauth/access_token
_FACEBOOK_VERIFY_TOKEN_ENDPOINT=https://graph.facebook.com/v8.0/me
_GOOGLE_ID=google
_GOOGLE_CLIENT_ID=1043919666765-opdjnqlgg0jl9rcnqnna8bsgd4nrhbem.apps.googleusercontent.com
_GOOGLE_CLIENT_SECRET=4Ns5zb8HKbbnaFbYqVttxBvK
_GOOGLE_REDIRECT_URI=http://localhost:9020/login/social/google
_GOOGLE_AUTH_ENDPOINT=https://accounts.google.com/o/oauth2/v2/auth
_GOOGLE_TOKEN_ENDPOINT=https://oauth2.googleapis.com/token
_GOOGLE_VERIFY_TOKEN_ENDPOINT=https://people.googleapis.com/v1/people/me?personFields=names,emailAddresses

.PHONY: base-env env env-test

base-env:
	@echo 'export HYDRA_ADMIN_URL="${_HYDRA_ADMIN_URL}"'
	@echo 'export FACEBOOK_ID="${_FACEBOOK_ID}"'
	@echo 'export FACEBOOK_CLIENT_ID="${_FACEBOOK_CLIENT_ID}"'
	@echo 'export FACEBOOK_CLIENT_SECRET="${_FACEBOOK_CLIENT_SECRET}"'
	@echo 'export FACEBOOK_REDIRECT_URI="${_FACEBOOK_REDIRECT_URI}"'
	@echo 'export FACEBOOK_AUTH_ENDPOINT="${_FACEBOOK_AUTH_ENDPOINT}"'
	@echo 'export FACEBOOK_TOKEN_ENDPOINT="${_FACEBOOK_TOKEN_ENDPOINT}"'
	@echo 'export FACEBOOK_VERIFY_TOKEN_ENDPOINT="${_FACEBOOK_VERIFY_TOKEN_ENDPOINT}"'
	@echo 'export GOOGLE_ID="${_GOOGLE_ID}"'
	@echo 'export GOOGLE_CLIENT_ID="${_GOOGLE_CLIENT_ID}"'
	@echo 'export GOOGLE_CLIENT_SECRET="${_GOOGLE_CLIENT_SECRET}"'
	@echo 'export GOOGLE_REDIRECT_URI="${_GOOGLE_REDIRECT_URI}"'
	@echo 'export GOOGLE_AUTH_ENDPOINT="${_GOOGLE_AUTH_ENDPOINT}"'
	@echo 'export GOOGLE_TOKEN_ENDPOINT="${_GOOGLE_TOKEN_ENDPOINT}"'
	@echo 'export GOOGLE_VERIFY_TOKEN_ENDPOINT="${_GOOGLE_VERIFY_TOKEN_ENDPOINT}"'

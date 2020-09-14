package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	admin2 "github.com/ory/hydra-client-go/client/admin"
)

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func handleSocialLoginCallbackGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	facebookAccessCode := r.URL.Query().Get("code")
	loginChallenge := r.URL.Query().Get("state")

	admin := getHydraAdmin()

	_, err := admin.Admin.GetLoginRequest(
		&admin2.GetLoginRequestParams{
			LoginChallenge: loginChallenge,
			Context:        DefaultContext,
		})

	if err != nil {
		handleError(w, err)
		return
	}

	accessTokenResponse, err := http.Get(
		fmt.Sprintf(
			FacebookAccessTokenEndpoint,
			FacebookClientID,
			FacebookClientSecret,
			FacebookRedirectURI,
			facebookAccessCode,
		),
	)
	if err != nil {
		handleError(w, err)
		return
	}

	scanner := bufio.NewScanner(accessTokenResponse.Body)
	scanner.Scan()
	accessTokenResponseString :=scanner.Bytes()

	accessTokenResponseStruct := &AccessTokenResponse{}
	err = json.Unmarshal(accessTokenResponseString, accessTokenResponseStruct)
	if err != nil {
		handleError(w, err)
		return
	}

	fmt.Print(accessTokenResponseStruct.AccessToken)
}




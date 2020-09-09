package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"runtime/debug"
	"time"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/julienschmidt/httprouter"
	hydra "github.com/ory/hydra-client-go/client"
	"golang.org/x/oauth2"
)

const Port = 9020
const HydraHost = "ory-hydra-example--hydra"
const HydraAdminPort = "9001"
const HydraPublicPort = "9000"
var HydraAdminURL = fmt.Sprintf("https://%s:%s", HydraHost, HydraAdminPort)
var HydraPublicURL = fmt.Sprintf("https://%s:%s", HydraHost, HydraPublicPort)
var DefaultContext = context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{Transport: &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}})

func main()  {
	r := httprouter.New()
	server := &http.Server{Addr: fmt.Sprintf(":%d", Port), Handler: r}
	r.GET("/login", handleLoginGet)
	r.POST("/login", handleLoginPost)

	r.GET("/consent", handleConsentGet)
	r.POST("/consent", handleConsentPost)

	err := server.ListenAndServe()
	panic(err)
}

func getHydraAdmin() *hydra.OryHydra {
	adminURL, _ := url.Parse(HydraAdminURL)
	skipTlsClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: time.Second * 1000,
	}
	transport := httptransport.NewWithClient(adminURL.Host, adminURL.Path, []string{adminURL.Scheme}, skipTlsClient) // todo fix skip tls
	admin := hydra.New(transport, nil)
	return admin
}


var errorTemplate = template.Must(template.New("").Parse(`<html>
<head></head>
<body>
<h1>Error</h1>
<h2>{{ .Error }}</h2>
</body>
</html>`))

func handleError(w http.ResponseWriter, err error) {
	debug.PrintStack()
	_ = errorTemplate.Execute(w, struct {
		Error string
	}{
		Error: err.Error(),
	})
	return
}

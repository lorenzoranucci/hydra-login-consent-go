package handler

import (
	"net/http"
)

type HydraConsentHandler struct {}

func NewHydraConsentHandler() *HydraConsentHandler {
	return &HydraConsentHandler{}
}

func (h *HydraConsentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.handleConsentGet(w, r)
	}

	/* todo implement for third party clients consent handling
	if r.Method == "POST" {
		h.handleLoginPost(w, r)
	}*/
}

func (h *HydraConsentHandler) handleConsentGet(w http.ResponseWriter, r *http.Request, ) {

}

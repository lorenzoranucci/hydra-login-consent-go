package handler

import (
	"encoding/json"
	"net/http"
)

func send(w http.ResponseWriter, data interface{}) {
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func fail(w http.ResponseWriter, e error, statusCode int) {
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(e); err != nil {
		//log error
	}
}

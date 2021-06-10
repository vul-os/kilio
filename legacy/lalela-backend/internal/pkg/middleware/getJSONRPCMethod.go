package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	lalelaException "lalela-backend/internal/pkg/exception"
)

func getMethod(r *http.Request) (string, error) {
	// Confirm that body of request has data
	if r.Body == nil {
		log.Error().Msg("http request body is nil")
		return "", lalelaException.ErrUnexpected{}
	}

	// Extract body of http Request
	var bodyBytes []byte
	bodyBytes, _ = ioutil.ReadAll(r.Body)

	// Reset body of request
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	// Retrieve id and method of json rpc request
	var req struct {
		// To unmarshal the received json
		Id     int    `json:"id"`
		Method string `json:"method"`
	}
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		log.Error().Msg("unable to unmarshall json rpc request")
		return "", lalelaException.ErrUnexpected{}
	}

	return req.Method, nil
}

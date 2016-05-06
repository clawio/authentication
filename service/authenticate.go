package service

import (
	"encoding/json"
	"net/http"

	"github.com/clawio/codes"
)

type (
	// AuthenticateRequest specifies the data received by the Authenticate endpoint.
	AuthenticateRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// AuthenticateResponse specifies the data returned from the Authenticate endpoint.
	AuthenticateResponse struct {
		AccessToken string `json:"access_token"`
	}
)

// Authenticate authenticates an user using an username and a password.
func (s *Service) Authenticate(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	authReq := &AuthenticateRequest{}
	if err := json.NewDecoder(r.Body).Decode(authReq); err != nil {
		e := codes.NewErr(codes.BadInputData, "")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(e)
		return
	}
	token, err := s.AuthenticationController.Authenticate(authReq.Username, authReq.Password)
	if err != nil {
		s.handleAuthenticateError(err, w)
		return
	}
	res := &AuthenticateResponse{AccessToken: token}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (s *Service) handleAuthenticateError(err error, w http.ResponseWriter) {
	e := codes.NewErr(codes.BadInputData, "user or password do not match")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(e)
	return
}

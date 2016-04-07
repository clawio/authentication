package service

import (
	"encoding/json"
	"net/http"

	"github.com/clawio/codes"
)

type (
	AuthenticateRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	AuthenticateResponse struct {
		Token string `json:"token"`
	}
)

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
	res := &AuthenticateResponse{Token: token}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (s *Service) handleAuthenticateError(err error, w http.ResponseWriter) {
	e := codes.NewErr(codes.BadInputData, "user or password do not match")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(e)
	return
}

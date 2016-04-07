package service

import (
	"encoding/json"
	"net/http"

	"github.com/clawio/codes"
	"github.com/gorilla/mux"
)

func (s *Service) Verify(w http.ResponseWriter, r *http.Request) {
	user, err := s.AuthenticationController.Verify(mux.Vars(r)["token"])
	if err != nil {
		s.handleVerifyError(err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
func (s *Service) handleVerifyError(err error, w http.ResponseWriter) {
	e := codes.NewErr(codes.InvalidToken, "")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(e)
	return
}

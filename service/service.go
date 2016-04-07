package service

import (
	"errors"
	"net/http"

	"github.com/NYTimes/gizmo/config"
	"github.com/clawio/authentication/authenticationcontroller"
	"github.com/prometheus/client_golang/prometheus"
)

type (

	// Service will implement server.Service and
	// handle all requests to the server.
	Service struct {
		Config                   *Config
		AuthenticationController authenticationcontroller.AuthenticationController
	}

	// Config is a struct to contain all the needed
	// configuration for our Service
	Config struct {
		Server                   *config.Server
		AuthenticationController *AuthenticationControllerConfig
	}

	// AuthenticationControllerConfig holds the configuration for
	// an AuthenticationController.
	AuthenticationControllerConfig struct {
		Type                   string
		SimpleDriver           string
		SimpleDSN              string
		SimpleJWTKey           string
		SimpleJWTSigningMethod string
	}
)

// New will instantiate and return
// a new Service that implements server.Service.
func New(cfg *Config) (*Service, error) {
	if cfg == nil {
		return nil, errors.New("config is nil")
	}
	if cfg.AuthenticationController == nil {
		return nil, errors.New("config.AuthenticationController is nil")
	}
	opts := &authenticationcontroller.SimpleAuthenticationControllerOptions{
		Driver:           cfg.AuthenticationController.SimpleDriver,
		DSN:              cfg.AuthenticationController.SimpleDSN,
		JWTKey:           cfg.AuthenticationController.SimpleJWTKey,
		JWTSigningMethod: cfg.AuthenticationController.SimpleJWTSigningMethod,
	}
	authenticationController, err := authenticationcontroller.NewSimpleAuthenticationController(opts)
	if err != nil {
		return nil, err
	}
	return &Service{
		Config: cfg,
		AuthenticationController: authenticationController,
	}, nil
}

// Prefix returns the string prefix used for all endpoints within
// this service.
func (s *Service) Prefix() string {
	return "/clawio/v1/auth"
}

// Middleware provides an http.Handler hook wrapped around all requests.
// In this implementation, we're using a GzipHandler middleware to
// compress our responses.
func (s *Service) Middleware(h http.Handler) http.Handler {
	return h
}

// Endpoints is a listing of all endpoints available in the MixedService.
func (s *Service) Endpoints() map[string]map[string]http.HandlerFunc {
	return map[string]map[string]http.HandlerFunc{
		"/metrics": map[string]http.HandlerFunc{
			"GET": func(w http.ResponseWriter, r *http.Request) {
				prometheus.Handler().ServeHTTP(w, r)
			},
		},
		"/authenticate": map[string]http.HandlerFunc{
			"POST": prometheus.InstrumentHandlerFunc("/authenticate", s.Authenticate),
		},
		"/verify/{token}": map[string]http.HandlerFunc{
			"GET": prometheus.InstrumentHandlerFunc("/verify", s.Verify),
		},
	}
}

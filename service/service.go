package service

import (
	"errors"
	"net/http"

	"github.com/NYTimes/gizmo/config"
	"github.com/clawio/authentication/authenticationcontroller"
	"github.com/clawio/authentication/authenticationcontroller/memory"
	"github.com/clawio/authentication/authenticationcontroller/simple"
	"github.com/clawio/authentication/lib"
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
		General                  *GeneralConfig
		AuthenticationController *AuthenticationControllerConfig
	}

	// GeneralConfig contains configuration parameters
	// for general parts of the service.
	GeneralConfig struct {
		BaseURL                  string
		JWTKey, JWTSigningMethod string
	}

	// AuthenticationControllerConfig holds the configuration for
	// an AuthenticationController.
	AuthenticationControllerConfig struct {
		Type string

		SimpleDriver string
		SimpleDSN    string

		MemoryUsers []*memory.User
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
	if cfg.General == nil {
		return nil, errors.New("config.General is nil")
	}

	var authenticationController authenticationcontroller.AuthenticationController
	switch cfg.AuthenticationController.Type {
	case "simple":
		a, err := getSimpleAuthenticationController(cfg)
		if err != nil {
			return nil, err
		}
		authenticationController = a
	case "memory":
		authenticationController = getMemoryAuthenticationController(cfg)
	default:
		return nil, errors.New("authenticationController type " + cfg.AuthenticationController.Type + " does not exist")
	}

	return &Service{
		Config: cfg,
		AuthenticationController: authenticationController,
	}, nil
}

func getSimpleAuthenticationController(cfg *Config) (authenticationcontroller.AuthenticationController, error) {
	authenticator := lib.NewAuthenticator(cfg.General.JWTKey, cfg.General.JWTSigningMethod)
	opts := &simple.Options{
		Driver:        cfg.AuthenticationController.SimpleDriver,
		DSN:           cfg.AuthenticationController.SimpleDSN,
		Authenticator: authenticator,
	}
	return simple.New(opts)
}
func getMemoryAuthenticationController(cfg *Config) authenticationcontroller.AuthenticationController {
	authenticator := lib.NewAuthenticator(cfg.General.JWTKey, cfg.General.JWTSigningMethod)
	opts := &memory.Options{
		Users:         cfg.AuthenticationController.MemoryUsers,
		Authenticator: authenticator,
	}
	return memory.New(opts)
}

// Prefix returns the string prefix used for all endpoints within
// this service.
func (s *Service) Prefix() string {
	return s.Config.General.BaseURL
}

// Middleware provides an http.Handler hook wrapped around all requests.
func (s *Service) Middleware(h http.Handler) http.Handler {
	return h
}

// Endpoints is a listing of all endpoints available in the MixedService.
func (s *Service) Endpoints() map[string]map[string]http.HandlerFunc {
	return map[string]map[string]http.HandlerFunc{
		"/metrics": {
			"GET": func(w http.ResponseWriter, r *http.Request) {
				prometheus.Handler().ServeHTTP(w, r)
			},
		},
		"/token": {
			"POST": prometheus.InstrumentHandlerFunc("/token", s.Token),
		},
	}
}

package service

import (
	"errors"
	"net/http"

	"github.com/NYTimes/gizmo/config"
	"github.com/clawio/authentication/authenticationcontroller"
	"github.com/clawio/authentication/authenticationcontroller/memory"
	"github.com/clawio/authentication/authenticationcontroller/simple"
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
		Type string

		SimpleDriver           string
		SimpleDSN              string
		SimpleJWTKey           string
		SimpleJWTSigningMethod string

		MemoryJWTKey           string
		MemoryJWTSigningMethod string
		MemoryUsers            []*memory.User
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
	opts := &simple.Options{
		Driver:           cfg.AuthenticationController.SimpleDriver,
		DSN:              cfg.AuthenticationController.SimpleDSN,
		JWTKey:           cfg.AuthenticationController.SimpleJWTKey,
		JWTSigningMethod: cfg.AuthenticationController.SimpleJWTSigningMethod,
	}
	return simple.New(opts)
}
func getMemoryAuthenticationController(cfg *Config) authenticationcontroller.AuthenticationController {
	opts := &memory.Options{
		Users:            cfg.AuthenticationController.MemoryUsers,
		JWTKey:           cfg.AuthenticationController.SimpleJWTKey,
		JWTSigningMethod: cfg.AuthenticationController.SimpleJWTSigningMethod,
	}
	return memory.New(opts)
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
		"/metrics": {
			"GET": func(w http.ResponseWriter, r *http.Request) {
				prometheus.Handler().ServeHTTP(w, r)
			},
		},
		"/authenticate": {
			"POST": prometheus.InstrumentHandlerFunc("/authenticate", s.Authenticate),
		},
		"/verify/{token}": {
			"GET": prometheus.InstrumentHandlerFunc("/verify", s.Verify),
		},
	}
}

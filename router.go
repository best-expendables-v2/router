package router

import (
	"errors"
	"github.com/best-expendables-v2/router/middleware"
	newrelic "github.com/newrelic/go-agent"
	"os"

	"github.com/best-expendables-v2/logger"

	"github.com/go-chi/chi"
)

type (
	// Configuration router configuration
	Configuration struct {
		// LoggerFactory using in ContextLogger middleware
		LoggerFactory logger.Factory

		// NewrelicApp
		NewrelicApp newrelic.Application

		// PanicHandler optional parameter
		// On nil panic returns only 500 status code
		PanicHandler middleware.PanicHandler

		AccessLog struct {
			// Disable access log middleware
			Disable bool
		}
	}
)

// New returns router with default list of middlewares
func New(cfg Configuration) (chi.Router, error) {
	if cfg.LoggerFactory == nil {
		return nil, errors.New("logger factory or newrelic is missing")
	}

	prefix, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	mux := chi.NewMux()
	mux.Use(
		middleware.RequestID(prefix),
		middleware.Authentication,
		middleware.ContextLogger(cfg.LoggerFactory),
		middleware.Newrelic(cfg.NewrelicApp),
		middleware.Recoverer(cfg.PanicHandler),
		middleware.Prometheus(),
	)

	if !cfg.AccessLog.Disable {
		mux.Use(middleware.AccessLog(middleware.AccessLogOptions{}))
	}

	return mux, nil
}

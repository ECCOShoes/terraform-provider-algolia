// Package client wraps the official Algolia Go SDK v4 for use by the Terraform
// provider. It exposes a thin constructor plus a couple of helpers so the
// resource code does not need to know the SDK's error and configuration
// details.
package client

import (
	"errors"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
)

// Config holds the credentials required to talk to an Algolia application.
type Config struct {
	// AppID is the Algolia application ID.
	AppID string
	// APIKey is an Algolia Admin API key with permission to manage index
	// settings and API keys.
	APIKey string
}

// New builds an Algolia Search API client. Authentication is handled by the
// SDK via the x-algolia-application-id and x-algolia-api-key headers.
func New(cfg Config) (*search.APIClient, error) {
	if cfg.AppID == "" {
		return nil, errors.New("app_id is required")
	}
	if cfg.APIKey == "" {
		return nil, errors.New("api_key is required")
	}
	return search.NewClient(cfg.AppID, cfg.APIKey)
}

// IsNotFound reports whether err is an Algolia API error with a 404 status.
func IsNotFound(err error) bool {
	var apiErr *search.APIError
	if errors.As(err, &apiErr) {
		return apiErr.Status == 404
	}
	return false
}

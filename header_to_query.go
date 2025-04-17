// Package header-to-query a plugin to convert headers to query parameters.
package header_to_query

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// Header represents a header mapping configuration
// for converting headers to query parameters.
type Header struct {
	Name       string `json:"name" yaml:"name"`
	Value      string `json:"value,omitempty" yaml:"value,omitempty"`
	KeepHeader bool   `json:"keepHeader,omitempty" yaml:"keepHeader,omitempty"`
}

// Config the plugin configuration.
type Config struct {
	Headers []Header `json:"headers,omitempty" yaml:"headers,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Headers: []Header{},
	}
}

// HeaderToQuery a plugin to convert headers to query parameters.
type HeaderToQuery struct {
	next     http.Handler
	headers  []Header
	name     string
}

// New created a new HeaderToQuery plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.Headers) == 0 {
		return nil, fmt.Errorf("headers cannot be empty")
	}

	return &HeaderToQuery{
		headers:  config.Headers,
		next:     next,
		name:     name,
	}, nil
}

func (p *HeaderToQuery) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	for _, h := range p.headers {
		values := req.Header.Values(h.Name)
		if len(values) == 0 {
			continue
		}
		// If Value is set, use as query param name, else use header name (lowercased)
		queryKey := h.Value
		if queryKey == "" {
			queryKey = h.Name
		}

		// Normalize query key after queryKey is set to avoid issues with misuse and invalid characters
		queryKey = normalizeKey(queryKey)
		for _, v := range values {
			q.Add(queryKey, v)
		}

		// Remove header if not kept
		if !h.KeepHeader {
			req.Header.Del(h.Name)
		}
	}
	// Set all query params at once
	req.URL.RawQuery = q.Encode()

	p.next.ServeHTTP(rw, req)
}

// normalizeKey converts header names to a suitable query key (lowercase, dashes to underscores, etc.)
func normalizeKey(s string) string {
	// Example: SERVICE-TAG -> service_tag
	res := strings.ReplaceAll(s, "-", "_")
	return strings.ToLower(res)
}

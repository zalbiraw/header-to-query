// Package headertoquery a plugin to convert headers to query parameters.
package headertoquery

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

// Header represents a header mapping configuration
// for converting headers to query parameters.
type Header struct {
	Name       string `json:"name" yaml:"name"`
	Key        string `json:"key,omitempty" yaml:"key,omitempty"`
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
	next    http.Handler
	headers []Header
	name    string
}

// New created a new HeaderToQuery plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.Headers) == 0 {
		return nil, errors.New("headers cannot be empty")
	}

	return &HeaderToQuery{
		headers: config.Headers,
		next:    next,
		name:    name,
	}, nil
}

func (p *HeaderToQuery) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	for _, h := range p.headers {
		values := r.Header.Values(h.Name)

		if len(values) == 0 {
			continue
		}

		// Remove all headers if keepHeader is false
		if !h.KeepHeader {
			r.Header.Del(h.Name)
		}

		// If Key is set, use as query param name, else use header name (lowercased)
		queryKey := h.Name
		if h.Key != "" {
			queryKey = h.Key
		}
		queryKey = normalizeKey(queryKey)

		for _, v := range values {
			q.Add(queryKey, v)
		}
	}
	r.URL.RawQuery = q.Encode()
	r.RequestURI = r.URL.RequestURI()

	clone := r.Clone(r.Context())
	clone.Body = r.Body

	p.next.ServeHTTP(rw, clone)
}

// normalizeKey converts header names to a suitable query key (lowercase, dashes to underscores, etc.)
func normalizeKey(s string) string {
	// Example: SERVICE-TAG -> service_tag
	res := strings.ReplaceAll(s, "-", "_")
	return strings.ToLower(res)
}

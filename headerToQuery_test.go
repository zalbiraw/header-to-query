package headertoquery_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zalbiraw/headertoquery"
)

func TestHeaderToQuery(t *testing.T) {
	cfg := headertoquery.CreateConfig()

	// Configure headers based on test data
	cfg.Headers = []headertoquery.Header{
		{
			Name: "SERVICE-TAG",
			Key:  "id",
		},
		{
			Name: "RANK",
		},
		{
			Name:       "GROUP",
			KeepHeader: true,
		},
	}

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := headertoquery.New(ctx, next, cfg, "header-to-query-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/get?test=test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set test headers
	req.Header.Add("SERVICE-TAG", "S117")
	req.Header.Add("SERVICE-TAG", "SPARTAN-117")
	req.Header.Add("SERVICE-TAG", "117")
	req.Header.Add("RANK", "Masterchief")
	req.Header.Add("GROUP", "UNSC")

	handler.ServeHTTP(recorder, req)

	// Assert headers
	assertHeaderNotExists(t, req, "SERVICE-TAG")
	assertHeaderNotExists(t, req, "RANK")
	assertHeaderEquals(t, req, "GROUP", "UNSC")

	// Assert query parameters
	assertQueryParamEquals(t, req, "id", "S117", "SPARTAN-117", "117")
	assertQueryParamEquals(t, req, "rank", "Masterchief")
	assertQueryParamEquals(t, req, "group", "UNSC")
	assertQueryParamEquals(t, req, "test", "test")
}

// assertHeaderNotExists checks that a header does not exist
func assertHeaderNotExists(t *testing.T, req *http.Request, header string) {
	t.Helper()
	if req.Header.Get(header) != "" {
		t.Errorf("header %q should not exist, got value %q", header, req.Header.Get(header))
	}
}

// assertHeaderEquals checks that a header exists with the expected value
func assertHeaderEquals(t *testing.T, req *http.Request, header, expected string) {
	t.Helper()
	if actual := req.Header.Get(header); actual != expected {
		t.Errorf("header %q: expected %q, got %q", header, expected, actual)
	}
}

// assertQueryParamEquals checks that a query parameter exists with the expected values (order-insensitive)
func assertQueryParamEquals(t *testing.T, req *http.Request, param string, expected ...string) {
	t.Helper()
	actual := req.URL.Query()[param]
	if !stringSlicesEqualIgnoreOrder(actual, expected) {
		t.Errorf("query parameter %q: expected values %v, got %v", param, expected, actual)
	}
}

// stringSlicesEqualIgnoreOrder checks if two string slices have the same elements, order-insensitive
func stringSlicesEqualIgnoreOrder(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	count := make(map[string]int)
	for _, v := range a {
		count[v]++
	}
	for _, v := range b {
		count[v]--
		if count[v] < 0 {
			return false
		}
	}
	return true
}

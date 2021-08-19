package internal

import (
	"fmt"
	"net/http"
	"testing"
)

func TestClientBuildURL(t *testing.T) {
	client := NewClient("https://example.com", "secret")
	actual := client.buildURL("/api/workflow")
	expected := "https://example.com/api/workflow"

	if actual != expected {
		t.Errorf("expected %q, got %q", actual, expected)
	}
}

func TestAddHeaders(t *testing.T) {
	test := func(t testing.TB, headers http.Header, key string, expected string) {
		t.Helper()

		actual := headers.Get(key)

		if actual != expected {
			t.Errorf("expected %v: %v, got %v", key, expected, actual)
		}
	}

	headers := http.Header{}
	accessKey := "secret"
	addHeaders(&headers, accessKey)

	test(t, headers, "Content-Type", "application/json")
	test(t, headers, "User-Agent", fmt.Sprintf("msgenctl/%v", Version))
	test(t, headers, "Ocp-Apim-Subscription-Key", accessKey)
}

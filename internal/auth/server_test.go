package auth

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestAuthServerLifecycle(t *testing.T) {
	var receivedToken string
	loginURL, err := StartAuthServer(func(token string) {
		receivedToken = token
	})
	if err != nil {
		t.Fatalf("StartAuthServer failed: %v", err)
	}
	if !strings.HasPrefix(loginURL, "http://127.0.0.1:") {
		t.Errorf("Unexpected login URL: %s", loginURL)
	}

	// Test GET /login page
	resp, err := http.Get(loginURL)
	if err != nil {
		t.Fatalf("Failed to GET %s: %v", loginURL, err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "Gotion Notion Login") {
		t.Errorf("Response body does not contain expected title")
	}

	// Test GET /callback with token
	callbackURL := strings.Replace(loginURL, "/login", "/callback?token=v02%3Atest_token_auth_123", 1)
	cbResp, err := http.Get(callbackURL)
	if err != nil {
		t.Fatalf("Failed to GET %s: %v", callbackURL, err)
	}
	defer cbResp.Body.Close()

	if cbResp.StatusCode != http.StatusOK {
		t.Errorf("Callback returned status %d; want 200", cbResp.StatusCode)
	}

	time.Sleep(50 * time.Millisecond)
	if receivedToken != "v02:test_token_auth_123" {
		t.Errorf("receivedToken = %q; want %q", receivedToken, "v02:test_token_auth_123")
	}

	StopAuthServer()
}

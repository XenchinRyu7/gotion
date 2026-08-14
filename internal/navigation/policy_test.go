package navigation

import (
	"testing"
)

func TestIsInternalURL(t *testing.T) {
	tests := []struct {
		url      string
		expected bool
	}{
		{"https://app.notion.com", true},
		{"https://app.notion.com/my-workspace/page-123", true},
		{"https://www.notion.so/login", true},
		{"https://notion.so/help", true},
		{"https://notion.com", true},
		{"https://custom.notion.site/docs", true},
		{"https://file.notionusercontent.com/f/u/123", true},
		{"https://accounts.google.com/o/oauth2/v2/auth", true},
		{"https://appleid.apple.com/auth/authorize", true},
		{"https://login.microsoftonline.com/common/oauth2/authorize", true},
		{"https://github.com/login/oauth/authorize", true},
		{"https://google.com/search?q=notion", false},
		{"https://youtube.com/watch?v=123", false},
		{"https://example.com", false},
		{"", false},
		{"invalid-url:::", false},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			got := IsInternalURL(tt.url)
			if got != tt.expected {
				t.Errorf("IsInternalURL(%q) = %v; want %v", tt.url, got, tt.expected)
			}
		})
	}
}

func TestGetStartURL(t *testing.T) {
	got := GetStartURL()
	if got != "https://app.notion.com" {
		t.Errorf("GetStartURL() = %q; want https://app.notion.com", got)
	}
}

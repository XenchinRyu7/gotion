package navigation

import (
	"net/url"
	"strings"
)

const DefaultStartURL = "https://app.notion.com"

// allowedSuffixes are domain suffixes that should be handled inside the Gotion desktop window.
var allowedSuffixes = []string{
	"notion.com",
	"notion.so",
	"notion.site",
	"notionusercontent.com",
}

// allowedAuthHosts are identity provider hosts required during Notion web authentication flows.
var allowedAuthHosts = []string{
	"accounts.google.com",
	"appleid.apple.com",
	"login.microsoftonline.com",
	"login.live.com",
	"github.com",
}

// GetStartURL returns the canonical initial URL for Notion.
func GetStartURL() string {
	return DefaultStartURL
}

// IsInternalURL checks whether the given URL belongs to Notion or an authorized authentication flow.
func IsInternalURL(rawURL string) bool {
	if rawURL == "" {
		return false
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	host := strings.ToLower(parsed.Hostname())
	if host == "" {
		return false
	}

	// Check Notion domains and subdomains
	for _, suffix := range allowedSuffixes {
		if host == suffix || strings.HasSuffix(host, "."+suffix) {
			return true
		}
	}

	// Check authentication provider hosts
	for _, authHost := range allowedAuthHosts {
		if host == authHost || strings.HasSuffix(host, "."+authHost) {
			return true
		}
	}

	return false
}

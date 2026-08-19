package script

import (
	"strings"
	"testing"
)

func TestGetInjectionScript(t *testing.T) {
	s := GetInjectionScript()
	if len(s) == 0 {
		t.Fatalf("GetInjectionScript returned empty string")
	}

	expectedSubstrings := []string{
		"gotion-mac-titlebar",
		"gotion-titlebar-windows",
		"gotion-btn-max { order: 2",
		"restore",
		"gotion_titlebar_style",
		"Title Bar: Mac",
		"gotion-traffic-lights",
		"gotion-btn-close",
		"gotion-btn-min",
		"gotion-btn-max",
		"gotion-resize-handle",
		"gotion-dropdown",
		"svg",
	}

	for _, sub := range expectedSubstrings {
		if !strings.Contains(s, sub) {
			t.Errorf("GetInjectionScript missing expected substring: %q", sub)
		}
	}
}

//go:build linux

package main

/*
#cgo webkit2_41 pkg-config: gtk+-3.0 webkit2gtk-4.1
#cgo !webkit2_41 pkg-config: gtk+-3.0 webkit2gtk-4.0
#include <gtk/gtk.h>
#include <webkit2/webkit2.h>

void enable_linux_persistent_cookies(const char* data_dir) {
    WebKitWebContext *ctx = webkit_web_context_get_default();
    if (!ctx) return;
    WebKitCookieManager *mgr = webkit_web_context_get_cookie_manager(ctx);
    if (!mgr) return;

    char cookie_path[1024];
    snprintf(cookie_path, sizeof(cookie_path), "%s/cookies.sqlite", data_dir);

    webkit_cookie_manager_set_persistent_storage(
        mgr,
        cookie_path,
        WEBKIT_COOKIE_PERSISTENT_STORAGE_SQLITE
    );
    webkit_cookie_manager_set_accept_policy(
        mgr,
        WEBKIT_COOKIE_POLICY_ACCEPT_ALWAYS
    );
}
*/
import "C"
import (
	"log"
	"os"
	"runtime/debug"
	"unsafe"
)

// initPlatformSettings initializes Linux-specific WebKitGTK persistent cookie storage
func (a *App) initPlatformSettings() {
	profileDir := a.userDataDir
	if profileDir != "" {
		_ = os.MkdirAll(profileDir, 0755)
		cDir := C.CString(profileDir)
		defer C.free(unsafe.Pointer(cDir))
		C.enable_linux_persistent_cookies(cDir)
		log.Printf("[Gotion Linux] Persistent SQLite cookie storage enabled at: %s/cookies.sqlite", profileDir)
	}
}

// sendNativeWheelZoom is a no-op fallback on Linux (handled by client-side CSS zoom).
func (a *App) sendNativeWheelZoom(delta int32) {
	log.Printf("[Gotion Go] sendNativeWheelZoom called with delta %d (Linux platform)", delta)
}

// trimProcessMemory triggers Go runtime GC & returns OS memory on Linux.
func (a *App) trimProcessMemory() {
	debug.FreeOSMemory()
}

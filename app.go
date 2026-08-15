package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"gotion/internal/auth"
	"gotion/internal/config"
	"gotion/internal/navigation"
	"gotion/internal/profile"
	"gotion/internal/script"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct manages application lifecycle and backend bindings.
type App struct {
	ctx         context.Context
	userDataDir string
	startURL    string
	windowState config.WindowState
	zoomSteps   int
}

// NewApp creates a new App instance with loaded configuration.
func NewApp() *App {
	userDataDir, err := profile.GetUserDataDir()
	if err != nil {
		log.Printf("[Gotion] Warning: Failed to get user data directory: %v", err)
	}

	savedState := config.LoadWindowState()

	return &App{
		userDataDir: userDataDir,
		startURL:    navigation.GetStartURL(),
		windowState: savedState,
	}
}

// startup is called when the app starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.initPlatformSettings()
	log.Printf("[Gotion] Application started successfully.")
	log.Printf("[Gotion] Target Start URL: %s", a.startURL)
	log.Printf("[Gotion] User Profile Directory: %s", a.userDataDir)

	// Restore window position or center
	if a.windowState.X >= 0 && a.windowState.Y >= 0 {
		runtime.WindowSetPosition(a.ctx, a.windowState.X, a.windowState.Y)
		runtime.WindowSetSize(a.ctx, a.windowState.Width, a.windowState.Height)
	} else {
		runtime.WindowCenter(a.ctx)
	}

	if a.windowState.IsMaximized {
		runtime.WindowMaximise(a.ctx)
	} else {
		runtime.WindowUnmaximise(a.ctx)
	}

	// Register Event listeners for Zoom actions triggered from external Notion DOM
	runtime.EventsOn(a.ctx, "zoom:in", func(data ...interface{}) {
		log.Printf("[Gotion Event] Received 'zoom:in' event from Frontend!")
		a.ZoomIn()
	})
	runtime.EventsOn(a.ctx, "zoom:out", func(data ...interface{}) {
		log.Printf("[Gotion Event] Received 'zoom:out' event from Frontend!")
		a.ZoomOut()
	})
	runtime.EventsOn(a.ctx, "zoom:reset", func(data ...interface{}) {
		log.Printf("[Gotion Event] Received 'zoom:reset' event from Frontend!")
		a.ResetZoom()
	})

	// Start watchdog to keep Mac titlebar and navigation hotkeys injected
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		injectionJS := script.GetInjectionScript()

		for {
			select {
			case <-a.ctx.Done():
				return
			case <-ticker.C:
				runtime.WindowExecJS(a.ctx, injectionJS)
			}
		}
	}()

	// Memory Trimmer: runs periodically to reclaim inactive/standby RAM pages
	go func() {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-a.ctx.Done():
				return
			case <-ticker.C:
				a.trimProcessMemory()
			}
		}
	}()
}

// domReady is called after the front-end has been loaded.
func (a *App) domReady(ctx context.Context) {
	log.Printf("[Gotion] DOM Ready.")
	runtime.WindowExecJS(a.ctx, script.GetInjectionScript())
}

// beforeClose is called before window closes to save state.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	a.SaveCurrentWindowState()
	return false
}

// shutdown is called when the app terminates.
func (a *App) shutdown(ctx context.Context) {
	a.SaveCurrentWindowState()
	log.Printf("[Gotion] Application shutting down cleanly.")
}

// SaveCurrentWindowState saves the current window geometry to config.
func (a *App) SaveCurrentWindowState() {
	if a.ctx == nil {
		return
	}

	isMaximized := runtime.WindowIsMaximised(a.ctx)
	x, y := runtime.WindowGetPosition(a.ctx)
	w, h := runtime.WindowGetSize(a.ctx)

	// If minimized or 0, avoid saving invalid dimensions
	if w <= 0 || h <= 0 {
		return
	}

	state := config.WindowState{
		X:           x,
		Y:           y,
		Width:       w,
		Height:      h,
		IsMaximized: isMaximized,
	}

	if err := config.SaveWindowState(state); err != nil {
		log.Printf("[Gotion] Warning: Failed to save window state: %v", err)
	} else {
		log.Printf("[Gotion] Saved window state: %+v", state)
	}
}

// Close closes the application window.
func (a *App) Close() {
	if a.ctx != nil {
		runtime.Quit(a.ctx)
	}
}

// Quit quits the application.
func (a *App) Quit() {
	if a.ctx != nil {
		runtime.Quit(a.ctx)
	}
}

// Minimise minimizes the application window.
func (a *App) Minimise() {
	if a.ctx != nil {
		runtime.WindowMinimise(a.ctx)
	}
}

// ToggleMaximise toggles between maximize and restore states reliably.
func (a *App) ToggleMaximise() {
	if a.ctx != nil {
		runtime.WindowToggleMaximise(a.ctx)
	}
}

// GetStartURL returns the canonical initial URL for Notion.
func (a *App) GetStartURL() string {
	return a.startURL
}

// GetUserDataDir returns the persistent WebView2 user data profile directory.
func (a *App) GetUserDataDir() string {
	return a.userDataDir
}

// OpenExternalURL opens a given external URL in the default browser.
func (a *App) OpenExternalURL(targetURL string) bool {
	if targetURL == "" {
		return false
	}
	log.Printf("[Gotion] Opening external URL in default browser: %s", targetURL)
	runtime.BrowserOpenURL(a.ctx, targetURL)
	return true
}

// IsInternalURL returns true if the URL belongs to Notion or its authorized auth providers.
func (a *App) IsInternalURL(targetURL string) bool {
	return navigation.IsInternalURL(targetURL)
}

// Reload reloads the current window/page.
func (a *App) Reload() {
	if a.ctx != nil {
		runtime.WindowReload(a.ctx)
	}
}

// GoBack navigates back in browser history.
func (a *App) GoBack() {
	if a.ctx != nil {
		runtime.WindowExecJS(a.ctx, "window.history.back();")
	}
}

// GoForward navigates forward in browser history.
func (a *App) GoForward() {
	if a.ctx != nil {
		runtime.WindowExecJS(a.ctx, "window.history.forward();")
	}
}

// HardReload forces a full reload of the application.
func (a *App) HardReload() {
	if a.ctx != nil {
		runtime.WindowReloadApp(a.ctx)
	}
}

// ZoomIn triggers client-side zoom in.
func (a *App) ZoomIn() {
	if a.ctx != nil {
		runtime.WindowExecJS(a.ctx, "if(window.__gotion_triggerZoomIn) window.__gotion_triggerZoomIn();")
	}
}

// ZoomOut triggers client-side zoom out.
func (a *App) ZoomOut() {
	if a.ctx != nil {
		runtime.WindowExecJS(a.ctx, "if(window.__gotion_triggerZoomOut) window.__gotion_triggerZoomOut();")
	}
}

// ResetZoom resets the native engine zoom back to default 100%.
func (a *App) ResetZoom() {
	if a.ctx != nil {
		runtime.WindowExecJS(a.ctx, "if(window.__gotion_triggerResetZoom) window.__gotion_triggerResetZoom();")
	}
}

const WHEEL_DELTA = 120

// GetAppVersion returns the Gotion client version.
func (a *App) GetAppVersion() string {
	return "0.1.0"
}

// StartBrowserLogin starts the local loopback auth bridge and opens the browser.
func (a *App) StartBrowserLogin() string {
	loginURL, err := auth.StartAuthServer(func(token string) {
		log.Printf("[Gotion Auth] Successfully received session token via browser auth!")
		if a.ctx != nil {
			runtime.EventsEmit(a.ctx, "auth:login_success", token)
			js := fmt.Sprintf(`(function(){
				document.cookie = "token_v2=" + encodeURIComponent(%q) + "; domain=.notion.so; path=/; max-age=31536000; secure";
				document.cookie = "token_v2=" + encodeURIComponent(%q) + "; domain=.notion.com; path=/; max-age=31536000; secure";
				window.location.href = "https://app.notion.com";
			})();`, token, token)
			runtime.WindowExecJS(a.ctx, js)
		}
	})
	if err != nil {
		log.Printf("[Gotion Auth] Error starting auth server: %v", err)
		return ""
	}
	if a.ctx != nil {
		runtime.BrowserOpenURL(a.ctx, loginURL)
	}
	return loginURL
}

// SetSessionToken saves and applies the Notion session token directly.
func (a *App) SetSessionToken(token string) bool {
	token = strings.TrimSpace(token)
	if token == "" {
		return false
	}
	token = strings.Trim(token, `"'`)
	if err := config.SaveSessionToken(token); err != nil {
		log.Printf("[Gotion Auth] Warning: Failed to save session token: %v", err)
	}
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "auth:login_success", token)
		js := fmt.Sprintf(`(function(){
			document.cookie = "token_v2=" + encodeURIComponent(%q) + "; domain=.notion.so; path=/; max-age=31536000; secure";
			document.cookie = "token_v2=" + encodeURIComponent(%q) + "; domain=.notion.com; path=/; max-age=31536000; secure";
			window.location.href = "https://app.notion.com";
		})();`, token, token)
		runtime.WindowExecJS(a.ctx, js)
	}
	return true
}

// GetSessionToken loads the saved session token from disk.
func (a *App) GetSessionToken() string {
	return config.LoadSessionToken()
}

// Logout clears the saved session and reloads back to the internal Gotion UI.
func (a *App) Logout() bool {
	_ = config.ClearSessionToken()
	if a.ctx != nil {
		runtime.WindowExecJS(a.ctx, `(function(){
			document.cookie = "token_v2=; domain=.notion.so; path=/; max-age=0; secure";
			document.cookie = "token_v2=; domain=.notion.com; path=/; max-age=0; secure";
		})();`)
		runtime.WindowReloadApp(a.ctx)
	}
	return true
}

// ReturnToInternalFrontend reloads the app to the root Wails frontend (dev or build).
func (a *App) ReturnToInternalFrontend() {
	if a.ctx != nil {
		runtime.WindowReloadApp(a.ctx)
	}
}




package main

import (
	"embed"
	"log"
	"os"
	"runtime/debug"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/src
var assets embed.FS

func main() {
	// Aggressive Chromium & V8 RAM trimming flags
	browserArgs := "--renderer-process-limit=1 --js-flags=\"--max-old-space-size=128 --optimize-for-size\" --enable-features=ServiceWorkerScopeCache,NetworkServiceInProcess,FontAccess --disable-features=Translate,OptimizationHints,MediaRouter,CalculateNativeWinOcclusion --disk-cache-size=314572800 --media-cache-size=104857600 --aggressive-cache-discard"
	_ = os.Setenv("WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS", browserArgs)

	// Optimize Go runtime memory footprint (trigger GC proactively)
	debug.SetGCPercent(30)

	// Create an instance of the app structure
	app := NewApp()

	// Configure Windows-specific WebView2 options
	windowsOptions := &windows.Options{
		WebviewUserDataPath:               app.GetUserDataDir(),
		Theme:                             windows.SystemDefault,
		BackdropType:                      windows.Auto,
		IsZoomControlEnabled:              true,
		DisablePinchZoom:                  false,
		DisableFramelessWindowDecorations: false, // Preserves Windows 11 rounded corners and aero shadow
		CustomTheme: &windows.ThemeSettings{
			DarkModeTitleBar:          windows.RGB(25, 25, 25),
			DarkModeTitleBarInactive:  windows.RGB(32, 32, 32),
			DarkModeTitleText:         windows.RGB(240, 240, 240),
			DarkModeTitleTextInactive: windows.RGB(140, 140, 140),
			DarkModeBorder:            windows.RGB(35, 35, 35),
			DarkModeBorderInactive:    windows.RGB(30, 30, 30),
		},
	}

	// Create application with options
	err := wails.Run(&options.App{
		Title:           "Gotion - Lightweight Notion Client",
		Width:           app.windowState.Width,
		Height:          app.windowState.Height,
		MinWidth:        800,
		MinHeight:       600,
		Frameless:       true, // Custom Mac-style Titlebar
		CSSDragProperty: "--wails-draggable",
		CSSDragValue:    "drag",
		DisableResize:   false,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour:       &options.RGBA{R: 25, G: 25, B: 25, A: 255},
		OnStartup:              app.startup,
		OnDomReady:             app.domReady,
		OnBeforeClose:          app.beforeClose,
		OnShutdown:             app.shutdown,
		Windows:                windowsOptions,
		BindingsAllowedOrigins: "https://*.notion.com,https://*.notion.so,https://notion.so,https://notion.com,https://app.notion.com",
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		log.Fatalf("[Gotion] Application error: %v", err)
	}
}

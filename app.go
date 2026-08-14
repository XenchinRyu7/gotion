package main

import (
	"context"
	"log"
	"runtime/debug"
	"syscall"
	"time"
	"unsafe"

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
		ticker := time.NewTicker(250 * time.Millisecond)
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

	// Hardcore Memory Trimmer: runs periodically to flush inactive/cached RAM to keep footprint ultra-low
	go func() {
		ticker := time.NewTicker(10 * time.Second)
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
	if a.ctx == nil {
		return
	}

	if runtime.WindowIsMaximised(a.ctx) {
		log.Printf("[Gotion] Restoring window to normal size...")
		runtime.WindowUnmaximise(a.ctx)

		// Explicit dimension restore as fallback for frameless Win32
		w := a.windowState.Width
		h := a.windowState.Height
		if w < 600 || h < 400 {
			w = 1280
			h = 800
		}
		runtime.WindowSetSize(a.ctx, w, h)
		if a.windowState.X >= 0 && a.windowState.Y >= 0 {
			runtime.WindowSetPosition(a.ctx, a.windowState.X, a.windowState.Y)
		} else {
			runtime.WindowCenter(a.ctx)
		}
	} else {
		log.Printf("[Gotion] Maximising window...")
		// Save current position and dimensions before maximizing
		x, y := runtime.WindowGetPosition(a.ctx)
		w, h := runtime.WindowGetSize(a.ctx)
		if w >= 600 && h >= 400 {
			a.windowState.X = x
			a.windowState.Y = y
			a.windowState.Width = w
			a.windowState.Height = h
		}
		runtime.WindowMaximise(a.ctx)
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

// OpenExternalURL opens a given external URL in the Windows default browser.
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

// ZoomIn triggers low-level native engine zoom in.
func (a *App) ZoomIn() {
	log.Printf("[Gotion Go] ===> ZoomIn() called from Frontend!")
	if a.ctx == nil {
		log.Printf("[Gotion Go] ERROR: a.ctx is nil")
		return
	}
	a.zoomSteps++
	if a.zoomSteps > 15 {
		a.zoomSteps = 15
	}
	log.Printf("[Gotion Go] zoomSteps is now: %d. Dispatching wheel delta: +%d", a.zoomSteps, WHEEL_DELTA)
	a.sendNativeWheelZoom(WHEEL_DELTA)
}

// ZoomOut triggers low-level native engine zoom out.
func (a *App) ZoomOut() {
	log.Printf("[Gotion Go] ===> ZoomOut() called from Frontend!")
	if a.ctx == nil {
		log.Printf("[Gotion Go] ERROR: a.ctx is nil")
		return
	}
	a.zoomSteps--
	if a.zoomSteps < -15 {
		a.zoomSteps = -15
	}
	log.Printf("[Gotion Go] zoomSteps is now: %d. Dispatching wheel delta: -%d", a.zoomSteps, WHEEL_DELTA)
	a.sendNativeWheelZoom(-WHEEL_DELTA)
}

// ResetZoom resets the native engine zoom back to default 100%.
func (a *App) ResetZoom() {
	log.Printf("[Gotion Go] ===> ResetZoom() called from Frontend! Current steps: %d", a.zoomSteps)
	if a.ctx == nil {
		log.Printf("[Gotion Go] ERROR: a.ctx is nil")
		return
	}
	if a.zoomSteps > 0 {
		for i := 0; i < a.zoomSteps; i++ {
			log.Printf("[Gotion Go] Reset step %d/%d (ZoomOut)", i+1, a.zoomSteps)
			a.sendNativeWheelZoom(-WHEEL_DELTA)
			time.Sleep(20 * time.Millisecond)
		}
	} else if a.zoomSteps < 0 {
		for i := 0; i < -a.zoomSteps; i++ {
			log.Printf("[Gotion Go] Reset step %d/%d (ZoomIn)", i+1, -a.zoomSteps)
			a.sendNativeWheelZoom(WHEEL_DELTA)
			time.Sleep(20 * time.Millisecond)
		}
	}
	a.zoomSteps = 0
	log.Printf("[Gotion Go] ResetZoom() completed. zoomSteps = 0")
}

// Win32 API declarations for native input synthesis
var (
	modUser32        = syscall.NewLazyDLL("user32.dll")
	procGetCursorPos = modUser32.NewProc("GetCursorPos")
	procSetCursorPos = modUser32.NewProc("SetCursorPos")
	procMouseEvent   = modUser32.NewProc("mouse_event")
	procKeybdEvent   = modUser32.NewProc("keybd_event")
)

const (
	VK_CONTROL        = 0x11
	KEYEVENTF_KEYUP   = 0x0002
	MOUSEEVENTF_WHEEL = 0x0800
	WHEEL_DELTA       = 120
)

type winPoint struct {
	X int32
	Y int32
}

// sendNativeWheelZoom sends a hardware-level Ctrl+MouseWheel event directly into the OS message pump.
func (a *App) sendNativeWheelZoom(delta int32) {
	var curPos winPoint
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&curPos)))

	x, y := runtime.WindowGetPosition(a.ctx)
	w, h := runtime.WindowGetSize(a.ctx)

	log.Printf("[Gotion Go] Window Rect: [x:%d, y:%d, w:%d, h:%d], Cursor: (%d, %d)", x, y, w, h, curPos.X, curPos.Y)

	cursorOutside := false
	if curPos.X < int32(x) || curPos.X > int32(x+w) || curPos.Y < int32(y) || curPos.Y > int32(y+h) {
		cursorOutside = true
		centerX := int32(x + w/2)
		centerY := int32(y + h/2)
		log.Printf("[Gotion Go] Cursor is outside window. Temporarily moving to center: (%d, %d)", centerX, centerY)
		procSetCursorPos.Call(uintptr(centerX), uintptr(centerY))
		time.Sleep(10 * time.Millisecond)
	}

	// 1. Press CTRL key down
	log.Printf("[Gotion Go] Sending KeyDown(VK_CONTROL)...")
	procKeybdEvent.Call(uintptr(VK_CONTROL), 0, 0, 0)
	time.Sleep(10 * time.Millisecond)

	// 2. Fire mouse wheel event (positive = zoom in, negative = zoom out)
	log.Printf("[Gotion Go] Sending MouseEvent(MOUSEEVENTF_WHEEL, delta=%d)...", delta)
	procMouseEvent.Call(
		uintptr(MOUSEEVENTF_WHEEL),
		0,
		0,
		uintptr(uint32(delta)),
		0,
	)
	time.Sleep(10 * time.Millisecond)

	// 3. Release CTRL key
	log.Printf("[Gotion Go] Sending KeyUp(VK_CONTROL)...")
	procKeybdEvent.Call(uintptr(VK_CONTROL), 0, uintptr(KEYEVENTF_KEYUP), 0)
	time.Sleep(10 * time.Millisecond)

	// 4. Restore cursor position if moved
	if cursorOutside {
		log.Printf("[Gotion Go] Restoring cursor back to (%d, %d)", curPos.X, curPos.Y)
		procSetCursorPos.Call(uintptr(curPos.X), uintptr(curPos.Y))
	}
	log.Printf("[Gotion Go] sendNativeWheelZoom finished successfully.")
}

// GetAppVersion returns the Gotion client version.
func (a *App) GetAppVersion() string {
	return "0.1.0"
}

// Win32 API declarations for process memory trimming
var (
	modKernel32                  = syscall.NewLazyDLL("kernel32.dll")
	modPsapi                     = syscall.NewLazyDLL("psapi.dll")
	procGetCurrentProcess        = modKernel32.NewProc("GetCurrentProcess")
	procGetCurrentProcessId      = modKernel32.NewProc("GetCurrentProcessId")
	procSetProcessWorkingSetSize = modKernel32.NewProc("SetProcessWorkingSetSize")
	procCreateToolhelp32Snapshot = modKernel32.NewProc("CreateToolhelp32Snapshot")
	procProcess32First           = modKernel32.NewProc("Process32FirstW")
	procProcess32Next            = modKernel32.NewProc("Process32NextW")
	procOpenProcess              = modKernel32.NewProc("OpenProcess")
	procCloseHandle              = modKernel32.NewProc("CloseHandle")
	procEmptyWorkingSet          = modPsapi.NewProc("EmptyWorkingSet")
)

const (
	TH32CS_SNAPPROCESS        = 0x00000002
	PROCESS_SET_QUOTA         = 0x0100
	PROCESS_QUERY_INFORMATION = 0x0400
)

type PROCESSENTRY32W struct {
	DwSize              uint32
	CntUsage            uint32
	Th32ProcessID       uint32
	Th32DefaultHeapID   uintptr
	Th32ModuleID        uint32
	CntThreads          uint32
	Th32ParentProcessID uint32
	PcPriClassBase      int32
	DwFlags             uint32
	SzExeFile           [260]uint16
}

// trimProcessMemory flushes unused working set pages to disk for Gotion and its child WebView2 processes.
func (a *App) trimProcessMemory() {
	// 1. Return Go runtime memory to OS
	debug.FreeOSMemory()

	// 2. Trim Gotion main process working set
	hMain, _, _ := procGetCurrentProcess.Call()
	if hMain != 0 {
		procSetProcessWorkingSetSize.Call(hMain, ^uintptr(0), ^uintptr(0))
		procEmptyWorkingSet.Call(hMain)
	}

	// 3. Find and trim all child msedgewebview2.exe processes
	myPid, _, _ := procGetCurrentProcessId.Call()
	hSnapshot, _, _ := procCreateToolhelp32Snapshot.Call(uintptr(TH32CS_SNAPPROCESS), 0)
	if hSnapshot == 0 || hSnapshot == ^uintptr(0) {
		return
	}
	defer procCloseHandle.Call(hSnapshot)

	var entry PROCESSENTRY32W
	entry.DwSize = uint32(unsafe.Sizeof(entry))

	ret, _, _ := procProcess32First.Call(hSnapshot, uintptr(unsafe.Pointer(&entry)))
	for ret != 0 {
		if entry.Th32ParentProcessID == uint32(myPid) {
			hChild, _, _ := procOpenProcess.Call(
				uintptr(PROCESS_SET_QUOTA|PROCESS_QUERY_INFORMATION),
				0,
				uintptr(entry.Th32ProcessID),
			)
			if hChild != 0 {
				procSetProcessWorkingSetSize.Call(hChild, ^uintptr(0), ^uintptr(0))
				procEmptyWorkingSet.Call(hChild)
				procCloseHandle.Call(hChild)
			}
		}
		ret, _, _ = procProcess32Next.Call(hSnapshot, uintptr(unsafe.Pointer(&entry)))
	}
}


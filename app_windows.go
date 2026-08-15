//go:build windows

package main

import (
	"log"
	"runtime/debug"
	"syscall"
	"time"
	"unsafe"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

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

// initPlatformSettings initializes Windows-specific WebView2 settings.
func (a *App) initPlatformSettings() {
	// Windows WebView2 handles user data path via windows.Options.WebviewUserDataPath
}

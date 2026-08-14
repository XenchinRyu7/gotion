<div align="center">

# 🪶 Gotion
### *The Ultra-Lightweight, Privacy-Respecting Desktop Client for Notion*

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Wails v2](https://img.shields.io/badge/Wails-v2-DF1A5A?style=flat&logo=wails)](https://wails.io)
[![WebView2](https://img.shields.io/badge/Runtime-Microsoft%20Edge%20WebView2-0078D7?style=flat&logo=microsoftedge)](https://developer.microsoft.com/en-us/microsoft-edge/webview2/)
[![Platform](https://img.shields.io/badge/Platform-Windows%2010%20%7C%2011-blue?style=flat&logo=windows)](https://microsoft.com)
[![Installer Size](https://img.shields.io/badge/Installer-~6.4%20MB-success?style=flat)](https://github.com)

A high-performance native Windows desktop wrapper for Notion built with **Go**, **Wails v2**, and **Microsoft Edge WebView2**. Zero bundled Chromium bloat, zero background telemetry overhead, and ultra-low RAM consumption (< 400 MB).

---

</div>

## Why Gotion?

The official Notion Desktop application uses **Electron**, which bundles an entire Node.js runtime and full Chromium instance into your PC. This often consumes **700 MB to 1.4 GB+ of RAM**, runs dozens of background processes, and tracks continuous telemetry events.

**Gotion** solves this by embedding Notion's genuine web workspace inside Windows' native **WebView2** runtime combined with Go native optimizations:

| Feature | Official Notion Client (Electron) | Gotion |
| :--- | :--- | :--- |
| **RAM Usage** | ~700 MB – 1.4 GB+ | **< 400 MB** *(with active Win32 Working Set Trimmer)* |
| **Installer Size** | ~110 MB+ | **~6.4 MB** *(NSIS Compressed)* |
| **Telemetry & Trackers** | Active (Segment, Intercom, Datadog RUM) | **Blocked & Mocked (0ms overhead)** |
| **Startup Speed** | Slow / Heavy Bundle Load | **Near-Instant** *(Persistent 300MB Disk Cache)* |
| **Sidebar** | Promo banners ("Notion Desktop") | **Clean & Distraction-Free** |
| **Window Frame** | Standard Windows Titlebar | **Mac-Style Minimalist Frameless Titlebar** |
| **Theme Sync** | Manual / Desynced Window Chrome | **Real-Time Dynamic Notion Light & Dark Sync** |

---

## Key Features

### 1. Frameless Titlebar & Window Controls
* **Mac-Style Traffic Lights:** Custom close (`#ff5f56`), minimize (`#ffbd2e`), and maximize (`#27c93f`) window controls with bidirectional maximize/restore toggling.
* **8-Direction Native Resizing:** Smooth resizing handles along all borders and corners (`N`, `S`, `E`, `W`, `NE`, `NW`, `SE`, `SW`).
* **Titlebar Counter-Scaling:** Keeps the title bar locked at an exact physical **38px** height during native zoom scaling.

### 2. Real-Time Dynamic Theme Adaptation
* Uses a 4-layer multi-tier luminance sampler (`detectNotionDarkTheme`) that directly inspects rendered Notion workspace elements (`.notion-frame`, `.notion-scroller`).
* **Light Theme:** Notion Warm Paper (`#f6f5f4`) surface with hairline borders (`#e6e6e6`) and charcoal typography (`#31302e`).
* **Dark Theme:** Notion Dark Surface (`#191919`) with dark charcoal borders (`#2c2c2c`) and soft off-white typography (`#e6e6e6`).
* Toggling theme inside Notion (`Ctrl + Shift + L` or via Appearance settings) changes the titlebar **instantly without page reload**.

### 3. Performance & Bloatware Stripper
* Intercepts and blocks network requests to heavy tracking and analytics endpoints (`api.segment.io`, `api.mixpanel.com`, `browser-intake-datadoghq.com`, `widget.intercom.io`, `amplitude.com`, `notion.so/api/v3/logUserEvent`, `notion.so/api/v3/ping`).
* Mocks tracking libraries (`window.analytics`, `window.Intercom`, `window.datadogRum`, `window.mixpanel`, `window.posthog`) with zero-cost no-op proxies.
* Strips unwanted promo banners (*"Aplikasi Notion - Notion Desktop"*) from the sidebar.

### 4. Hardcore V8 Memory Cap & Win32 Working Set Trimmer
* Restricts V8 JavaScript heap to 128 MB max with bytecode size optimization (`--js-flags="--max-old-space-size=128 --optimize-for-size"`).
* Limits Chromium renderers to a single process (`--renderer-process-limit=1`).
* Employs a background Win32 memory trimmer calling `EmptyWorkingSet` and `SetProcessWorkingSetSize` every 10 seconds across the main process and all child WebView2 processes, continuously reclaiming standby RAM.

### 5. Aggressive Disk Caching & Persistent Profile
* Allocates a dedicated 300 MB Chromium Disk Cache in `%APPDATA%\Gotion\Profile`.
* Large Notion JavaScript bundles and IndexedDB stores are preserved locally for lightning-fast subsequent launches.

---

## Keyboard Shortcuts Reference

| Shortcut | Action |
| :--- | :--- |
| <kbd>Alt</kbd> + <kbd>←</kbd> | Navigate Back in History |
| <kbd>Alt</kbd> + <kbd>→</kbd> | Navigate Forward in History |
| <kbd>Ctrl</kbd> + <kbd>R</kbd> / <kbd>F5</kbd> | Reload Page |
| <kbd>Ctrl</kbd> + <kbd>+</kbd> / <kbd>=</kbd> | Zoom In (Native Engine Zoom) |
| <kbd>Ctrl</kbd> + <kbd>-</kbd> | Zoom Out (Native Engine Zoom) |
| <kbd>Ctrl</kbd> + <kbd>0</kbd> | Reset Zoom to 100% |
| <kbd>Ctrl</kbd> + <kbd>MouseWheel</kbd> | Smooth Low-Level Native Zoom |
| <kbd>F11</kbd> | Maximize / Restore Window |
| <kbd>Ctrl</kbd> + <kbd>Shift</kbd> + <kbd>L</kbd> | Toggle Notion Light / Dark Mode |
| <kbd>Alt</kbd> + <kbd>F4</kbd> | Quit Application |

---

## Requirements & Building

### Prerequisites
* **Windows 10 / 11** (64-bit)
* **Go 1.21+**
* **Wails CLI v2** (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)
* **NSIS** *(optional, for generating the Windows installer)*

### Development
Run Gotion in live development mode with hot reload:
```powershell
wails dev
```

### Standalone Production Binary
Compile an ultra-optimized standalone executable (`build\bin\gotion.exe`):
```powershell
wails build -ldflags "-s -w"
```

### NSIS Setup Installer (~6.4 MB)
Compile the standalone executable and package it into a Windows NSIS Setup Installer (`build\bin\gotion-amd64-installer.exe`):
```powershell
wails build -nsis -ldflags "-s -w"
```

---

## Project Structure

```
gotion/
├── app.go                  # Lifecycle, Win32 zoom synthesis, working set memory trimmer
├── main.go                 # Wails options, Chromium & V8 arguments, window config
├── wails.json              # Wails project config and product metadata
├── internal/
│   ├── config/             # Window state persistence (dimensions, position, maximized state)
│   ├── navigation/         # URL routing policies & auth host whitelist
│   ├── profile/            # User profile data directory resolver (%APPDATA%\Gotion\Profile)
│   └── script/             # Injected JS & CSS (theme observer, tracker stripper, titlebar)
├── build/
│   ├── appicon.png         # Main application icon
│   ├── windows/            # Windows manifest, info.json, and icon.ico
│   │   └── installer/      # NSIS installer script (project.nsi)
│   └── bin/                # Compiled binaries and installer output
└── frontend/               # Frontend asset server directory
```

---

## Security & Privacy

* **Strict Navigation Isolation:** Notion workspaces and identity providers (Google, Apple, Microsoft, GitHub) stay inside the desktop client; external outbound links are intercepted and routed to the user's default system browser.
* **Telemetry Stripping:** Telemetry calls are blocked locally on the client before leaving the machine.
* **Encrypted Storage:** WebView2 credentials and session cookies are stored securely in the user's local Windows DPAPI-protected profile directory.

---

## License & Disclaimer

This project is open-source under the [MIT License](LICENSE).

> **Disclaimer:** Gotion is an unofficial third-party desktop client and is not affiliated with, maintained, sponsored, or endorsed by Notion Labs, Inc. Notion and the Notion logo are registered trademarks of Notion Labs, Inc.

<div align="center">

# 🪶 Gotion
### *The Ultra-Lightweight, Cross-Platform Desktop Client for Notion*

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Wails v2](https://img.shields.io/badge/Wails-v2-DF1A5A?style=flat&logo=wails)](https://wails.io)
[![Runtime](https://img.shields.io/badge/Runtime-WebView2%20%7C%20WebKitGTK%20%7C%20WKWebView-0078D7?style=flat)](https://wails.io)
[![Platforms](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-blue?style=flat)](https://github.com)
[![Installer Size](https://img.shields.io/badge/Installer-~6.4%20MB-success?style=flat)](https://github.com)

A high-performance, native desktop wrapper for Notion built with **Go**, **Wails v2**, and native operating system web engines (**Microsoft Edge WebView2** on Windows, **WebKitGTK** on Linux, and **WKWebView** on macOS). Zero bundled Chromium bloat, zero background telemetry overhead, and ultra-low memory consumption (< 400 MB).

---

</div>

## Overview

The official Notion Desktop client relies on **Electron**, packaging an entire Node.js runtime and full Chromium browser instance into the application. This typically consumes **700 MB to 1.4 GB+ of RAM**, runs multiple redundant background helper processes, and introduces noticeable startup latency.

<table>
  <tr>
    <td width="50%">
      <img src="https://github.com/user-attachments/assets/f6cf9433-9962-448b-90b4-3f2e1b27cd4b" alt="1" style="width: 100%; height: auto;" />
    </td>
    <td width="50%">
      <img src="https://github.com/user-attachments/assets/972e6f70-3b2e-4938-abd1-7040c4d26cee" alt="2" style="width: 100%; height: auto;" />
    </td>
  </tr>
</table>



**Gotion** replaces the heavy Electron shell with a lightweight Go backend that embeds the genuine Notion web application (`https://app.notion.com`) inside your operating system's native web engine:

```
┌────────────────────────────────────────────────────────────────────────┐
│                        Gotion Desktop Shell                            │
│     (Go Runtime + Wails v2 + Native Frameless Window Architecture)     │
├────────────────────────────────────────────────────────────────────────┤
│                      Native OS Web Engine                              │
│   Windows: WebView2   |   Linux: WebKit2GTK   |   macOS: WKWebView     │
├────────────────────────────────────────────────────────────────────────┤
│                       Notion Web Workspace                             │
│       (https://app.notion.com + Persistent Local Profile Cache)        │
└────────────────────────────────────────────────────────────────────────┘
```

---

## Performance Comparison

| Metric / Feature | Official Notion Client (Electron) | Gotion |
| :--- | :--- | :--- |
| **Memory Footprint (Idle / Active)** | ~700 MB – 1.4 GB+ | **< 400 MB** *(with active Win32 Working Set Trimming)* |
| **Binary / Installer Size** | ~110 MB+ | **~6.4 MB** *(Windows NSIS Compressed)* |
| **Background Telemetry** | Active third-party analytics | **Cleaned & Minimalist** |
| **Startup Time** | Slow / Cold bundle parsing | **Near-Instant** *(Persistent 300MB Local Disk Cache)* |
| **Multi-Account Support** | Supported | **Supported** *(Direct workspace switching)* |
| **Window Frame** | Standard OS Window Chrome | **Custom Frameless Titlebar with Mac/Windows Style Controls** |
| **Theme Synchronization** | Manual / Desynced Window Titlebar | **Real-Time Dynamic Light & Dark Theme Detection** |
| **Supported Platforms** | Windows, macOS | **Windows 10/11, Linux (AppImage/Deb/Rpm), macOS** |

---

## Core Features

### 1. Frameless Window with 8-Directional Native Resizing
* **Switchable Titlebar Styles:** Choose Mac-style traffic lights or Windows-style right-aligned controls from `Actions -> Title Bar`.
* **Mac-Style Traffic Lights:** Custom Close (`#ff5f56`), Minimize (`#ffbd2e`), and Maximize/Restore (`#27c93f`) window controls with bidirectional maximize toggling.
* **Persistent Titlebar Preference:** Selected style is saved in the WebView profile and restored on next launch.
* **8-Direction Resizing Handles:** Smooth resizing along all edges and corners (`N`, `S`, `E`, `W`, `NE`, `NW`, `SE`, `SW`).
* **Titlebar Counter-Scaling:** Automatically computes `1 / zoomRatio` to ensure the custom 38px titlebar visually retains its physical dimensions during native engine zoom.

### 2. Real-Time Dynamic Theme Adaptation
* Features a 4-layer multi-tier luminance sampler (`detectNotionDarkTheme`) that directly inspects rendered workspace containers (`.notion-frame`, `.notion-scroller`):
  * **Light Theme:** Notion Warm Paper (`#f6f5f4`) surface with hairline borders (`#e6e6e6`) and charcoal typography (`#31302e`).
  * **Dark Theme:** Notion Dark Surface (`#191919`) with dark hairline borders (`#2c2c2c`) and soft off-white typography (`#e6e6e6`).
* Toggling themes inside Notion (`Ctrl + Shift + L` or via Settings -> Appearance) updates the entire title bar instantly without requiring an application reload.

### 3. High-Performance Caching & Memory Optimization
* **Dedicated Disk Caching:** Allocates 300 MB disk cache and 100 MB media cache in persistent profile storage, ensuring Notion's large bundle files load locally on startup.
* **Working Set Memory Trimming:** Periodically flushes inactive standby memory pages out of physical RAM via native Win32 `EmptyWorkingSet` APIs on Windows.
* **Go Runtime Tuning:** Configures proactive garbage collection (`debug.SetGCPercent(30)`) to maintain a lean memory footprint.

### 4. Clean Workspace Interface
* Automatically removes promotional download banners (*"Aplikasi Notion - Notion Desktop"*) from the sidebar to provide a clean, native workspace feel.

### 5. Flexible Authentication System
* **Native In-App OAuth:** Supports direct login via Google, Apple, Microsoft, GitHub, Okta, Auth0, and SAML identity providers on Windows and macOS.
* **Local Loopback Auth Bridge:** Built-in loopback authentication server (`internal/auth/server.go`) listening on `http://127.0.0.1:28795/login` for direct `token_v2` token import, export, and external browser authentication.

> [!WARNING]
> **Linux Authentication Requirement:**
> Due to strict bot detection and embedded browser restrictions enforced by Google/Apple OAuth on WebKitGTK (`403: disallowed_useragent`), authentication on **Linux** is supported exclusively via **"Continue with Email" (Magic Link / OTP Code)** or direct session token import (`token_v2`) via the built-in Auth Bridge. Google SSO and Apple Sign-In work out-of-the-box on Windows and macOS.

### 6. Safe Navigation Routing
* Notion workspaces and identity provider authentication domains remain inside the application window.
* External outbound links are automatically intercepted and routed to the operating system's default web browser.

---

## Keyboard Shortcuts Reference

| Shortcut | Action |
| :--- | :--- |
| <kbd>Alt</kbd> + <kbd>←</kbd> | Navigate Back in History |
| <kbd>Alt</kbd> + <kbd>→</kbd> | Navigate Forward in History |
| <kbd>Ctrl</kbd> + <kbd>R</kbd> / <kbd>F5</kbd> | Reload Page |
| <kbd>Ctrl</kbd> + <kbd>+</kbd> / <kbd>=</kbd> | Zoom In |
| <kbd>Ctrl</kbd> + <kbd>-</kbd> | Zoom Out |
| <kbd>Ctrl</kbd> + <kbd>0</kbd> | Reset Zoom to 100% |
| <kbd>Ctrl</kbd> + <kbd>MouseWheel</kbd> | Smooth Native Engine Zoom |
| <kbd>F11</kbd> | Maximize / Restore Window |
| <kbd>Ctrl</kbd> + <kbd>Shift</kbd> + <kbd>L</kbd> | Toggle Notion Light / Dark Theme |
| <kbd>Alt</kbd> + <kbd>F4</kbd> | Quit Application |

---

## Cross-Platform Storage & Profiles

Gotion isolates user data into persistent profile directories so login sessions, cookies, and local caches survive restarts:

| Platform | Profile & Cache Directory | Configuration Directory |
| :--- | :--- | :--- |
| **Windows** | `%LOCALAPPDATA%\Gotion\profile` | `%LOCALAPPDATA%\Gotion\config` |
| **Linux** | `~/.config/Gotion/profile` | `~/.config/Gotion/config` |
| **macOS** | `~/Library/Application Support/Gotion/profile` | `~/Library/Application Support/Gotion/config` |

---

## Building and Packaging

### Prerequisites

#### Windows
* **Windows 10 / 11** (64-bit)
* **Go 1.21+**
* **Wails CLI v2** (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)
* **NSIS** *(required for compiling the `.exe` setup installer)*

#### Linux
* **Go 1.21+**
* **Wails CLI v2**
* **GTK3 & WebKitGTK Development Headers**:
  * **Debian / Ubuntu / Linux Mint:**
    ```bash
    sudo apt update
    sudo apt install build-essential pkg-config libgtk-3-dev libwebkit2gtk-4.0-dev
    ```
    *(For modern distributions with WebKitGTK 4.1: `libwebkit2gtk-4.1-dev`)*
  * **Arch Linux / Manjaro:**
    ```bash
    sudo pacman -S base-devel gtk3 webkit2gtk
    ```
  * **Fedora / RHEL:**
    ```bash
    sudo dnf install gcc-c++ gtk3-devel webkit2gtk4.0-devel
    ```

#### macOS
* **macOS 11.0+** (Apple Silicon or Intel)
* **Xcode Command Line Tools** (`xcode-select --install`)
* **Go 1.21+** and **Wails CLI v2**

---

### Local Build Commands

#### Windows
```powershell
# Build standalone executable (build/bin/gotion.exe)
wails build -ldflags "-s -w"

# Build executable + compressed NSIS setup installer (build/bin/gotion-amd64-installer.exe)
wails build -nsis -ldflags "-s -w"
```

#### Linux
```bash
# Build native Linux binary (build/bin/gotion)
wails build -platform linux/amd64 -ldflags "-s -w"
```

#### macOS
```bash
# Build Universal macOS Application Bundle (build/bin/gotion.app)
wails build -platform darwin/universal -ldflags "-s -w"
```

---

## Automated Multi-Platform Releases (GitHub Actions)

Gotion includes a production-ready CI/CD pipeline in [`.github/workflows/release.yml`](.github/workflows/release.yml). Pushing a semantic version tag (e.g. `v1.0.0`) or triggering the workflow manually automatically compiles and publishes packages for all major operating systems:

* **Windows:**
  * `gotion-windows-amd64.exe` (Standalone Portable Executable)
  * `gotion-windows-amd64-installer.exe` (NSIS Setup Wizard)
* **Linux:**
  * `gotion-1.0.0-x86_64.AppImage` (Universal Portable AppImage)
  * `gotion_1.0.0_amd64.deb` (Debian, Ubuntu, Pop!_OS, Mint)
  * `gotion-1.0.0.x86_64.rpm` (Fedora, RHEL, openSUSE)
  * `gotion-linux-amd64.tar.gz` (Standard ELF Tarball)
* **macOS:**
  * `gotion-1.0.0-macos-universal.dmg` (Drag-and-Drop Disk Image)
  * `gotion-macos-universal.zip` (Universal Application Bundle)
* **Checksums:**
  * `SHA256SUMS.txt` (Automated verification digests)

---

## Project Structure

```
gotion/
├── .github/workflows/      # Automated multi-platform CI/CD release pipeline (release.yml)
├── app.go                  # Core application lifecycle, window state, zoom and auth bindings
├── app_windows.go          # Windows-specific memory trimmer and native mouse input synthesis
├── app_linux.go            # Linux-specific WebKitGTK persistent cookie storage with SQLite
├── app_other.go            # Fallback handlers for macOS and other platforms
├── main.go                 # Application entrypoint, Wails initialization, WebView2 options
├── wails.json              # Wails project config, metadata, and packaging flags
├── icon.png                # Master high-resolution application icon
├── internal/
│   ├── auth/               # Loopback authentication server & token_v2 bridge (port 28795)
│   ├── config/             # Window geometry state and session token persistence
│   ├── navigation/         # URL routing policy, authentication domain whitelist, and link delegation
│   ├── profile/            # Cross-platform persistent profile directory resolver
│   └── script/             # Injected client script (dynamic theme observer, custom titlebar, shortcuts)
├── build/
│   ├── appicon.png         # Resampled 512x512 application icon
│   ├── windows/            # Windows manifest, info.json metadata, and multi-resolution icon.ico
│   │   └── installer/      # Nullsoft Scriptable Install System (NSIS) setup script (project.nsi)
│   └── bin/                # Compiled binary outputs and installer executables
└── frontend/               # Local asset server directory (fallback splash screen and loader)
```

---

## Security & Privacy Architecture

* **Isolated Outbound Traffic:** Outbound links that do not belong to Notion or authorized identity providers are blocked from executing in the internal webview and delegated to the default system browser.
* **Local Session Storage:** Session cookies and local storage are persisted strictly on the user's local machine inside the user-owned config/profile directory.
* **Zero Telemetry Proxying:** No proxy servers or intermediate logging services are used; all network traffic communicates directly with Notion's official endpoints (`*.notion.com` / `*.notion.so`).

---

## License & Disclaimer

Distributed under the [MIT License](LICENSE).

> **Disclaimer:** Gotion is an independent, unofficial desktop client and is not affiliated with, sponsored by, maintained by, or endorsed by Notion Labs, Inc. "Notion" is a registered trademark of Notion Labs, Inc.

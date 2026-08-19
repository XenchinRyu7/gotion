package script

// GetInjectionScript returns the JavaScript code to inject into Notion pages.
// It renders a flat, edge-to-edge title bar that visually belongs to Notion's design system:
// warm monochrome surface, hairline border, modern Lucide-style minimalist vector SVGs,
// flat inline controls with subtle hover states, clean font-scaling zoom, and window controls.
func GetInjectionScript() string {
	return `
(function() {
    var host = window.location.hostname.toLowerCase();
    var isNotionDomain = host === "app.notion.com" || host.endsWith(".notion.com") ||
                         host === "notion.so" || host.endsWith(".notion.so") ||
                         host === "notion.site" || host.endsWith(".notion.site") ||
                         host === "notionusercontent.com" || host.endsWith(".notionusercontent.com");

    // 0. User-Agent compatibility fix for Google OAuth & Notion Web
    try {
        var chromeUA = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36";
        if (navigator.userAgent !== chromeUA) {
            Object.defineProperty(navigator, 'userAgent', { get: function() { return chromeUA; }, configurable: true });
            Object.defineProperty(navigator, 'appVersion', { get: function() { return "5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"; }, configurable: true });
            Object.defineProperty(navigator, 'vendor', { get: function() { return "Google Inc."; }, configurable: true });
        }
    } catch (e) {}

    // 1. Popup & Window.Open Interceptor (Prevents Notion "Browser is blocking pop-ups" warning)
    if (isNotionDomain && !window.__gotion_popup_interceptor__) {
        window.__gotion_popup_interceptor__ = true;

        function isNotionOnlyUrl(rawUrl) {
            if (!rawUrl) return false;
            try {
                var u = new URL(rawUrl, window.location.href);
                var h = u.hostname.toLowerCase();
                if (h === "app.notion.com" || h.endsWith(".notion.com") ||
                    h === "notion.so" || h.endsWith(".notion.so") ||
                    h === "notion.site" || h.endsWith(".notion.site") ||
                    h === "notionusercontent.com" || h.endsWith(".notionusercontent.com")) {
                    return true;
                }
            } catch(e) {}
            return false;
        }

        function showGotionToast(msg) {
            var existing = document.getElementById("gotion-auth-toast");
            if (existing) existing.remove();

            var toast = document.createElement("div");
            toast.id = "gotion-auth-toast";
            toast.innerHTML = '<span style="font-size:15px;margin-right:8px;">🌐</span><span>' + msg + '</span>';
            toast.setAttribute("style", "position:fixed;bottom:24px;left:50%;transform:translateX(-50%);background:#222;color:#f0f0f0;padding:10px 18px;border-radius:8px;border:1px solid #3a3a3a;box-shadow:0 8px 30px rgba(0,0,0,0.6);font-size:13px;z-index:2147483647;display:flex;align-items:center;font-family:-apple-system,BlinkMacSystemFont,sans-serif;pointer-events:none;");
            (document.body || document.documentElement).appendChild(toast);
            setTimeout(function() {
                if (toast && toast.parentNode) toast.parentNode.removeChild(toast);
            }, 6000);
        }

        var isNavigating = false;
        function handleTargetNavigation(targetUrl) {
            if (!targetUrl || targetUrl === "about:blank") return;
            if (isNavigating) return;
            isNavigating = true;
            setTimeout(function() { isNavigating = false; }, 1000);

            if (isNotionOnlyUrl(targetUrl)) {
                window.location.href = targetUrl;
            } else {
                showGotionToast("Login opened in browser popup. (Or use 'Continue with email' to sign in directly here)");
                if (window.go && window.go.main && window.go.main.App && window.go.main.App.OpenExternalURL) {
                    window.go.main.App.OpenExternalURL(targetUrl);
                } else if (window.runtime && window.runtime.BrowserOpenURL) {
                    window.runtime.BrowserOpenURL(targetUrl);
                }
            }
        }

        function createMockWindow(initialUrl) {
            var mock = {
                closed: false,
                name: "",
                opener: window,
                parent: window,
                top: window,
                focus: function() {},
                blur: function() {},
                close: function() { this.closed = true; },
                postMessage: function(msg, origin) {
                    window.postMessage(msg, origin || "*");
                },
                location: {
                    _href: initialUrl || "",
                    get href() { return this._href; },
                    set href(val) {
                        this._href = val;
                        handleTargetNavigation(val);
                    },
                    replace: function(val) {
                        this.href = val;
                    },
                    assign: function(val) {
                        this.href = val;
                    },
                    toString: function() { return this._href; }
                },
                document: {
                    write: function() {},
                    writeln: function() {},
                    open: function() {},
                    close: function() {}
                }
            };
            if (initialUrl && initialUrl !== "about:blank") {
                handleTargetNavigation(initialUrl);
            }
            return mock;
        }

        window.open = function(url, target, features) {
            return createMockWindow(url);
        };
    }

    // Only inject Notion titlebar and styles if we are on a Notion domain
    if (!isNotionDomain) {
        return;
    }

    var styleId = "gotion-notion-styles";
    var existingStyle = document.getElementById(styleId);
    if (!existingStyle) {
        var style = document.createElement("style");
        style.id = styleId;
        style.textContent = [
            "/* Clean Sidebar: Hide Notion Desktop download promo banners */",
            "a[href*='/desktop'],",
            "a[href*='notion.so/desktop'] {",
            "    display: none !important;",
            "}",
            "/* Fix Notion database table cell font clipping & cramped row height on Linux WebKitGTK */",
            ".notion-table-view,",
            ".notion-collection-item,",
            ".notion-table-view-cell,",
            "[data-block-id] {",
            "    font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Liberation Sans', sans-serif !important;",
            "}",
            ".notion-table-view-cell {",
            "    overflow: visible !important;",
            "}",
            ".notion-table-view-cell > div,",
            ".notion-table-view-cell [contenteditable] {",
            "    line-height: 1.35 !important;",
            "    min-height: 24px !important;",
            "    display: flex !important;",
            "    align-items: center !important;",
            "}",
            ".notion-table-view-row {",
            "    min-height: 34px !important;",
            "}",
            "#gotion-mac-titlebar, #gotion-mac-titlebar.gotion-theme-light {",
            "    --gt-bg: #f6f5f4 !important;",
            "    --gt-border: #e6e6e6 !important;",
            "    --gt-text: #31302e !important;",
            "    --gt-muted: #615d59 !important;",
            "    --gt-faint: #9b9a97 !important;",
            "    --gt-hover: rgba(55, 53, 47, 0.08) !important;",
            "    --gt-active: rgba(55, 53, 47, 0.14) !important;",
            "    --gt-popover-bg: #ffffff !important;",
            "    --gt-popover-shadow: 0 4px 16px rgba(0, 0, 0, 0.08) !important;",
            "    --gt-popover-border: #e6e6e6 !important;",
            "    --gt-traffic-border-close: #e0443e !important;",
            "    --gt-traffic-border-min: #dea123 !important;",
            "    --gt-traffic-border-max: #1aab29 !important;",
            "    --gt-blue: #0075de !important;",
            "}",
            "#gotion-mac-titlebar.gotion-theme-dark {",
            "    --gt-bg: #191919 !important;",
            "    --gt-border: #2c2c2c !important;",
            "    --gt-text: #e6e6e6 !important;",
            "    --gt-muted: #999999 !important;",
            "    --gt-faint: #6e6e6e !important;",
            "    --gt-hover: rgba(255, 255, 255, 0.08) !important;",
            "    --gt-active: rgba(255, 255, 255, 0.14) !important;",
            "    --gt-popover-bg: #222222 !important;",
            "    --gt-popover-shadow: 0 4px 20px rgba(0, 0, 0, 0.5) !important;",
            "    --gt-popover-border: #2c2c2c !important;",
            "    --gt-traffic-border-close: #b03a35 !important;",
            "    --gt-traffic-border-min: #b07f18 !important;",
            "    --gt-traffic-border-max: #167a21 !important;",
            "    --gt-blue: #0075de !important;",
            "}",
             "#gotion-mac-titlebar {",
            "    position: fixed !important;",
            "    top: 0 !important;",
            "    left: 0 !important;",
            "    width: 100vw !important;",
            "    height: 38px !important;",
            "    background-color: var(--gt-bg) !important;",
            "    color: var(--gt-text) !important;",
            "    display: flex !important;",
            "    align-items: center !important;",
            "    justify-content: space-between !important;",
            "    padding: 0 12px !important;",
            "    z-index: 2147483640 !important;",
            "    border-bottom: 1px solid var(--gt-border) !important;",
            "    user-select: none !important;",
            "    -webkit-user-select: none !important;",
            "    --wails-draggable: drag !important;",
            "    box-sizing: border-box !important;",
            "    font-family: Inter, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, sans-serif !important;",
            "    cursor: default !important;",
             "    transition: background-color 0.15s ease, border-color 0.15s ease !important;",
             "}",
             "#gotion-mac-titlebar.gotion-titlebar-windows {",
             "    justify-content: flex-start !important;",
             "}",
             "#gotion-mac-titlebar.gotion-titlebar-windows .gotion-left-group {",
             "    order: 3 !important;",
             "    margin-left: auto !important;",
             "    gap: 0 !important;",
             "}",
             "#gotion-mac-titlebar.gotion-titlebar-windows .gotion-traffic-lights {",
             "    flex-direction: row !important;",
             "    gap: 0 !important;",
             "    padding: 0 !important;",
             "}",
             "#gotion-mac-titlebar.gotion-titlebar-windows .gotion-btn-min { order: 1 !important; }",
             "#gotion-mac-titlebar.gotion-titlebar-windows .gotion-btn-max { order: 2 !important; }",
             "#gotion-mac-titlebar.gotion-titlebar-windows .gotion-btn-close { order: 3 !important; }",
             "#gotion-mac-titlebar.gotion-titlebar-windows .gotion-btn-traffic {",
             "    width: 46px !important;",
             "    height: 38px !important;",
             "    border: none !important;",
             "    border-radius: 0 !important;",
             "    background: transparent !important;",
             "    color: var(--gt-text) !important;",
             "}",
             "#gotion-mac-titlebar.gotion-titlebar-windows .gotion-btn-traffic:hover {",
             "    background: var(--gt-hover) !important;",
             "    filter: none !important;",
             "}",
             "#gotion-mac-titlebar.gotion-titlebar-windows .gotion-btn-close:hover {",
             "    background: #c42b1c !important;",
             "    color: #ffffff !important;",
             "}",
             "#gotion-mac-titlebar.gotion-titlebar-windows .gotion-btn-traffic svg {",
             "    opacity: 1 !important;",
             "    width: 10px !important;",
             "    height: 10px !important;",
             "}",
             "#gotion-mac-titlebar.gotion-titlebar-windows .gotion-btn-close {",
             "    background: transparent !important;",
             "}",
             "#gotion-mac-titlebar.gotion-titlebar-windows .gotion-btn-min,",
             "#gotion-mac-titlebar.gotion-titlebar-windows .gotion-btn-max {",
             "    background: transparent !important;",
             "}",
             "#gotion-mac-titlebar.gotion-titlebar-windows .gotion-titlebar-center {",
             "    order: 2 !important;",
             "}",
             "#gotion-mac-titlebar.gotion-titlebar-windows .gotion-titlebar-right {",
             "    order: 1 !important;",
             "    margin-left: 12px !important;",
             "}",
             "#gotion-mac-titlebar.gotion-titlebar-windows .gotion-btn-close svg line { stroke: currentColor !important; }",
             "#gotion-mac-titlebar.gotion-titlebar-windows .gotion-btn-min svg line { stroke: currentColor !important; }",
             "#gotion-mac-titlebar.gotion-titlebar-windows .gotion-btn-max svg polygon { fill: currentColor !important; }",
            "#notion-app {",
            "    position: absolute !important;",
            "    top: 38px !important;",
            "    left: 0 !important;",
            "    right: 0 !important;",
            "    bottom: 0 !important;",
            "    height: calc(100vh - 38px) !important;",
            "    width: 100vw !important;",
            "}",
            ".gotion-left-group {",
            "    display: flex !important;",
            "    align-items: center !important;",
            "    gap: 12px !important;",
            "    --wails-draggable: no-drag !important;",
            "}",
            ".gotion-traffic-lights {",
            "    display: flex !important;",
            "    align-items: center !important;",
            "    gap: 8px !important;",
            "    padding: 4px 2px !important;",
            "    --wails-draggable: no-drag !important;",
            "}",
            ".gotion-btn-traffic {",
            "    width: 11px !important;",
            "    height: 11px !important;",
            "    border-radius: 50% !important;",
            "    border: none !important;",
            "    cursor: pointer !important;",
            "    padding: 0 !important;",
            "    margin: 0 !important;",
            "    display: flex !important;",
            "    align-items: center !important;",
            "    justify-content: center !important;",
            "    outline: none !important;",
            "    transition: transform 0.1s ease, filter 0.15s ease !important;",
            "}",
            ".gotion-btn-traffic:active {",
            "    transform: scale(0.9) !important;",
            "}",
            ".gotion-btn-close {",
            "    background-color: #ff5f56 !important;",
            "    border: 1px solid var(--gt-traffic-border-close) !important;",
            "}",
            ".gotion-btn-min {",
            "    background-color: #ffbd2e !important;",
            "    border: 1px solid var(--gt-traffic-border-min) !important;",
            "}",
            ".gotion-btn-max {",
            "    background-color: #27c93f !important;",
            "    border: 1px solid var(--gt-traffic-border-max) !important;",
            "}",
            ".gotion-traffic-lights .gotion-btn-traffic svg {",
            "    opacity: 0 !important;",
            "    transition: opacity 0.15s ease !important;",
            "}",
            ".gotion-traffic-lights:hover .gotion-btn-traffic svg {",
            "    opacity: 0.8 !important;",
            "}",
            ".gotion-titlebar-center {",
            "    position: absolute !important;",
            "    left: 50% !important;",
            "    transform: translateX(-50%) !important;",
            "    display: flex !important;",
            "    align-items: center !important;",
            "    gap: 6px !important;",
            "    pointer-events: none !important;",
            "}",
            ".gotion-app-title {",
            "    font-size: 13px !important;",
            "    font-weight: 500 !important;",
            "    color: var(--gt-text) !important;",
            "    letter-spacing: -0.1px !important;",
            "}",
            "/* Flat inline right controls - No Card, No Pill container */",
            ".gotion-titlebar-right {",
            "    display: flex !important;",
            "    align-items: center !important;",
            "    gap: 3px !important;",
            "    --wails-draggable: no-drag !important;",
            "}",
            ".gotion-nav-btn {",
            "    width: 28px !important;",
            "    height: 28px !important;",
            "    border-radius: 5px !important;",
            "    background: transparent !important;",
            "    border: none !important;",
            "    color: var(--gt-muted) !important;",
            "    cursor: pointer !important;",
            "    display: flex !important;",
            "    align-items: center !important;",
            "    justify-content: center !important;",
            "    transition: background-color 0.12s ease, color 0.12s ease !important;",
            "    padding: 0 !important;",
            "    outline: none !important;",
            "}",
            ".gotion-nav-btn:hover {",
            "    background-color: var(--gt-hover) !important;",
            "    color: var(--gt-text) !important;",
            "}",
            ".gotion-nav-btn:active {",
            "    background-color: var(--gt-active) !important;",
            "}",
            ".gotion-dropdown {",
            "    position: relative !important;",
            "    display: inline-flex !important;",
            "    align-items: center !important;",
            "    --wails-draggable: no-drag !important;",
            "}",
            ".gotion-dropdown-btn {",
            "    height: 28px !important;",
            "    padding: 0 8px !important;",
            "    border-radius: 5px !important;",
            "    background: transparent !important;",
            "    border: none !important;",
            "    color: var(--gt-text) !important;",
            "    font-size: 12px !important;",
            "    font-weight: 500 !important;",
            "    cursor: pointer !important;",
            "    display: flex !important;",
            "    align-items: center !important;",
            "    gap: 4px !important;",
            "    outline: none !important;",
            "    transition: background-color 0.12s ease !important;",
            "}",
            ".gotion-dropdown-btn:hover {",
            "    background-color: var(--gt-hover) !important;",
            "}",
            ".gotion-dropdown-btn:active {",
            "    background-color: var(--gt-active) !important;",
            "}",
            ".gotion-dropdown-content {",
            "    display: none !important;",
            "    position: absolute !important;",
            "    right: 0 !important;",
            "    top: 32px !important;",
            "    background-color: var(--gt-popover-bg) !important;",
            "    min-width: 190px !important;",
            "    box-shadow: var(--gt-popover-shadow) !important;",
            "    border: 1px solid var(--gt-popover-border) !important;",
            "    border-radius: 6px !important;",
            "    z-index: 2147483647 !important;",
            "    padding: 4px !important;",
            "    animation: gotionFadeIn 0.12s ease-out !important;",
            "}",
            "@keyframes gotionFadeIn {",
            "    from { opacity: 0; transform: translateY(-3px); }",
            "    to { opacity: 1; transform: translateY(0); }",
            "}",
            ".gotion-dropdown.open .gotion-dropdown-content {",
            "    display: block !important;",
            "}",
            ".gotion-menu-item {",
            "    display: flex !important;",
            "    align-items: center !important;",
            "    justify-content: space-between !important;",
            "    width: 100% !important;",
            "    padding: 6px 8px !important;",
            "    background: transparent !important;",
            "    border: none !important;",
            "    color: var(--gt-text) !important;",
            "    font-size: 12px !important;",
            "    font-weight: 400 !important;",
            "    border-radius: 4px !important;",
            "    text-align: left !important;",
            "    cursor: pointer !important;",
            "    outline: none !important;",
            "    box-sizing: border-box !important;",
            "    transition: background-color 0.1s ease !important;",
            "}",
            ".gotion-menu-item-left {",
            "    display: flex !important;",
            "    align-items: center !important;",
            "    gap: 8px !important;",
            "}",
            ".gotion-menu-item:hover {",
            "    background-color: var(--gt-hover) !important;",
            "}",
            ".gotion-menu-item:active {",
            "    background-color: var(--gt-active) !important;",
            "}",
            ".gotion-menu-divider {",
            "    height: 1px !important;",
            "    background-color: var(--gt-border) !important;",
            "    margin: 4px 0 !important;",
            "}",
            ".gotion-menu-shortcut {",
            "    font-size: 10px !important;",
            "    color: var(--gt-faint) !important;",
            "    font-family: inherit !important;",
            "}",
            "/* Window Resize Handles (Layered on top of titlebar) */",
            ".gotion-resize-handle {",
            "    position: fixed !important;",
            "    z-index: 2147483647 !important;",
            "    background: transparent !important;",
            "    user-select: none !important;",
            "    --wails-draggable: no-drag !important;",
            "}",
            ".gotion-resize-top { top: 0 !important; left: 10px !important; right: 10px !important; height: 6px !important; cursor: n-resize !important; }",
            ".gotion-resize-bottom { bottom: 0 !important; left: 10px !important; right: 10px !important; height: 6px !important; cursor: s-resize !important; }",
            ".gotion-resize-left { top: 10px !important; bottom: 10px !important; left: 0 !important; width: 6px !important; cursor: w-resize !important; }",
            ".gotion-resize-right { top: 10px !important; bottom: 10px !important; right: 0 !important; width: 6px !important; cursor: e-resize !important; }",
            ".gotion-resize-top-left { top: 0 !important; left: 0 !important; width: 10px !important; height: 10px !important; cursor: nw-resize !important; }",
            ".gotion-resize-top-right { top: 0 !important; right: 0 !important; width: 10px !important; height: 10px !important; cursor: ne-resize !important; }",
            ".gotion-resize-bottom-left { bottom: 0 !important; left: 0 !important; width: 10px !important; height: 10px !important; cursor: sw-resize !important; }",
            ".gotion-resize-bottom-right { bottom: 0 !important; right: 0 !important; width: 10px !important; height: 10px !important; cursor: se-resize !important; }"
        ].join("\n");
        (document.head || document.documentElement).appendChild(style);
    }

    function invokeNative(msg) {
        console.log("[Gotion IPC] invokeNative:", msg);
        // 1. WebKitGTK (Linux) & WKWebView (macOS)
        try {
            if (window.webkit && window.webkit.messageHandlers && window.webkit.messageHandlers.external && window.webkit.messageHandlers.external.postMessage) {
                window.webkit.messageHandlers.external.postMessage(msg);
                return;
            }
        } catch(e) {}

        // 2. WebView2 (Windows)
        try {
            if (window.chrome && window.chrome.webview && window.chrome.webview.postMessage) {
                window.chrome.webview.postMessage(msg);
                return;
            }
        } catch(e) {}

        // 3. Wails legacy invoke
        try {
            if (window.WailsInvoke) {
                window.WailsInvoke(msg);
                return;
            }
        } catch(e) {}
    }

    function doClose() {
        console.log("[Gotion JS] doClose called");
        invokeNative("Q");
        if (window.go && window.go.main && window.go.main.App && window.go.main.App.Close) {
            window.go.main.App.Close();
        } else if (window.runtime && window.runtime.Quit) {
            window.runtime.Quit();
        }
    }

    function doMinimise() {
        console.log("[Gotion JS] doMinimise called");
        invokeNative("Wm");
        if (window.go && window.go.main && window.go.main.App && window.go.main.App.Minimise) {
            window.go.main.App.Minimise();
        } else if (window.runtime && window.runtime.WindowMinimise) {
            window.runtime.WindowMinimise();
        }
    }

    function doToggleMaximise() {
        console.log("[Gotion JS] doToggleMaximise called");
        invokeNative("Wt");
        if (window.go && window.go.main && window.go.main.App && window.go.main.App.ToggleMaximise) {
            window.go.main.App.ToggleMaximise();
        } else if (window.runtime && window.runtime.WindowToggleMaximise) {
            window.runtime.WindowToggleMaximise();
        }
    }

    // Modern minimalist vector SVGs (Lucide style, 1.8px stroke, rounded ends)
    var icons = {
        back: '<svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="m15 18-6-6 6-6"/></svg>',
        forward: '<svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="m9 18 6-6-6-6"/></svg>',
        reload: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12a9 9 0 0 0-9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"/><path d="M3 3v5h5"/><path d="M3 12a9 9 0 0 0 9 9 9.75 9.75 0 0 0 6.74-2.74L21 16"/><path d="M16 21h5v-5"/></svg>',
        chevronDown: '<svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="m6 9 6 6 6-6"/></svg>',
        home: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/><polyline points="9 22 9 12 15 12 15 22"/></svg>',
        zoomIn: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/><line x1="11" y1="8" x2="11" y2="14"/><line x1="8" y1="11" x2="14" y2="11"/></svg>',
        zoomOut: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/><line x1="8" y1="11" x2="14" y2="11"/></svg>',
        resetZoom: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"/><path d="M3 3v5h5"/></svg>',
         maximize: '<svg width="12" height="12" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="1" stroke-linejoin="round"><rect x="1.5" y="1.5" width="9" height="9"/></svg>',
         restore: '<svg width="12" height="12" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="1" stroke-linejoin="round"><path d="M3.5 3.5h6v6h-6z"/><path d="M2.5 8.5h-1v-6h6v1"/></svg>',
        user: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>',
        quit: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg>'
    };

    // Zoom Manager (Client-side CSS zoom level)
     var zoomLevels = [0.6, 0.7, 0.8, 0.85, 0.9, 0.95, 1.0, 1.05, 1.1, 1.15, 1.25, 1.35, 1.5, 1.75, 2.0];
    var zoomIndex = 6; // default 1.0

    try {
        var savedZ = localStorage.getItem("gotion_zoom_factor");
        if (savedZ) {
            var f = parseFloat(savedZ);
            for (var zi = 0; zi < zoomLevels.length; zi++) {
                if (Math.abs(zoomLevels[zi] - f) < 0.01) {
                    zoomIndex = zi;
                    break;
                }
            }
        }
    } catch(e) {}

    function applyZoom() {
        var factor = zoomLevels[zoomIndex];
        console.log("[Gotion Zoom] Applying zoom level:", factor);
        
        var target = document.getElementById("notion-app") || document.body;
        if (target) {
            target.style.zoom = factor;
        }

        var bar = document.getElementById("gotion-mac-titlebar");
        if (bar) {
            bar.style.zoom = 1.0;
        }

        showGotionToast("Zoom: " + Math.round(factor * 100) + "%");

        try {
            localStorage.setItem("gotion_zoom_factor", factor.toString());
        } catch(e) {}
    }

    function triggerZoomIn() {
        if (zoomIndex < zoomLevels.length - 1) {
            zoomIndex++;
            applyZoom();
        }
    }

    function triggerZoomOut() {
        if (zoomIndex > 0) {
            zoomIndex--;
            applyZoom();
        }
    }

    function triggerResetZoom() {
        zoomIndex = 6; // 1.0
        applyZoom();
    }

    window.__gotion_triggerZoomIn = triggerZoomIn;
    window.__gotion_triggerZoomOut = triggerZoomOut;
     window.__gotion_triggerResetZoom = triggerResetZoom;

     var titlebarStyle = "mac";
     try {
         var savedTitlebarStyle = localStorage.getItem("gotion_titlebar_style");
         if (savedTitlebarStyle === "windows") titlebarStyle = "windows";
     } catch(e) {}

     function applyTitlebarStyle() {
         var bar = document.getElementById("gotion-mac-titlebar");
         if (!bar) return;
         bar.classList.toggle("gotion-titlebar-windows", titlebarStyle === "windows");
         bar.setAttribute("data-titlebar-style", titlebarStyle);
         var switchItem = document.getElementById("gotion-titlebar-switch");
         if (switchItem) {
             switchItem.querySelector("span").textContent = "Title Bar: " + (titlebarStyle === "windows" ? "Windows" : "Mac");
         }
     }

     function toggleTitlebarStyle() {
         titlebarStyle = titlebarStyle === "windows" ? "mac" : "windows";
         try { localStorage.setItem("gotion_titlebar_style", titlebarStyle); } catch(e) {}
         applyTitlebarStyle();
         showGotionToast("Title bar: " + (titlebarStyle === "windows" ? "Windows" : "Mac"));
     }

     function syncMaximiseIcon() {
         var btn = document.querySelector(".gotion-btn-max");
         if (!btn) return;
         var isMaximised = false;
         try {
             if (window.runtime && window.runtime.WindowIsMaximised) {
                 isMaximised = window.runtime.WindowIsMaximised();
             }
         } catch(e) {}
         if (isMaximised && typeof isMaximised.then === "function") {
             isMaximised.then(function(value) {
                 btn.innerHTML = value ? icons.restore : icons.maximize;
             });
         } else {
             btn.innerHTML = isMaximised ? icons.restore : icons.maximize;
         }
     }

    // Apply saved zoom on startup
    setTimeout(applyZoom, 200);

    // 2. Ensure Titlebar DOM is present
    if (!document.getElementById("gotion-mac-titlebar")) {
        var bar = document.createElement("div");
        bar.id = "gotion-mac-titlebar";
        bar.setAttribute("style", "--wails-draggable: drag;");

        // Double-click detection & native window drag
        var lastTitlebarClickTime = 0;
        bar.addEventListener("mousedown", function(e) {
            if (e.target.closest("button") || e.target.closest(".gotion-dropdown") || e.target.closest(".gotion-no-drag")) {
                return;
            }
            if (e.button !== 0) return;

            var now = Date.now();
            if (now - lastTitlebarClickTime < 350) {
                lastTitlebarClickTime = 0;
                e.preventDefault();
                e.stopPropagation();
                doToggleMaximise();
                return;
            }
            lastTitlebarClickTime = now;

            // Trigger GTK / Windows native window drag
            invokeNative("drag");
        });

        bar.addEventListener("dblclick", function(e) {
            if (e.target.closest("button") || e.target.closest(".gotion-dropdown") || e.target.closest(".gotion-no-drag")) {
                return;
            }
            e.preventDefault();
            e.stopPropagation();
            doToggleMaximise();
        });

        // Left Group (Traffic Lights)
        var leftGroup = document.createElement("div");
        leftGroup.className = "gotion-left-group gotion-no-drag";

         var trafficLights = document.createElement("div");
        trafficLights.className = "gotion-traffic-lights";

        // Close
        var btnClose = document.createElement("button");
         btnClose.className = "gotion-btn-traffic gotion-btn-close";
        btnClose.title = "Close (Ctrl+Q / Alt+F4)";
         btnClose.innerHTML = '<svg width="6" height="6" viewBox="0 0 6 6"><line x1="1" y1="1" x2="5" y2="5" stroke="#4d0000" stroke-width="1.2" stroke-linecap="round"/><line x1="5" y1="1" x2="1" y2="5" stroke="#4d0000" stroke-width="1.2" stroke-linecap="round"/></svg>';
        btnClose.onclick = function(e) {
            e.preventDefault();
            e.stopPropagation();
            doClose();
        };

        // Minimize
        var btnMin = document.createElement("button");
        btnMin.className = "gotion-btn-traffic gotion-btn-min";
        btnMin.title = "Minimize";
        btnMin.innerHTML = '<svg width="6" height="6" viewBox="0 0 6 6"><line x1="1" y1="3" x2="5" y2="3" stroke="#664400" stroke-width="1.2" stroke-linecap="round"/></svg>';
        btnMin.onclick = function(e) {
            e.preventDefault();
            e.stopPropagation();
            doMinimise();
        };

        // Maximize
        var btnMax = document.createElement("button");
        btnMax.className = "gotion-btn-traffic gotion-btn-max";
        btnMax.title = "Maximize / Restore";
         btnMax.innerHTML = icons.maximize;
         btnMax.onclick = function(e) {
            e.preventDefault();
            e.stopPropagation();
             doToggleMaximise();
             setTimeout(syncMaximiseIcon, 100);
         };

        trafficLights.appendChild(btnClose);
        trafficLights.appendChild(btnMin);
        trafficLights.appendChild(btnMax);
        leftGroup.appendChild(trafficLights);

        // Center Title
        var center = document.createElement("div");
        center.className = "gotion-titlebar-center";
        center.innerHTML = '<span class="gotion-app-title">Gotion - Lightweight Notion Client</span>';

        // Right Navigation & Flat Actions Group
        var rightGroup = document.createElement("div");
        rightGroup.className = "gotion-titlebar-right gotion-no-drag";

        var btnBack = document.createElement("button");
        btnBack.className = "gotion-nav-btn";
        btnBack.title = "Back (Alt + ←)";
        btnBack.innerHTML = icons.back;
        btnBack.onclick = function(e) {
            e.stopPropagation();
            window.history.back();
        };

        var btnFwd = document.createElement("button");
        btnFwd.className = "gotion-nav-btn";
        btnFwd.title = "Forward (Alt + →)";
        btnFwd.innerHTML = icons.forward;
        btnFwd.onclick = function(e) {
            e.stopPropagation();
            window.history.forward();
        };

        var btnReload = document.createElement("button");
        btnReload.className = "gotion-nav-btn";
        btnReload.title = "Reload Page (Ctrl + R)";
        btnReload.innerHTML = icons.reload;
        btnReload.onclick = function(e) {
            e.stopPropagation();
            window.location.reload();
        };

        // Dropdown Menu Container (Flat Trigger)
        var dropdown = document.createElement("div");
        dropdown.className = "gotion-dropdown";

        var dropdownBtn = document.createElement("button");
        dropdownBtn.className = "gotion-dropdown-btn";
        dropdownBtn.innerHTML = '<span>Actions</span>' + icons.chevronDown;
        dropdownBtn.onclick = function(e) {
            e.stopPropagation();
            dropdown.classList.toggle("open");
        };

        var dropdownContent = document.createElement("div");
        dropdownContent.className = "gotion-dropdown-content";

        function createMenuItem(iconSvg, label, shortcut, action) {
            var item = document.createElement("button");
            item.className = "gotion-menu-item";
            item.innerHTML = '<div class="gotion-menu-item-left">' + iconSvg + '<span>' + label + '</span></div>' + (shortcut ? '<span class="gotion-menu-shortcut">' + shortcut + '</span>' : '');
            item.onclick = function(e) {
                e.stopPropagation();
                dropdown.classList.remove("open");
                action();
            };
            return item;
        }

        dropdownContent.appendChild(createMenuItem(icons.back, "Go Back", "Alt+←", function() { window.history.back(); }));
        dropdownContent.appendChild(createMenuItem(icons.forward, "Go Forward", "Alt+→", function() { window.history.forward(); }));
        dropdownContent.appendChild(createMenuItem(icons.reload, "Reload Page", "Ctrl+R", function() { window.location.reload(); }));
        dropdownContent.appendChild(createMenuItem(icons.home, "Notion Home", "", function() { window.location.href = "https://app.notion.com"; }));

        var divider1 = document.createElement("div");
        divider1.className = "gotion-menu-divider";
        dropdownContent.appendChild(divider1);

        dropdownContent.appendChild(createMenuItem(icons.zoomIn, "Zoom In (+)", "Ctrl +", function() {
            triggerZoomIn();
        }));
        dropdownContent.appendChild(createMenuItem(icons.zoomOut, "Zoom Out (-)", "Ctrl -", function() {
            triggerZoomOut();
        }));
        dropdownContent.appendChild(createMenuItem(icons.resetZoom, "Reset Zoom", "Ctrl 0", function() {
            triggerResetZoom();
        }));

        var divider2 = document.createElement("div");
        divider2.className = "gotion-menu-divider";
        dropdownContent.appendChild(divider2);

         dropdownContent.appendChild(createMenuItem(icons.maximize, "Maximize / Restore", "F11", function() { doToggleMaximise(); }));
         var titlebarSwitch = createMenuItem(icons.maximize, "Title Bar: Mac", "", toggleTitlebarStyle);
         titlebarSwitch.id = "gotion-titlebar-switch";
         dropdownContent.appendChild(titlebarSwitch);
        dropdownContent.appendChild(createMenuItem(icons.user, "Switch Account / Login", "", function() {
            if (window.go && window.go.main && window.go.main.App && window.go.main.App.Logout) {
                window.go.main.App.Logout();
            } else if (window.runtime && window.runtime.WindowReloadApp) {
                window.runtime.WindowReloadApp();
            }
        }));
        dropdownContent.appendChild(createMenuItem(icons.quit, "Quit Gotion", "Alt+F4", function() { doClose(); }));

        dropdown.appendChild(dropdownBtn);
        dropdown.appendChild(dropdownContent);

        // Close dropdown when clicking outside
        document.addEventListener("click", function(e) {
            if (!dropdown.contains(e.target)) {
                dropdown.classList.remove("open");
            }
        });

        rightGroup.appendChild(btnBack);
        rightGroup.appendChild(btnFwd);
        rightGroup.appendChild(btnReload);
        rightGroup.appendChild(dropdown);

        bar.appendChild(leftGroup);
        bar.appendChild(center);
        bar.appendChild(rightGroup);

         document.documentElement.appendChild(bar);
         applyTitlebarStyle();
         syncMaximiseIcon();
     }

     applyTitlebarStyle();
     syncMaximiseIcon();

    // 3. Inject 8 Native Window Resize Handles
    var resizeHandles = [
        { cls: "gotion-resize-top", edge: "n-resize" },
        { cls: "gotion-resize-bottom", edge: "s-resize" },
        { cls: "gotion-resize-left", edge: "w-resize" },
        { cls: "gotion-resize-right", edge: "e-resize" },
        { cls: "gotion-resize-top-left", edge: "nw-resize" },
        { cls: "gotion-resize-top-right", edge: "ne-resize" },
        { cls: "gotion-resize-bottom-left", edge: "sw-resize" },
        { cls: "gotion-resize-bottom-right", edge: "se-resize" }
    ];

    resizeHandles.forEach(function(h) {
        if (!document.querySelector("." + h.cls)) {
            var handle = document.createElement("div");
            handle.className = "gotion-resize-handle " + h.cls;
            handle.addEventListener("mousedown", function(e) {
                if (e.button === 0) {
                    e.preventDefault();
                    e.stopPropagation();
                    invokeNative("resize:" + h.edge);
                }
            });
            document.documentElement.appendChild(handle);
        }
    });

    // 4. Global Shortcuts & Link Delegation
    if (!window.__gotion_global_shortcuts__) {
        window.__gotion_global_shortcuts__ = true;

        window.addEventListener("keydown", function(e) {
            // Zoom shortcuts: Ctrl + '=', Ctrl + '+', Ctrl + '-'
            if ((e.ctrlKey || e.metaKey) && (e.key === "=" || e.key === "+" || e.code === "Equal" || e.code === "NumpadAdd")) {
                e.preventDefault();
                e.stopImmediatePropagation();
                triggerZoomIn();
                return;
            }
            if ((e.ctrlKey || e.metaKey) && (e.key === "-" || e.key === "_" || e.code === "Minus" || e.code === "NumpadSubtract")) {
                e.preventDefault();
                e.stopImmediatePropagation();
                triggerZoomOut();
                return;
            }
            if ((e.ctrlKey || e.metaKey) && (e.key === "0" || e.code === "Digit0" || e.code === "Numpad0")) {
                e.preventDefault();
                e.stopImmediatePropagation();
                triggerResetZoom();
                return;
            }

            if (e.altKey && (e.key === "ArrowLeft" || e.code === "ArrowLeft")) {
                e.preventDefault();
                window.history.back();
                return;
            }
            if (e.altKey && (e.key === "ArrowRight" || e.code === "ArrowRight")) {
                e.preventDefault();
                window.history.forward();
                return;
            }
            if (((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === "r" && !e.shiftKey) || e.key === "F5") {
                e.preventDefault();
                window.location.reload();
                return;
            }
            if ((e.ctrlKey || e.metaKey) && e.shiftKey && e.key.toLowerCase() === "r") {
                e.preventDefault();
                window.location.href = "https://app.notion.com";
                return;
            }
        }, true);

        window.addEventListener("mouseup", function(e) {
            if (e.button === 3) {
                e.preventDefault();
                window.history.back();
            } else if (e.button === 4) {
                e.preventDefault();
                window.history.forward();
            }
        }, true);

        document.addEventListener("click", function(e) {
            var target = e.target.closest("a");
            if (!target || !target.href) return;
            var href = target.href;
            if (!href.startsWith("http://") && !href.startsWith("https://")) return;

            var isInternal = false;
            try {
                var url = new URL(href);
                var host = url.hostname.toLowerCase();
                if (host === "app.notion.com" || host.endsWith(".notion.com") ||
                    host === "notion.so" || host.endsWith(".notion.so") ||
                    host === "notion.site" || host.endsWith(".notion.site") ||
                    host === "notionusercontent.com" || host.endsWith(".notionusercontent.com")) {
                    isInternal = true;
                }
            } catch (err) {}

            if (!isInternal) {
                e.preventDefault();
                e.stopPropagation();
                if (window.go && window.go.main && window.go.main.App && window.go.main.App.OpenExternalURL) {
                    window.go.main.App.OpenExternalURL(href);
                } else if (window.runtime && window.runtime.BrowserOpenURL) {
                    window.runtime.BrowserOpenURL(href);
                } else {
                    window.open(href, "_blank");
                }
            }
        }, true);
    }

    // 5. Real-time Dynamic Theme Synchronization (Multi-tier Luminance & Class Sampler)
    function detectNotionDarkTheme() {
        var html = document.documentElement;
        var body = document.body;
        var app = document.getElementById("notion-app");
        var appInner = document.querySelector(".notion-app-inner");

        // Layer 1: Direct class & data-theme inspection
        var targets = [html, body, app, appInner].filter(Boolean);
        for (var i = 0; i < targets.length; i++) {
            var el = targets[i];
            if (el.classList.contains("dark") || el.classList.contains("notion-dark-theme") || el.getAttribute("data-theme") === "dark") {
                return true;
            }
            if (el.classList.contains("light") || el.classList.contains("notion-light-theme") || el.getAttribute("data-theme") === "light") {
                return false;
            }
        }

        // Layer 2: Notion CSS custom properties check
        for (var j = 0; j < targets.length; j++) {
            try {
                var cs = window.getComputedStyle(targets[j]);
                var themeBg = cs.getPropertyValue("--theme--bg") || cs.getPropertyValue("--color-bg-default") || cs.getPropertyValue("--bg-color");
                if (themeBg) {
                    themeBg = themeBg.trim().toLowerCase();
                    if (themeBg === "#191919" || themeBg === "#202020" || themeBg.indexOf("rgb(25") !== -1 || themeBg.indexOf("rgb(32") !== -1) {
                        return true;
                    }
                    if (themeBg === "#ffffff" || themeBg === "#fff" || themeBg.indexOf("rgb(255") !== -1 || themeBg.indexOf("rgb(246") !== -1) {
                        return false;
                    }
                }
            } catch (err) {}
        }

        // Layer 3: Direct RGB Luminance Sampling on active Notion workspace containers
        var sampleSelectors = [
            ".notion-frame",
            ".notion-page-content",
            ".notion-scroller",
            ".notion-sidebar",
            ".notion-topbar",
            ".notion-app-inner",
            "#notion-app"
        ];
        for (var k = 0; k < sampleSelectors.length; k++) {
            try {
                var elem = document.querySelector(sampleSelectors[k]);
                if (elem) {
                    var bg = window.getComputedStyle(elem).backgroundColor;
                    if (bg && bg !== "transparent" && bg !== "rgba(0, 0, 0, 0)") {
                        var rgb = bg.match(/\d+/g);
                        if (rgb && rgb.length >= 3) {
                            var r = parseInt(rgb[0], 10);
                            var g = parseInt(rgb[1], 10);
                            var b = parseInt(rgb[2], 10);
                            // Standard ITU-R BT.601 perceptual luminance formula
                            var lum = (r * 299 + g * 587 + b * 114) / 1000;
                            return lum < 128;
                        }
                    }
                }
            } catch (err) {}
        }

        // Layer 4: Default fallback to OS preference
        return window.matchMedia && window.matchMedia("(prefers-color-scheme: dark)").matches;
    }

    function syncGotionTheme() {
        var bar = document.getElementById("gotion-mac-titlebar");
        if (!bar) return;

        var isDark = detectNotionDarkTheme();
        if (isDark) {
            bar.classList.add("gotion-theme-dark");
            bar.classList.remove("gotion-theme-light");
        } else {
            bar.classList.add("gotion-theme-light");
            bar.classList.remove("gotion-theme-dark");
        }
    }

    if (!window.__gotion_theme_observer__) {
        window.__gotion_theme_observer__ = true;
        if (window.MutationObserver && document.body) {
            var themeObserver = new MutationObserver(function() {
                syncGotionTheme();
            });
            themeObserver.observe(document.body, { attributes: true, attributeFilter: ["class", "data-theme"] });
            if (document.documentElement) {
                themeObserver.observe(document.documentElement, { attributes: true, attributeFilter: ["class", "data-theme"] });
            }
        }
        if (window.matchMedia) {
            window.matchMedia("(prefers-color-scheme: dark)").addEventListener("change", syncGotionTheme);
        }
        setInterval(syncGotionTheme, 500);
    }
    syncGotionTheme();

    // 6. Automatic Titlebar Counter-Scaling (Keeps Titlebar at fixed physical 38px height during Native Zoom)
    if (!window.__gotion_base_dpr__) {
        window.__gotion_base_dpr__ = window.devicePixelRatio || 1;
    }

    function updateTitlebarScale() {
        var baseDPR = window.__gotion_base_dpr__ || 1;
        var currentDPR = window.devicePixelRatio || 1;
        var zoomRatio = currentDPR / baseDPR;

        var scaleStyle = document.getElementById("gotion-scale-style");
        if (!scaleStyle) {
            scaleStyle = document.createElement("style");
            scaleStyle.id = "gotion-scale-style";
            (document.head || document.documentElement).appendChild(scaleStyle);
        }

        if (Math.abs(zoomRatio - 1) < 0.01) {
            scaleStyle.textContent = "";
        } else {
            var inv = 1 / zoomRatio;
            var topOffset = 38 * inv;
            scaleStyle.textContent = [
                "#gotion-mac-titlebar {",
                "    transform: scale(" + inv + ") !important;",
                "    transform-origin: top left !important;",
                "    width: calc(100vw * " + zoomRatio + ") !important;",
                "}",
                "#notion-app {",
                "    top: " + topOffset + "px !important;",
                "    height: calc(100vh - " + topOffset + "px) !important;",
                "}"
            ].join("\n");
        }
    }

    window.addEventListener("resize", updateTitlebarScale);
    updateTitlebarScale();
})();
`
}

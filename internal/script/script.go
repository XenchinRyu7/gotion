package script

// GetInjectionScript returns the JavaScript code to inject into Notion pages.
// It renders a flat, edge-to-edge title bar that visually belongs to Notion's design system:
// warm monochrome surface, hairline border, modern Lucide-style minimalist vector SVGs,
// flat inline controls with subtle hover states, clean font-scaling zoom, and window controls.
func GetInjectionScript() string {
	return `
(function() {
    // 0. Performance & Bloatware Stripper (Save ~50-80MB RAM & CPU cycles)
    if (!window.__gotion_bloatware_stripped__) {
        window.__gotion_bloatware_stripped__ = true;

        // Mock tracking libraries as no-op dummies so Notion doesn't crash or lag
        var noop = function() {};
        var dummyTracker = new Proxy(noop, {
            get: function(target, prop) {
                if (prop === "on" || prop === "off" || prop === "once" || prop === "ready") {
                    return function(cb) { if (typeof cb === 'function') cb(); return dummyTracker; };
                }
                return noop;
            },
            apply: function() { return dummyTracker; }
        });

        try {
            window.Intercom = noop;
            window.analytics = dummyTracker;
            window.datadogRum = dummyTracker;
            window.mixpanel = dummyTracker;
            window.posthog = dummyTracker;
        } catch (e) {}

        // Block heavy tracking and telemetry endpoints
        var blockedEndpoints = [
            "api.segment.io",
            "api.mixpanel.com",
            "browser-intake-datadoghq.com",
            "widget.intercom.io",
            "amplitude.com",
            "api/v3/logUserEvent",
            "api/v3/ping"
        ];

        // Intercept window.fetch
        var originalFetch = window.fetch;
        if (originalFetch) {
            window.fetch = function(input, init) {
                var url = typeof input === "string" ? input : (input && input.url ? input.url : "");
                if (url) {
                    for (var i = 0; i < blockedEndpoints.length; i++) {
                        if (url.indexOf(blockedEndpoints[i]) !== -1) {
                            return Promise.resolve(new Response("{}", { status: 200, headers: { "Content-Type": "application/json" } }));
                        }
                    }
                }
                return originalFetch.apply(this, arguments);
            };
        }

        // Intercept XMLHttpRequest
        var originalXhrOpen = XMLHttpRequest.prototype.open;
        XMLHttpRequest.prototype.open = function(method, url) {
            this.__url = url;
            return originalXhrOpen.apply(this, arguments);
        };
        var originalXhrSend = XMLHttpRequest.prototype.send;
        XMLHttpRequest.prototype.send = function(body) {
            var url = this.__url || "";
            if (url) {
                for (var i = 0; i < blockedEndpoints.length; i++) {
                    if (url.indexOf(blockedEndpoints[i]) !== -1) {
                        return;
                    }
                }
            }
            return originalXhrSend.apply(this, arguments);
        };
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
        if (window.chrome && window.chrome.webview && window.chrome.webview.postMessage) {
            window.chrome.webview.postMessage(msg);
        } else if (window.WailsInvoke) {
            window.WailsInvoke(msg);
        }
    }

    function triggerToggleMaximise() {
        if (window.go && window.go.main && window.go.main.App && window.go.main.App.ToggleMaximise) {
            window.go.main.App.ToggleMaximise();
        } else {
            invokeNative("Wt");
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
        maximize: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><rect width="18" height="18" x="3" y="3" rx="2"/><path d="M9 3v18"/><path d="M15 3v18"/></svg>',
        quit: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg>'
    };

    function emitWailsEvent(name) {
        invokeNative('EE{"name":"' + name + '"}');
    }

    function triggerZoomIn() {
        console.log("[Gotion JS] triggerZoomIn called.");
        emitWailsEvent("zoom:in");
        if (window.go && window.go.main && window.go.main.App && window.go.main.App.ZoomIn) {
            window.go.main.App.ZoomIn();
        }
    }

    function triggerZoomOut() {
        console.log("[Gotion JS] triggerZoomOut called.");
        emitWailsEvent("zoom:out");
        if (window.go && window.go.main && window.go.main.App && window.go.main.App.ZoomOut) {
            window.go.main.App.ZoomOut();
        }
    }

    function triggerResetZoom() {
        console.log("[Gotion JS] triggerResetZoom called.");
        emitWailsEvent("zoom:reset");
        if (window.go && window.go.main && window.go.main.App && window.go.main.App.ResetZoom) {
            window.go.main.App.ResetZoom();
        }
    }

    var lastHeaderClickTime = 0;

    // 2. Ensure Titlebar DOM is present
    if (!document.getElementById("gotion-mac-titlebar")) {
        var bar = document.createElement("div");
        bar.id = "gotion-mac-titlebar";
        bar.setAttribute("style", "--wails-draggable: drag;");

        // Double-click toggles maximize
        bar.addEventListener("dblclick", function(e) {
            if (e.target.closest("button") || e.target.closest(".gotion-dropdown") || e.target.closest(".gotion-no-drag")) {
                return;
            }
            triggerToggleMaximise();
        });

        // Left mouse down initiates drag
        bar.addEventListener("mousedown", function(e) {
            if (e.target.closest("button") || e.target.closest(".gotion-dropdown") || e.target.closest(".gotion-no-drag")) {
                return;
            }
            if (e.button !== 0) return;

            var now = Date.now();
            if (now - lastHeaderClickTime < 350) {
                lastHeaderClickTime = 0;
                triggerToggleMaximise();
                return;
            }
            lastHeaderClickTime = now;
            invokeNative("drag");
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
        btnClose.innerHTML = '<svg width="6" height="6" viewBox="0 0 6 6"><line x1="1" y1="1" x2="5" y2="5" stroke="#4d0000" stroke-width="1" stroke-linecap="round"/><line x1="5" y1="1" x2="1" y2="5" stroke="#4d0000" stroke-width="1" stroke-linecap="round"/></svg>';
        btnClose.onclick = function(e) {
            e.stopPropagation();
            invokeNative("Q");
        };

        // Minimize
        var btnMin = document.createElement("button");
        btnMin.className = "gotion-btn-traffic gotion-btn-min";
        btnMin.title = "Minimize";
        btnMin.innerHTML = '<svg width="6" height="6" viewBox="0 0 6 6"><line x1="1" y1="3" x2="5" y2="3" stroke="#664400" stroke-width="1" stroke-linecap="round"/></svg>';
        btnMin.onclick = function(e) {
            e.stopPropagation();
            invokeNative("Wm");
        };

        // Maximize
        var btnMax = document.createElement("button");
        btnMax.className = "gotion-btn-traffic gotion-btn-max";
        btnMax.title = "Maximize / Restore";
        btnMax.innerHTML = '<svg width="6" height="6" viewBox="0 0 6 6"><polygon points="1,1 5,1 1,5" fill="#004d00"/><polygon points="5,5 1,5 5,1" fill="#004d00"/></svg>';
        btnMax.onclick = function(e) {
            e.stopPropagation();
            triggerToggleMaximise();
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

        dropdownContent.appendChild(createMenuItem(icons.maximize, "Maximize / Restore", "F11", function() { triggerToggleMaximise(); }));
        dropdownContent.appendChild(createMenuItem(icons.quit, "Quit Gotion", "Alt+F4", function() { invokeNative("Q"); }));

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
    }

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
                    host === "notionusercontent.com" || host.endsWith(".notionusercontent.com") ||
                    host === "accounts.google.com" || host === "appleid.apple.com" ||
                    host === "login.microsoftonline.com" || host === "github.com") {
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

    // 6. Clean Sidebar: Strip Notion Desktop App download promo from Notion sidebar
    function cleanSidebarPromos() {
        var sidebar = document.querySelector(".notion-sidebar") || document.querySelector("#notion-app");
        if (!sidebar) return;

        var items = sidebar.querySelectorAll("a, div[role='button'], div");
        for (var i = 0; i < items.length; i++) {
            var el = items[i];
            var text = (el.textContent || "").trim();
            if (text === "Notion Desktop" || text === "Aplikasi Notion" || text.indexOf("Download Notion") !== -1 || (el.href && el.href.indexOf("/desktop") !== -1)) {
                var parent = el.closest("div[role='button']") || el.closest("a") || el;
                if (parent && parent.style.display !== "none") {
                    parent.style.setProperty("display", "none", "important");
                }
            }
        }
    }

    if (!window.__gotion_theme_observer__) {
        window.__gotion_theme_observer__ = true;
        if (window.MutationObserver && document.body) {
            var themeObserver = new MutationObserver(function() {
                syncGotionTheme();
                cleanSidebarPromos();
            });
            themeObserver.observe(document.body, { attributes: true, attributeFilter: ["class", "style", "data-theme"], subtree: true });
            if (document.documentElement) {
                themeObserver.observe(document.documentElement, { attributes: true, attributeFilter: ["class", "style", "data-theme"] });
            }
        }
        if (window.matchMedia) {
            window.matchMedia("(prefers-color-scheme: dark)").addEventListener("change", syncGotionTheme);
        }
        setInterval(function() {
            syncGotionTheme();
            cleanSidebarPromos();
        }, 250);
    }
    syncGotionTheme();
    cleanSidebarPromos();

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

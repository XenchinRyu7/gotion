package auth

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"gotion/internal/config"
)

const DefaultAuthPort = 28795

type AuthServer struct {
	mu        sync.Mutex
	server    *http.Server
	port      int
	onSuccess func(token string)
	running   bool
}

var globalAuthServer = &AuthServer{}

// StartAuthServer starts the local loopback auth bridge server.
func StartAuthServer(onSuccess func(token string)) (string, error) {
	globalAuthServer.mu.Lock()
	defer globalAuthServer.mu.Unlock()

	globalAuthServer.onSuccess = onSuccess

	if globalAuthServer.running {
		return fmt.Sprintf("http://127.0.0.1:%d/login", globalAuthServer.port), nil
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", DefaultAuthPort))
	if err != nil {
		// Fallback to random available port
		listener, err = net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return "", fmt.Errorf("failed to bind auth server: %w", err)
		}
	}

	port := listener.Addr().(*net.TCPAddr).Port
	globalAuthServer.port = port

	mux := http.NewServeMux()
	mux.HandleFunc("/login", handleLoginBridge)
	mux.HandleFunc("/callback", handleCallback)
	mux.HandleFunc("/status", handleStatus)

	server := &http.Server{
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	globalAuthServer.server = server
	globalAuthServer.running = true

	go func() {
		log.Printf("[Gotion Auth] Local bridge server started at http://127.0.0.1:%d", port)
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Printf("[Gotion Auth] Server error: %v", err)
		}
	}()

	return fmt.Sprintf("http://127.0.0.1:%d/login", port), nil
}

// StopAuthServer stops the running auth bridge server.
func StopAuthServer() {
	globalAuthServer.mu.Lock()
	defer globalAuthServer.mu.Unlock()

	if globalAuthServer.running && globalAuthServer.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = globalAuthServer.server.Shutdown(ctx)
		globalAuthServer.running = false
		log.Printf("[Gotion Auth] Local bridge server stopped.")
	}
}

func handleLoginBridge(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Gotion - Notion Login Bridge</title>
    <style>
        :root {
            --bg: #191919;
            --surface: #222222;
            --border: #333333;
            --text: #f0f0f0;
            --muted: #999999;
            --accent: #2383e2;
            --accent-hover: #1c6ec0;
            --success: #2ea043;
        }
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            background-color: var(--bg);
            color: var(--text);
            display: flex;
            align-items: center;
            justify-content: center;
            min-height: 100vh;
            padding: 20px;
        }
        .card {
            background: var(--surface);
            border: 1px solid var(--border);
            border-radius: 12px;
            padding: 32px;
            max-width: 520px;
            width: 100%;
            box-shadow: 0 8px 32px rgba(0,0,0,0.4);
        }
        .logo {
            display: flex;
            align-items: center;
            gap: 12px;
            margin-bottom: 24px;
        }
        .logo svg { border-radius: 8px; }
        h1 { font-size: 20px; font-weight: 600; }
        p { color: var(--muted); font-size: 14px; line-height: 1.5; margin-bottom: 20px; }
        .step-box {
            background: rgba(255,255,255,0.03);
            border: 1px solid var(--border);
            border-radius: 8px;
            padding: 16px;
            margin-bottom: 16px;
        }
        .step-title {
            font-size: 14px;
            font-weight: 600;
            margin-bottom: 8px;
            display: flex;
            align-items: center;
            gap: 8px;
        }
        .step-num {
            background: var(--accent);
            color: white;
            width: 22px;
            height: 22px;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 12px;
        }
        .btn {
            display: inline-flex;
            align-items: center;
            justify-content: center;
            gap: 8px;
            width: 100%;
            padding: 10px 16px;
            border-radius: 6px;
            font-size: 14px;
            font-weight: 500;
            text-decoration: none;
            cursor: pointer;
            border: none;
            transition: all 0.15s ease;
        }
        .btn-primary { background: var(--accent); color: white; }
        .btn-primary:hover { background: var(--accent-hover); }
        .btn-bookmarklet {
            background: #2a2a2a;
            border: 1px dashed #555;
            color: #4da3ff;
            cursor: grab;
            margin-top: 8px;
        }
        .input-group {
            display: flex;
            gap: 8px;
            margin-top: 10px;
        }
        input[type="text"] {
            flex: 1;
            padding: 10px 12px;
            background: #151515;
            border: 1px solid var(--border);
            border-radius: 6px;
            color: var(--text);
            font-size: 13px;
            font-family: monospace;
            outline: none;
        }
        input[type="text"]:focus { border-color: var(--accent); }
        .hint { font-size: 12px; color: var(--muted); margin-top: 6px; }
        .badge { background: #2ea04333; color: #3fb950; font-size: 11px; padding: 2px 6px; border-radius: 4px; }
    </style>
</head>
<body>
    <div class="card">
        <div class="logo">
            <svg width="32" height="32" viewBox="0 0 100 100" fill="none">
                <rect width="100" height="100" rx="22" fill="#333"/>
                <path d="M28 72V28L52 60V28H72V72L48 40V72H28Z" fill="#FFF"/>
            </svg>
            <div>
                <h1>Gotion Notion Login</h1>
                <span class="badge">Desktop Authentication Bridge</span>
            </div>
        </div>

        <p>Follow the steps below to authenticate your Notion account in this browser and sync it seamlessly to Gotion Desktop.</p>

        <div class="step-box">
            <div class="step-title">
                <div class="step-num">1</div>
                <span>Log in to Notion Web</span>
            </div>
            <p style="margin-bottom: 12px; font-size: 13px;">Open Notion in your browser and log in with your Google, Apple, or Email account.</p>
            <a href="https://www.notion.so/login" target="_blank" class="btn btn-primary">Open Notion Login in New Tab ↗</a>
        </div>

        <div class="step-box">
            <div class="step-title">
                <div class="step-num">2</div>
                <span>Sync Session to Gotion</span>
            </div>
            <p style="margin-bottom: 8px; font-size: 13px;">After logging in on Notion, paste your <code>token_v2</code> cookie or submit below:</p>
            <form action="/callback" method="GET">
                <div class="input-group">
                    <input type="text" name="token" placeholder="Paste token_v2 or session token here..." required autofocus />
                    <button type="submit" class="btn btn-primary" style="width: auto;">Send to Gotion</button>
                </div>
            </form>
        </div>
    </div>
</body>
</html>`

	_, _ = w.Write([]byte(html))
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	token := strings.TrimSpace(r.URL.Query().Get("token"))
	if token == "" {
		_ = r.ParseForm()
		token = strings.TrimSpace(r.FormValue("token"))
	}

	if token == "" {
		http.Error(w, "Missing token parameter", http.StatusBadRequest)
		return
	}

	// Clean token
	token = strings.Trim(token, `"'`)
	if err := config.SaveSessionToken(token); err != nil {
		log.Printf("[Gotion Auth] Warning: Failed to save session token: %v", err)
	}

	if globalAuthServer.onSuccess != nil {
		globalAuthServer.onSuccess(token)
	}

	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Login Successful - Gotion</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            background: #191919;
            color: #f0f0f0;
            display: flex;
            align-items: center;
            justify-content: center;
            min-height: 100vh;
            margin: 0;
        }
        .box {
            background: #222;
            border: 1px solid #333;
            border-radius: 12px;
            padding: 36px;
            text-align: center;
            max-width: 420px;
            box-shadow: 0 8px 32px rgba(0,0,0,0.5);
        }
        .icon { font-size: 48px; margin-bottom: 16px; }
        h1 { font-size: 20px; margin-bottom: 12px; }
        p { color: #999; font-size: 14px; line-height: 1.5; }
    </style>
</head>
<body>
    <div class="box">
        <div class="icon">🎉</div>
        <h1>Successfully Authenticated!</h1>
        <p>Gotion Desktop is now logged in and loading your workspace.</p>
        <p style="margin-top: 16px; font-size: 12px; color: #666;">You can safely close this browser tab and return to Gotion.</p>
    </div>
    <script>
        setTimeout(function() {
            window.close();
        }, 4000);
    </script>
</body>
</html>`

	_, _ = w.Write([]byte(html))
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	token := config.LoadSessionToken()
	if token != "" {
		_, _ = w.Write([]byte(`{"authenticated": true}`))
	} else {
		_, _ = w.Write([]byte(`{"authenticated": false}`))
	}
}

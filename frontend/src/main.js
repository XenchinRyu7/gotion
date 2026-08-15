// Gotion Bootstrap & Navigation Handler

const NOTION_APP_URL = "https://app.notion.com";

// Store initial internal Wails URL (handles difference between dev http://localhost:... and build wails://...)
try {
    sessionStorage.setItem("gotion_internal_url", window.location.href);
} catch (e) {}

function closeWindow(e) {
    if (e) e.stopPropagation();
    if (window.runtime && window.runtime.Quit) {
        window.runtime.Quit();
    } else if (window.go && window.go.main && window.go.main.App) {
        window.go.main.App.Close();
    }
}

function minimizeWindow(e) {
    if (e) e.stopPropagation();
    if (window.runtime && window.runtime.WindowMinimise) {
        window.runtime.WindowMinimise();
    } else if (window.go && window.go.main && window.go.main.App) {
        window.go.main.App.Minimise();
    }
}

function maximizeWindow(e) {
    if (e) e.stopPropagation();
    if (window.runtime && window.runtime.WindowToggleMaximise) {
        window.runtime.WindowToggleMaximise();
    } else if (window.go && window.go.main && window.go.main.App) {
        window.go.main.App.ToggleMaximise();
    }
}

function navigateToNotion() {
    const loader = document.getElementById("loader");
    const errorContainer = document.getElementById("error-container");
    const statusText = document.getElementById("status-text");

    if (loader) loader.classList.remove("hidden");
    if (errorContainer) errorContainer.classList.add("hidden");
    if (statusText) statusText.innerText = "Connecting to Notion...";

    if (!navigator.onLine) {
        showError("You appear to be offline. Please check your internet connection.");
        return;
    }

    try {
        window.location.replace(NOTION_APP_URL);
    } catch (e) {
        showError("Failed to connect to Notion: " + e.message);
    }
}

function showError(message) {
    const loader = document.getElementById("loader");
    const errorContainer = document.getElementById("error-container");
    const errorDetails = document.getElementById("error-details");

    if (loader) loader.classList.add("hidden");
    if (errorContainer) errorContainer.classList.remove("hidden");
    if (errorDetails && message) errorDetails.innerText = message;
}

function retryConnection() {
    navigateToNotion();
}

window.addEventListener("online", () => {
    navigateToNotion();
});

window.addEventListener("offline", () => {
    showError("Network connection lost. Please reconnect.");
});

// Run immediate navigation on DOM load
window.addEventListener("DOMContentLoaded", () => {
    setTimeout(navigateToNotion, 60);
});

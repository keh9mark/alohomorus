package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
//go:embed backend/main.bin
var assets embed.FS

const (
	pythonServerURL   = "http://127.0.0.1:8000"
	secretHeader      = "X-App-Token"
	secretHeaderValue = "your_very_secret_token_12345"
	serverTimeout     = 15 * time.Second
)

var (
	pythonProcess *exec.Cmd
	binPath       string
)

func getBinaryPath() (string, error) {
	if binPath != "" {
		return binPath, nil
	}

	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}
	exeDir := filepath.Dir(exePath)
	binPath = filepath.Join(exeDir, "main.bin")
	return binPath, nil
}

func extractBinary() error {
	binData, err := fs.ReadFile(assets, "backend/main.bin")
	if err != nil {
		return fmt.Errorf("failed to read embedded binary: %w", err)
	}

	path, err := getBinaryPath()
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, binData, 0755); err != nil {
		return fmt.Errorf("failed to write binary file: %w", err)
	}

	return nil
}

func isServerReady() bool {
	client := http.Client{Timeout: 1 * time.Second}
	req, err := http.NewRequest("GET", pythonServerURL+"/health", nil)
	if err != nil {
		return false
	}

	// Добавляем секретный заголовок
	req.Header.Set(secretHeader, secretHeaderValue)

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// Считаем сервер готовым, если он отвечает (даже с 403)
	return resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusForbidden
}

func waitForServerReady() bool {
	timeout := time.After(serverTimeout)
	tick := time.NewTicker(500 * time.Millisecond)
	defer tick.Stop()

	for {
		select {
		case <-timeout:
			log.Println("Server startup timed out")
			return false
		case <-tick.C:
			if isServerReady() {
				return true
			}
		}
	}
}

func startPythonServer() error {
	if err := extractBinary(); err != nil {
		return fmt.Errorf("binary extraction failed: %w", err)
	}

	path, err := getBinaryPath()
	if err != nil {
		return fmt.Errorf("failed to get binary path: %w", err)
	}

	pythonProcess = exec.Command(path)

	// Настройка атрибутов процесса
	if runtime.GOOS != "windows" {
		pythonProcess.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	}

	pythonProcess.Stdout = os.Stdout
	pythonProcess.Stderr = os.Stderr

	if err := pythonProcess.Start(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	if !waitForServerReady() {
		stopPythonServer()
		return fmt.Errorf("server didn't become ready in time")
	}

	log.Printf("Python server started with PID %d", pythonProcess.Process.Pid)
	return nil
}

func stopPythonServer() {
	if pythonProcess == nil || pythonProcess.Process == nil {
		return
	}

	pid := pythonProcess.Process.Pid
	log.Printf("Stopping Python server (PID %d)", pid)

	var err error
	if runtime.GOOS == "windows" {
		cmd := exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(pid))
		err = cmd.Run()
	} else {
		err = syscall.Kill(-pid, syscall.SIGKILL)
	}

	if err != nil {
		log.Printf("Error stopping process: %v", err)
	} else {
		log.Println("Python server stopped successfully")
	}

	if path, err := getBinaryPath(); err == nil {
		if removeErr := os.Remove(path); removeErr != nil && !os.IsNotExist(removeErr) {
			log.Printf("Error removing binary file: %v", removeErr)
		}
	}

	pythonProcess = nil
}

func createPythonProxy() http.Handler {
	target, _ := url.Parse(pythonServerURL)
	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Header.Set(secretHeader, secretHeaderValue)
		if len(req.URL.Path) >= 4 && req.URL.Path[:4] == "/api" {
			req.URL.Path = req.URL.Path[4:]
		}
	}

	return proxy
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting application...")

	if err := startPythonServer(); err != nil {
		log.Printf("Warning: %v. Continuing with frontend only...", err)
	}
	defer stopPythonServer()

	pythonProxy := createPythonProxy()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) >= 4 && r.URL.Path[:4] == "/api" {
			pythonProxy.ServeHTTP(w, r)
			return
		}
		http.FileServer(http.FS(assets)).ServeHTTP(w, r)
	})

	err := wails.Run(&options.App{
		Title:  "Secure Wails App",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: handler,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			log.Println("Wails app started")
		},
		OnShutdown: func(ctx context.Context) {
			log.Println("Wails app shutting down")
			stopPythonServer()
		},
		Bind: []interface{}{},
	})

	if err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}

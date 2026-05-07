package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

// Wails needs embedded frontend assets for production builds.
//go:embed all:frontend/dist
var assets embed.FS

type App struct {
	ctx context.Context

	server     *http.Server
	serverOnce sync.Once
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.startHTTPServer()
}

func (a *App) shutdown(ctx context.Context) {
	if a.server == nil {
		return
	}

	stopCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_ = a.server.Shutdown(stopCtx)
}

func (a *App) startHTTPServer() {
	a.serverOnce.Do(func() {
		mux := http.NewServeMux()

		mux.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
			allowLocalDev(w)

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			if r.Method != http.MethodGet {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}

			fmt.Println("GET /api/hello")

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			_ = json.NewEncoder(w).Encode(map[string]string{
				"message": "Привіт із Go backend!!!",
				"time":    time.Now().Format(time.RFC3339),
			})
		})

		a.server = &http.Server{
			Addr:    "127.0.0.1:34116",
			Handler: mux,
		}

		go func() {
			log.Println("Go HTTP backend: http://127.0.0.1:34116")
			if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Printf("HTTP backend error: %v", err)
			}
		}()
	})
}

func allowLocalDev(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "Wails Vue Go Hello",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:  app.startup,
		OnShutdown: app.shutdown,
	})
	if err != nil {
		log.Fatal(err)
	}
}

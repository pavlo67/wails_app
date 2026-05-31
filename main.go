package main

import (
	"context"
	"embed"
	"log"
	"net/http"
	"sync"
	"time"

	"wails-vue-go/backend/appcore"
	"wails-vue-go/backend/logadapter"
	"wails-vue-go/backend/modules/filepanel"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

// Wails needs embedded frontend assets for production builds.
//go:embed all:frontend/dist
var assets embed.FS

type App struct {
	ctx        context.Context
	server     *http.Server
	serverOnce sync.Once
}

func NewApp() *App { return &App{} }

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
		logger := logadapter.New("wails_app")
		registry := appcore.NewRegistry()

		filePanelModule, err := filepanel.NewFromStartup(logger)
		if err != nil {
			log.Fatalf("filepanel init failed: %v", err)
		}
		registry.Register(filePanelModule)

		mux := http.NewServeMux()
		registry.RegisterHTTP(mux)

		a.server = &http.Server{
			Addr:    "127.0.0.1:34116",
			Handler: appcore.WithLocalDevCORS(mux),
		}

		go func() {
			log.Println("Go HTTP backend: http://127.0.0.1:34116")
			if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Printf("HTTP backend error: %v", err)
			}
		}()
	})
}

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "Wails File Panel",
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

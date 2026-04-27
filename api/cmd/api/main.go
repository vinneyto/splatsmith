package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/vinneyto/splatra/api/internal/app"
	"github.com/vinneyto/splatra/api/internal/httpapi"
)

func main() {
	configPath := flag.String("config", "./config/standalone.yaml", "path to YAML config")
	flag.Parse()

	cfg, err := app.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	runtime, err := app.BuildRuntime(cfg)
	if err != nil {
		log.Fatalf("build runtime: %v", err)
	}
	defer func() {
		if err := runtime.Close(); err != nil {
			log.Printf("runtime close error: %v", err)
		}
	}()

	apiModule := httpapi.NewModule(cfg.API, httpapi.Dependencies{
		Mode:                string(runtime.Mode),
		AuthService:         runtime.AuthService,
		JobViewer:           runtime.JobViewer,
		DefaultResultURLTTL: time.Duration(runtime.ResultURLTTL) * time.Second,
	})

	srv := &http.Server{
		Addr:              apiModule.ListenAddr(),
		Handler:           apiModule.Handler(),
		ReadHeaderTimeout: 10 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("http shutdown error: %v", err)
		}
	}()

	log.Printf("splatra api started on %s (mode=%s)", srv.Addr, runtime.Mode)
	log.Printf("docs: http://localhost%s/docs | openapi: http://localhost%s/openapi.json", srv.Addr, srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen and serve: %v", err)
	}

	log.Printf("splatra api stopped")
}

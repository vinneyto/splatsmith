package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/vinneyto/ariadne/api/internal/app"
)

func main() {
	configPath := flag.String("config", "./config/standalone.yaml", "path to YAML config")
	token := flag.String("token", "dev-token", "token for startup auth probe")
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

	identity, err := runtime.AuthService.Authenticate(context.Background(), "Bearer "+*token)
	if err != nil {
		log.Printf("startup auth probe failed (expected for aws stubs): %v", err)
	} else {
		log.Printf("startup auth probe success: user_id=%s email=%s", identity.UserID, identity.Email)
	}

	fmt.Printf("api bootstrap is ready (mode=%s)\n", runtime.Mode)
}
